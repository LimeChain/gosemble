package timestamp

import (
	"bytes"
	"errors"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/hooks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

var (
	errTimestampUpdatedOnce   = errors.New("Timestamp must be updated only once in the block")
	errTimestampMinimumPeriod = errors.New("Timestamp must increment by at least <MinimumPeriod> between sequential blocks")
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
			Arguments:  sc.NewVaryingData(sc.Compact{Number: sc.NewU64(0)}),
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

func (c callSet) DecodeArgs(buffer *bytes.Buffer) (primitives.Call, error) {
	compact, err := sc.DecodeCompact[sc.U64](buffer)
	if err != nil {
		return nil, err
	}
	c.Arguments = sc.NewVaryingData(compact)
	return c, nil
}

func (c callSet) Encode(buffer *bytes.Buffer) error {
	return c.Callable.Encode(buffer)
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
	return primitives.PaysYes
}

func (c callSet) Dispatch(origin primitives.RuntimeOrigin, args sc.VaryingData) (primitives.PostDispatchInfo, error) {
	valueTs, ok := args[0].(sc.Compact)
	if !ok {
		return primitives.PostDispatchInfo{}, errors.New("couldn't dispatch call set timestamp compact value")
	}
	return primitives.PostDispatchInfo{}, c.set(origin, sc.U64(valueTs.ToBigInt().Uint64()))
}

func (c callSet) Docs() string {
	return "Set the current time."
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
func (c callSet) set(origin primitives.RuntimeOrigin, now sc.U64) error {
	if !origin.IsNoneOrigin() {
		return primitives.NewDispatchErrorBadOrigin()
	}

	didUpdate := c.storage.DidUpdate.Exists()
	if didUpdate {
		return primitives.NewDispatchErrorOther(sc.Str(errTimestampUpdatedOnce.Error()))
	}

	previousTimestamp, err := c.storage.Now.Get()
	if err != nil {
		return primitives.NewDispatchErrorOther(sc.Str(err.Error()))
	}

	if !(previousTimestamp == 0 ||
		now >= previousTimestamp+c.constants.MinimumPeriod) {
		return primitives.NewDispatchErrorOther(sc.Str(errTimestampMinimumPeriod.Error()))
	}

	c.storage.Now.Put(now)
	c.storage.DidUpdate.Put(true)

	return c.onTimestampSet.OnTimestampSet(now)
}
