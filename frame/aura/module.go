package aura

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/aura"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/hooks"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

var (
	KeyTypeId = [4]byte{'a', 'u', 'r', 'a'}
)

type Module[N sc.Numeric] struct {
	primitives.DefaultProvideInherent
	hooks.DefaultDispatchModule[N]
	Index     sc.U8
	Config    *Config
	Storage   *storage
	Constants *consts
}

func New[N sc.Numeric](index sc.U8, config *Config) Module[N] {
	storage := newStorage()
	constants := newConstants(config.MinimumPeriod)

	return Module[N]{
		Index:     index,
		Config:    config,
		Storage:   storage,
		Constants: constants,
	}
}

func (m Module[N]) GetIndex() sc.U8 {
	return m.Index
}

func (m Module[N]) Functions() map[sc.U8]primitives.Call {
	return map[sc.U8]primitives.Call{}
}

func (m Module[N]) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (m Module[N]) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (m Module[N]) KeyType() primitives.PublicKeyType {
	return m.Config.KeyType
}

func (m Module[N]) KeyTypeId() [4]byte {
	return KeyTypeId
}

func (m Module[N]) OnInitialize(_ N) primitives.Weight {
	slot := m.currentSlotFromDigests()

	if slot.HasValue {
		newSlot := slot.Value

		currentSlot := m.Storage.CurrentSlot.Get()

		if currentSlot >= newSlot {
			log.Critical("Slot must increase")
		}

		m.Storage.CurrentSlot.Put(newSlot)

		totalAuthorities := m.Storage.Authorities.DecodeLen()
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

		return constants.DbWeight.ReadsWrites(2, 1)
	} else {
		return constants.DbWeight.Reads(1)
	}
}

func (m Module[N]) OnTimestampSet(now sc.U64) {
	slotDuration := m.SlotDuration()
	if slotDuration == 0 {
		log.Critical("Aura slot duration cannot be zero.")
	}

	timestampSlot := now / slotDuration

	currentSlot := m.Storage.CurrentSlot.Get()
	if currentSlot != timestampSlot {
		log.Critical("Timestamp slot must match `CurrentSlot`")
	}
}

func (m Module[N]) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return m.metadataTypes(), primitives.MetadataModule{
		Name: "Aura",
		Storage: sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
			Prefix: "Aura",
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
		}),
		Call:      sc.NewOption[sc.Compact](nil),
		Event:     sc.NewOption[sc.Compact](nil),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{},
		Error:     sc.NewOption[sc.Compact](nil),
		Index:     m.Index,
	}
}

func (m Module[N]) metadataTypes() sc.Sequence[primitives.MetadataType] {
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
	}
}

func (m Module[N]) currentSlotFromDigests() sc.Option[slot] {
	digest := m.Config.SystemDigest()

	for keyDigest, dig := range digest {
		if keyDigest == primitives.DigestTypePreRuntime {
			for _, digestItem := range dig {
				if reflect.DeepEqual(sc.FixedSequenceU8ToBytes(digestItem.Engine), aura.EngineId[:]) {
					buffer := &bytes.Buffer{}
					buffer.Write(sc.SequenceU8ToBytes(digestItem.Payload))

					return sc.NewOption[slot](sc.DecodeU64(buffer))
				}
			}
		}
	}

	return sc.NewOption[slot](nil)
}

func (m Module[N]) SlotDuration() sc.U64 {
	return m.Constants.MinimumPeriod.Mul(2)
}
