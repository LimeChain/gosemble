package mocks

import (
	"bytes"

	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type RuntimeDecoder struct {
	mock.Mock
}

func (m *RuntimeDecoder) DecodeBlock(buffer *bytes.Buffer) (types.Block, error) {
	args := m.Called(buffer)

	return args.Get(0).(types.Block), nil
}

func (m *RuntimeDecoder) DecodeUncheckedExtrinsic(buffer *bytes.Buffer) (types.UncheckedExtrinsic, error) {
	args := m.Called(buffer)

	return args.Get(0).(types.UncheckedExtrinsic), nil
}

func (m *RuntimeDecoder) DecodeCall(buffer *bytes.Buffer) (primitives.Call, error) {
	args := m.Called(buffer)

	return args.Get(0).(primitives.Call), nil
}
