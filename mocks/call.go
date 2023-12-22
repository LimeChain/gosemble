package mocks

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type Call struct {
	mock.Mock
}

func (m *Call) Encode(buffer *bytes.Buffer) error {
	m.Called(buffer)
	return nil
}

func (m *Call) Bytes() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func (m *Call) ModuleIndex() sc.U8 {
	args := m.Called()

	return args.Get(0).(sc.U8)
}

func (m *Call) FunctionIndex() sc.U8 {
	args := m.Called()

	return args.Get(0).(sc.U8)
}

func (m *Call) Args() sc.VaryingData {
	args := m.Called()

	return args.Get(0).(sc.VaryingData)
}

func (m *Call) Dispatch(origin types.RuntimeOrigin, a sc.VaryingData) (types.PostDispatchInfo, error) {
	args := m.Called(origin, a)

	if args.Get(1) == nil {
		return args.Get(0).(types.PostDispatchInfo), nil
	}

	return args.Get(0).(types.PostDispatchInfo), args.Get(1).(error)
}

func (m *Call) BaseWeight() types.Weight {
	args := m.Called()

	return args.Get(0).(types.Weight)
}

func (m *Call) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	args := m.Called(baseWeight)

	return args.Get(0).(types.DispatchClass)
}

func (m *Call) PaysFee(baseWeight types.Weight) types.Pays {
	args := m.Called(baseWeight)

	return args.Get(0).(types.Pays)
}

func (m *Call) WeighData(baseWeight types.Weight) types.Weight {
	args := m.Called(baseWeight)

	return args.Get(0).(types.Weight)
}

func (m *Call) DecodeArgs(buffer *bytes.Buffer) (types.Call, error) {
	args := m.Called(buffer)

	if args.Get(1) == nil {
		return args.Get(0).(types.Call), nil
	}

	return args.Get(0).(types.Call), args.Get(1).(error)
}
