package mocks

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type Block struct {
	mock.Mock
}

func (m *Block) Encode(buffer *bytes.Buffer) {
	m.Called(buffer)
}

func (m *Block) Bytes() []byte {
	args := m.Called()

	return args.Get(0).([]byte)
}

func (m *Block) Header() types.Header {
	args := m.Called()

	return args.Get(0).(types.Header)
}

func (m *Block) Extrinsics() sc.Sequence[types.UncheckedExtrinsic] {
	args := m.Called()

	return args.Get(0).(sc.Sequence[types.UncheckedExtrinsic])
}
