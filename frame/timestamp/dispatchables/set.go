package timestamp

import (
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

func (_ FnSet) Index() sc.U8 {
	return timestamp.FunctionSetIndex
}

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

func (_ FnSet) WeightInfo(baseWeight types.Weight, target sc.Sequence[sc.U8]) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ FnSet) ClassifyDispatch(baseWeight types.Weight, target sc.Sequence[sc.U8]) types.DispatchClass {
	return types.NewDispatchClassMandatory()
}

func (_ FnSet) PaysFee(baseWeight types.Weight, target sc.Sequence[sc.U8]) types.Pays {
	return types.NewPaysYes()
}

func (fn FnSet) Dispatch(origin types.RuntimeOrigin, now sc.U64) (ok sc.Empty, err types.DispatchError) {
	return set(origin, now)
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
func set(origin types.RuntimeOrigin, now sc.U64) (ok sc.Empty, err types.DispatchError) {
	ok, err = EnsureNone(origin)
	if err != nil {
		return ok, err
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

	return ok, err
}

// Ensure that the origin `o` represents an unsigned extrinsic. Returns `Ok` or an `Err` otherwise.
func EnsureNone(o types.RuntimeOrigin) (ok sc.Empty, err types.DispatchError) { // [OuterOrigin, AccountId] \ OuterOrigin: Into<Result<RawOrigin<AccountId>, OuterOrigin>>
	if o.IsNoneOrigin() {
		ok = sc.Empty{}
	} else {
		err = types.NewDispatchErrorBadOrigin()
	}
	return ok, err
}
