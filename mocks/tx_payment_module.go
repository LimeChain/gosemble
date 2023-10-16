package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type TransactionPaymentModule struct {
	mock.Mock
}

func (m *TransactionPaymentModule) CreateInherent(inherent types.InherentData) sc.Option[types.Call] {
	args := m.Called(inherent)
	return args.Get(0).(sc.Option[types.Call])
}

func (m *TransactionPaymentModule) CheckInherent(call types.Call, data types.InherentData) types.FatalError {
	args := m.Called(call, data)
	return args.Get(0).(types.FatalError)
}

func (m *TransactionPaymentModule) InherentIdentifier() [8]byte {
	args := m.Called()
	return args.Get(0).([8]byte)
}

func (m *TransactionPaymentModule) IsInherent(call types.Call) bool {
	args := m.Called(call)
	return args.Bool(0)
}

func (m *TransactionPaymentModule) OnInitialize(n sc.U64) types.Weight {
	args := m.Called(n)
	return args.Get(0).(types.Weight)
}

func (m *TransactionPaymentModule) OnRuntimeUpgrade() types.Weight {
	args := m.Called()
	return args.Get(0).(types.Weight)
}

func (m *TransactionPaymentModule) OnFinalize(n sc.U64) {
	m.Called(n)
}

func (m *TransactionPaymentModule) OnIdle(n sc.U64, remainingWeight types.Weight) types.Weight {
	args := m.Called(n, remainingWeight)
	return args.Get(0).(types.Weight)
}

func (m *TransactionPaymentModule) OffchainWorker(n sc.U64) {
	m.Called(n)
}

func (m *TransactionPaymentModule) GetIndex() sc.U8 {
	args := m.Called()
	return args.Get(0).(sc.U8)
}

func (m *TransactionPaymentModule) Functions() map[sc.U8]types.Call {
	args := m.Called()
	return args.Get(0).(map[sc.U8]types.Call)
}

func (m *TransactionPaymentModule) PreDispatch(call types.Call) (sc.Empty, types.TransactionValidityError) {
	args := m.Called(call)
	return args.Get(0).(sc.Empty), args.Get(1).(types.TransactionValidityError)
}

func (m *TransactionPaymentModule) ValidateUnsigned(source types.TransactionSource, call types.Call) (types.ValidTransaction, types.TransactionValidityError) {
	args := m.Called(source, call)
	return args.Get(0).(types.ValidTransaction), args.Get(1).(types.TransactionValidityError)
}

func (m *TransactionPaymentModule) Metadata() (sc.Sequence[types.MetadataType], types.MetadataModule) {
	args := m.Called()
	return args.Get(0).(sc.Sequence[types.MetadataType]), args.Get(1).(types.MetadataModule)
}

func (m *TransactionPaymentModule) ComputeFee(len sc.U32, info types.DispatchInfo, tip types.Balance) types.Balance {
	args := m.Called(len, info, tip)
	return args.Get(0).(types.Balance)
}

func (m *TransactionPaymentModule) ComputeFeeDetails(len sc.U32, info types.DispatchInfo, tip types.Balance) types.FeeDetails {
	args := m.Called(len, info, tip)
	return args.Get(0).(types.FeeDetails)
}

func (m *TransactionPaymentModule) ComputeActualFee(len sc.U32, info types.DispatchInfo, postInfo types.PostDispatchInfo, tip types.Balance) types.Balance {
	args := m.Called(len, info, postInfo, tip)
	return args.Get(0).(types.Balance)
}

func (m *TransactionPaymentModule) OperationalFeeMultiplier() sc.U8 {
	args := m.Called()
	return args.Get(0).(sc.U8)
}
