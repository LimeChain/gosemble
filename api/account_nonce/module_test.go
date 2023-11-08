package account_nonce

import (
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	mockSystem      *mocks.SystemModule
	mockMemoryUtils *mocks.MemoryTranslator
)

func Test_Module_Name(t *testing.T) {
	target := setup()

	assert.Equal(t, ApiModuleName, target.Name())
}

func Test_Module_Item(t *testing.T) {
	target := setup()

	hexName := common.MustBlake2b8([]byte(ApiModuleName))
	expect := types.NewApiItem(hexName, apiVersion)

	result := target.Item()

	assert.Equal(t, expect, result)
}

func Test_Module_AccountNonce(t *testing.T) {
	target := setup()

	publicKey := constants.OneAddressAccountId
	nonce := sc.U32(5)
	accountInfo := types.AccountInfo{
		Nonce: nonce,
	}
	expect := int64(7)

	mockMemoryUtils.On("GetWasmMemorySlice", int32(0), int32(1)).Return(publicKey.Bytes())
	mockSystem.On("Get", publicKey).Return(accountInfo, nil)
	mockMemoryUtils.On("BytesToOffsetAndSize", nonce.Bytes()).Return(expect)

	result := target.AccountNonce(0, 1)

	assert.Equal(t, expect, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", int32(0), int32(1))
	mockSystem.AssertCalled(t, "Get", publicKey)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", nonce.Bytes())
}

func setup() Module[types.Ed25519Signer] {
	mockSystem = new(mocks.SystemModule)
	mockMemoryUtils = new(mocks.MemoryTranslator)

	target := New[types.Ed25519Signer](mockSystem)
	target.memUtils = mockMemoryUtils

	return target
}
