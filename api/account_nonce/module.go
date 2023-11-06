package account_nonce

import (
	"bytes"

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

type Module struct {
	systemModule system.Module
	memUtils     utils.WasmMemoryTranslator
}

func New(systemModule system.Module) Module {
	return Module{
		systemModule: systemModule,
		memUtils:     utils.NewMemoryTranslator(),
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() types.ApiItem {
	hash, err := hashing.MustBlake2b8([]byte(ApiModuleName))
	if err != nil {
		log.Critical(err.Error())
	}
	return types.NewApiItem(hash, apiVersion)
}

// AccountNonce returns the account nonce of given AccountId.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded AccountId.
// Returns a pointer-size of the SCALE-encoded nonce of the AccountId.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-accountnonceapi-account-nonce)
func (m Module) AccountNonce(dataPtr int32, dataLen int32) int64 {
	b := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	publicKey, err := types.DecodePublicKey(buffer)
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
