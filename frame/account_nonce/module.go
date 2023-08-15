package account_nonce

import (
	"bytes"

	"github.com/LimeChain/gosemble/frame/system/module"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

type Module struct {
	systemModule module.SystemModule
}

func New(systemModule module.SystemModule) Module {
	return Module{systemModule}
}

// AccountNonce returns the account nonce of given AccountId.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded AccountId.
// Returns a pointer-size of the SCALE-encoded nonce of the AccountId.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-accountnonceapi-account-nonce)
func (m Module) AccountNonce(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	publicKey := types.DecodePublicKey(buffer)
	nonce := m.systemModule.Get(publicKey).Nonce

	return utils.BytesToOffsetAndSize(nonce.Bytes())
}
