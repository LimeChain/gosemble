package module

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/support"
	"github.com/LimeChain/gosemble/primitives/storage"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type testCall struct {
	primitives.Callable
}

func newTestCall(moduleId, functionId sc.U8) primitives.Call {
	call := testCall{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
	}

	return call
}

func (c testCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(sc.DecodeSequence[sc.U8](buffer))
	return c
}

func (c testCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c testCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c testCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c testCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c testCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ testCall) BaseWeight(args ...any) primitives.Weight {
	return primitives.WeightFromParts(1_000_000, 0)
}

func (_ testCall) WeightInfo(baseWeight primitives.Weight) primitives.Weight {
	return primitives.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ testCall) ClassifyDispatch(baseWeight primitives.Weight) primitives.DispatchClass {
	return primitives.NewDispatchClassNormal()
}

func (_ testCall) PaysFee(baseWeight primitives.Weight) primitives.Pays {
	return primitives.NewPaysYes()
}

func (_ testCall) Dispatch(origin primitives.RuntimeOrigin, _ sc.VaryingData) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo] {
	storage.Set([]byte("testvalue"), []byte{1})

	support.WithStorageLayer(func() (ok primitives.PostDispatchInfo, err primitives.DispatchError) {
		storage.Set([]byte("testvalue"), []byte{2})
		return ok, primitives.NewDispatchErrorOther("revert")
	})

	return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{Ok: primitives.PostDispatchInfo{}}
}
