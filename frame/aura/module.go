package aura

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/hooks"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	errSlotMustIncrease      = "Slot must increase"
	errSlotDurationZero      = "Aura slot duration cannot be zero."
	errTimestampSlotMismatch = "Timestamp slot must match `CurrentSlot`"
)

var (
	EngineId  = [4]byte{'a', 'u', 'r', 'a'}
	KeyTypeId = [4]byte{'a', 'u', 'r', 'a'}
)

type AuraModule interface {
	primitives.Module

	KeyType() primitives.PublicKeyType
	KeyTypeId() [4]byte
	OnTimestampSet(now sc.U64)
	SlotDuration() sc.U64
	GetAuthorities() sc.Option[sc.Sequence[sc.U8]]
}

type Module struct {
	primitives.DefaultInherentProvider
	hooks.DefaultDispatchModule
	index     sc.U8
	config    *Config
	storage   *storage
	constants *consts
}

func New(index sc.U8, config *Config) Module {
	storage := newStorage()

	return Module{
		index:     index,
		config:    config,
		storage:   storage,
		constants: newConstants(config.DbWeight, config.MinimumPeriod),
	}
}

func (m Module) GetIndex() sc.U8 {
	return m.index
}

func (m Module) name() sc.Str {
	return "Aura"
}

func (m Module) Functions() map[sc.U8]primitives.Call {
	return map[sc.U8]primitives.Call{}
}

func (m Module) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (m Module) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (m Module) KeyType() primitives.PublicKeyType {
	return m.config.KeyType
}

func (m Module) KeyTypeId() [4]byte {
	return KeyTypeId
}

func (m Module) OnInitialize(_ sc.U64) (primitives.Weight, error) {
	slot, err := m.currentSlotFromDigests()
	if err != nil {
		return primitives.Weight{}, err
	}

	if slot.HasValue {
		newSlot := slot.Value

		currentSlot := m.storage.CurrentSlot.Get()

		if currentSlot >= newSlot {
			log.Critical(errSlotMustIncrease)
		}

		m.storage.CurrentSlot.Put(newSlot)

		totalAuthorities := m.storage.Authorities.DecodeLen()
		if totalAuthorities.HasValue {
			_ = currentSlot % totalAuthorities.Value

			// TODO: implement once  Session module is added
			/*
				if T::DisabledValidators::is_disabled(authority_index as u32) {
							panic!(
								"Validator with index {:?} is disabled and should not be attempting to author blocks.",
								authority_index,
							);
						}
			*/
		}

		return m.constants.DbWeight.ReadsWrites(2, 1), nil
	} else {
		return m.constants.DbWeight.Reads(1), nil
	}
}

func (m Module) OnTimestampSet(now sc.U64) {
	slotDuration := m.SlotDuration()
	if slotDuration == 0 {
		log.Critical(errSlotDurationZero)
	}

	timestampSlot := now / slotDuration

	currentSlot := m.storage.CurrentSlot.Get()
	if currentSlot != timestampSlot {
		log.Critical(errTimestampSlotMismatch)
	}
}

func (m Module) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	dataV14 := primitives.MetadataModuleV14{
		Name:      m.name(),
		Storage:   m.metadataStorage(),
		Call:      sc.NewOption[sc.Compact](nil),
		CallDef:   sc.NewOption[primitives.MetadataDefinitionVariant](nil),
		Event:     sc.NewOption[sc.Compact](nil),
		EventDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{},
		Error:     sc.NewOption[sc.Compact](nil),
		ErrorDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
		Index:     m.index,
	}
	return m.metadataTypes(), primitives.MetadataModule{
		Version:   primitives.ModuleVersion14,
		ModuleV14: dataV14,
	}
}

func (m Module) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithParams(
			metadata.TypesAuraStorageAuthorities,
			"BoundedVec<T::AuthorityId, T::MaxAuthorities>",
			sc.Sequence[sc.Str]{"bounded_collection", "bounded_vec", "BoundedVec"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionField(metadata.TypesSequencePubKeys),
				}), sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesAuthorityId, "T"),
				primitives.NewMetadataEmptyTypeParameter("S"),
			}),

		primitives.NewMetadataTypeWithPath(metadata.TypesAuthorityId,
			"sp_consensus_aura sr25519 app_sr25519 Public",
			sc.Sequence[sc.Str]{"sp_consensus_aura", "sr25519", "app_sr25519", "Public"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionField(metadata.TypesSr25519PubKey)})),

		primitives.NewMetadataTypeWithPath(metadata.TypesSr25519PubKey,
			"sp_core sr25519 Public",
			sc.Sequence[sc.Str]{"sp_core", "sr25519", "Public"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence32U8)})),

		primitives.NewMetadataType(metadata.TypesSequencePubKeys,
			"[]PublicKey",
			primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesAuthorityId))),

		primitives.NewMetadataTypeWithPath(metadata.TypesAuraSlot,
			"sp_consensus_slots Slot",
			sc.Sequence[sc.Str]{"sp_consensus_slots", "Slot"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU64),
				})),

		// type 924
		primitives.NewMetadataType(metadata.TypesTupleSequenceU8KeyTypeId, "(Seq<U8>, KeyTypeId)",
			primitives.NewMetadataTypeDefinitionTuple(sc.Sequence[sc.Compact]{sc.ToCompact(metadata.TypesSequenceU8), sc.ToCompact(metadata.TypesKeyTypeId)})),

		// type 923
		primitives.NewMetadataType(metadata.TypesSequenceTupleSequenceU8KeyTypeId, "[]byte TupleSequenceU8KeyTypeId", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesTupleSequenceU8KeyTypeId))),

		// type 922
		primitives.NewMetadataTypeWithParam(metadata.TypesOptionTupleSequenceU8KeyTypeId, "Option<TupleSequenceU8KeyTypeId>", sc.Sequence[sc.Str]{"Option"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"None",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					""),
				primitives.NewMetadataDefinitionVariant(
					"Some",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceTupleSequenceU8KeyTypeId),
					},
					1,
					""),
			}),
			primitives.NewMetadataTypeParameter(metadata.TypesSequenceTupleSequenceU8KeyTypeId, "T")),
	}
}

func (m Module) metadataStorage() sc.Option[primitives.MetadataModuleStorage] {
	return sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
		Prefix: m.name(),
		Items: sc.Sequence[primitives.MetadataModuleStorageEntry]{
			primitives.NewMetadataModuleStorageEntry(
				"Authorities",
				primitives.MetadataModuleStorageEntryModifierDefault,
				primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesAuraStorageAuthorities)),
				"The current authority set."),
			primitives.NewMetadataModuleStorageEntry(
				"CurrentSlot",
				primitives.MetadataModuleStorageEntryModifierDefault,
				primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesAuraSlot)),
				"The current slot of this block.   This will be set in `on_initialize`."),
		},
	})
}

func (m Module) currentSlotFromDigests() (sc.Option[slot], error) {
	digest := m.config.SystemDigest()

	for keyDigest, dig := range digest {
		if keyDigest == primitives.DigestTypePreRuntime {
			for _, digestItem := range dig {
				if reflect.DeepEqual(sc.FixedSequenceU8ToBytes(digestItem.Engine), EngineId[:]) {
					buffer := &bytes.Buffer{}
					buffer.Write(sc.SequenceU8ToBytes(digestItem.Payload))

					decodeResult, err := sc.DecodeU64(buffer)
					if err != nil {
						return sc.Option[slot]{}, err
					}

					return sc.NewOption[slot](decodeResult), nil
				}
			}
		}
	}

	return sc.NewOption[slot](nil), nil
}

func (m Module) SlotDuration() sc.U64 {
	return m.constants.MinimumPeriod * 2
}

func (m Module) GetAuthorities() sc.Option[sc.Sequence[sc.U8]] {
	return m.storage.Authorities.GetBytes()
}
