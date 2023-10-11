package mocks

import (
	"bytes"

	"github.com/LimeChain/gosemble/execution/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type RuntimeDecoder struct {
	mock.Mock
}

func (m *RuntimeDecoder) DecodeBlock(buffer *bytes.Buffer) types.Block {
	args := m.Called(buffer)

	return args.Get(0).(types.Block)
}

func (m *RuntimeDecoder) DecodeUncheckedExtrinsic(buffer *bytes.Buffer) types.UncheckedExtrinsic {
	args := m.Called(buffer)

	return args.Get(0).(types.UncheckedExtrinsic)
}

func (m *RuntimeDecoder) DecodeCall(buffer *bytes.Buffer) primitives.Call {
	args := m.Called(buffer)

	return args.Get(0).(primitives.Call)
}
