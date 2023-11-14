package account_nonce

import (
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
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

func Test_Module_Metadata(t *testing.T) {
	target := setup()

	expect := types.RuntimeApiMetadata{
		Name: ApiModuleName,
		Methods: sc.Sequence[types.RuntimeApiMethodMetadata]{
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
		},
		Docs: sc.Sequence[sc.Str]{" The API to query account nonce."},
	}

	assert.Equal(t, expect, target.Metadata())
}

func setup() Module[types.Ed25519PublicKey] {
	mockSystem = new(mocks.SystemModule)
	mockMemoryUtils = new(mocks.MemoryTranslator)

	target := New[types.Ed25519PublicKey](mockSystem)
	target.memUtils = mockMemoryUtils

	return target
}
