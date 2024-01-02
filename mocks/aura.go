package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type AuraModule struct {
	mock.Mock
}

func (m *AuraModule) GetIndex() sc.U8 {
	args := m.Called()
	return args.Get(0).(sc.U8)
}

func (m *AuraModule) Functions() map[sc.U8]primitives.Call {
	args := m.Called()
	return args.Get(0).(map[sc.U8]primitives.Call)
}

func (m *AuraModule) PreDispatch(call primitives.Call) (sc.Empty, error) {
	args := m.Called(call)
	return args.Get(0).(sc.Empty), args.Get(1).(error)
}

func (m *AuraModule) ValidateUnsigned(txSource primitives.TransactionSource, call primitives.Call) (primitives.ValidTransaction, error) {
	args := m.Called(txSource, call)
	return args.Get(0).(primitives.ValidTransaction), args.Get(1).(error)
}

func (m *AuraModule) KeyType() primitives.PublicKeyType {
	args := m.Called()
	return args.Get(0).(primitives.PublicKeyType)
}

func (m *AuraModule) KeyTypeId() [4]byte {
	args := m.Called()
	return args.Get(0).([4]byte)
}

func (m *AuraModule) OnInitialize(n sc.U64) (primitives.Weight, error) {
	args := m.Called(n)
	if args.Get(1) == nil {
		return args.Get(0).(primitives.Weight), nil
	}
	return args.Get(0).(primitives.Weight), args.Get(1).(error)
}

func (m *AuraModule) OnTimestampSet(now sc.U64) error {
	args := m.Called(now)
	if args.Error(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (m *AuraModule) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	args := m.Called()
	return args.Get(0).(sc.Sequence[primitives.MetadataType]), args.Get(1).(primitives.MetadataModule)
}

func (m *AuraModule) SlotDuration() sc.U64 {
	args := m.Called()
	return args.Get(0).(sc.U64)
}

func (m *AuraModule) GetAuthorities() (sc.Option[sc.Sequence[sc.U8]], error) {
	args := m.Called()
	if args.Get(1) == nil {
		return args.Get(0).(sc.Option[sc.Sequence[sc.U8]]), nil
	}
	return args.Get(0).(sc.Option[sc.Sequence[sc.U8]]), args.Error(1)
}

func (m *AuraModule) CreateInherent(inherent types.InherentData) (sc.Option[types.Call], error) {
	args := m.Called(inherent)
	if args.Get(1) == nil {
		return args.Get(0).(sc.Option[types.Call]), nil
	}
	return args.Get(0).(sc.Option[types.Call]), args.Error(1)
}

func (m *AuraModule) CheckInherent(call types.Call, data types.InherentData) error {
	args := m.Called(call, data)
	return args.Get(0).(error)
}

func (m *AuraModule) InherentIdentifier() [8]byte {
	args := m.Called()
	return args.Get(0).([8]byte)
}

func (m *AuraModule) IsInherent(call types.Call) bool {
	args := m.Called(call)
	return args.Get(0).(bool)
}

func (m *AuraModule) OnRuntimeUpgrade() primitives.Weight {
	args := m.Called()
	return args.Get(0).(primitives.Weight)
}

func (m *AuraModule) OnFinalize(n sc.U64) error {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (m *AuraModule) OnIdle(n sc.U64, remainingWeight primitives.Weight) primitives.Weight {
	args := m.Called(n, remainingWeight)
	return args.Get(0).(primitives.Weight)
}

func (m *AuraModule) OffchainWorker(n sc.U64) {
	m.Called(n)
}
