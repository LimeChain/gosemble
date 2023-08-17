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
	"github.com/LimeChain/gosemble/primitives/log"
	storage_root "github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/trie"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	functionRemarkIndex = 0
)

type Module[N sc.Numeric] struct {
	primitives.DefaultProvideInherent
	hooks.DefaultDispatchModule[N]
	Index     sc.U8
	Config    *Config
	Storage   *storage[N]
	Constants *consts
	functions map[sc.U8]primitives.Call
}

func New[N sc.Numeric](index sc.U8, config *Config) Module[N] {
	functions := make(map[sc.U8]primitives.Call)
	storage := newStorage[N]()
	constants := newConstants(config.BlockHashCount, config.BlockWeights, config.BlockLength, config.Version)

	functions[functionRemarkIndex] = newRemarkCall(index, functionRemarkIndex)
	// TODO: add more dispatchables

	return Module[N]{
		Index:     index,
		Config:    config,
		Storage:   storage,
		Constants: constants,
		functions: functions,
	}
}

func (m Module[N]) GetIndex() sc.U8 {
	return m.Index
}

func (m Module[N]) Functions() map[sc.U8]primitives.Call {
	return m.functions
}

func (m Module[N]) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (m Module[N]) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (m Module[N]) Initialize(blockNumber N, parentHash primitives.Blake2bHash, digest primitives.Digest) {
	m.Storage.ExecutionPhase.Put(primitives.NewExtrinsicPhaseInitialization())
	m.Storage.ExtrinsicIndex.Put(sc.U32(0))
	m.Storage.BlockNumber.Put(blockNumber)
	m.Storage.Digest.Put(digest)
	m.Storage.ParentHash.Put(parentHash)
	m.Storage.BlockHash.Put(N(blockNumber-1), parentHash)
	m.Storage.BlockWeight.Clear()
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
func (m Module[N]) RegisterExtraWeightUnchecked(weight primitives.Weight, class primitives.DispatchClass) {
	currentWeight := m.Storage.BlockWeight.Get()
	currentWeight.Accrue(weight, class)
	m.Storage.BlockWeight.Put(currentWeight)
}

func (m Module[N]) NoteFinishedInitialize() {
	m.Storage.ExecutionPhase.Put(primitives.NewExtrinsicPhaseApply(sc.U32(0)))
}

// NoteExtrinsic - what the extrinsic data of the current extrinsic index is.
//
// This is required to be called before applying an extrinsic. The data will used
// in [`finalize`] to calculate the correct extrinsics root.
func (m Module[N]) NoteExtrinsic(encodedExt []byte) {
	extrinsicIndex := m.Storage.ExtrinsicIndex.Get()
	m.Storage.ExtrinsicData.Put(extrinsicIndex, sc.BytesToSequenceU8(encodedExt))
}

// NoteAppliedExtrinsic - To be called immediately after an extrinsic has been applied.
//
// Emits an `ExtrinsicSuccess` or `ExtrinsicFailed` event depending on the outcome.
// The emitted event contains the post-dispatch corrected weight including
// the base-weight for its dispatch class.
func (m Module[N]) NoteAppliedExtrinsic(r *primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo], info primitives.DispatchInfo) {
	baseWeight := m.Constants.BlockWeights.Get(info.Class).BaseExtrinsic // TODO: convert to be a const from module
	info.Weight = primitives.ExtractActualWeight(r, &info).SaturatingAdd(baseWeight)
	info.PaysFee = primitives.ExtractActualPaysFee(r, &info)

	if r.HasError {
		// log.Trace(fmt.Sprintf("Extrinsic failed at block(%d): {%v}", m.Storage.BlockNumber.Get(), r.Err))
		log.Trace("Extrinsic failed at block(" + strconv.Itoa(int(m.Storage.BlockNumber.Get())) + "): {}")
		m.DepositEvent(NewEventExtrinsicFailed(r.Err.Error, info))
	} else {
		m.DepositEvent(NewEventExtrinsicSuccess(info))
	}

	nextExtrinsicIndex := m.Storage.ExtrinsicIndex.Get() + sc.U32(1)
	m.Storage.ExtrinsicIndex.Put(nextExtrinsicIndex)

	m.Storage.ExecutionPhase.Put(primitives.NewExtrinsicPhaseApply(nextExtrinsicIndex))
}

func (m Module[N]) Finalize() primitives.Header[N] {
	m.Storage.ExecutionPhase.Clear()
	m.Storage.AllExtrinsicsLen.Clear()

	blockNumber := m.Storage.BlockNumber.Get()
	parentHash := m.Storage.ParentHash.Get()
	digest := m.Storage.Digest.Get()
	extrinsicCount := m.Storage.ExtrinsicCount.Take()

	var extrinsics []byte

	for i := 0; i < int(extrinsicCount); i++ {
		sci := sc.U32(i)

		extrinsic := m.Storage.ExtrinsicData.TakeBytes(sci)
		extrinsics = append(extrinsics, extrinsic...)
	}

	buf := &bytes.Buffer{}
	extrinsicsRootBytes := trie.Blake2256OrderedRoot(
		append(sc.ToCompact(uint64(extrinsicCount)).Bytes(), extrinsics...),
		constants.StorageVersion)
	buf.Write(extrinsicsRootBytes)
	extrinsicsRoot := primitives.DecodeH256(buf)
	buf.Reset()

	// saturating_sub
	toRemove := blockNumber - N(m.Constants.BlockHashCount-1)
	if toRemove > blockNumber {
		toRemove = 0
	}

	if toRemove != 0 {
		m.Storage.BlockHash.Remove(N(toRemove))
	}

	storageRootBytes := storage_root.Root(int32(m.Constants.Version.StateVersion))
	buf.Write(storageRootBytes)
	storageRoot := primitives.DecodeH256(buf)
	buf.Reset()

	return primitives.Header[N]{
		ExtrinsicsRoot: extrinsicsRoot,
		StateRoot:      storageRoot,
		ParentHash:     parentHash,
		Number:         blockNumber,
		Digest:         digest,
	}
}

func (m Module[N]) NoteFinishedExtrinsics() {
	extrinsicIndex := m.Storage.ExtrinsicIndex.Take()
	m.Storage.ExtrinsicCount.Put(extrinsicIndex)
	m.Storage.ExecutionPhase.Put(primitives.NewExtrinsicPhaseFinalization())
}

func (m Module[N]) ResetEvents() {
	m.Storage.Events.Clear()
	m.Storage.EventCount.Clear()
	m.Storage.EventTopics.Clear(sc.U32(math.MaxUint32))
}

func (m Module[N]) Get(key primitives.PublicKey) primitives.AccountInfo {
	return m.Storage.Account.Get(key)
}

func (m Module[N]) CanDecProviders(who primitives.Address32) bool {
	acc := m.Get(who.FixedSequence)

	return acc.Consumers == 0 || acc.Providers > 1
}

// DepositEvent deposits an event into block's event record.
func (m Module[N]) DepositEvent(event primitives.Event) {
	m.depositEventIndexed([]primitives.H256{}, event)
}

func (m Module[N]) Mutate(who primitives.Address32, f func(who *primitives.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	accountInfo := m.Get(who.FixedSequence)

	result := f(&accountInfo)
	if !result.HasError {
		m.Storage.Account.Put(who.FixedSequence, accountInfo)
	}

	return result
}

func (m Module[N]) TryMutateExists(who primitives.Address32, f func(who *primitives.AccountData) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
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

func (m Module[N]) incProviders(who primitives.Address32) primitives.IncRefStatus {
	result := m.Mutate(who, func(a *primitives.AccountInfo) sc.Result[sc.Encodable] {
		if a.Providers == 0 && a.Sufficients == 0 {
			a.Providers = 1
			m.onCreatedAccount(who)

			return sc.Result[sc.Encodable]{
				HasError: false,
				Value:    primitives.IncRefStatusCreated,
			}
		} else {
			// saturating_add
			newProviders := a.Providers + 1
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

func (m Module[N]) decProviders(who primitives.Address32) (primitives.DecRefStatus, primitives.DispatchError) {
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

		account.Providers -= 1
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
func (m Module[N]) depositEventIndexed(topics []primitives.H256, event primitives.Event) {
	blockNumber := m.Storage.BlockNumber.Get()
	if blockNumber == 0 {
		return
	}

	eventRecord := primitives.EventRecord{
		Phase:  m.Storage.ExecutionPhase.Get(),
		Event:  event,
		Topics: topics,
	}

	oldEventCount := m.Storage.EventCount.Get()
	newEventCount := oldEventCount + 1 // checked_add
	if newEventCount < oldEventCount {
		return
	}

	m.Storage.EventCount.Put(newEventCount)

	m.Storage.Events.Append(eventRecord)

	topicValue := sc.NewVaryingData(blockNumber, oldEventCount)
	for _, topic := range topics {
		m.Storage.EventTopics.Append(topic, topicValue)
	}
}

func (m Module[N]) onCreatedAccount(who primitives.Address32) {
	// hook on creating new account, currently not used in Substrate
	//T::OnNewAccount::on_new_account(&who);
	m.DepositEvent(NewEventNewAccount(who.FixedSequence))
}

func (m Module[N]) onKilledAccount(who primitives.Address32) {
	m.DepositEvent(NewEventKilledAccount(who.FixedSequence))
}

// TODO: Check difference with TryMutateExists
func (m Module[N]) AccountTryMutateExists(who primitives.Address32, f func(who *primitives.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	account := m.Get(who.FixedSequence)

	result := f(&account)

	if !result.HasError {
		m.Storage.Account.Put(who.FixedSequence, account)
	}

	return result
}

func (m Module[N]) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	metadataModule := primitives.MetadataModule{
		Name: "System",
		Storage: sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
			Prefix: "System",
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
		}),
		Call:  sc.NewOption[sc.Compact](sc.ToCompact(metadata.SystemCalls)),
		Event: sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesSystemEvent)),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{
			primitives.NewMetadataModuleConstant(
				"BlockWeights",
				sc.ToCompact(metadata.TypesBlockWeights),
				sc.BytesToSequenceU8(m.Constants.BlockWeights.Bytes()),
				"Block & extrinsics weights: base values and limits.",
			),
			primitives.NewMetadataModuleConstant(
				"BlockLength",
				sc.ToCompact(metadata.TypesBlockLength),
				sc.BytesToSequenceU8(m.Constants.BlockLength.Bytes()),
				"The maximum length of a block (in bytes).",
			),
			primitives.NewMetadataModuleConstant(
				"BlockHashCount",
				sc.ToCompact(metadata.PrimitiveTypesU32),
				sc.BytesToSequenceU8(m.Constants.BlockHashCount.Bytes()),
				"Maximum number of block number to block hash mappings to keep (oldest pruned first).",
			),
			primitives.NewMetadataModuleConstant(
				"DbWeight",
				sc.ToCompact(metadata.TypesDbWeight),
				sc.BytesToSequenceU8(constants.DbWeight.Bytes()),
				"The weight of runtime database operations the runtime can invoke.",
			),
			primitives.NewMetadataModuleConstant(
				"Version",
				sc.ToCompact(metadata.TypesRuntimeVersion),
				sc.BytesToSequenceU8(m.Constants.Version.Bytes()),
				"Get the chain's current version.",
			),
		},
		Error: sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesSystemErrors)),
		Index: m.Index,
	}

	return m.metadataTypes(), metadataModule
}

func (m Module[N]) metadataTypes() sc.Sequence[primitives.MetadataType] {
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

		primitives.NewMetadataTypeWithPath(metadata.CheckNonZeroSender, "CheckNonZeroSender", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_non_zero_sender", "CheckNonZeroSender"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataTypeWithPath(metadata.CheckSpecVersion, "CheckSpecVersion", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_spec_version", "CheckSpecVersion"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataTypeWithPath(metadata.CheckTxVersion, "CheckTxVersion", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_tx_version", "CheckTxVersion"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataTypeWithPath(metadata.CheckGenesis, "CheckGenesis", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_genesis", "CheckGenesis"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataTypeWithPath(metadata.CheckMortality, "CheckMortality", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_mortality", "CheckMortality"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesEra, "Era")})),
		primitives.NewMetadataTypeWithPath(metadata.CheckNonce, "CheckNonce", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_nonce", "CheckNonce"}, primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU32))),
		primitives.NewMetadataTypeWithPath(metadata.CheckWeight, "CheckWeight", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_weight", "CheckWeight"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),

		primitives.NewMetadataTypeWithPath(metadata.TypesEra, "Era", sc.Sequence[sc.Str]{"sp_runtime", "generic", "era", "Era"}, primitives.NewMetadataTypeDefinitionVariant(primitives.EraTypeDefinition())),
	}
}
