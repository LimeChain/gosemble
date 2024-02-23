package system

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

// TODO: implement
type callDoTask struct {
	primitives.Callable
}

func newCallDoTask(moduleId sc.U8, functionId sc.U8) primitives.Call {
	call := callDoTask{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
			Arguments:  sc.NewVaryingData(),
		},
	}

	return call
}

func (c callDoTask) DecodeArgs(buffer *bytes.Buffer) (primitives.Call, error) {
	c.Arguments = sc.NewVaryingData()
	return c, nil
}

func (c callDoTask) Encode(buffer *bytes.Buffer) error {
	return c.Callable.Encode(buffer)
}

func (c callDoTask) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callDoTask) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callDoTask) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callDoTask) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (c callDoTask) BaseWeight() primitives.Weight {
	return callDoTaskWeight(primitives.RuntimeDbWeight{})
}

func (_ callDoTask) WeighData(baseWeight primitives.Weight) primitives.Weight {
	return primitives.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callDoTask) ClassifyDispatch(baseWeight primitives.Weight) primitives.DispatchClass {
	return primitives.NewDispatchClassNormal()
}

func (_ callDoTask) PaysFee(baseWeight primitives.Weight) primitives.Pays {
	return primitives.PaysYes
}

func (_ callDoTask) Dispatch(origin primitives.RuntimeOrigin, _ sc.VaryingData) (primitives.PostDispatchInfo, error) {
	return primitives.PostDispatchInfo{}, nil
}

func (_ callDoTask) Docs() string {
	return "Do some specified task."
}
