package account_nonce

import (
	"bytes"

	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

func AccountNonce(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)

	buffer := &bytes.Buffer{}
	buffer.Write(b)

	publicKey := types.DecodePublicKey(buffer)
	nonce := system.StorageAccountNonce(publicKey)

	return utils.BytesToOffsetAndSize(nonce.Bytes())
}
