package mocks

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type UncheckedExtrinsic struct {
	mock.Mock
}

func (uxt *UncheckedExtrinsic) Encode(buffer *bytes.Buffer) {
	uxt.Called(buffer)
}

func (uxt *UncheckedExtrinsic) Bytes() []byte {
	args := uxt.Called()
	return args.Get(0).([]byte)
}

func (uxt *UncheckedExtrinsic) Signature() sc.Option[primitives.ExtrinsicSignature] {
	args := uxt.Called()
	return args.Get(0).(sc.Option[primitives.ExtrinsicSignature])
}

func (uxt *UncheckedExtrinsic) Function() primitives.Call {
	args := uxt.Called()
	return args.Get(0).(primitives.Call)
}

func (uxt *UncheckedExtrinsic) Extra() primitives.SignedExtra {
	args := uxt.Called()
	return args.Get(0).(primitives.SignedExtra)
}

func (uxt *UncheckedExtrinsic) IsSigned() sc.Bool {
	args := uxt.Called()
	return args.Get(0).(sc.Bool)
}

func (uxt *UncheckedExtrinsic) Check(lookup primitives.AccountIdLookup) (sc.Option[primitives.Address32], primitives.TransactionValidityError) {
	args := uxt.Called(lookup)

	var arg0 sc.Option[primitives.Address32]
	var arg1 primitives.TransactionValidityError

	if args.Get(0) != nil {
		arg0 = args.Get(0).(sc.Option[primitives.Address32])
	}

	if args.Get(1) != nil {
		arg1 = args.Get(1).(primitives.TransactionValidityError)
	}

	return arg0, arg1
}