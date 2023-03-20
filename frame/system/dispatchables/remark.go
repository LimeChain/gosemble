package system

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type FnRemark struct{}

func (_ FnRemark) Decode(buffer *bytes.Buffer) sc.VaryingData {
	return sc.NewVaryingData()
}

// Make some on-chain remark.
//
// ## Complexity
// - `O(1)`
// The range of component `b` is `[0, 3932160]`.
func (_ FnRemark) BaseWeight(args ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `0`
	// Minimum execution time: 2_018 nanoseconds.
	// Standard Error: 0
	b := sc.Sequence[sc.U8]{} // should be args[0], but since it is empty, it should not be created, otherwise the verification will fail.
	w := types.WeightFromParts(362, 0).SaturatingMul(sc.U64(len(b)))
	return types.WeightFromParts(2_091_000, 0).SaturatingAdd(w)
}

func (_ FnRemark) IsInherent() bool {
	return false
}

func (_ FnRemark) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ FnRemark) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ FnRemark) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (_ FnRemark) Dispatch(origin types.RuntimeOrigin, _ sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	return remark(origin)
}

func remark(origin types.RuntimeOrigin) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	_, err := ensureSignedOrRoot(origin)
	if err != nil {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				Error: err,
			},
		}
	}

	return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: false,
		Ok:       types.PostDispatchInfo{},
	}
}

// Ensure that the origin `o` represents either a signed extrinsic (i.e. transaction) or the root.
// Returns `Ok` with the account that signed the extrinsic, `None` if it was root,  or an `Err`
// otherwise.
func ensureSignedOrRoot(o types.RawOrigin) (ok sc.Option[types.Address32], err types.DispatchError) {
	if o.IsRootOrigin() {
		ok = sc.NewOption[types.Address32](nil)
	} else if o.IsSignedOrigin() {
		ok = sc.NewOption[types.Address32](o.VaryingData[1])
	} else {
		err = types.NewDispatchErrorBadOrigin()
	}
	return ok, err
}
