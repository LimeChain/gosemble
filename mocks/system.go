package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type SystemModule struct {
	mock.Mock
}

func (m *SystemModule) CreateInherent(inherent primitives.InherentData) (sc.Option[primitives.Call], error) {
	args := m.Called(inherent)
	return args.Get(0).(sc.Option[primitives.Call]), nil
}

func (m *SystemModule) CheckInherent(call primitives.Call, data primitives.InherentData) primitives.FatalError {
	args := m.Called(call, data)
	return args.Get(0).(primitives.FatalError)
}

func (m *SystemModule) InherentIdentifier() [8]byte {
	args := m.Called()
	return args.Get(0).([8]byte)
}

func (m *SystemModule) IsInherent(call primitives.Call) bool {
	args := m.Called(call)
	return args.Get(0).(bool)
}

func (m *SystemModule) OnInitialize(n sc.U64) (primitives.Weight, error) {
	args := m.Called(n)
	return args.Get(0).(primitives.Weight), nil
}

func (m *SystemModule) OnRuntimeUpgrade() primitives.Weight {
	args := m.Called()
	return args.Get(0).(primitives.Weight)
}

func (m *SystemModule) OnFinalize(n sc.U64) error {
	m.Called(n)
	return nil
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

func (m *SystemModule) PreDispatch(call primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	args := m.Called(call)
	return args.Get(0).(sc.Empty), args.Get(1).(primitives.TransactionValidityError)
}

func (m *SystemModule) ValidateUnsigned(txSource primitives.TransactionSource, call primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	args := m.Called(txSource, call)
	return args.Get(0).(primitives.ValidTransaction), args.Get(1).(primitives.TransactionValidityError)
}

func (m *SystemModule) Initialize(blockNumber sc.U64, parentHash primitives.Blake2bHash, digest primitives.Digest) {
	m.Called(blockNumber, parentHash, digest)
}

func (m *SystemModule) RegisterExtraWeightUnchecked(weight primitives.Weight, class primitives.DispatchClass) error {
	m.Called(weight, class)
	return nil
}

func (m *SystemModule) NoteFinishedInitialize() {
	m.Called()
}

func (m *SystemModule) NoteExtrinsic(encodedExt []byte) error {
	m.Called(encodedExt)
	return nil
}

func (m *SystemModule) NoteAppliedExtrinsic(r *primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo], info primitives.DispatchInfo) error {
	m.Called(r, info)
	return nil
}

func (m *SystemModule) Finalize() (primitives.Header, error) {
	args := m.Called()
	return args.Get(0).(primitives.Header), nil
}

func (m *SystemModule) NoteFinishedExtrinsics() error {
	m.Called()
	return nil
}

func (m *SystemModule) ResetEvents() {
	m.Called()
}

func (m *SystemModule) Get(key primitives.PublicKey) (primitives.AccountInfo, error) {
	args := m.Called(key)
	return args.Get(0).(primitives.AccountInfo), nil
}

func (m *SystemModule) CanDecProviders(who primitives.Address32) (bool, error) {
	args := m.Called(who)
	return args.Get(0).(bool), nil
}

func (m *SystemModule) DepositEvent(event primitives.Event) {
	m.Called(event)
}

func (m *SystemModule) Mutate(who primitives.Address32, f func(who *primitives.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	args := m.Called(who, f)
	return args.Get(0).(sc.Result[sc.Encodable])
}

func (m *SystemModule) TryMutateExists(who primitives.Address32, f func(who *primitives.AccountData) sc.Result[sc.Encodable]) (sc.Result[sc.Encodable], error) {
	args := m.Called(who, f)
	return args.Get(0).(sc.Result[sc.Encodable]), nil
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

func (m *SystemModule) BlockWeights() types.BlockWeights {
	args := m.Called()
	return args.Get(0).(types.BlockWeights)
}

func (m *SystemModule) BlockLength() types.BlockLength {
	args := m.Called()
	return args.Get(0).(types.BlockLength)
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

func (m *SystemModule) StorageDigest() (types.Digest, error) {
	args := m.Called()
	return args.Get(0).(types.Digest), nil
}

func (m *SystemModule) StorageBlockWeight() (primitives.ConsumedWeight, error) {
	args := m.Called()
	return args.Get(0).(primitives.ConsumedWeight), nil
}

func (m *SystemModule) StorageBlockWeightSet(weight primitives.ConsumedWeight) {
	m.Called(weight)
}

func (m *SystemModule) StorageBlockHash(key sc.U64) (types.Blake2bHash, error) {
	args := m.Called(key)
	return args.Get(0).(types.Blake2bHash), nil
}

func (m *SystemModule) StorageBlockHashSet(key sc.U64, value types.Blake2bHash) {
	m.Called(key, value)
}

func (m *SystemModule) StorageBlockHashExists(key sc.U64) bool {
	args := m.Called(key)

	return args.Get(0).(bool)
}

func (m *SystemModule) StorageBlockNumber() (sc.U64, error) {
	args := m.Called()
	return args.Get(0).(sc.U64), nil
}

func (m *SystemModule) StorageBlockNumberSet(blockNumber sc.U64) {
	m.Called(blockNumber)
}

func (m *SystemModule) StorageLastRuntimeUpgrade() (types.LastRuntimeUpgradeInfo, error) {
	args := m.Called()
	return args.Get(0).(types.LastRuntimeUpgradeInfo), nil
}

func (m *SystemModule) StorageLastRuntimeUpgradeSet(lrui types.LastRuntimeUpgradeInfo) {
	m.Called(lrui)
}

func (m *SystemModule) StorageAccount(key types.PublicKey) (types.AccountInfo, error) {
	args := m.Called(key)
	return args.Get(0).(types.AccountInfo), nil
}

func (m *SystemModule) StorageAccountSet(key types.PublicKey, value types.AccountInfo) {
	m.Called(key, value)
}

func (m *SystemModule) StorageAllExtrinsicsLen() (sc.U32, error) {
	args := m.Called()
	return args.Get(0).(sc.U32), nil
}

func (m *SystemModule) StorageAllExtrinsicsLenSet(value sc.U32) {
	m.Called(value)
}
