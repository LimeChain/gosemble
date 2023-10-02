package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/support"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type SystemModule struct {
	mock.Mock
}

func (m *SystemModule) CreateInherent(inherent primitives.InherentData) sc.Option[primitives.Call] {
	args := m.Called(inherent)
	return args.Get(0).(sc.Option[primitives.Call])
}

func (m *SystemModule) CheckInherent(call primitives.Call, data primitives.InherentData) error {
	args := m.Called(call, data)
	return args.Error(0)
}

func (m *SystemModule) InherentIdentifier() [8]byte {
	args := m.Called()
	return args.Get(0).([8]byte)
}

func (m *SystemModule) IsInherent(call primitives.Call) bool {
	args := m.Called(call)
	return args.Get(0).(bool)
}

func (m *SystemModule) OnInitialize(n sc.U64) primitives.Weight {
	args := m.Called(n)
	return args.Get(0).(primitives.Weight)
}

func (m *SystemModule) OnRuntimeUpgrade() primitives.Weight {
	args := m.Called()
	return args.Get(0).(primitives.Weight)
}

func (m *SystemModule) OnFinalize(n sc.U64) {
	m.Called(n)
}

func (m *SystemModule) OnIdle(n sc.U64, remainingWeight primitives.Weight) primitives.Weight {
	args := m.Called(n, remainingWeight)
	return args.Get(0).(primitives.Weight)
}

func (m *SystemModule) OffchainWorker(n sc.U64) {
	m.Called(n)
}

func (m *SystemModule) GetIndex() sc.U8 {
	args := m.Called()
	return args.Get(0).(sc.U8)
}

func (m *SystemModule) Functions() map[sc.U8]primitives.Call {
	args := m.Called()
	return args.Get(0).(map[sc.U8]primitives.Call)
}

func (m *SystemModule) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	args := m.Called()
	return args.Get(0).(sc.Empty), args.Get(1).(primitives.TransactionValidityError)
}

func (m *SystemModule) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	args := m.Called()
	return args.Get(0).(primitives.ValidTransaction), args.Get(1).(primitives.TransactionValidityError)
}

func (m *SystemModule) Initialize(blockNumber sc.U64, parentHash primitives.Blake2bHash, digest primitives.Digest) {
	m.Called(blockNumber, parentHash, digest)
}

func (m *SystemModule) RegisterExtraWeightUnchecked(weight primitives.Weight, class primitives.DispatchClass) {
	m.Called(weight, class)
}

func (m *SystemModule) NoteFinishedInitialize() {
	m.Called()
}

func (m *SystemModule) NoteExtrinsic(encodedExt []byte) {
	m.Called(encodedExt)
}

func (m *SystemModule) NoteAppliedExtrinsic(r *primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo], info primitives.DispatchInfo) {
	m.Called(r, info)
}

func (m *SystemModule) Finalize() primitives.Header {
	args := m.Called()
	return args.Get(0).(primitives.Header)
}

func (m *SystemModule) NoteFinishedExtrinsics() {
	m.Called()
}

func (m *SystemModule) ResetEvents() {
	m.Called()
}

func (m *SystemModule) Get(key primitives.PublicKey) primitives.AccountInfo {
	args := m.Called(key)
	return args.Get(0).(primitives.AccountInfo)
}

func (m *SystemModule) CanDecProviders(who primitives.Address32) bool {
	args := m.Called(who)
	return args.Get(0).(bool)
}

func (m *SystemModule) DepositEvent(event primitives.Event) {
	m.Called(event)
}

func (m *SystemModule) Mutate(who primitives.Address32, f func(who *primitives.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	args := m.Called(who, f)
	return args.Get(0).(sc.Result[sc.Encodable])
}

func (m *SystemModule) TryMutateExists(who primitives.Address32, f func(who *primitives.AccountData) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	args := m.Called(who, f)
	return args.Get(0).(sc.Result[sc.Encodable])
}

func (m *SystemModule) AccountTryMutateExists(who primitives.Address32, f func(who *primitives.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	args := m.Called(who, f)
	return args.Get(0).(sc.Result[sc.Encodable])
}

func (m *SystemModule) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	args := m.Called()
	return args.Get(0).(sc.Sequence[primitives.MetadataType]), args.Get(1).(primitives.MetadataModule)
}

func (m *SystemModule) RuntimeUpgrade() bool {
	args := m.Called()
	return args.Get(0).(bool)
}

func (m *SystemModule) BlockWeights() system.BlockWeights {
	args := m.Called()
	return args.Get(0).(system.BlockWeights)
}

func (m *SystemModule) BlockLength() system.BlockLength {
	args := m.Called()
	return args.Get(0).(system.BlockLength)
}

func (m *SystemModule) Version() types.RuntimeVersion {
	args := m.Called()
	return args.Get(0).(types.RuntimeVersion)
}

func (m *SystemModule) DbWeight() types.RuntimeDbWeight {
	args := m.Called()
	return args.Get(0).(types.RuntimeDbWeight)
}

func (m *SystemModule) BlockHashCount() sc.U64 {
	args := m.Called()
	return args.Get(0).(sc.U64)
}

func (m *SystemModule) StorageDigest() support.StorageValue[types.Digest] {
	args := m.Called()
	return args.Get(0).(support.StorageValue[types.Digest])
}

func (m *SystemModule) StorageBlockWeight() support.StorageValue[primitives.ConsumedWeight] {
	args := m.Called()
	return args.Get(0).(support.StorageValue[primitives.ConsumedWeight])
}

func (m *SystemModule) StorageBlockHash() support.StorageMap[sc.U64, types.Blake2bHash] {
	args := m.Called()
	return args.Get(0).(support.StorageMap[sc.U64, types.Blake2bHash])
}

func (m *SystemModule) StorageBlockNumber() support.StorageValue[sc.U64] {
	args := m.Called()
	return args.Get(0).(support.StorageValue[sc.U64])
}

func (m *SystemModule) StorageLastRuntimeUpgrade() support.StorageValue[types.LastRuntimeUpgradeInfo] {
	args := m.Called()
	return args.Get(0).(support.StorageValue[types.LastRuntimeUpgradeInfo])
}

func (m *SystemModule) StorageAccount() support.StorageMap[types.PublicKey, types.AccountInfo] {
	args := m.Called()
	return args.Get(0).(support.StorageMap[types.PublicKey, types.AccountInfo])
}

func (m *SystemModule) StorageAllExtrinsicsLen() support.StorageValue[sc.U32] {
	args := m.Called()
	return args.Get(0).(support.StorageValue[sc.U32])
}
