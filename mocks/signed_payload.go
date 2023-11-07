package mocks

import (
	"bytes"

	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type SignedPayload struct {
	mock.Mock
}

func (sp *SignedPayload) Encode(buffer *bytes.Buffer) error {
	sp.Called(buffer)
	return nil
}

func (sp *SignedPayload) Bytes() []byte {
	args := sp.Called()
	return args.Get(0).([]byte)
}

func (sp *SignedPayload) AdditionalSigned() primitives.AdditionalSigned {
	args := sp.Called()
	return args.Get(0).(primitives.AdditionalSigned)
}

func (sp *SignedPayload) Call() primitives.Call {
	args := sp.Called()
	return args.Get(0).(primitives.Call)
}

func (sp *SignedPayload) Extra() primitives.SignedExtra {
	args := sp.Called()
	return args.Get(0).(primitives.SignedExtra)
}
