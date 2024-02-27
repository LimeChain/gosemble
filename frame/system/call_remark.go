package system

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

// callRemark makes an on-chain remark.
// Can be executed by any origin.
type callRemark struct {
	primitives.Callable
}

func newCallRemark(moduleId sc.U8, functionId sc.U8) primitives.Call {
	call := callRemark{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
			Arguments:  sc.NewVaryingData(sc.Sequence[sc.U8]{}),
		},
	}

	return call
}

func (c callRemark) DecodeArgs(buffer *bytes.Buffer) (primitives.Call, error) {
	args, err := sc.DecodeSequence[sc.U8](buffer)
	if err != nil {
		return nil, err
	}
	c.Arguments = sc.NewVaryingData(args)
	return c, nil
}

func (c callRemark) Encode(buffer *bytes.Buffer) error {
	return c.Callable.Encode(buffer)
}

func (c callRemark) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callRemark) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callRemark) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callRemark) Args() sc.VaryingData {
	return c.Callable.Args()
}

// Make some on-chain remark.
//
// ## Complexity
// - `O(1)`
// The range of component `b` is `[0, 3932160]`.
func (c callRemark) BaseWeight() primitives.Weight {
	b := sc.Sequence[sc.U8]{}
	if c.Arguments != nil {
		b = c.Arguments[0].(sc.Sequence[sc.U8])
	}

	return callRemarkWeight(primitives.RuntimeDbWeight{}, sc.U64(len(b)))
}

func (_ callRemark) WeighData(baseWeight primitives.Weight) primitives.Weight {
	return primitives.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callRemark) ClassifyDispatch(baseWeight primitives.Weight) primitives.DispatchClass {
	return primitives.NewDispatchClassNormal()
}

func (_ callRemark) PaysFee(baseWeight primitives.Weight) primitives.Pays {
	return primitives.PaysYes
}

func (_ callRemark) Dispatch(origin primitives.RuntimeOrigin, _ sc.VaryingData) (primitives.PostDispatchInfo, error) {
	return primitives.PostDispatchInfo{}, remark(origin)
}

// remark makes some on-chain remark.
func remark(origin primitives.RuntimeOrigin) error {
	_, err := EnsureSignedOrRoot(origin)
	return err
}

// EnsureSignedOrRoot ensures the origin represents either a signed extrinsic or the root.
// Returns an empty Option if the origin is `Root`.
// Returns an Option with the signer if the origin is signed.
// Returns a `BadOrigin` error if neither of the above.
func EnsureSignedOrRoot(origin primitives.RawOrigin) (sc.Option[primitives.AccountId], error) {
	if origin.IsRootOrigin() {
		return sc.NewOption[primitives.AccountId](nil), nil
	} else if origin.IsSignedOrigin() {
		return sc.NewOption[primitives.AccountId](origin.VaryingData[1]), nil
	}

	return sc.Option[primitives.AccountId]{}, primitives.NewDispatchErrorBadOrigin()
}

func (_ callRemark) Docs() string {
	return "Make some on-chain remark."
}
