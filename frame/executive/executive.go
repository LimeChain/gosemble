package executive

import (
	"fmt"
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
	InitializeBlock(header primitives.Header) error
	ExecuteBlock(block primitives.Block) error
	ApplyExtrinsic(uxt primitives.UncheckedExtrinsic) error
	FinalizeBlock() (primitives.Header, error)
	ValidateTransaction(source primitives.TransactionSource, uxt primitives.UncheckedExtrinsic, blockHash primitives.Blake2bHash) (primitives.ValidTransaction, error)
	OffchainWorker(header primitives.Header) error
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
func (m module) InitializeBlock(header primitives.Header) error {
	log.Trace("init_block")
	m.system.ResetEvents()

	weight := primitives.WeightZero()
	upgrade, err := m.runtimeUpgrade()
	if err != nil {
		return err
	}
	if upgrade {
		weight = weight.SaturatingAdd(m.executeOnRuntimeUpgrade())
	}

	m.system.Initialize(header.Number, header.ParentHash, extractPreRuntimeDigest(header.Digest))

	onInit, err := m.runtimeExtrinsic.OnInitialize(header.Number)
	if err != nil {
		return err
	}

	weight = weight.SaturatingAdd(onInit)
	weight = weight.SaturatingAdd(m.system.BlockWeights().BaseBlock)
	// use in case of dynamic weight calculation
	err = m.system.RegisterExtraWeightUnchecked(weight, primitives.NewDispatchClassMandatory())
	if err != nil {
		return err
	}

	m.system.NoteFinishedInitialize()
	return nil
}

func (m module) ExecuteBlock(block primitives.Block) error {
	log.Trace(fmt.Sprintf("execute_block %v", block.Header().Number))

	err := m.InitializeBlock(block.Header())
	if err != nil {
		return err
	}

	err = m.initialChecks(block)
	if err != nil {
		return err
	}

	m.executeExtrinsicsWithBookKeeping(block)

	header := block.Header()
	err = m.finalChecks(&header)
	if err != nil {
		return err
	}
	return nil
}

// ApplyExtrinsic applies extrinsic outside the block execution function.
//
// This doesn't attempt to validate anything regarding the block, but it builds a list of uxt
// hashes.
func (m module) ApplyExtrinsic(uxt primitives.UncheckedExtrinsic) error {
	encoded := uxt.Bytes()
	encodedLen := sc.ToCompact(len(encoded))

	log.Trace("apply_extrinsic")

	// Verify that the signature is good.
	checked, err := uxt.Check()
	if err != nil {
		return err
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
		_, isDispatchErr := err.(primitives.DispatchError)
		if !isDispatchErr {
			return err
		}

		// Mandatory(inherents) are not allowed to fail.
		//
		// The entire block should be discarded if an inherent fails to apply. Otherwise
		// it may open an attack vector.
		if isMandatoryDispatch(dispatchInfo) {
			return primitives.NewTransactionValidityError(primitives.NewInvalidTransactionBadMandatory())
		}
	}

	if err := m.system.NoteAppliedExtrinsic(res, err, dispatchInfo); err != nil {
		log.Critical(err.Error())
	}

	return err
}

func (m module) FinalizeBlock() (primitives.Header, error) {
	log.Trace("finalize_block")
	err := m.system.NoteFinishedExtrinsics()
	if err != nil {
		return primitives.Header{}, err
	}
	blockNumber, err := m.system.StorageBlockNumber()
	if err != nil {
		return primitives.Header{}, err
	}

	err = m.idleAndFinalizeHook(blockNumber)
	if err != nil {
		return primitives.Header{}, err
	}

	return m.system.Finalize()
}

// ValidateTransaction checks a given signed transaction for validity. This doesn't execute any
// side-effects; it merely checks whether the transaction would panic if it were included or
// not.
//
// Changes made to storage should be discarded.
func (m module) ValidateTransaction(source primitives.TransactionSource, uxt primitives.UncheckedExtrinsic, blockHash primitives.Blake2bHash) (primitives.ValidTransaction, error) {
	log.Trace("validate_transaction")
	currentBlockNumber, err := m.system.StorageBlockNumber()
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	m.system.Initialize(currentBlockNumber+1, blockHash, primitives.Digest{})

	log.Trace("using_encoded")
	encodedLen := sc.ToCompact(len(uxt.Bytes()))

	log.Trace("check")
	checked, errCheck := uxt.Check()
	if errCheck != nil {
		return primitives.ValidTransaction{}, errCheck
	}

	log.Trace("dispatch_info")
	dispatchInfo := primitives.GetDispatchInfo(checked.Function())
	if isMandatoryDispatch(dispatchInfo) {
		return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionMandatoryValidation())
	}

	log.Trace("validate")
	unsignedValidator := extrinsic.NewUnsignedValidatorForChecked(m.runtimeExtrinsic)
	return checked.Validate(unsignedValidator, source, &dispatchInfo, encodedLen)
}

func (m module) OffchainWorker(header primitives.Header) error {
	m.system.Initialize(header.Number, header.ParentHash, header.Digest)

	hash := m.hashing.Blake256(header.Bytes())
	blockHash, err := primitives.NewBlake2bHash(sc.BytesToSequenceU8(hash)...)
	if err != nil {
		return err
	}

	m.system.StorageBlockHashSet(header.Number, blockHash)

	m.runtimeExtrinsic.OffchainWorker(header.Number)

	return nil
}

func (m module) idleAndFinalizeHook(blockNumber sc.U64) error {
	weight, err := m.system.StorageBlockWeight()
	if err != nil {
		return err
	}

	maxWeight := m.system.BlockWeights().MaxBlock

	total, totalErr := weight.Total()
	if totalErr != nil {
		log.Critical(totalErr.Error())
	}
	remainingWeight := maxWeight.SaturatingSub(total)

	if remainingWeight.AllGt(primitives.WeightZero()) {
		usedWeight := m.runtimeExtrinsic.OnIdle(blockNumber, remainingWeight)
		err = m.system.RegisterExtraWeightUnchecked(usedWeight, primitives.NewDispatchClassMandatory())
		if err != nil {
			return err
		}
	}

	err = m.runtimeExtrinsic.OnFinalize(blockNumber)
	if err != nil {
		return err
	}
	return nil
}

func (m module) executeExtrinsicsWithBookKeeping(block primitives.Block) {
	for _, ext := range block.Extrinsics() {
		if err := m.ApplyExtrinsic(ext); err != nil {
			log.Critical(err.Error())
		}
	}

	m.system.NoteFinishedExtrinsics()

	m.idleAndFinalizeHook(block.Header().Number)
}

func (m module) initialChecks(block primitives.Block) error {
	log.Trace("initial_checks")

	header := block.Header()
	blockNumber := header.Number

	if blockNumber > 0 {
		storageParentHash, err := m.system.StorageBlockHash(blockNumber - 1)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(storageParentHash, header.ParentHash) {
			log.Critical("parent hash should be valid")
		}
	}

	inherentsAreFirst := m.runtimeExtrinsic.EnsureInherentsAreFirst(block)
	if inherentsAreFirst >= 0 {
		log.Critical(fmt.Sprintf("invalid inherent position for extrinsic at index [%d]", inherentsAreFirst))
	}
	return nil
}

func (m module) runtimeUpgrade() (sc.Bool, error) {
	last, err := m.system.StorageLastRuntimeUpgrade()
	if err != nil {
		return false, err
	}

	if m.system.Version().SpecVersion > last.SpecVersion ||
		last.SpecName != m.system.Version().SpecName {

		current := primitives.LastRuntimeUpgradeInfo{
			SpecVersion: m.system.Version().SpecVersion,
			SpecName:    m.system.Version().SpecName,
		}
		m.system.StorageLastRuntimeUpgradeSet(current)

		return true, nil
	}

	return false, nil
}

func (m module) finalChecks(header *primitives.Header) error {
	newHeader, err := m.system.Finalize()
	if err != nil {
		return err
	}

	if len(header.Digest.Sequence) != len(newHeader.Digest.Sequence) {
		log.Critical("Number of digest must match the calculated")
	}

	for i, digest := range header.Digest.Sequence {
		otherDigest := newHeader.Digest.Sequence[i]
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
	return nil
}

// executeOnRuntimeUpgrade - Execute all `OnRuntimeUpgrade` of this runtime, and return the aggregate weight.
func (m module) executeOnRuntimeUpgrade() primitives.Weight {
	weight := m.onRuntimeUpgrade.OnRuntimeUpgrade()

	return weight.SaturatingAdd(m.runtimeExtrinsic.OnRuntimeUpgrade())
}

func extractPreRuntimeDigest(digest primitives.Digest) primitives.Digest {
	return digest.OnlyPreRuntimes()
}

func isMandatoryDispatch(dispatchInfo primitives.DispatchInfo) sc.Bool {
	isMandatory, err := dispatchInfo.Class.Is(primitives.DispatchClassMandatory)
	if err != nil {
		log.Critical(err.Error())
	}
	return isMandatory
}
