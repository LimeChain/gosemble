package system

import (
	"bytes"
	"math"
	"reflect"
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/frame/support"
	"github.com/LimeChain/gosemble/hooks"
	"github.com/LimeChain/gosemble/primitives/log"
	storage_root "github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/trie"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	functionRemarkIndex = iota
)

type Module interface {
	types.InherentProvider
	hooks.DispatchModule

	GetIndex() sc.U8
	Functions() map[sc.U8]primitives.Call
	PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError)
	ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError)
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
	Mutate(who primitives.Address32, f func(who *primitives.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable]
	TryMutateExists(who primitives.Address32, f func(who *primitives.AccountData) sc.Result[sc.Encodable]) sc.Result[sc.Encodable]
	AccountTryMutateExists(who primitives.Address32, f func(who *primitives.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable]
	Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule)

	BlockWeights() BlockWeights
	BlockLength() BlockLength
	Version() types.RuntimeVersion
	DbWeight() types.RuntimeDbWeight
	BlockHashCount() sc.U64

	StorageDigest() support.StorageValue[types.Digest]
	StorageBlockWeight() support.StorageValue[primitives.ConsumedWeight]
	StorageBlockHash() support.StorageMap[sc.U64, types.Blake2bHash]
	StorageBlockNumber() support.StorageValue[sc.U64]
	StorageLastRuntimeUpgrade() support.StorageValue[types.LastRuntimeUpgradeInfo]
	StorageAccount() support.StorageMap[types.PublicKey, types.AccountInfo]
	StorageAllExtrinsicsLen() support.StorageValue[sc.U32]
}

type module struct {
	primitives.DefaultInherentProvider
	hooks.DefaultDispatchModule
	Index     sc.U8
	Config    *Config
	storage   *storage
	constants *consts
	functions map[sc.U8]primitives.Call
}

func New(index sc.U8, config *Config) Module {
	functions := make(map[sc.U8]primitives.Call)
	storage := newStorage()
	constants := newConstants(config.BlockHashCount, config.BlockWeights, config.BlockLength, config.DbWeight, config.Version)

	functions[functionRemarkIndex] = newCallRemark(index, functionRemarkIndex)
	// TODO: add more dispatchables

	return module{
		Index:     index,
		Config:    config,
		storage:   storage,
		constants: constants,
		functions: functions,
	}
}

func (m module) name() sc.Str {
	return "System"
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

func (m module) Initialize(blockNumber sc.U64, parentHash primitives.Blake2bHash, digest primitives.Digest) {
	m.storage.ExecutionPhase.Put(primitives.NewExtrinsicPhaseInitialization())
	m.storage.ExtrinsicIndex.Put(sc.U32(0))
	m.StorageBlockNumber().Put(blockNumber)
	m.storage.Digest.Put(digest)
	m.storage.ParentHash.Put(parentHash)
	m.StorageBlockHash().Put(blockNumber-1, parentHash)
	m.StorageBlockWeight().Clear()
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
	currentWeight := m.StorageBlockWeight().Get()
	currentWeight.Accrue(weight, class)
	m.StorageBlockWeight().Put(currentWeight)
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
		blockNum := m.StorageBlockNumber().Get()
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

	blockNumber := m.StorageBlockNumber().Get()
	parentHash := m.storage.ParentHash.Get()
	digest := m.StorageDigest().Get()
	extrinsicCount := m.storage.ExtrinsicCount.Take()

	var extrinsics []byte

	for i := 0; i < int(extrinsicCount); i++ {
		sci := sc.U32(i)

		extrinsic := m.storage.ExtrinsicData.TakeBytes(sci)
		extrinsics = append(extrinsics, extrinsic...)
	}

	buf := &bytes.Buffer{}
	extrinsicsRootBytes := trie.Blake2256OrderedRoot(
		append(sc.ToCompact(uint64(extrinsicCount)).Bytes(), extrinsics...),
		constants.StorageVersion)
	buf.Write(extrinsicsRootBytes)
	extrinsicsRoot := primitives.DecodeH256(buf)
	buf.Reset()

	toRemove := sc.SaturatingSubU64(blockNumber, m.constants.BlockHashCount)
	toRemove = sc.SaturatingSubU64(toRemove, 1)

	if toRemove > blockNumber {
		toRemove = 0
	}

	if toRemove != 0 {
		m.StorageBlockHash().Remove(toRemove)
	}

	storageRootBytes := storage_root.Root(int32(m.constants.Version.StateVersion))
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

func (m module) Mutate(who primitives.Address32, f func(who *primitives.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	accountInfo := m.Get(who.FixedSequence)

	result := f(&accountInfo)
	if !result.HasError {
		m.storage.Account.Put(who.FixedSequence, accountInfo)
	}

	return result
}

func (m module) TryMutateExists(who primitives.Address32, f func(who *primitives.AccountData) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
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

	isProviding := !reflect.DeepEqual(someData, primitives.AccountData{})

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

	m.Mutate(who, func(a *primitives.AccountInfo) sc.Result[sc.Encodable] {
		if someData != nil {
			a.Data = *someData
		} else {
			a.Data = primitives.AccountData{}
		}

		return sc.Result[sc.Encodable]{}
	})

	return result
}

func (m module) incProviders(who primitives.Address32) primitives.IncRefStatus {
	result := m.Mutate(who, func(a *primitives.AccountInfo) sc.Result[sc.Encodable] {
		if a.Providers == 0 && a.Sufficients == 0 {
			a.Providers = 1
			m.onCreatedAccount(who)

			return sc.Result[sc.Encodable]{
				HasError: false,
				Value:    primitives.IncRefStatusCreated,
			}
		} else {
			newProviders := sc.SaturatingAddU32(a.Providers, 1)
			if newProviders < a.Providers {
				newProviders = math.MaxUint32
			}

			return sc.Result[sc.Encodable]{
				HasError: false,
				Value:    primitives.IncRefStatusExisted,
			}
		}
	})

	return result.Value.(primitives.IncRefStatus)
}

func (m module) decProviders(who primitives.Address32) (primitives.DecRefStatus, primitives.DispatchError) {
	result := m.AccountTryMutateExists(who, func(account *primitives.AccountInfo) sc.Result[sc.Encodable] {
		if account.Providers == 0 {
			log.Warn("Logic error: Unexpected underflow in reducing provider")

			account.Providers = 1
		}

		if account.Providers == 1 && account.Consumers == 0 && account.Sufficients == 0 {
			return sc.Result[sc.Encodable]{
				HasError: false,
				Value:    primitives.DecRefStatusReaped,
			}
		}

		if account.Providers == 1 && account.Consumers > 0 {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value:    primitives.NewDispatchErrorConsumerRemaining(),
			}
		}

		account.Providers = account.Providers - 1
		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    primitives.DecRefStatusExists,
		}
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
	blockNumber := m.StorageBlockNumber().Get()
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

// TODO: Check difference with TryMutateExists
func (m module) AccountTryMutateExists(who primitives.Address32, f func(who *primitives.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	account := m.Get(who.FixedSequence)

	result := f(&account)

	if !result.HasError {
		m.storage.Account.Put(who.FixedSequence, account)
	}

	return result
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
						0,
						"The name of specification does not match between the current runtime and the new runtime."),
					primitives.NewMetadataDefinitionVariant(
						"SpecVersionNeedsToIncrease",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						1,
						"The specification version is not allowed to decrease between the current runtime and the new runtime."),
					primitives.NewMetadataDefinitionVariant(
						"FailedToExtractRuntimeVersion",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						2,
						"Failed to extract the runtime version from the new runtime.  Either calling `Core_version` or decoding `RuntimeVersion` failed."),
					primitives.NewMetadataDefinitionVariant(
						"NonDefaultComposite",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						3,
						"Suicide called when the account has non-default composite data."),
					primitives.NewMetadataDefinitionVariant(
						"NonZeroRefCount",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						4,
						"There is a non-zero reference count preventing the account from being purged."),
					primitives.NewMetadataDefinitionVariant(
						"CallFiltered",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						5,
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

func (m module) BlockWeights() BlockWeights {
	return m.constants.BlockWeights
}

func (m module) BlockLength() BlockLength {
	return m.constants.BlockLength
}

func (m module) Version() types.RuntimeVersion {
	return m.constants.Version
}

func (m module) BlockHashCount() sc.U64 {
	return m.constants.BlockHashCount
}

func (m module) DbWeight() types.RuntimeDbWeight {
	return m.constants.DbWeight
}

func (m module) StorageDigest() support.StorageValue[types.Digest] {
	return m.storage.Digest
}

func (m module) StorageBlockWeight() support.StorageValue[primitives.ConsumedWeight] {
	return m.storage.BlockWeight
}

func (m module) StorageBlockHash() support.StorageMap[sc.U64, types.Blake2bHash] {
	return m.storage.BlockHash
}

func (m module) StorageBlockNumber() support.StorageValue[sc.U64] {
	return m.storage.BlockNumber
}

func (m module) StorageLastRuntimeUpgrade() support.StorageValue[types.LastRuntimeUpgradeInfo] {
	return m.storage.LastRuntimeUpgrade
}

func (m module) StorageAccount() support.StorageMap[types.PublicKey, types.AccountInfo] {
	return m.storage.Account
}

func (m module) StorageAllExtrinsicsLen() support.StorageValue[sc.U32] {
	return m.storage.AllExtrinsicsLen
}
