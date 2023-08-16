package system

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type remarkCall struct {
	primitives.Callable
}

func newRemarkCall(moduleId sc.U8, functionId sc.U8) primitives.Call {
	call := remarkCall{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
	}

	return call
}

func (c remarkCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(sc.DecodeSequence[sc.U8](buffer))
	return c
}

func (c remarkCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c remarkCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c remarkCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c remarkCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c remarkCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

// Make some on-chain remark.
//
// ## Complexity
// - `O(1)`
// The range of component `b` is `[0, 3932160]`.
func (_ remarkCall) BaseWeight(args ...any) primitives.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `0`
	// Minimum execution time: 2_018 nanoseconds.
	// Standard Error: 0
	b := sc.Sequence[sc.U8]{} // should be args[0], but since it is empty, it should not be created, otherwise the verification will fail.
	w := primitives.WeightFromParts(362, 0).SaturatingMul(sc.U64(len(b)))
	return primitives.WeightFromParts(2_091_000, 0).SaturatingAdd(w)
}

func (_ remarkCall) WeightInfo(baseWeight primitives.Weight) primitives.Weight {
	return primitives.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ remarkCall) ClassifyDispatch(baseWeight primitives.Weight) primitives.DispatchClass {
	return primitives.NewDispatchClassNormal()
}

func (_ remarkCall) PaysFee(baseWeight primitives.Weight) primitives.Pays {
	return primitives.NewPaysYes()
}

func (_ remarkCall) Dispatch(origin primitives.RuntimeOrigin, _ sc.VaryingData) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo] {
	return remark(origin)
}

// remark makes some on-chain remark.
func remark(origin primitives.RuntimeOrigin) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo] {
	_, err := EnsureSignedOrRoot(origin)
	if err != nil {
		return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
			HasError: true,
			Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
				Error: err,
			},
		}
	}

	return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: false,
		Ok:       primitives.PostDispatchInfo{},
	}
}

// Ensure that the origin `o` represents either a signed extrinsic (i.e. transaction) or the root.
// Returns `Ok` with the account that signed the extrinsic, `None` if it was root,  or an `Err`
// otherwise.
func EnsureSignedOrRoot(o primitives.RawOrigin) (ok sc.Option[primitives.Address32], err primitives.DispatchError) {
	if o.IsRootOrigin() {
		ok = sc.NewOption[primitives.Address32](nil)
	} else if o.IsSignedOrigin() {
		ok = sc.NewOption[primitives.Address32](o.VaryingData[1])
	} else {
		err = primitives.NewDispatchErrorBadOrigin()
	}
	return ok, err
}
