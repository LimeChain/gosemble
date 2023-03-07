package executive

import (
	"fmt"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/execution/extrinsic"
	"github.com/LimeChain/gosemble/frame/aura"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/crypto"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
)

// InitializeBlock initialises a block with the given header,
// starting the execution of a particular block.
func InitializeBlock(header types.Header) {
	system.ResetEvents()

	if runtimeUpgrade() {
		// TODO: weight
		/*
			weight = weight.saturating_add(Self::execute_on_runtime_upgrade());
		*/
	}

	system.Initialize(header.Number, header.ParentHash, extractPreRuntimeDigest(header.Digest))

	// TODO: weight + on_initialize
	/*
		weight = weight.saturating_add(<AllPalletsWithSystem as OnInitialize<
					System::BlockNumber,
				>>::on_initialize(*block_number));
				weight = weight.saturating_add(
					<System::BlockWeights as frame_support::traits::Get<_>>::get().base_block,
				);
				<frame_system::Pallet<System>>::register_extra_weight_unchecked(
					weight,
					DispatchClass::Mandatory,
				);
	*/
	aura.OnInitialize()

	system.NoteFinishedInitialize()
}

// ApplyExtrinsic applies extrinsic outside the block execution function.
//
// This doesn't attempt to validate anything regarding the block, but it builds a list of uxt
// hashes.
func ApplyExtrinsic(uxt types.UncheckedExtrinsic) (ok types.DispatchOutcome, err types.TransactionValidityError) { // types.ApplyExtrinsicResult
	encoded := uxt.Bytes()
	encodedLen := sc.ToCompact(len(encoded))

	log.Info("apply_extrinsic")

	// Verify that the signature is good.
	xt, err := extrinsic.Unchecked(uxt).Check(types.DefaultAccountIdLookup())
	if err != nil {
		return ok, err
	}

	// We don't need to make sure to `note_extrinsic` only after we know it's going to be
	// executed to prevent it from leaking in storage since at this point, it will either
	// execute or panic (and revert storage changes).
	system.NoteExtrinsic(encoded)

	// AUDIT: Under no circumstances may this function panic from here onwards.

	// Decode parameters and dispatch
	dispatchInfo := extrinsic.GetDispatchInfo(xt)

	unsignedValidator := extrinsic.UnsignedValidatorForChecked{}
	res, err := extrinsic.Checked(xt).Apply(unsignedValidator, &dispatchInfo, encodedLen)

	// Mandatory(inherents) are not allowed to fail.
	//
	// The entire block should be discarded if an inherent fails to apply. Otherwise
	// it may open an attack vector.
	if res.HasError && dispatchInfo.Class.Is(types.DispatchClassMandatory) {
		return ok, types.NewTransactionValidityError(types.NewInvalidTransactionBadMandatory())
	}

	system.NoteAppliedExtrinsic(&res, dispatchInfo)

	if err != nil {
		return ok, err
	}

	return types.NewDispatchOutcome(nil), err
}

func ExecuteBlock(block types.Block) {
	InitializeBlock(block.Header)

	initialChecks(block)

	crypto.ExtCryptoStartBatchVerify()
	executeExtrinsicsWithBookKeeping(block)
	if crypto.ExtCryptoFinishBatchVerify() != 1 {
		log.Critical("Signature verification failed")
	}

	finalChecks(&block.Header)
}

func executeExtrinsicsWithBookKeeping(block types.Block) {
	for _, ext := range block.Extrinsics {
		_, err := ApplyExtrinsic(ext)
		if err != nil {
			log.Critical(string(err[0].Bytes()))
		}
	}

	system.NoteFinishedExtrinsics()
	system.IdleAndFinalizeHook(block.Header.Number)
}

func initialChecks(block types.Block) {
	header := block.Header
	blockNumber := header.Number

	if blockNumber > 0 {
		storageParentHash := system.StorageGetBlockHash(blockNumber - 1)

		if !reflect.DeepEqual(storageParentHash, header.ParentHash) {
			log.Critical("parent hash should be valid")
		}
	}

	inherentsAreFirst := system.EnsureInherentsAreFirst(block)

	if inherentsAreFirst >= 0 {
		log.Critical(fmt.Sprintf("invalid inherent position for extrinsic at index [%d]", inherentsAreFirst))
	}
}

func runtimeUpgrade() sc.Bool {
	systemHash := hashing.Twox128(constants.KeySystem)
	lastRuntimeUpgradeHash := hashing.Twox128(constants.KeyLastRuntimeUpgrade)

	keyLru := append(systemHash, lastRuntimeUpgradeHash...)
	lrupi := storage.GetDecode(keyLru, types.DecodeLastRuntimeUpgradeInfo)

	if constants.RuntimeVersion.SpecVersion > sc.U32(lrupi.SpecVersion.ToBigInt().Int64()) ||
		lrupi.SpecName != constants.RuntimeVersion.SpecName {

		valueLru := append(
			sc.ToCompact(constants.RuntimeVersion.SpecVersion).Bytes(),
			constants.RuntimeVersion.SpecName.Bytes()...)
		storage.Set(keyLru, valueLru)

		return true
	}

	return false
}

func extractPreRuntimeDigest(digest types.Digest) types.Digest {
	result := types.Digest{}
	for k, v := range digest {
		if k == types.DigestTypePreRuntime {
			result[k] = v
		}
	}

	return result
}

func finalChecks(header *types.Header) {
	newHeader := system.Finalize()

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

// Check a given signed transaction for validity. This doesn't execute any
// side-effects; it merely checks whether the transaction would panic if it were included or
// not.
//
// Changes made to storage should be discarded.
func ValidateTransaction(source types.TransactionSource, uxt types.UncheckedExtrinsic, blockHash types.Blake2bHash) (ok types.ValidTransaction, err types.TransactionValidityError) {
	currentBlockNumber := system.StorageGetBlockNumber()
	system.Initialize(currentBlockNumber+1, blockHash, types.Digest{})

	log.Trace("validate_transaction")

	log.Trace("using_encoded")
	encodedLen := sc.ToCompact(len(uxt.Bytes()))

	log.Trace("check")
	xt, err := extrinsic.Unchecked(uxt).Check(types.DefaultAccountIdLookup())
	if err != nil {
		return ok, err
	}

	log.Trace("dispatch_info")
	dispatchInfo := extrinsic.GetDispatchInfo(xt) // xt.GetDispatchInfo()

	if dispatchInfo.Class.Is(types.DispatchClassMandatory) {
		return ok, types.NewTransactionValidityError(types.NewInvalidTransactionMandatoryValidation())
	}

	log.Trace("validate")
	unsignedValidator := extrinsic.UnsignedValidatorForChecked{}
	return extrinsic.Checked(xt).Validate(unsignedValidator, source, &dispatchInfo, encodedLen)
}
