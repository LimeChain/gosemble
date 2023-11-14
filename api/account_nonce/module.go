package account_nonce

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "AccountNonceApi"
	apiVersion    = 1
)

type Module[T types.PublicKey] struct {
	systemModule system.Module
	memUtils     utils.WasmMemoryTranslator
}

func New[T types.PublicKey](systemModule system.Module) Module[T] {
	return Module[T]{
		systemModule: systemModule,
		memUtils:     utils.NewMemoryTranslator(),
	}
}

func (m Module[T]) Name() string {
	return ApiModuleName
}

func (m Module[T]) Item() types.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return types.NewApiItem(hash, apiVersion)
}

// AccountNonce returns the account nonce of given AccountId.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded AccountId.
// Returns a pointer-size of the SCALE-encoded nonce of the AccountId.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-accountnonceapi-account-nonce)
func (m Module[T]) AccountNonce(dataPtr int32, dataLen int32) int64 {
	b := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	publicKey, err := types.DecodeAccountId[T](buffer)
	if err != nil {
		log.Critical(err.Error())
	}
	account, err := m.systemModule.Get(publicKey)
	if err != nil {
		log.Critical(err.Error())
	}
	nonce := account.Nonce

	return m.memUtils.BytesToOffsetAndSize(nonce.Bytes())
}

func (m Module[T]) Metadata() types.RuntimeApiMetadata {
	methods := sc.Sequence[types.RuntimeApiMethodMetadata]{
		types.RuntimeApiMethodMetadata{
			Name: "account_nonce",
			Inputs: sc.Sequence[types.RuntimeApiMethodParamMetadata]{
				types.RuntimeApiMethodParamMetadata{
					Name: "account",
					Type: sc.ToCompact(metadata.TypesAddress32),
				},
			},
			Output: sc.ToCompact(metadata.PrimitiveTypesU32),
			Docs:   sc.Sequence[sc.Str]{" Get current account nonce of given `AccountId`."},
		},
	}

	return types.RuntimeApiMetadata{
		Name:    ApiModuleName,
		Methods: methods,
		Docs:    sc.Sequence[sc.Str]{" The API to query account nonce."},
	}
}
