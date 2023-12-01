package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type testCall struct {
	Callable
}

func (c testCall) Dispatch(origin RuntimeOrigin, args sc.VaryingData) DispatchResultWithPostInfo[PostDispatchInfo] {
	return DispatchResultWithPostInfo[PostDispatchInfo]{
		HasError: false,
		Ok:       PostDispatchInfo{},
		Err:      DispatchErrorWithPostInfo[PostDispatchInfo]{},
	}
}

func (c testCall) BaseWeight() Weight {
	return WeightFromParts(3, 4)
}

func (c testCall) ClassifyDispatch(baseWeight Weight) DispatchClass {
	return NewDispatchClassNormal()
}

func (c testCall) PaysFee(baseWeight Weight) Pays {
	return PaysYes
}

func (c testCall) WeighData(baseWeight Weight) Weight {
	return baseWeight
}

func (c testCall) DecodeArgs(buffer *bytes.Buffer) (Call, error) {
	return c, nil
}
