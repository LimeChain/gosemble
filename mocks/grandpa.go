package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type GrandpaModule struct {
	mock.Mock
}

func (m *GrandpaModule) GetIndex() sc.U8 {
	args := m.Called()
	return args.Get(0).(sc.U8)
}

func (m *GrandpaModule) Functions() map[sc.U8]primitives.Call {
	args := m.Called()
	return args.Get(0).(map[sc.U8]primitives.Call)
}

func (m *GrandpaModule) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	args := m.Called()
	return args.Get(0).(sc.Empty), args.Get(1).(primitives.TransactionValidityError)
}

func (m *GrandpaModule) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	args := m.Called()
	return args.Get(0).(primitives.ValidTransaction), args.Get(1).(primitives.TransactionValidityError)
}

func (m *GrandpaModule) KeyType() primitives.PublicKeyType {
	args := m.Called()
	return args.Get(0).(primitives.PublicKeyType)
}

func (m *GrandpaModule) KeyTypeId() [4]byte {
	args := m.Called()
	return args.Get(0).([4]byte)
}

func (m *GrandpaModule) OnInitialize(_ sc.U64) primitives.Weight {
	args := m.Called()
	return args.Get(0).(primitives.Weight)
}

func (m *GrandpaModule) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	args := m.Called()
	return args.Get(0).(sc.Sequence[primitives.MetadataType]), args.Get(1).(primitives.MetadataModule)
}

func (m *GrandpaModule) Authorities() sc.Sequence[primitives.Authority] {
	args := m.Called()
	return args.Get(0).(sc.Sequence[primitives.Authority])
}

func (m *GrandpaModule) CreateInherent(inherent types.InherentData) sc.Option[types.Call] {
	args := m.Called()
	return args.Get(0).(sc.Option[types.Call])
}

func (m *GrandpaModule) CheckInherent(call types.Call, data types.InherentData) error {
	args := m.Called()
	return args.Get(0).(error)
}

func (m *GrandpaModule) InherentIdentifier() [8]byte {
	args := m.Called()
	return args.Get(0).([8]byte)
}

func (m *GrandpaModule) IsInherent(call types.Call) bool {
	args := m.Called()
	return args.Get(0).(bool)
}

func (m *GrandpaModule) OnRuntimeUpgrade() primitives.Weight {
	args := m.Called()
	return args.Get(0).(primitives.Weight)
}

func (m *GrandpaModule) OnFinalize(n sc.U64) {
	m.Called()
}

func (m *GrandpaModule) OnIdle(n sc.U64, remainingWeight primitives.Weight) primitives.Weight {
	args := m.Called()
	return args.Get(0).(primitives.Weight)
}

func (m *GrandpaModule) OffchainWorker(n sc.U64) {
	m.Called()
}
