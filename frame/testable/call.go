package testable

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/support"
	"github.com/LimeChain/gosemble/primitives/io"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type callTest struct {
	primitives.Callable
}

func newCallTest(moduleId, functionId sc.U8) primitives.Call {
	call := callTest{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
	}

	return call
}

func (c callTest) Metadata() sc.Sequence[primitives.RuntimeApiMethodParamMetadata] {
	return sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{}
}

func (c callTest) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(sc.DecodeSequence[sc.U8](buffer))
	return c
}

func (c callTest) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c callTest) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callTest) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callTest) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callTest) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ callTest) BaseWeight() primitives.Weight {
	return primitives.WeightFromParts(1_000_000, 0)
}

func (_ callTest) WeighData(baseWeight primitives.Weight) primitives.Weight {
	return primitives.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callTest) ClassifyDispatch(baseWeight primitives.Weight) primitives.DispatchClass {
	return primitives.NewDispatchClassNormal()
}

func (_ callTest) PaysFee(baseWeight primitives.Weight) primitives.Pays {
	return primitives.NewPaysYes()
}

func (_ callTest) Dispatch(origin primitives.RuntimeOrigin, _ sc.VaryingData) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo] {
	storage := io.NewStorage()
	storage.Set([]byte("testvalue"), []byte{1})

	transactional := support.NewTransactional[primitives.PostDispatchInfo]()
	transactional.WithStorageLayer(func() (primitives.PostDispatchInfo, primitives.DispatchError) {
		storage.Set([]byte("testvalue"), []byte{2})
		return primitives.PostDispatchInfo{}, primitives.NewDispatchErrorOther("revert")
	})

	return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{Ok: primitives.PostDispatchInfo{}}
}
