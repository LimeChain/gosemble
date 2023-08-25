package executive

import (
	"reflect"
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/extrinsic"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/hooks"
	"github.com/LimeChain/gosemble/primitives/crypto"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Module[N sc.Numeric] struct {
	system           system.Module[N]
	runtimeExtrinsic extrinsic.RuntimeExtrinsic[N]
	onRuntimeUpgrade hooks.OnRuntimeUpgrade
}

func New[N sc.Numeric](systemModule system.Module[N], runtimeExtrinsic extrinsic.RuntimeExtrinsic[N], onRuntimeUpgrade hooks.OnRuntimeUpgrade) Module[N] {
	return Module[N]{
		system:           systemModule,
		runtimeExtrinsic: runtimeExtrinsic,
		onRuntimeUpgrade: onRuntimeUpgrade,
	}
}

// InitializeBlock initialises a block with the given header,
// starting the execution of a particular block.
func (m Module[N]) InitializeBlock(header primitives.Header[N]) {
	log.Trace("init_block")
	m.system.ResetEvents()

	weight := primitives.WeightZero()
	if m.runtimeUpgrade() {
		weight = weight.SaturatingAdd(m.executeOnRuntimeUpgrade())
	}

	m.system.Initialize(header.Number, header.ParentHash, extractPreRuntimeDigest(header.Digest))

	weight = weight.SaturatingAdd(m.runtimeExtrinsic.OnInitialize(header.Number))
	weight = weight.SaturatingAdd(m.system.Constants.BlockWeights.BaseBlock)
	// use in case of dynamic weight calculation
	m.system.RegisterExtraWeightUnchecked(weight, primitives.NewDispatchClassMandatory())

	m.system.NoteFinishedInitialize()
}

func (m Module[N]) ExecuteBlock(block types.Block[N]) {
	// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
	// log.Trace(fmt.Sprintf("execute_block %v", block.Header.Number))
	log.Trace("execute_block " + strconv.Itoa(int(sc.To[sc.U64](block.Header.Number))))

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
func (m Module[N]) ApplyExtrinsic(uxt types.UncheckedExtrinsic) (primitives.DispatchOutcome, primitives.TransactionValidityError) {
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
	log.Trace("get_dispatch_info: weight ref time " + strconv.Itoa(int(dispatchInfo.Weight.RefTime)))

	unsignedValidator := extrinsic.NewUnsignedValidatorForChecked(m.runtimeExtrinsic)
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

func (m Module[N]) FinalizeBlock() primitives.Header[N] {
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
func (m Module[N]) ValidateTransaction(source primitives.TransactionSource, uxt types.UncheckedExtrinsic, blockHash primitives.Blake2bHash) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	currentBlockNumber := m.system.Storage.BlockNumber.Get()
	blockNumber := currentBlockNumber.Add(sc.NewNumeric[N](1)).(N)
	m.system.Initialize(blockNumber, blockHash, primitives.Digest{})

	log.Trace("validate_transaction")

	log.Trace("using_encoded")
	encodedLen := sc.ToCompact(len(uxt.Bytes()))

	log.Trace("check")
	xt, err := extrinsic.Unchecked(uxt).Check(primitives.DefaultAccountIdLookup())
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	log.Trace("dispatch_info")
	dispatchInfo := primitives.GetDispatchInfo(xt.Function)

	if dispatchInfo.Class.Is(primitives.DispatchClassMandatory) {
		return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionMandatoryValidation())
	}

	log.Trace("validate")
	unsignedValidator := extrinsic.NewUnsignedValidatorForChecked(m.runtimeExtrinsic)
	return extrinsic.Checked(xt).Validate(unsignedValidator, source, &dispatchInfo, encodedLen)
}

func (m Module[N]) OffchainWorker(header primitives.Header[N]) {
	m.system.Initialize(header.Number, header.ParentHash, header.Digest)

	hash := hashing.Blake256(header.Bytes())
	blockHash := primitives.NewBlake2bHash(sc.BytesToSequenceU8(hash)...)

	m.system.Storage.BlockHash.Put(N(header.Number), blockHash)

	m.runtimeExtrinsic.OffchainWorker(header.Number)
}

func (m Module[N]) idleAndFinalizeHook(blockNumber N) {
	weight := m.system.Storage.BlockWeight.Get()

	maxWeight := m.system.Constants.BlockWeights.MaxBlock
	remainingWeight := maxWeight.SaturatingSub(weight.Total())

	if remainingWeight.AllGt(primitives.WeightZero()) {
		usedWeight := m.runtimeExtrinsic.OnIdle(blockNumber, remainingWeight)
		m.system.RegisterExtraWeightUnchecked(usedWeight, primitives.NewDispatchClassMandatory())
	}

	m.runtimeExtrinsic.OnFinalize(blockNumber)
}

func (m Module[N]) executeExtrinsicsWithBookKeeping(block types.Block[N]) {
	for _, ext := range block.Extrinsics {
		_, err := m.ApplyExtrinsic(ext)
		if err != nil {
			log.Critical(string(err[0].Bytes()))
		}
	}

	m.system.NoteFinishedExtrinsics()

	m.idleAndFinalizeHook(block.Header.Number)
}

func (m Module[N]) initialChecks(block types.Block[N]) {
	log.Trace("initial_checks")

	header := block.Header
	blockNumber := header.Number

	if blockNumber.Gt(sc.NewNumeric[N](0)) {
		storageParentHash := m.system.Storage.BlockHash.Get(blockNumber.Sub(sc.NewNumeric[N](1)).(N))

		if !reflect.DeepEqual(storageParentHash, header.ParentHash) {
			log.Critical("parent hash should be valid")
		}
	}

	inherentsAreFirst := m.runtimeExtrinsic.EnsureInherentsAreFirst(block)
	if inherentsAreFirst >= 0 {
		// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
		// log.Critical(fmt.Sprintf("invalid inherent position for extrinsic at index [%d]", inherentsAreFirst))
		log.Critical("invalid inherent position for extrinsic at index " + strconv.Itoa(int(inherentsAreFirst)))
	}
}

func (m Module[N]) runtimeUpgrade() sc.Bool {
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

func (m Module[N]) finalChecks(header *primitives.Header[N]) {
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

// executeOnRuntimeUpgrade - Execute all `OnRuntimeUpgrade` of this runtime, and return the aggregate weight.
func (m Module[N]) executeOnRuntimeUpgrade() primitives.Weight {
	weight := m.onRuntimeUpgrade.OnRuntimeUpgrade()

	return weight.SaturatingAdd(m.runtimeExtrinsic.OnRuntimeUpgrade())
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
