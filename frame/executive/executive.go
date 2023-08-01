package executive

import (
	"fmt"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/extrinsic"
	"github.com/LimeChain/gosemble/execution/inherent"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/aura"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/frame/system/module"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/primitives/crypto"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Module struct {
	system module.SystemModule
}

func New(systemModule module.SystemModule) Module {
	return Module{
		system: systemModule,
	}
}

// InitializeBlock initialises a block with the given header,
// starting the execution of a particular block.
func (m Module) InitializeBlock(header primitives.Header) {
	log.Trace("init_block")
	m.system.ResetEvents()

	weight := primitives.WeightZero()
	if m.runtimeUpgrade() {
		weight = weight.SaturatingAdd(executeOnRuntimeUpgrade())
	}

	m.system.Initialize(header.Number, header.ParentHash, extractPreRuntimeDigest(header.Digest))

	// TODO: accumulate the weight from all pallets that have on_initialize
	weight = weight.SaturatingAdd(aura.OnInitialize())
	weight = weight.SaturatingAdd(system.DefaultBlockWeights().BaseBlock)
	// use in case of dynamic weight calculation
	m.system.RegisterExtraWeightUnchecked(weight, primitives.NewDispatchClassMandatory())

	m.system.NoteFinishedInitialize()
}

func (m Module) ExecuteBlock(block types.Block) {
	log.Trace(fmt.Sprintf("execute_block %v", block.Header.Number))

	m.InitializeBlock(block.Header)

	m.initialChecks(block)

	crypto.ExtCryptoStartBatchVerify()
	m.executeExtrinsicsWithBookKeeping(block)
	if crypto.ExtCryptoFinishBatchVerify() != 1 {
		log.Critical("Signature verification failed")
	}

	m.finalChecks(&block.Header)
}

// ApplyExtrinsic applies extrinsic outside the block execution function.
//
// This doesn't attempt to validate anything regarding the block, but it builds a list of uxt
// hashes.
func (m Module) ApplyExtrinsic(uxt types.UncheckedExtrinsic) (primitives.DispatchOutcome, primitives.TransactionValidityError) {
	encoded := uxt.Bytes()
	encodedLen := sc.ToCompact(len(encoded))

	log.Trace("apply_extrinsic")

	// Verify that the signature is good.
	xt, err := extrinsic.Unchecked(uxt).Check(primitives.DefaultAccountIdLookup())
	if err != nil {
		return primitives.DispatchOutcome{}, err
	}

	// We don't need to make sure to `note_extrinsic` only after we know it's going to be
	// executed to prevent it from leaking in storage since at this point, it will either
	// execute or panic (and revert storage changes).
	m.system.NoteExtrinsic(encoded)

	// AUDIT: Under no circumstances may this function panic from here onwards.

	// Decode parameters and dispatch
	dispatchInfo := primitives.GetDispatchInfo(xt.Function)
	log.Trace("get_dispatch_info: weight ref time " + dispatchInfo.Weight.RefTime.String())

	unsignedValidator := extrinsic.UnsignedValidatorForChecked{}
	res, err := extrinsic.Checked(xt).Apply(unsignedValidator, &dispatchInfo, encodedLen)
	if err != nil {
		return primitives.DispatchOutcome{}, err
	}

	// Mandatory(inherents) are not allowed to fail.
	//
	// The entire block should be discarded if an inherent fails to apply. Otherwise
	// it may open an attack vector.
	if res.HasError && dispatchInfo.Class.Is(primitives.DispatchClassMandatory) {
		return primitives.DispatchOutcome{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionBadMandatory())
	}

	m.system.NoteAppliedExtrinsic(&res, dispatchInfo)

	if res.HasError {
		return primitives.NewDispatchOutcome(res.Err.Error), nil
	}

	return primitives.NewDispatchOutcome(nil), nil
}

func (m Module) FinalizeBlock() primitives.Header {
	log.Trace("finalize_block")
	m.system.NoteFinishedExtrinsics()
	blockNumber := m.system.Storage.BlockNumber.Get()

	m.idleAndFinalizeHook(blockNumber)

	return m.system.Finalize()
}

// ValidateTransaction checks a given signed transaction for validity. This doesn't execute any
// side-effects; it merely checks whether the transaction would panic if it were included or
// not.
//
// Changes made to storage should be discarded.
func (m Module) ValidateTransaction(source primitives.TransactionSource, uxt types.UncheckedExtrinsic, blockHash primitives.Blake2bHash) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	currentBlockNumber := m.system.Storage.BlockNumber.Get()
	m.system.Initialize(currentBlockNumber+1, blockHash, primitives.Digest{})

	log.Trace("validate_transaction")

	log.Trace("using_encoded")
	encodedLen := sc.ToCompact(len(uxt.Bytes()))

	log.Trace("check")
	xt, err := extrinsic.Unchecked(uxt).Check(primitives.DefaultAccountIdLookup())
	if err != nil {
		return ok, err
	}

	log.Trace("dispatch_info")
	dispatchInfo := primitives.GetDispatchInfo(xt.Function)

	if dispatchInfo.Class.Is(primitives.DispatchClassMandatory) {
		return ok, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionMandatoryValidation())
	}

	log.Trace("validate")
	unsignedValidator := extrinsic.UnsignedValidatorForChecked{}
	return extrinsic.Checked(xt).Validate(unsignedValidator, source, &dispatchInfo, encodedLen)
}

func (m Module) OffchainWorker(header primitives.Header) {
	m.system.Initialize(header.Number, header.ParentHash, header.Digest)

	hash := hashing.Blake256(header.Bytes())
	blockHash := primitives.NewBlake2bHash(sc.BytesToSequenceU8(hash)...)

	m.system.Storage.BlockHash.Put(header.Number, blockHash)

	// TODO:
	/*
		<AllPalletsWithSystem as OffchainWorker<System::BlockNumber>>::offchain_worker(*header.number(),)
	*/
}

func (m Module) idleAndFinalizeHook(blockNumber primitives.BlockNumber) {
	weight := m.system.Storage.BlockWeight.Get()

	maxWeight := system.DefaultBlockWeights().MaxBlock
	remainingWeight := maxWeight.SaturatingSub(weight.Total())

	if remainingWeight.AllGt(primitives.WeightZero()) {
		// TODO: call on_idle hook for each pallet
		usedWeight := onIdle(blockNumber, remainingWeight)
		m.system.RegisterExtraWeightUnchecked(usedWeight, primitives.NewDispatchClassMandatory())
	}

	// Each pallet (babe, grandpa) has its own on_finalize that has to be implemented once it is supported
	// TODO:
	timestamp.OnFinalize()
}

func (m Module) executeExtrinsicsWithBookKeeping(block types.Block) {
	for _, ext := range block.Extrinsics {
		_, err := m.ApplyExtrinsic(ext)
		if err != nil {
			log.Critical(string(err[0].Bytes()))
		}
	}

	m.system.NoteFinishedExtrinsics()

	m.idleAndFinalizeHook(block.Header.Number)
}

func (m Module) initialChecks(block types.Block) {
	log.Trace("initial_checks")

	header := block.Header
	blockNumber := header.Number

	if blockNumber > 0 {
		storageParentHash := m.system.Storage.BlockHash.Get(blockNumber - 1)

		if !reflect.DeepEqual(storageParentHash, header.ParentHash) {
			log.Critical("parent hash should be valid")
		}
	}

	inherentsAreFirst := inherent.EnsureInherentsAreFirst(block)

	if inherentsAreFirst >= 0 {
		log.Critical(fmt.Sprintf("invalid inherent position for extrinsic at index [%d]", inherentsAreFirst))
	}
}

func (m Module) runtimeUpgrade() sc.Bool {
	last := m.system.Storage.LastRuntimeUpgrade.Get()

	if m.system.Constants.Version.SpecVersion > sc.U32(last.SpecVersion.ToBigInt().Int64()) ||
		last.SpecName != m.system.Constants.Version.SpecName {

		current := primitives.LastRuntimeUpgradeInfo{
			SpecVersion: sc.ToCompact(m.system.Constants.Version.SpecVersion),
			SpecName:    m.system.Constants.Version.SpecName,
		}
		m.system.Storage.LastRuntimeUpgrade.Put(current)

		return true
	}

	return false
}

func (m Module) finalChecks(header *primitives.Header) {
	newHeader := m.system.Finalize()

	if len(header.Digest) != len(newHeader.Digest) {
		log.Critical("Number of digest must match the calculated")
	}

	for key, digest := range header.Digest {
		otherDigest := newHeader.Digest[key]
		if !reflect.DeepEqual(digest, otherDigest) {
			log.Critical("digest item must match that calculated")
		}
	}

	if !reflect.DeepEqual(header.StateRoot, newHeader.StateRoot) {
		log.Critical("Storage root must match that calculated")
	}

	if !reflect.DeepEqual(header.ExtrinsicsRoot, newHeader.ExtrinsicsRoot) {
		log.Critical("Transaction trie must be valid")
	}
}

// Execute all `OnRuntimeUpgrade` of this runtime, and return the aggregate weight.
func executeOnRuntimeUpgrade() primitives.Weight {
	// TODO: ex: balances
	// call on_runtime_upgrade hook for all modules that implement it
	return onRuntimeUpgrade()
}

func extractPreRuntimeDigest(digest primitives.Digest) primitives.Digest {
	result := primitives.Digest{}
	for k, v := range digest {
		if k == primitives.DigestTypePreRuntime {
			result[k] = v
		}
	}

	return result
}
