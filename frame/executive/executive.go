package executive

import (
	"bytes"
	"fmt"
	"github.com/LimeChain/gosemble/execution/extrinsic"
	"github.com/LimeChain/gosemble/primitives/crypto"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
)

// InitializeBlock initialises a block with the given header,
// starting the execution of a particular block.
func InitializeBlock(header types.Header) {
	system.ResetEvents()

	if runtimeUpgrade() {
		// TODO: weight
	}

	system.Initialize(header.Number, header.ParentHash, extractPreRuntimeDigest(header.Digest))

	// TODO: weight

	system.NoteFinishedInitialize()
}

// ApplyExtrinsic applies extrinsic outside the block execution function.
//
// This doesn't attempt to validate anything regarding the block, but it builds a list of uxt
// hashes.
func ApplyExtrinsic(uxt types.UncheckedExtrinsic) (ok types.DispatchOutcome, err types.TransactionValidityError) { // types.ApplyExtrinsicResult
	// sp_io.InitTracing()
	encoded := uxt.Bytes()
	encodedLen := sc.ToCompact(uint64(len(encoded)))
	// sp_tracing.EnterSpan(sp_tracing.InfoSpan("apply_extrinsic", hexdisplay.From(&encoded)))

	// Verify that the signature is good.
	xt, err := uxt.Check(types.DefaultAccountIdLookup())
	if err != nil {
		return ok, err
	}

	// We don't need to make sure to `note_extrinsic` only after we know it's going to be
	// executed to prevent it from leaking in storage since at this point, it will either
	// execute or panic (and revert storage changes).
	system.NoteExtrinsic(encoded)

	// AUDIT: Under no circumstances may this function panic from here onwards.

	// Decode parameters and dispatch
	dispatchInfo := xt.GetDispatchInfo()
	res, err := extrinsic.ApplyUnsignedValidator(xt, &dispatchInfo, encodedLen)

	// Mandatory(inherents) are not allowed to fail.
	//
	// The entire block should be discarded if an inherent fails to apply. Otherwise
	// it may open an attack vector.
	if res.HasError && (dispatchInfo.Class == types.MandatoryDispatch) {
		return ok, types.NewTransactionValidityError(types.NewInvalidTransaction(types.BadMandatoryError))
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
		panic("Signature verification failed")
	}

	finalChecks(&block.Header)
}

func executeExtrinsicsWithBookKeeping(block types.Block) {
	for _, ext := range block.Extrinsics {
		_, err := ApplyExtrinsic(ext)
		if err != nil {
			panic(string(err[0].Bytes()))
		}
	}

	system.NoteFinishedExtrinsics()
	system.IdleAndFinalizeHook(block.Header.Number)
}

func initialChecks(block types.Block) {
	header := block.Header

	blockNumber := header.Number

	if blockNumber > 0 {
		systemHash := hashing.Twox128(constants.KeySystem)
		previousBlock := blockNumber - 1
		blockNumHash := hashing.Twox64(previousBlock.Bytes())

		blockNumKey := append(systemHash, hashing.Twox128(constants.KeyBlockHash)...)
		blockNumKey = append(blockNumKey, blockNumHash...)
		blockNumKey = append(blockNumKey, previousBlock.Bytes()...)

		previousHash := storage.Get(blockNumKey)

		storageParentHash := types.Blake2bHash{}
		if len(previousHash) > 1 {
			buf := &bytes.Buffer{}

			buf.Write(previousHash[1:]) // Remove option byte
			bytesSequence := sc.DecodeSequence[sc.U8](buf)
			buf.Reset()

			buf.Write(sc.SequenceU8ToBytes(bytesSequence))

			storageParentHash = types.DecodeBlake2bHash(buf)
			buf.Reset()
		} else {
			panic("storage parent hash not found")
		}

		if !reflect.DeepEqual(storageParentHash.FixedSequence, header.ParentHash.FixedSequence) {
			panic("parent hash should be valid")
		}
	}

	inherentsAreFirst := system.EnsureInherentsAreFirst(block)

	if inherentsAreFirst >= 0 {
		panic(fmt.Sprintf("invalid inherent position for extrinsic at index [%d]", inherentsAreFirst))
	}
}

func runtimeUpgrade() sc.Bool {
	systemHash := hashing.Twox128(constants.KeySystem)
	lastRuntimeUpgradeHash := hashing.Twox128(constants.KeyLastRuntimeUpgrade)

	keyLru := append(systemHash, lastRuntimeUpgradeHash...)
	last := storage.Get(keyLru)

	lrupi := types.LastRuntimeUpgradeInfo{}
	if len(last) > 1 {
		buf := &bytes.Buffer{}
		buf.Write(last[1:])

		bytesSequence := sc.DecodeSequence[sc.U8](buf)
		buf.Reset()
		buf.Write(sc.SequenceU8ToBytes(bytesSequence))

		lrupi = types.DecodeLastRuntimeUpgradeInfo(buf)
	}

	if constants.RuntimeVersion.SpecVersion > sc.U32(lrupi.SpecVersion.ToBigInt().Int64()) ||
		lrupi.SpecName != constants.RuntimeVersion.SpecName {

		valueLru := append(
			sc.ToCompact(uint64(constants.RuntimeVersion.SpecVersion)).Bytes(),
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
		panic("Number of digest must match the calculated")
	}

	for key, digest := range header.Digest {
		otherDigest := newHeader.Digest[key]
		if !reflect.DeepEqual(digest, otherDigest) {
			panic("digest item must match that calculated")
		}
	}

	if !reflect.DeepEqual(header.StateRoot, newHeader.StateRoot) {
		panic("Storage root must match that calculated")
	}

	if !reflect.DeepEqual(header.ExtrinsicsRoot, newHeader.ExtrinsicsRoot) {
		panic("Transaction trie must be valid")
	}
}
