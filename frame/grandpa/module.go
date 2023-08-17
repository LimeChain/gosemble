package grandpa

import (
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/hooks"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

var (
	AuthorityVersion sc.U8 = 1
	EngineId               = [4]byte{'f', 'r', 'n', 'k'}
	KeyTypeId              = [4]byte{'g', 'r', 'a', 'n'}
)

type Module[N sc.Numeric] struct {
	primitives.DefaultProvideInherent
	hooks.DefaultDispatchModule[N]
	Index   sc.U8
	storage *storage
}

func New[N sc.Numeric](index sc.U8) Module[N] {
	return Module[N]{
		Index:   index,
		storage: newStorage(),
	}
}

func (gm Module[N]) KeyType() primitives.PublicKeyType {
	return primitives.PublicKeyEd25519
}

func (gm Module[N]) KeyTypeId() [4]byte {
	return KeyTypeId
}

func (gm Module[N]) GetIndex() sc.U8 {
	return gm.Index
}

func (gm Module[N]) Functions() map[sc.U8]primitives.Call {
	return map[sc.U8]primitives.Call{}
}

func (gm Module[N]) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (gm Module[N]) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (gm Module[N]) Authorities() sc.Sequence[primitives.Authority] {
	versionedAuthorityList := gm.storage.Authorities.Get()

	authorities := versionedAuthorityList.AuthorityList
	if versionedAuthorityList.Version != AuthorityVersion {
		// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
		// log.Warn(fmt.Sprintf("unknown Grandpa authorities version: [%d]", versionedAuthorityList.Version))
		log.Warn("unknown Grandpa authorities version: [" + strconv.Itoa(int(versionedAuthorityList.Version)) + "]")
		return sc.Sequence[primitives.Authority]{}
	}

	return authorities
}

func (gm Module[N]) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return gm.metadataTypes(), primitives.MetadataModule{
		Name:      "Grandpa",
		Storage:   sc.Option[primitives.MetadataModuleStorage]{},
		Call:      sc.NewOption[sc.Compact](nil),
		Event:     sc.NewOption[sc.Compact](nil),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{},
		Error:     sc.NewOption[sc.Compact](nil),
		Index:     gm.Index,
	}
}

func (gm Module[N]) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithParams(metadata.GrandpaCalls, "Grandpa calls", sc.Sequence[sc.Str]{"pallet_grandpa", "pallet", "Call"}, primitives.NewMetadataTypeDefinitionVariant(
			// TODO: types
			sc.Sequence[primitives.MetadataDefinitionVariant]{}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataEmptyTypeParameter("T"),
				primitives.NewMetadataEmptyTypeParameter("I"),
			}),
	}
}
