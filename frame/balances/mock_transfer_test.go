package balances

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type mockAccountMutator struct {
	mock.Mock
}

func (m *mockAccountMutator) ensureCanWithdraw(who types.AccountId[types.PublicKey], amount sc.U128, reasons types.Reasons, newBalance sc.U128) error {
	args := m.Called(who, amount, reasons, newBalance)
	err, _ := args[0].(error)
	return err
}

func (m *mockAccountMutator) tryMutateAccountWithDust(who types.AccountId[types.PublicKey], f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	args := m.Called(who, f)

	return args[0].(sc.Result[sc.Encodable])
}

func (m *mockAccountMutator) tryMutateAccount(who types.AccountId[types.PublicKey], f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	args := m.Called(who, f)

	return args[0].(sc.Result[sc.Encodable])
}
