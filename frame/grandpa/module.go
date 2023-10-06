package grandpa

import (
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/hooks"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	name = sc.Str("Grandpa")
)

const (
	PauseFailedError = iota
	ResumeFailedError
	ChangePendingError
	TooSoonError
	InvalidKeyOwnershipProofError
	InvalidEquivocationProofError
	DuplicateOffenceReportError
)

var (
	AuthorityVersion sc.U8 = 1
	EngineId               = [4]byte{'f', 'r', 'n', 'k'}
	KeyTypeId              = [4]byte{'g', 'r', 'a', 'n'}
)

type GrandpaModule interface {
	primitives.Module

	KeyType() primitives.PublicKeyType
	KeyTypeId() [4]byte
	Authorities() sc.Sequence[primitives.Authority]
}

type Module struct {
	primitives.DefaultInherentProvider
	hooks.DefaultDispatchModule
	Index   sc.U8
	storage *storage
}

func New(index sc.U8) Module {
	return Module{
		Index:   index,
		storage: newStorage(),
	}
}

func (m Module) KeyType() primitives.PublicKeyType {
	return primitives.PublicKeyEd25519
}

func (m Module) KeyTypeId() [4]byte {
	return KeyTypeId
}

func (m Module) GetIndex() sc.U8 {
	return m.Index
}

func (m Module) name() sc.Str {
	return name
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

func (m Module) Authorities() sc.Sequence[primitives.Authority] {
	versionedAuthorityList := m.storage.Authorities.Get()

	authorities := versionedAuthorityList.AuthorityList
	if versionedAuthorityList.Version != AuthorityVersion {
		// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
		// log.Warn(fmt.Sprintf("unknown Grandpa authorities version: [%d]", versionedAuthorityList.Version))
		log.Warn("unknown Grandpa authorities version: [" + strconv.Itoa(int(versionedAuthorityList.Version)) + "]")
		return sc.Sequence[primitives.Authority]{}
	}

	return authorities
}

func (m Module) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return m.metadataTypes(), primitives.MetadataModule{
		Name:      m.name(),
		Storage:   sc.Option[primitives.MetadataModuleStorage]{},
		Call:      sc.NewOption[sc.Compact](nil),
		CallDef:   sc.NewOption[primitives.MetadataDefinitionVariant](nil),
		Event:     sc.NewOption[sc.Compact](nil),
		EventDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{},
		Error:     sc.NewOption[sc.Compact](nil),
		Index:     m.Index,
	}
}

func (m Module) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithParams(metadata.GrandpaCalls, "Grandpa calls", sc.Sequence[sc.Str]{"pallet_grandpa", "pallet", "Call"}, primitives.NewMetadataTypeDefinitionVariant(
			// TODO: types
			sc.Sequence[primitives.MetadataDefinitionVariant]{}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataEmptyTypeParameter("T"),
				primitives.NewMetadataEmptyTypeParameter("I"),
			}),
		primitives.NewMetadataTypeWithParams(metadata.TypesGrandpaErrors, "The `Error` enum of this pallet.", sc.Sequence[sc.Str]{"pallet_grandpa", "pallet", "Error"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant("PauseFailed", sc.Sequence[primitives.MetadataTypeDefinitionField]{}, PauseFailedError, ""),
				primitives.NewMetadataDefinitionVariant("ResumeFailed", sc.Sequence[primitives.MetadataTypeDefinitionField]{}, ResumeFailedError, ""),
				primitives.NewMetadataDefinitionVariant("ChangePending", sc.Sequence[primitives.MetadataTypeDefinitionField]{}, ChangePendingError, ""),
				primitives.NewMetadataDefinitionVariant("TooSoon", sc.Sequence[primitives.MetadataTypeDefinitionField]{}, TooSoonError, ""),
				primitives.NewMetadataDefinitionVariant("InvalidKeyOwnershipProof", sc.Sequence[primitives.MetadataTypeDefinitionField]{}, InvalidKeyOwnershipProofError, ""),
				primitives.NewMetadataDefinitionVariant("InvalidEquivocationProof", sc.Sequence[primitives.MetadataTypeDefinitionField]{}, InvalidEquivocationProofError, ""),
				primitives.NewMetadataDefinitionVariant("DuplicateOffenceReport", sc.Sequence[primitives.MetadataTypeDefinitionField]{}, DuplicateOffenceReportError, ""),
			}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataEmptyTypeParameter("T"),
			}),
		primitives.NewMetadataTypeWithPath(metadata.TypesGrandpaAppPublic, "sp_consensus_grandpa app Public", sc.Sequence[sc.Str]{"sp_consensus_grandpa", "app", "Public"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesEd25519PubKey),
			})),
		primitives.NewMetadataType(metadata.TypesTupleGrandaAppPublicU64, "(GrandpaAppPublic, U64)",
			primitives.NewMetadataTypeDefinitionTuple(sc.Sequence[sc.Compact]{sc.ToCompact(metadata.TypesGrandpaAppPublic), sc.ToCompact(metadata.PrimitiveTypesU64)})),
		primitives.NewMetadataType(metadata.TypesSequenceTupleGrandpaAppPublic, "[]byte (GrandpaAppPublic, U64)", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesTupleGrandaAppPublicU64))),
	}
}
