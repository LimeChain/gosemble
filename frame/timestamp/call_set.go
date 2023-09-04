package timestamp

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/hooks"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	errTimestampUpdatedOnce   = "Timestamp must be updated only once in the block"
	errTimestampMinimumPeriod = "Timestamp must increment by at least <MinimumPeriod> between sequential blocks"
)

type callSet struct {
	storage        *storage
	onTimestampSet hooks.OnTimestampSet[sc.U64]
	constants      *consts
	primitives.Callable
}

func newCallSet(moduleId sc.U8, functionId sc.U8, storage *storage, constants *consts, onTimestampSet hooks.OnTimestampSet[sc.U64]) primitives.Call {
	call := callSet{
		storage:   storage,
		constants: constants,
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		onTimestampSet: onTimestampSet,
	}

	return call
}

func newCallSetWithArgs(moduleId sc.U8, functionId sc.U8, args sc.VaryingData) primitives.Call {
	call := callSet{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
			Arguments:  args,
		},
	}

	return call
}

func (c callSet) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(sc.DecodeCompact(buffer))
	return c
}

func (c callSet) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c callSet) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callSet) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callSet) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callSet) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (c callSet) BaseWeight() primitives.Weight {
	// Storage: Timestamp Now (r:1 w:1)
	// Proof: Timestamp Now (max_values: Some(1), max_size: Some(8), added: 503, mode: MaxEncodedLen)
	// Storage: Babe CurrentSlot (r:1 w:0)
	// Proof: Babe CurrentSlot (max_values: Some(1), max_size: Some(8), added: 503, mode: MaxEncodedLen)
	// TODO: Consensus algorithm affects weight values.
	// Proof Size summary in bytes:
	//  Measured:  `312`
	//  Estimated: `1006`
	// Minimum execution time: 9_106 nanoseconds.
	r := c.constants.DbWeight.Reads(2)
	w := c.constants.DbWeight.Writes(1)
	return primitives.WeightFromParts(9_258_000, 1006).SaturatingAdd(r).SaturatingAdd(w)
}

func (_ callSet) WeighData(baseWeight primitives.Weight) primitives.Weight {
	return primitives.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callSet) ClassifyDispatch(baseWeight primitives.Weight) primitives.DispatchClass {
	return primitives.NewDispatchClassMandatory()
}

func (_ callSet) PaysFee(baseWeight primitives.Weight) primitives.Pays {
	return primitives.NewPaysYes()
}

func (c callSet) Dispatch(origin primitives.RuntimeOrigin, args sc.VaryingData) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo] {
	compactTs := args[0].(sc.Compact)
	return c.set(origin, sc.U64(compactTs.ToBigInt().Uint64()))
}

// set sets the current time.
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
func (c callSet) set(origin primitives.RuntimeOrigin, now sc.U64) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo] {
	if !origin.IsNoneOrigin() {
		return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
			HasError: true,
			Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
				Error: primitives.NewDispatchErrorBadOrigin(),
			},
		}
	}

	didUpdate := c.storage.DidUpdate.Exists()
	if didUpdate {
		log.Critical(errTimestampUpdatedOnce)
	}

	previousTimestamp := c.storage.Now.Get()

	if !(previousTimestamp.Eq(sc.U64(0)) ||
		now.Gte(previousTimestamp.Add(c.constants.MinimumPeriod))) {
		log.Critical(errTimestampMinimumPeriod)
	}

	c.storage.Now.Put(now)
	c.storage.DidUpdate.Put(true)

	c.onTimestampSet.OnTimestampSet(now)

	return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: false,
		Ok:       primitives.PostDispatchInfo{},
	}
}
