package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type OnChargeTransaction struct {
	mock.Mock
}

func (ct *OnChargeTransaction) CorrectAndDepositFee(who primitives.AccountId[types.PublicKey], correctedFee types.Balance, tip types.Balance, alreadyWithdrawn sc.Option[types.Balance]) error {
	args := ct.Called(who, correctedFee, tip, alreadyWithdrawn)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(error)
}

func (ct *OnChargeTransaction) WithdrawFee(who primitives.AccountId[types.PublicKey], call primitives.Call, info *types.DispatchInfo, fee types.Balance, tip types.Balance) (
	sc.Option[types.Balance], error) {

	args := ct.Called(who, call, info, fee, tip)

	if args.Get(0).(sc.Option[types.Balance]).HasValue == true {
		return args.Get(0).(sc.Option[types.Balance]), nil
	}

	return sc.NewOption[types.Balance](nil), args.Get(1).(error)
}
