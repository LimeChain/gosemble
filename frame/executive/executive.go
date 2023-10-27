package executive

import (
	"reflect"
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/extrinsic"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/io"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Module interface {
	InitializeBlock(header primitives.Header)
	ExecuteBlock(block primitives.Block)
	ApplyExtrinsic(uxt primitives.UncheckedExtrinsic) (primitives.DispatchOutcome, primitives.TransactionValidityError)
	FinalizeBlock() primitives.Header
	ValidateTransaction(source primitives.TransactionSource, uxt primitives.UncheckedExtrinsic, blockHash primitives.Blake2bHash) (primitives.ValidTransaction, primitives.TransactionValidityError)
	OffchainWorker(header primitives.Header)
}

type module struct {
	system           system.Module
	onRuntimeUpgrade primitives.OnRuntimeUpgrade
	runtimeExtrinsic extrinsic.RuntimeExtrinsic
	hashing          io.Hashing
}

func New(systemModule system.Module, runtimeExtrinsic extrinsic.RuntimeExtrinsic, onRuntimeUpgrade primitives.OnRuntimeUpgrade) Module {
	return module{
		system:           systemModule,
		onRuntimeUpgrade: onRuntimeUpgrade,
		runtimeExtrinsic: runtimeExtrinsic,
		hashing:          io.NewHashing(),
	}
}

// InitializeBlock initialises a block with the given header,
// starting the execution of a particular block.
func (m module) InitializeBlock(header primitives.Header) {
	log.Trace("init_block")
	m.system.ResetEvents()

	weight := primitives.WeightZero()
	if m.runtimeUpgrade() {
		weight = weight.SaturatingAdd(m.executeOnRuntimeUpgrade())
	}

	m.system.Initialize(header.Number, header.ParentHash, extractPreRuntimeDigest(header.Digest))

	weight = weight.SaturatingAdd(m.runtimeExtrinsic.OnInitialize(header.Number))
	weight = weight.SaturatingAdd(m.system.BlockWeights().BaseBlock)
	// use in case of dynamic weight calculation
	m.system.RegisterExtraWeightUnchecked(weight, primitives.NewDispatchClassMandatory())

	m.system.NoteFinishedInitialize()
}

func (m module) ExecuteBlock(block primitives.Block) {
	// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
	// log.Trace(fmt.Sprintf("execute_block %v", block.Header.Number))
	log.Trace("execute_block " + strconv.Itoa(int(block.Header().Number)))

	m.InitializeBlock(block.Header())

	m.initialChecks(block)

	m.executeExtrinsicsWithBookKeeping(block)

	header := block.Header()
	m.finalChecks(&header)
}

// ApplyExtrinsic applies extrinsic outside the block execution function.
//
// This doesn't attempt to validate anything regarding the block, but it builds a list of uxt
// hashes.
func (m module) ApplyExtrinsic(uxt primitives.UncheckedExtrinsic) (primitives.DispatchOutcome, primitives.TransactionValidityError) {
	encoded := uxt.Bytes()
	encodedLen := sc.ToCompact(len(encoded))

	log.Trace("apply_extrinsic")

	// Verify that the signature is good.
	checked, err := uxt.Check(primitives.DefaultAccountIdLookup())
	if err != nil {
		return primitives.DispatchOutcome{}, err
	}

	// We don't need to make sure to `note_extrinsic` only after we know it's going to be
	// executed to prevent it from leaking in storage since at this point, it will either
	// execute or panic (and revert storage changes).
	m.system.NoteExtrinsic(encoded)

	// AUDIT: Under no circumstances may this function panic from here onwards.

	// Decode parameters and dispatch
	dispatchInfo := primitives.GetDispatchInfo(checked.Function())
	log.Trace("get_dispatch_info: weight ref time " + strconv.Itoa(int(dispatchInfo.Weight.RefTime)))

	unsignedValidator := extrinsic.NewUnsignedValidatorForChecked(m.runtimeExtrinsic)
	res, err := checked.Apply(unsignedValidator, &dispatchInfo, encodedLen)
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

func (m module) FinalizeBlock() primitives.Header {
	log.Trace("finalize_block")
	m.system.NoteFinishedExtrinsics()
	blockNumber := m.system.StorageBlockNumber()

	m.idleAndFinalizeHook(blockNumber)

	return m.system.Finalize()
}

// ValidateTransaction checks a given signed transaction for validity. This doesn't execute any
// side-effects; it merely checks whether the transaction would panic if it were included or
// not.
//
// Changes made to storage should be discarded.
func (m module) ValidateTransaction(source primitives.TransactionSource, uxt primitives.UncheckedExtrinsic, blockHash primitives.Blake2bHash) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	log.Trace("validate_transaction")
	currentBlockNumber := m.system.StorageBlockNumber()

	m.system.Initialize(currentBlockNumber+1, blockHash, primitives.Digest{})

	log.Trace("using_encoded")
	encodedLen := sc.ToCompact(len(uxt.Bytes()))

	log.Trace("check")
	checked, err := uxt.Check(primitives.DefaultAccountIdLookup())
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	log.Trace("dispatch_info")
	dispatchInfo := primitives.GetDispatchInfo(checked.Function())

	if dispatchInfo.Class.Is(primitives.DispatchClassMandatory) {
		return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionMandatoryValidation())
	}

	log.Trace("validate")
	unsignedValidator := extrinsic.NewUnsignedValidatorForChecked(m.runtimeExtrinsic)
	return checked.Validate(unsignedValidator, source, &dispatchInfo, encodedLen)
}

func (m module) OffchainWorker(header primitives.Header) {
	m.system.Initialize(header.Number, header.ParentHash, header.Digest)

	hash := m.hashing.Blake256(header.Bytes())
	blockHash := primitives.NewBlake2bHash(sc.BytesToSequenceU8(hash)...)

	m.system.StorageBlockHashSet(header.Number, blockHash)

	m.runtimeExtrinsic.OffchainWorker(header.Number)
}

func (m module) idleAndFinalizeHook(blockNumber sc.U64) {
	weight := m.system.StorageBlockWeight()

	maxWeight := m.system.BlockWeights().MaxBlock
	remainingWeight := maxWeight.SaturatingSub(weight.Total())

	if remainingWeight.AllGt(primitives.WeightZero()) {
		usedWeight := m.runtimeExtrinsic.OnIdle(blockNumber, remainingWeight)
		m.system.RegisterExtraWeightUnchecked(usedWeight, primitives.NewDispatchClassMandatory())
	}

	m.runtimeExtrinsic.OnFinalize(blockNumber)
}

func (m module) executeExtrinsicsWithBookKeeping(block primitives.Block) {
	for _, ext := range block.Extrinsics() {
		_, err := m.ApplyExtrinsic(ext)
		if err != nil {
			log.Critical(string(err[0].Bytes()))
		}
	}

	m.system.NoteFinishedExtrinsics()

	m.idleAndFinalizeHook(block.Header().Number)
}

func (m module) initialChecks(block primitives.Block) {
	log.Trace("initial_checks")

	header := block.Header()
	blockNumber := header.Number

	if blockNumber > 0 {
		storageParentHash := m.system.StorageBlockHash(blockNumber - 1)

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

func (m module) runtimeUpgrade() sc.Bool {
	last := m.system.StorageLastRuntimeUpgrade()

	if m.system.Version().SpecVersion > last.SpecVersion ||
		last.SpecName != m.system.Version().SpecName {

		current := primitives.LastRuntimeUpgradeInfo{
			SpecVersion: m.system.Version().SpecVersion,
			SpecName:    m.system.Version().SpecName,
		}
		m.system.StorageLastRuntimeUpgradeSet(current)

		return true
	}

	return false
}

func (m module) finalChecks(header *primitives.Header) {
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
func (m module) executeOnRuntimeUpgrade() primitives.Weight {
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
