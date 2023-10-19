package mocks

import (
	"bytes"

	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type ExtrinsicSignature struct {
	mock.Mock
}

func (es *ExtrinsicSignature) Encode(buffer *bytes.Buffer) {
	es.Called(buffer)
}

func (es *ExtrinsicSignature) Bytes() []byte {
	args := es.Called()

	return args.Get(0).([]byte)
}

func (es *ExtrinsicSignature) DecodeExtrinsicSignature() types.ExtrinsicSignature {
	args := es.Called()

	return args.Get(0).(types.ExtrinsicSignature)
}
