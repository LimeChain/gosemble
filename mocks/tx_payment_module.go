package mocks

import (
	sc "github.com/LimeChain/goscale"
	tx_types "github.com/LimeChain/gosemble/frame/transaction_payment/types"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type TransactionPaymentModule struct {
	mock.Mock
}

func (m *TransactionPaymentModule) CreateInherent(inherent types.InherentData) (sc.Option[types.Call], error) {
	args := m.Called(inherent)
	if args.Get(1) == nil {
		return args.Get(0).(sc.Option[types.Call]), nil
	}
	return args.Get(0).(sc.Option[types.Call]), args.Get(1).(error)
}

func (m *TransactionPaymentModule) CheckInherent(call types.Call, data types.InherentData) error {
	args := m.Called(call, data)
	return args.Get(0).(error)
}

func (m *TransactionPaymentModule) InherentIdentifier() [8]byte {
	args := m.Called()
	return args.Get(0).([8]byte)
}

func (m *TransactionPaymentModule) IsInherent(call types.Call) bool {
	args := m.Called(call)
	return args.Bool(0)
}

func (m *TransactionPaymentModule) OnInitialize(n sc.U64) (types.Weight, error) {
	args := m.Called(n)
	if args.Get(1) == nil {
		return args.Get(0).(types.Weight), nil
	}
	return args.Get(0).(types.Weight), args.Get(1).(error)
}

func (m *TransactionPaymentModule) OnRuntimeUpgrade() types.Weight {
	args := m.Called()
	return args.Get(0).(types.Weight)
}

func (m *TransactionPaymentModule) OnFinalize(n sc.U64) error {
	m.Called(n)
	return nil
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

func (m *TransactionPaymentModule) PreDispatch(call types.Call) (sc.Empty, error) {
	args := m.Called(call)
	return args.Get(0).(sc.Empty), args.Get(1).(error)
}

func (m *TransactionPaymentModule) ValidateUnsigned(source types.TransactionSource, call types.Call) (types.ValidTransaction, error) {
	args := m.Called(source, call)
	return args.Get(0).(types.ValidTransaction), args.Get(1).(error)
}

func (m *TransactionPaymentModule) Metadata(mdGenerator *types.MetadataGenerator) types.MetadataModule {
	args := m.Called(mdGenerator)
	return args.Get(0).(types.MetadataModule)
}

func (m *TransactionPaymentModule) ComputeFee(len sc.U32, info types.DispatchInfo, tip types.Balance) (types.Balance, error) {
	args := m.Called(len, info, tip)
	if args.Get(1) == nil {
		return args.Get(0).(types.Balance), nil
	}
	return args.Get(0).(types.Balance), args.Get(1).(error)
}

func (m *TransactionPaymentModule) ComputeFeeDetails(len sc.U32, info types.DispatchInfo, tip types.Balance) (tx_types.FeeDetails, error) {
	args := m.Called(len, info, tip)
	if args.Get(1) == nil {
		return args.Get(0).(tx_types.FeeDetails), nil
	}
	return args.Get(0).(tx_types.FeeDetails), args.Get(1).(error)
}

func (m *TransactionPaymentModule) ComputeActualFee(len sc.U32, info types.DispatchInfo, postInfo types.PostDispatchInfo, tip types.Balance) (types.Balance, error) {
	args := m.Called(len, info, postInfo, tip)
	if args.Get(1) == nil {
		return args.Get(0).(types.Balance), nil
	}
	return args.Get(0).(types.Balance), args.Get(1).(error)
}

func (m *TransactionPaymentModule) OperationalFeeMultiplier() sc.U8 {
	args := m.Called()
	return args.Get(0).(sc.U8)
}
