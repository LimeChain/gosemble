package grandpa

import (
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/hooks"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "GrandpaApi"
	apiVersion    = 3
)

var (
	AuthorityVersion sc.U8 = 1
	EngineId               = [4]byte{'f', 'r', 'n', 'k'}
	KeyTypeId              = [4]byte{'g', 'r', 'a', 'n'}
)

type Module struct {
	primitives.DefaultProvideInherent
	hooks.DefaultDispatchModule[sc.U32]
	Index   sc.U8
	storage *storage
}

func NewModule(index sc.U8) Module {
	return Module{
		Index:   index,
		storage: newStorage(),
	}
}

func (gm Module) Name() string {
	return ApiModuleName
}

func (gm Module) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

func (gm Module) KeyType() primitives.PublicKeyType {
	return primitives.PublicKeyEd25519
}

func (gm Module) KeyTypeId() [4]byte {
	return KeyTypeId
}

func (gm Module) GetIndex() sc.U8 {
	return gm.Index
}

func (gm Module) Functions() map[sc.U8]primitives.Call {
	return map[sc.U8]primitives.Call{}
}

func (gm Module) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (gm Module) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (gm Module) Authorities() int64 {
	versionedAuthorityList := gm.storage.Authorities.Get()

	authorities := versionedAuthorityList.AuthorityList
	if versionedAuthorityList.Version != AuthorityVersion {
		// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
		// log.Warn(fmt.Sprintf("unknown Grandpa authorities version: [%d]", versionedAuthorityList.Version))
		log.Warn("unknown Grandpa authorities version: [" + strconv.Itoa(int(versionedAuthorityList.Version)) + "]")
		authorities = sc.Sequence[primitives.Authority]{}
	}

	return utils.BytesToOffsetAndSize(authorities.Bytes())
}

func (gm Module) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
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

func (gm Module) metadataTypes() sc.Sequence[primitives.MetadataType] {
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
