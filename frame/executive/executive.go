package executive

import (
	"bytes"

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

func runtimeUpgrade() sc.Bool {
	systemHash := hashing.Twox128(constants.KeySystem)
	lastRuntimeUpgradeHash := hashing.Twox128(constants.KeyLastRuntimeUpgrade)

	keyLru := append(systemHash, lastRuntimeUpgradeHash...)
	last := storage.Get(keyLru)

	buf := &bytes.Buffer{}
	buf.Write(last)

	lrupi, err := types.DecodeLastRuntimeUpgradeInfo(buf)
	if err != nil {
		panic(err)
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

// Apply extrinsic outside of the block execution function.
//
// This doesn't attempt to validate anything regarding the block, but it builds a list of uxt
// hashes.
func ApplyExtrinsic(uxt types.UncheckedExtrinsic) types.ApplyExtrinsicResult {
	// sp_io.InitTracing()
	encoded := uxt.Bytes()
	encodedLen := sc.ToCompact(uint64(len(encoded)))
	// sp_tracing.EnterSpan(sp_tracing.InfoSpan("apply_extrinsic", hexdisplay.From(&encoded)))

	// Verify that the signature is good.
	xt, err := uxt.Check() // TODO: args: (&Default::default())
	if err != nil {
		return types.NewApplyExtrinsicResult(err)
	}

	// We don't need to make sure to `note_extrinsic` only after we know it's going to be
	// executed to prevent it from leaking in storage since at this point, it will either
	// execute or panic (and revert storage changes).
	system.NoteExtrinsic(encoded) // system.PalletSystem

	// AUDIT: Under no circumstances may this function panic from here onwards.

	// Decode parameters and dispatch
	dispatchInfo := xt.GetDispatchInfo()
	res, err := xt.ApplyUnsignedValidator(&dispatchInfo, encodedLen)

	// Mandatory(inherents) are not allowed to fail.
	//
	// The entire block should be discarded if an inherent fails to apply. Otherwise
	// it may open an attack vector.
	if err != nil && (dispatchInfo.Class == types.MandatoryDispatch) {
		return types.NewApplyExtrinsicResult(
			types.NewTransactionValidityError(
				types.NewInvalidTransaction(types.BadMandatoryError),
			),
		)
	}

	system.NoteAppliedExtrinsic(&res, dispatchInfo) // system.PalletSystem

	if err != nil {
		return types.NewApplyExtrinsicResult(err)
	}

	return types.NewApplyExtrinsicResult(types.NewDispatchOutcome(nil))
}
