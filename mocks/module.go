package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type Module struct {
	mock.Mock
}

func (m *Module) GetIndex() sc.U8 {
	args := m.Called()

	return args.Get(0).(sc.U8)
}

func (m *Module) Functions() map[sc.U8]types.Call {
	args := m.Called()

	return args.Get(0).(map[sc.U8]types.Call)
}

func (m *Module) PreDispatch(call types.Call) (sc.Empty, types.TransactionValidityError) {
	args := m.Called(call)

	if args.Get(1) == nil {
		return args.Get(0).(sc.Empty), nil
	}

	return args.Get(0).(sc.Empty), args.Get(1).(types.TransactionValidityError)
}

func (m *Module) ValidateUnsigned(txSource types.TransactionSource, call types.Call) (types.ValidTransaction, types.TransactionValidityError) {
	args := m.Called(txSource, call)

	if args.Get(1) == nil {
		return args.Get(0).(types.ValidTransaction), nil
	}

	return args.Get(0).(types.ValidTransaction), args.Get(1).(types.TransactionValidityError)
}

func (m *Module) Metadata() (sc.Sequence[types.MetadataType], types.MetadataModule) {
	args := m.Called()

	return args.Get(0).(sc.Sequence[types.MetadataType]), args.Get(1).(types.MetadataModule)
}

func (m *Module) CreateInherent(inherent types.InherentData) sc.Option[types.Call] {
	args := m.Called(inherent)

	return args.Get(0).(sc.Option[types.Call])
}

func (m *Module) CheckInherent(call types.Call, data types.InherentData) types.FatalError {
	args := m.Called(call, data)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(types.FatalError)
}

func (m *Module) InherentIdentifier() [8]byte {
	args := m.Called()

	return args.Get(0).([8]byte)
}

func (m *Module) IsInherent(call types.Call) bool {
	args := m.Called(call)

	return args.Get(0).(bool)
}

func (m *Module) OnFinalize(n sc.U64) {
	m.Called(n)
}

func (m *Module) OnIdle(n sc.U64, remainingWeight types.Weight) types.Weight {
	args := m.Called(n, remainingWeight)

	return args.Get(0).(types.Weight)
}

func (m *Module) OffchainWorker(n sc.U64) {
	m.Called(n)
}

func (m *Module) OnInitialize(n sc.U64) types.Weight {
	args := m.Called(n)

	return args.Get(0).(types.Weight)
}

func (m *Module) OnRuntimeUpgrade() types.Weight {
	args := m.Called()

	return args.Get(0).(types.Weight)
}