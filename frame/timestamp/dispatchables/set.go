package dispatchables

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/frame/aura"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
)

type FnSet struct{}

// Storage: Timestamp Now (r:1 w:1)
// Proof: Timestamp Now (max_values: Some(1), max_size: Some(8), added: 503, mode: MaxEncodedLen)
// Storage: Babe CurrentSlot (r:1 w:0)
// Proof: Babe CurrentSlot (max_values: Some(1), max_size: Some(8), added: 503, mode: MaxEncodedLen)
func (_ FnSet) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `312`
	//  Estimated: `1006`
	// Minimum execution time: 9_106 nanoseconds.
	r := constants.DbWeight.Reads(2)
	w := constants.DbWeight.Writes(1)
	return types.WeightFromParts(9_258_000, 1006).SaturatingAdd(r).SaturatingAdd(w)
}

func (_ FnSet) Decode(buffer *bytes.Buffer) sc.VaryingData {
	return sc.NewVaryingData(
		sc.DecodeCompact(buffer),
	)
}

func (_ FnSet) IsInherent() bool {
	return true
}

func (_ FnSet) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ FnSet) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassMandatory()
}

func (_ FnSet) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (fn FnSet) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	compactTs := args[0].(sc.Compact)
	return set(origin, sc.U64(compactTs.ToBigInt().Uint64()))
}

// Set the current time.
//
// This call should be invoked exactly once per block. It will panic at the finalization
// phase, if this call hasn't been invoked by that time.
//
// The timestamp should be greater than the previous one by the amount specified by
// `MinimumPeriod`.
//
// The dispatch origin for this call must be `Inherent`.
//
// ## Complexity
//   - `O(1)` (Note that implementations of `OnTimestampSet` must also be `O(1)`)
//   - 1 storage read and 1 storage mutation (codec `O(1)`). (because of `DidUpdate::take` in
//     `on_finalize`)
//   - 1 event handler `on_timestamp_set`. Must be `O(1)`.
func set(origin types.RuntimeOrigin, now sc.U64) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	if !origin.IsNoneOrigin() {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				Error: types.NewDispatchErrorBadOrigin(),
			},
		}
	}

	timestampHash := hashing.Twox128(constants.KeyTimestamp)
	didUpdateHash := hashing.Twox128(constants.KeyDidUpdate)

	didUpdate := storage.Exists(append(timestampHash, didUpdateHash...))

	if didUpdate == 1 {
		log.Critical("Timestamp must be updated only once in the block")
	}

	nowHash := hashing.Twox128(constants.KeyNow)
	previousTimestamp := storage.GetDecode(append(timestampHash, nowHash...), sc.DecodeU64)

	if !(previousTimestamp == 0 || now >= previousTimestamp+timestamp.MinimumPeriod) {
		log.Critical("Timestamp must increment by at least <MinimumPeriod> between sequential blocks")
	}

	storage.Set(append(timestampHash, nowHash...), now.Bytes())
	storage.Set(append(timestampHash, didUpdateHash...), sc.Bool(true).Bytes())

	// TODO: Every consensus that uses the timestamp must implement
	// <T::OnTimestampSet as OnTimestampSet<_>>::on_timestamp_set(now)

	// TODO:
	// timestamp module should not depend on the aura module
	aura.OnTimestampSet(now)

	return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: false,
		Ok:       types.PostDispatchInfo{},
	}
}
