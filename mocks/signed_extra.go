package mocks

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type SignedExtra struct {
	mock.Mock
}

func (m *SignedExtra) Encode(buffer *bytes.Buffer) error {
	m.Called(buffer)
	return nil
}

func (m *SignedExtra) Bytes() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func (m *SignedExtra) Decode(buffer *bytes.Buffer) {
	m.Called(buffer)
}

func (m *SignedExtra) AdditionalSigned() (types.AdditionalSigned, types.TransactionValidityError) {
	args := m.Called()

	if args.Get(1) != nil {
		return args.Get(0).(types.AdditionalSigned), args.Get(1).(types.TransactionValidityError)
	}

	return args.Get(0).(types.AdditionalSigned), nil
}

func (m *SignedExtra) Validate(who types.AccountId[types.SignerAddress], call types.Call, info *types.DispatchInfo, length sc.Compact) (types.ValidTransaction, types.TransactionValidityError) {
	args := m.Called(who, call, info, length)

	if args.Get(1) != nil {
		return args.Get(0).(types.ValidTransaction), args.Get(1).(types.TransactionValidityError)
	}

	return args.Get(0).(types.ValidTransaction), nil
}

func (m *SignedExtra) ValidateUnsigned(call types.Call, info *types.DispatchInfo, length sc.Compact) (types.ValidTransaction, types.TransactionValidityError) {
	args := m.Called(call, info, length)

	if args.Get(1) != nil {
		return args.Get(0).(types.ValidTransaction), args.Get(1).(types.TransactionValidityError)
	}

	return args.Get(0).(types.ValidTransaction), nil
}

func (m *SignedExtra) PreDispatch(who types.AccountId[types.SignerAddress], call types.Call, info *types.DispatchInfo, length sc.Compact) (sc.Sequence[types.Pre], types.TransactionValidityError) {
	args := m.Called(who, call, info, length)

	if args.Get(1) != nil {
		return args.Get(0).(sc.Sequence[types.Pre]), args.Get(1).(types.TransactionValidityError)
	}

	return args.Get(0).(sc.Sequence[types.Pre]), nil
}

func (m *SignedExtra) PreDispatchUnsigned(call types.Call, info *types.DispatchInfo, length sc.Compact) types.TransactionValidityError {
	args := m.Called(call, info, length)

	return args.Get(0).(types.TransactionValidityError)
}

func (m *SignedExtra) PostDispatch(pre sc.Option[sc.Sequence[types.Pre]], info *types.DispatchInfo, postInfo *types.PostDispatchInfo, length sc.Compact, result *types.DispatchResult) types.TransactionValidityError {
	args := m.Called(pre, info, postInfo, length, result)

	return args.Get(0).(types.TransactionValidityError)
}

func (m *SignedExtra) Metadata() (sc.Sequence[types.MetadataType], sc.Sequence[types.MetadataSignedExtension]) {
	args := m.Called()

	return args.Get(0).(sc.Sequence[types.MetadataType]), args.Get(1).(sc.Sequence[types.MetadataSignedExtension])
}
