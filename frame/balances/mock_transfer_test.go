package balances

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type mockAccountMutator struct {
	mock.Mock
}

func (m *mockAccountMutator) ensureCanWithdraw(who types.AccountId, amount sc.U128, reasons types.Reasons, newBalance sc.U128) error {
	args := m.Called(who, amount, reasons, newBalance)

	if args[0] != nil {
		return args[0].(error)
	}

	return nil
}

func (m *mockAccountMutator) tryMutateAccountWithDust(who types.AccountId, f func(who *types.AccountData, bool bool) (sc.Encodable, error)) (sc.Encodable, error) {
	args := m.Called(who, f)

	if args[1] != nil {
		return args[0].(sc.Encodable), args[1].(error)
	}

	return args[0].(sc.Encodable), nil
}

func (m *mockAccountMutator) tryMutateAccount(who types.AccountId, f func(who *types.AccountData, bool bool) (sc.Encodable, error)) (sc.Encodable, error) {
	args := m.Called(who, f)

	if args[1] != nil {
		return args[0].(sc.Encodable), args[1].(error)
	}

	return args[0].(sc.Encodable), nil
}
