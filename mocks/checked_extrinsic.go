package mocks

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type CheckedExtrinsic struct {
	mock.Mock
}

func (c *CheckedExtrinsic) Apply(validator primitives.UnsignedValidator, info *primitives.DispatchInfo, length sc.Compact) (primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo], error) {
	args := c.Called(validator, info, length)

	var arg0 primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]
	var arg1 error

	if args.Get(0) != nil {
		arg0 = args.Get(0).(primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo])
	}

	if args.Get(1) != nil {
		arg1 = args.Get(1).(error)
	}

	return arg0, arg1
}

func (c *CheckedExtrinsic) Function() primitives.Call {
	args := c.Called()
	return args.Get(0).(primitives.Call)
}

func (c *CheckedExtrinsic) Validate(validator primitives.UnsignedValidator, source primitives.TransactionSource, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, error) {
	args := c.Called(validator, source, info, length)

	var arg0 primitives.ValidTransaction
	var arg1 error

	if args.Get(0) != nil {
		arg0 = args.Get(0).(primitives.ValidTransaction)
	}

	if args.Get(1) != nil {
		arg1 = args.Get(1).(error)
	}

	return arg0, arg1
}
