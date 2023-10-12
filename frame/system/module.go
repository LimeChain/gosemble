package system

import (
	"bytes"
	"math"
	"reflect"
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/hooks"
	"github.com/LimeChain/gosemble/primitives/io"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	functionRemarkIndex = iota
)

const (
	name = sc.Str("System")
)

type Module interface {
	primitives.Module

	Initialize(blockNumber sc.U64, parentHash primitives.Blake2bHash, digest primitives.Digest)
	RegisterExtraWeightUnchecked(weight primitives.Weight, class primitives.DispatchClass)
	NoteFinishedInitialize()
	NoteExtrinsic(encodedExt []byte)
	NoteAppliedExtrinsic(r *primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo], info primitives.DispatchInfo)
	Finalize() primitives.Header
	NoteFinishedExtrinsics()
	ResetEvents()
	Get(key primitives.PublicKey) primitives.AccountInfo
	CanDecProviders(who primitives.Address32) bool
	DepositEvent(event primitives.Event)
	TryMutateExists(who primitives.Address32, f func(who *primitives.AccountData) sc.Result[sc.Encodable]) sc.Result[sc.Encodable]

	BlockHashCount() sc.U64
	BlockLength() types.BlockLength
	BlockWeights() types.BlockWeights
	DbWeight() types.RuntimeDbWeight
	Version() types.RuntimeVersion

	StorageDigest() types.Digest

	StorageBlockWeight() primitives.ConsumedWeight
	StorageBlockWeightSet(weight primitives.ConsumedWeight)

	StorageBlockHash(key sc.U64) types.Blake2bHash
	StorageBlockHashSet(key sc.U64, value types.Blake2bHash)
	StorageBlockHashExists(key sc.U64) bool

	StorageBlockNumber() sc.U64

	StorageLastRuntimeUpgrade() types.LastRuntimeUpgradeInfo
	StorageLastRuntimeUpgradeSet(lrui types.LastRuntimeUpgradeInfo)

	StorageAccount(key types.PublicKey) types.AccountInfo
	StorageAccountSet(key types.PublicKey, value types.AccountInfo)

	StorageAllExtrinsicsLen() sc.U32
	StorageAllExtrinsicsLenSet(value sc.U32)
}

type module struct {
	primitives.DefaultInherentProvider
	hooks.DefaultDispatchModule
	Index     sc.U8
	Config    *Config
	storage   *storage
	constants *consts
	functions map[sc.U8]primitives.Call
	trie      io.Trie
	ioStorage io.Storage
}

func New(index sc.U8, config *Config) Module {
	functions := make(map[sc.U8]primitives.Call)
	constants := newConstants(config.BlockHashCount, config.BlockWeights, config.BlockLength, config.DbWeight, config.Version)

	functions[functionRemarkIndex] = newCallRemark(index, functionRemarkIndex)
	// TODO: add more dispatchables

	return module{
		Index:     index,
		Config:    config,
		storage:   newStorage(),
		constants: constants,
		functions: functions,
		trie:      io.NewTrie(),
		ioStorage: io.NewStorage(),
	}
}

func (m module) name() sc.Str {
	return name
}

func (m module) GetIndex() sc.U8 {
	return m.Index
}

func (m module) Functions() map[sc.U8]primitives.Call {
	return m.functions
}

func (m module) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (m module) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (m module) BlockHashCount() sc.U64 {
	return m.constants.BlockHashCount
}

func (m module) BlockLength() types.BlockLength {
	return m.constants.BlockLength
}

func (m module) BlockWeights() types.BlockWeights {
	return m.constants.BlockWeights
}

func (m module) DbWeight() types.RuntimeDbWeight {
	return m.constants.DbWeight
}

func (m module) Version() types.RuntimeVersion {
	return m.constants.Version
}

func (m module) StorageDigest() types.Digest {
	return m.storage.Digest.Get()
}

func (m module) StorageBlockWeight() primitives.ConsumedWeight {
	return m.storage.BlockWeight.Get()
}

func (m module) StorageBlockWeightSet(weight primitives.ConsumedWeight) {
	m.storage.BlockWeight.Put(weight)
}

func (m module) StorageBlockHash(key sc.U64) types.Blake2bHash {
	return m.storage.BlockHash.Get(key)
}

func (m module) StorageBlockHashSet(key sc.U64, value types.Blake2bHash) {
	m.storage.BlockHash.Put(key, value)
}

func (m module) StorageBlockHashExists(key sc.U64) bool {
	return m.storage.BlockHash.Exists(key)
}

func (m module) StorageBlockNumber() sc.U64 {
	return m.storage.BlockNumber.Get()
}

func (m module) StorageBlockNumberSet(blockNumber sc.U64) {
	m.storage.BlockNumber.Put(blockNumber)
}

func (m module) StorageLastRuntimeUpgrade() types.LastRuntimeUpgradeInfo {
	return m.storage.LastRuntimeUpgrade.Get()
}

func (m module) StorageLastRuntimeUpgradeSet(lrui types.LastRuntimeUpgradeInfo) {
	m.storage.LastRuntimeUpgrade.Put(lrui)
}

func (m module) StorageAccount(key types.PublicKey) types.AccountInfo {
	return m.storage.Account.Get(key)
}

func (m module) StorageAccountSet(key types.PublicKey, value types.AccountInfo) {
	m.storage.Account.Put(key, value)
}

func (m module) StorageAllExtrinsicsLen() sc.U32 {
	return m.storage.AllExtrinsicsLen.Get()
}

func (m module) StorageAllExtrinsicsLenSet(value sc.U32) {
	m.storage.AllExtrinsicsLen.Put(value)
}

func (m module) Initialize(blockNumber sc.U64, parentHash primitives.Blake2bHash, digest primitives.Digest) {
	m.storage.ExecutionPhase.Put(primitives.NewExtrinsicPhaseInitialization())
	m.storage.ExtrinsicIndex.Put(sc.U32(0))
	m.storage.BlockNumber.Put(blockNumber)
	m.storage.Digest.Put(digest)
	m.storage.ParentHash.Put(parentHash)
	m.storage.BlockHash.Put(blockNumber-1, parentHash)
	m.storage.BlockWeight.Clear()
}

// RegisterExtraWeightUnchecked - Inform the system pallet of some additional weight that should be accounted for, in the
// current block.
//
// NOTE: use with extra care; this function is made public only be used for certain pallets
// that need it. A runtime that does not have dynamic calls should never need this and should
// stick to static weights. A typical use case for this is inner calls or smart contract calls.
// Furthermore, it only makes sense to use this when it is presumably  _cheap_ to provide the
// argument `weight`; In other words, if this function is to be used to account for some
// unknown, user provided call's weight, it would only make sense to use it if you are sure you
// can rapidly compute the weight of the inner call.
//
// Even more dangerous is to note that this function does NOT take any action, if the new sum
// of block weight is more than the block weight limit. This is what the _unchecked_.
//
// Another potential use-case could be for the `on_initialize` and `on_finalize` hooks.
func (m module) RegisterExtraWeightUnchecked(weight primitives.Weight, class primitives.DispatchClass) {
	currentWeight := m.storage.BlockWeight.Get()
	currentWeight.Accrue(weight, class)
	m.storage.BlockWeight.Put(currentWeight)
}

func (m module) NoteFinishedInitialize() {
	m.storage.ExecutionPhase.Put(primitives.NewExtrinsicPhaseApply(sc.U32(0)))
}

// NoteExtrinsic - what the extrinsic data of the current extrinsic index is.
//
// This is required to be called before applying an extrinsic. The data will used
// in [`finalize`] to calculate the correct extrinsics root.
func (m module) NoteExtrinsic(encodedExt []byte) {
	extrinsicIndex := m.storage.ExtrinsicIndex.Get()
	m.storage.ExtrinsicData.Put(extrinsicIndex, sc.BytesToSequenceU8(encodedExt))
}

// NoteAppliedExtrinsic - To be called immediately after an extrinsic has been applied.
//
// Emits an `ExtrinsicSuccess` or `ExtrinsicFailed` event depending on the outcome.
// The emitted event contains the post-dispatch corrected weight including
// the base-weight for its dispatch class.
func (m module) NoteAppliedExtrinsic(r *primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo], info primitives.DispatchInfo) {
	baseWeight := m.BlockWeights().Get(info.Class).BaseExtrinsic
	info.Weight = primitives.ExtractActualWeight(r, &info).SaturatingAdd(baseWeight)
	info.PaysFee = primitives.ExtractActualPaysFee(r, &info)

	if r.HasError {
		// log.Trace(fmt.Sprintf("Extrinsic failed at block(%d): {%v}", m.Storage.BlockNumber.Get(), r.Err))
		blockNum := m.StorageBlockNumber()
		log.Trace("Extrinsic failed at block(" + strconv.Itoa(int(blockNum)) + "): {}")

		m.DepositEvent(newEventExtrinsicFailed(m.Index, r.Err.Error, info))
	} else {
		m.DepositEvent(newEventExtrinsicSuccess(m.Index, info))
	}

	nextExtrinsicIndex := m.storage.ExtrinsicIndex.Get() + 1
	m.storage.ExtrinsicIndex.Put(nextExtrinsicIndex)
	m.storage.ExecutionPhase.Put(primitives.NewExtrinsicPhaseApply(nextExtrinsicIndex))
}

func (m module) Finalize() primitives.Header {
	m.storage.ExecutionPhase.Clear()
	m.storage.AllExtrinsicsLen.Clear()

	blockNumber := m.StorageBlockNumber()
	parentHash := m.storage.ParentHash.Get()
	digest := m.StorageDigest()
	extrinsicCount := m.storage.ExtrinsicCount.Take()

	var extrinsics []byte

	for i := 0; i < int(extrinsicCount); i++ {
		sci := sc.U32(i)

		extrinsic := m.storage.ExtrinsicData.TakeBytes(sci)
		extrinsics = append(extrinsics, extrinsic...)
	}

	buf := &bytes.Buffer{}
	extrinsicsRootBytes := m.trie.Blake2256OrderedRoot(
		append(sc.ToCompact(uint64(extrinsicCount)).Bytes(), extrinsics...),
		constants.StorageVersion)
	buf.Write(extrinsicsRootBytes)
	extrinsicsRoot := primitives.DecodeH256(buf)
	buf.Reset()

	toRemove := sc.SaturatingSubU64(blockNumber, m.constants.BlockHashCount)
	toRemove = sc.SaturatingSubU64(toRemove, 1)
	if toRemove != 0 {
		m.storage.BlockHash.Remove(toRemove)
	}

	storageRootBytes := m.ioStorage.Root(int32(m.constants.Version.StateVersion))
	buf.Write(storageRootBytes)
	storageRoot := primitives.DecodeH256(buf)
	buf.Reset()

	return primitives.Header{
		ExtrinsicsRoot: extrinsicsRoot,
		StateRoot:      storageRoot,
		ParentHash:     parentHash,
		Number:         blockNumber,
		Digest:         digest,
	}
}

func (m module) NoteFinishedExtrinsics() {
	extrinsicIndex := m.storage.ExtrinsicIndex.Take()
	m.storage.ExtrinsicCount.Put(extrinsicIndex)
	m.storage.ExecutionPhase.Put(primitives.NewExtrinsicPhaseFinalization())
}

func (m module) ResetEvents() {
	m.storage.Events.Clear()
	m.storage.EventCount.Clear()
	m.storage.EventTopics.Clear(sc.U32(math.MaxUint32))
}

func (m module) Get(key primitives.PublicKey) primitives.AccountInfo {
	return m.storage.Account.Get(key)
}

func (m module) CanDecProviders(who primitives.Address32) bool {
	acc := m.Get(who.FixedSequence)

	return acc.Consumers == 0 || acc.Providers > 1
}

// DepositEvent deposits an event into block's event record.
func (m module) DepositEvent(event primitives.Event) {
	m.depositEventIndexed([]primitives.H256{}, event)
}

func (m module) TryMutateExists(who primitives.Address32, f func(*primitives.AccountData) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	account := m.Get(who.FixedSequence)
	wasProviding := false
	if !reflect.DeepEqual(account.Data, primitives.AccountData{}) {
		wasProviding = true
	}

	someData := &primitives.AccountData{}
	if wasProviding {
		someData = &account.Data
	}

	result := f(someData)
	if result.HasError {
		return result
	}

	isProviding := !reflect.DeepEqual(*someData, primitives.AccountData{})

	if !wasProviding && isProviding {
		m.incProviders(who)
	} else if wasProviding && !isProviding {
		status, err := m.decProviders(who)
		if err != nil {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value:    err,
			}
		}
		if status == primitives.DecRefStatusExists {
			return result
		}
	} else if !wasProviding && !isProviding {
		return result
	}

	m.storage.Account.Mutate(who.FixedSequence, func(a *primitives.AccountInfo) sc.Result[sc.Encodable] {
		return mutateAccount(a, someData)
	})

	return result
}

func (m module) incProviders(who primitives.Address32) primitives.IncRefStatus {
	result := m.storage.Account.Mutate(who.FixedSequence, func(account *primitives.AccountInfo) sc.Result[sc.Encodable] {
		return m.incrementProviders(who, account)
	})

	return result.Value.(primitives.IncRefStatus)
}

func (m module) decrementProviders(who primitives.Address32, maybeAccount *sc.Option[primitives.AccountInfo]) sc.Result[sc.Encodable] {
	if maybeAccount.HasValue {
		account := &maybeAccount.Value

		if account.Providers == 0 {
			log.Warn("Logic error: Unexpected underflow in reducing provider")
			account.Providers = 1
		}

		if account.Providers == 1 && account.Consumers == 0 && account.Sufficients == 0 {
			m.onKilledAccount(who)
			// No providers left (and no consumers) and no sufficients. Account dead.
			return sc.Result[sc.Encodable]{
				HasError: false,
				Value:    primitives.DecRefStatusReaped,
			}
		}
		if account.Providers == 1 && account.Consumers > 0 {
			// Cannot remove last provider if there are consumers.
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value:    primitives.NewDispatchErrorConsumerRemaining(),
			}
		}
		// Account will continue to exist as there is either > 1 provider or
		// > 0 sufficients.
		account.Providers = account.Providers - 1
		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    primitives.DecRefStatusExists,
		}
	} else {
		log.Warn("Logic error: Account already dead when reducing provider")
		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    primitives.DecRefStatusReaped,
		}
	}
}

func (m module) incrementProviders(who primitives.Address32, account *primitives.AccountInfo) sc.Result[sc.Encodable] {
	if account.Providers == 0 && account.Sufficients == 0 {
		account.Providers = 1
		m.onCreatedAccount(who)

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    primitives.IncRefStatusCreated,
		}
	} else {
		account.Providers = sc.SaturatingAddU32(account.Providers, 1)

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    primitives.IncRefStatusExisted,
		}
	}
}

func (m module) decProviders(who primitives.Address32) (primitives.DecRefStatus, primitives.DispatchError) {
	result := m.storage.Account.TryMutateExists(who.FixedSequence, func(maybeAccount *sc.Option[primitives.AccountInfo]) sc.Result[sc.Encodable] {
		return m.decrementProviders(who, maybeAccount)
	})

	if result.HasError {
		return sc.U8(0), result.Value.(primitives.DispatchError)
	}

	return result.Value.(primitives.DecRefStatus), nil
}

// depositEventIndexed Deposits an event into this block's event record adding this event
// to the corresponding topic indexes.
//
// This will update storage entries that correspond to the specified topics.
// It is expected that light-clients could subscribe to this topics.
//
// NOTE: Events not registered at the genesis block and quietly omitted.
func (m module) depositEventIndexed(topics []primitives.H256, event primitives.Event) {
	blockNumber := m.StorageBlockNumber()
	if blockNumber == 0 {
		return
	}

	eventRecord := primitives.EventRecord{
		Phase:  m.storage.ExecutionPhase.Get(),
		Event:  event,
		Topics: topics,
	}

	oldEventCount := m.storage.EventCount.Get()
	newEventCount, err := sc.CheckedAddU32(oldEventCount, 1)
	if err != nil {
		return
	}

	m.storage.EventCount.Put(newEventCount)

	m.storage.Events.Append(eventRecord)

	topicValue := sc.NewVaryingData(blockNumber, oldEventCount)
	for _, topic := range topics {
		m.storage.EventTopics.Append(topic, topicValue)
	}
}

func (m module) onCreatedAccount(who primitives.Address32) {
	// hook on creating new account, currently not used in Substrate
	//T::OnNewAccount::on_new_account(&who);
	m.DepositEvent(newEventNewAccount(m.Index, who.FixedSequence))
}

func (m module) onKilledAccount(who primitives.Address32) {
	m.DepositEvent(newEventKilledAccount(m.Index, who.FixedSequence))
}

func (m module) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	metadataModule := primitives.MetadataModule{
		Name:    m.name(),
		Storage: m.metadataStorage(),
		Call:    sc.NewOption[sc.Compact](sc.ToCompact(metadata.SystemCalls)),
		CallDef: sc.NewOption[primitives.MetadataDefinitionVariant](
			primitives.NewMetadataDefinitionVariantStr(
				m.name(),
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.SystemCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<System, Runtime>"),
				},
				m.Index,
				"Call.System"),
		),
		Event: sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesSystemEvent)),
		EventDef: sc.NewOption[primitives.MetadataDefinitionVariant](
			primitives.NewMetadataDefinitionVariantStr(
				m.name(),
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSystemEvent, "frame_system::Event<Runtime>"),
				},
				m.Index,
				"Events.System"),
		),
		Constants: m.metadataConstants(),
		Error:     sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesSystemErrors)),
		Index:     m.Index,
	}

	return m.metadataTypes(), metadataModule
}

func (m module) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithPath(metadata.TypesPhase,
			"frame_system Phase",
			sc.Sequence[sc.Str]{"frame_system", "Phase"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"ApplyExtrinsic",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32),
						},
						primitives.PhaseApplyExtrinsic,
						"Phase.ApplyExtrinsic"),
					primitives.NewMetadataDefinitionVariant(
						"Finalization",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						primitives.PhaseFinalization,
						"Phase.Finalization"),
					primitives.NewMetadataDefinitionVariant(
						"Initialization",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						primitives.PhaseInitialization,
						"Phase.Initialization"),
				})),
		primitives.NewMetadataType(metadata.TypesSystemEventStorage,
			"Vec<Box<EventRecord<T::RuntimeEvent, T::Hash>>>",
			primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesEventRecord))),

		primitives.NewMetadataType(metadata.TypesVecBlockNumEventIndex, "Vec<BlockNumber, EventIndex>",
			primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesTupleU32U32))),

		primitives.NewMetadataTypeWithParam(metadata.TypesPerDispatchClassWeight, "PerDispatchClass[Weight]", sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "normal", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "operational", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "mandatory", "T"),
			},
		),
			primitives.NewMetadataTypeParameter(metadata.TypesWeight, "T"),
		),
		primitives.NewMetadataTypeWithPath(metadata.TypesWeightPerClass, "WeightPerClass", sc.Sequence[sc.Str]{"frame_system", "limits", "WeightsPerClass"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "base_extrinsic", "Weight"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesOptionWeight, "max_extrinsic", "Option<Weight>"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesOptionWeight, "max_total", "Option<Weight>"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesOptionWeight, "reserved", "Option<Weight>"),
			})),
		primitives.NewMetadataTypeWithParam(metadata.TypesPerDispatchClassWeightsPerClass, "PerDispatchClass<WeightPerClass>", sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeightPerClass, "normal", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeightPerClass, "operational", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeightPerClass, "mandatory", "T"),
			}),
			primitives.NewMetadataTypeParameter(metadata.TypesWeightPerClass, "T")),

		primitives.NewMetadataTypeWithPath(metadata.TypesBlockWeights,
			"BlockWeights",
			sc.Sequence[sc.Str]{"frame_system", "limits", "BlockWeights"}, primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "base_block", "Weight"),
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "max_block", "Weight"),
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesPerDispatchClassWeightsPerClass, "per_class", "PerDispatchClass<WeightPerClass>"),
				})),

		primitives.NewMetadataTypeWithPath(metadata.TypesDbWeight, "sp_weights RuntimeDbWeight", sc.Sequence[sc.Str]{"sp_weights", "RuntimeDbWeight"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU64), // read
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU64), // write
			})),

		primitives.NewMetadataTypeWithPath(metadata.TypesBlockLength,
			"frame_system limits BlockLength",
			sc.Sequence[sc.Str]{"frame_system", "limits", "BlockLength"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesPerDispatchClassU32, "max", "PerDispatchClass<u32>"), // max
			})),

		primitives.NewMetadataTypeWithParams(metadata.TypesEventRecord,
			"frame_system EventRecord",
			sc.Sequence[sc.Str]{"frame_system", "EventRecord"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesPhase, "phase", "Phase"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesRuntimeEvent, "event", "E"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesVecTopics, "topics", "Vec<T>"),
			}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesRuntimeEvent, "E"),
				primitives.NewMetadataTypeParameter(metadata.TypesH256, "T"),
			}),
		primitives.NewMetadataTypeWithPath(metadata.TypesSystemEvent,
			"frame_system pallet Event",
			sc.Sequence[sc.Str]{"frame_system", "pallet", "Event"}, primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"ExtrinsicSuccess",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesDispatchInfo, "dispatch_info", "DispatchInfo"),
						},
						EventExtrinsicSuccess,
						"Event.ExtrinsicSuccess"),
					primitives.NewMetadataDefinitionVariant(
						"ExtrinsicFailed",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesDispatchError, "dispatch_error", "DispatchError"),
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesDispatchInfo, "dispatch_info", "DispatchInfo"),
						},
						EventExtrinsicFailed,
						"Events.ExtrinsicFailed"),
					primitives.NewMetadataDefinitionVariant(
						"CodeUpdated",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						EventCodeUpdated,
						"Events.CodeUpdated"),
					primitives.NewMetadataDefinitionVariant(
						"NewAccount",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
						},
						EventNewAccount,
						"Events.NewAccount"),
					primitives.NewMetadataDefinitionVariant(
						"KilledAccount",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
						},
						EventKilledAccount,
						"Events.KilledAccount"),
					primitives.NewMetadataDefinitionVariant(
						"Remarked",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "sender", "T::AccountId"),
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesH256, "hash", "T::Hash"),
						},
						EventRemarked,
						"Events.Remarked"),
				})),

		primitives.NewMetadataTypeWithPath(metadata.TypesLastRuntimeUpgradeInfo,
			"LastRuntimeUpgradeInfo",
			sc.Sequence[sc.Str]{"frame_system", "LastRuntimeUpgradeInfo"}, primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionField(metadata.TypesCompactU32),
					primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesString),
				})),

		primitives.NewMetadataTypeWithPath(metadata.TypesSystemErrors,
			"frame_system pallet Error",
			sc.Sequence[sc.Str]{"frame_system", "pallet", "Error"}, primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"InvalidSpecName",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorInvalidSpecName,
						"The name of specification does not match between the current runtime and the new runtime."),
					primitives.NewMetadataDefinitionVariant(
						"SpecVersionNeedsToIncrease",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorSpecVersionNeedsToIncrease,
						"The specification version is not allowed to decrease between the current runtime and the new runtime."),
					primitives.NewMetadataDefinitionVariant(
						"FailedToExtractRuntimeVersion",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorFailedToExtractRuntimeVersion,
						"Failed to extract the runtime version from the new runtime.  Either calling `Core_version` or decoding `RuntimeVersion` failed."),
					primitives.NewMetadataDefinitionVariant(
						"NonDefaultComposite",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorNonDefaultComposite,
						"Suicide called when the account has non-default composite data."),
					primitives.NewMetadataDefinitionVariant(
						"NonZeroRefCount",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorNonZeroRefCount,
						"There is a non-zero reference count preventing the account from being purged."),
					primitives.NewMetadataDefinitionVariant(
						"CallFiltered",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorCallFiltered,
						"The origin filter prevent the call to be dispatched."),
				})),

		primitives.NewMetadataTypeWithParam(metadata.SystemCalls,
			"System calls",
			sc.Sequence[sc.Str]{"frame_system", "pallet", "Call"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"remark",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),
						},
						functionRemarkIndex,
						"Make some on-chain remark."),
				}),
			primitives.NewMetadataEmptyTypeParameter("T")),

		primitives.NewMetadataTypeWithPath(metadata.TypesEra, "Era", sc.Sequence[sc.Str]{"sp_runtime", "generic", "era", "Era"}, primitives.NewMetadataTypeDefinitionVariant(primitives.EraTypeDefinition())),
	}
}

func (m module) metadataStorage() sc.Option[primitives.MetadataModuleStorage] {
	return sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
		Prefix: m.name(),
		Items: sc.Sequence[primitives.MetadataModuleStorageEntry]{
			primitives.NewMetadataModuleStorageEntry(
				"Account",
				primitives.MetadataModuleStorageEntryModifierDefault,
				primitives.NewMetadataModuleStorageEntryDefinitionMap(
					sc.Sequence[primitives.MetadataModuleStorageHashFunc]{primitives.MetadataModuleStorageHashFuncMultiBlake128Concat},
					sc.ToCompact(metadata.TypesAddress32),
					sc.ToCompact(metadata.TypesAccountInfo)),
				"The full account information for a particular account ID."),
			primitives.NewMetadataModuleStorageEntry(
				"ExtrinsicCount",
				primitives.MetadataModuleStorageEntryModifierOptional,
				primitives.NewMetadataModuleStorageEntryDefinitionPlain(
					sc.ToCompact(metadata.PrimitiveTypesU32)),
				"Total extrinsics count for the current block."),
			primitives.NewMetadataModuleStorageEntry(
				"BlockWeight",
				primitives.MetadataModuleStorageEntryModifierDefault,
				primitives.NewMetadataModuleStorageEntryDefinitionPlain(
					sc.ToCompact(metadata.TypesPerDispatchClassWeight)),
				"The current weight for the block."),
			primitives.NewMetadataModuleStorageEntry(
				"AllExtrinsicsLen",
				primitives.MetadataModuleStorageEntryModifierOptional,
				primitives.NewMetadataModuleStorageEntryDefinitionPlain(
					sc.ToCompact(metadata.PrimitiveTypesU32)),
				"Total length (in bytes) for all extrinsics put together, for the current block."),
			primitives.NewMetadataModuleStorageEntry(
				"BlockHash",
				primitives.MetadataModuleStorageEntryModifierDefault,
				primitives.NewMetadataModuleStorageEntryDefinitionMap(
					sc.Sequence[primitives.MetadataModuleStorageHashFunc]{primitives.MetadataModuleStorageHashFuncMultiXX64},
					sc.ToCompact(metadata.PrimitiveTypesU32),
					sc.ToCompact(metadata.TypesFixedSequence32U8)),
				"Map of block numbers to block hashes."),
			primitives.NewMetadataModuleStorageEntry(
				"ExtrinsicData",
				primitives.MetadataModuleStorageEntryModifierDefault,
				primitives.NewMetadataModuleStorageEntryDefinitionMap(
					sc.Sequence[primitives.MetadataModuleStorageHashFunc]{primitives.MetadataModuleStorageHashFuncMultiXX64},
					sc.ToCompact(metadata.PrimitiveTypesU32),
					sc.ToCompact(metadata.TypesSequenceU8)),
				"Extrinsics data for the current block (maps an extrinsic's index to its data)."),
			primitives.NewMetadataModuleStorageEntry(
				"Number",
				primitives.MetadataModuleStorageEntryModifierDefault,
				primitives.NewMetadataModuleStorageEntryDefinitionPlain(
					sc.ToCompact(metadata.PrimitiveTypesU32)),
				"The current block number being processed. Set by `execute_block`."),
			primitives.NewMetadataModuleStorageEntry(
				"ParentHash",
				primitives.MetadataModuleStorageEntryModifierDefault,
				primitives.NewMetadataModuleStorageEntryDefinitionPlain(
					sc.ToCompact(metadata.TypesFixedSequence32U8)),
				"Hash of the previous block."),
			primitives.NewMetadataModuleStorageEntry(
				"Digest",
				primitives.MetadataModuleStorageEntryModifierDefault,
				primitives.NewMetadataModuleStorageEntryDefinitionPlain(
					sc.ToCompact(metadata.TypesDigest)),
				"Digest of the current block, also part of the block header."),
			primitives.NewMetadataModuleStorageEntry(
				"Events",
				primitives.MetadataModuleStorageEntryModifierDefault,
				primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesSystemEventStorage)),
				"Events deposited for the current block.   NOTE: The item is unbound and should therefore never be read on chain."),
			primitives.NewMetadataModuleStorageEntry(
				"EventTopics",
				primitives.MetadataModuleStorageEntryModifierDefault,
				primitives.NewMetadataModuleStorageEntryDefinitionMap(
					sc.Sequence[primitives.MetadataModuleStorageHashFunc]{primitives.MetadataModuleStorageHashFuncMultiBlake128Concat},
					sc.ToCompact(metadata.TypesH256),
					sc.ToCompact(metadata.TypesVecBlockNumEventIndex)), "Mapping between a topic (represented by T::Hash) and a vector of indexes  of events in the `<Events<T>>` list."),
			primitives.NewMetadataModuleStorageEntry(
				"EventCount",
				primitives.MetadataModuleStorageEntryModifierDefault,
				primitives.NewMetadataModuleStorageEntryDefinitionPlain(
					sc.ToCompact(metadata.PrimitiveTypesU32)),
				"The number of events in the `Events<T>` list."),
			primitives.NewMetadataModuleStorageEntry(
				"LastRuntimeUpgrade",
				primitives.MetadataModuleStorageEntryModifierOptional,
				primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesLastRuntimeUpgradeInfo)),
				"Stores the `spec_version` and `spec_name` of when the last runtime upgrade happened."),
			primitives.NewMetadataModuleStorageEntry(
				"ExecutionPhase",
				primitives.MetadataModuleStorageEntryModifierOptional,
				primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesPhase)),
				"The execution phase of the block."),
		},
	})
}

func (m module) metadataConstants() sc.Sequence[primitives.MetadataModuleConstant] {
	return sc.Sequence[primitives.MetadataModuleConstant]{
		primitives.NewMetadataModuleConstant(
			"BlockWeights",
			sc.ToCompact(metadata.TypesBlockWeights),
			sc.BytesToSequenceU8(m.BlockWeights().Bytes()),
			"Block & extrinsics weights: base values and limits.",
		),
		primitives.NewMetadataModuleConstant(
			"BlockLength",
			sc.ToCompact(metadata.TypesBlockLength),
			sc.BytesToSequenceU8(m.BlockLength().Bytes()),
			"The maximum length of a block (in bytes).",
		),
		primitives.NewMetadataModuleConstant(
			"BlockHashCount",
			sc.ToCompact(metadata.PrimitiveTypesU32),
			sc.BytesToSequenceU8(m.BlockHashCount().Bytes()),
			"Maximum number of block number to block hash mappings to keep (oldest pruned first).",
		),
		primitives.NewMetadataModuleConstant(
			"DbWeight",
			sc.ToCompact(metadata.TypesDbWeight),
			sc.BytesToSequenceU8(m.DbWeight().Bytes()),
			"The weight of runtime database operations the runtime can invoke.",
		),
		primitives.NewMetadataModuleConstant(
			"Version",
			sc.ToCompact(metadata.TypesRuntimeVersion),
			sc.BytesToSequenceU8(m.Version().Bytes()),
			"Get the chain's current version.",
		),
	}
}

func mutateAccount(account *primitives.AccountInfo, data *primitives.AccountData) sc.Result[sc.Encodable] {
	if data != nil {
		account.Data = *data
	} else {
		account.Data = primitives.AccountData{}
	}

	return sc.Result[sc.Encodable]{}
}
