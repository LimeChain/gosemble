package session_keys

import (
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	dataPtr    = int32(0)
	dataLen    = int32(1)
	ptrAndSize = int64(2)

	seed         = sc.NewOption[sc.Sequence[sc.U8]](sc.BytesToSequenceU8([]byte("test-seed")))
	keyTypeId    = [4]byte{'t', 'e', 's', 't'}
	publicKey    = []byte("test-public-key")
	seqPublicKey = sc.BytesToSequenceU8(publicKey)
)

var (
	mockCrypto      *mocks.IoCrypto
	mockMemoryUtils *mocks.MemoryTranslator
	mockSessionKey  *mocks.AuraModule
)

func Test_Module_Name(t *testing.T) {
	target := setup()

	result := target.Name()

	assert.Equal(t, ApiModuleName, result)
}

func Test_Module_Item(t *testing.T) {
	target := setup()

	hexName := common.MustBlake2b8([]byte(ApiModuleName))
	expect := primitives.NewApiItem(hexName, apiVersion)

	result := target.Item()

	assert.Equal(t, expect, result)
}

func Test_Module_GenerateSessionKeys_Ed25519(t *testing.T) {
	target := setup()

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(seed.Bytes())
	mockSessionKey.On("KeyType").Return(primitives.PublicKeyEd25519)
	mockSessionKey.On("KeyTypeId").Return(keyTypeId)
	mockCrypto.On("Ed25519Generate", keyTypeId[:], seed.Bytes()).Return(publicKey)
	mockMemoryUtils.On("BytesToOffsetAndSize", seqPublicKey.Bytes()).Return(ptrAndSize)

	result := target.GenerateSessionKeys(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockSessionKey.AssertCalled(t, "KeyType")
	mockSessionKey.AssertCalled(t, "KeyTypeId")
	mockCrypto.AssertCalled(t, "Ed25519Generate", keyTypeId[:], seed.Bytes())
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", seqPublicKey.Bytes())
}

func Test_Module_GenerateSessionKeys_Sr25519(t *testing.T) {
	target := setup()

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(seed.Bytes())
	mockSessionKey.On("KeyType").Return(primitives.PublicKeySr25519)
	mockSessionKey.On("KeyTypeId").Return(keyTypeId)
	mockCrypto.On("Sr25519Generate", keyTypeId[:], seed.Bytes()).Return(publicKey)
	mockMemoryUtils.On("BytesToOffsetAndSize", seqPublicKey.Bytes()).Return(ptrAndSize)

	result := target.GenerateSessionKeys(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockSessionKey.AssertCalled(t, "KeyType")
	mockSessionKey.AssertCalled(t, "KeyTypeId")
	mockCrypto.AssertCalled(t, "Sr25519Generate", keyTypeId[:], seed.Bytes())
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", seqPublicKey.Bytes())
}

func Test_Module_GenerateSessionKeys_Ecdsa(t *testing.T) {
	target := setup()

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(seed.Bytes())
	mockSessionKey.On("KeyType").Return(primitives.PublicKeyEcdsa)
	mockSessionKey.On("KeyTypeId").Return(keyTypeId)
	mockCrypto.On("EcdsaGenerate", keyTypeId[:], seed.Bytes()).Return(publicKey)
	mockMemoryUtils.On("BytesToOffsetAndSize", seqPublicKey.Bytes()).Return(ptrAndSize)

	result := target.GenerateSessionKeys(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockSessionKey.AssertCalled(t, "KeyType")
	mockSessionKey.AssertCalled(t, "KeyTypeId")
	mockCrypto.AssertCalled(t, "EcdsaGenerate", keyTypeId[:], seed.Bytes())
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", seqPublicKey.Bytes())
}

func Test_Module_GenerateSessionKeys_InvalidPublicKeyType(t *testing.T) {
	target := setup()
	invalidKeyType := sc.U8(5)

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(seed.Bytes())
	mockSessionKey.On("KeyType").Return(invalidKeyType)
	mockSessionKey.On("KeyTypeId").Return(keyTypeId)
	mockCrypto.On("EcdsaGenerate", keyTypeId[:], seed.Bytes()).Return(publicKey)
	mockMemoryUtils.On("BytesToOffsetAndSize", seqPublicKey.Bytes()).Return(ptrAndSize)

	assert.PanicsWithValue(t, "invalid public key type", func() {
		target.GenerateSessionKeys(dataPtr, dataLen)
	})

	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockSessionKey.AssertCalled(t, "KeyType")
	mockSessionKey.AssertNotCalled(t, "KeyTypeId")
	mockCrypto.AssertNotCalled(t, "EcdsaGenerate", keyTypeId[:], seed.Bytes())
	mockMemoryUtils.AssertNotCalled(t, "BytesToOffsetAndSize", seqPublicKey.Bytes())
}

func Test_Module_DecodeSessionKeys(t *testing.T) {
	target := setup()

	key := common.MustHexToBytes("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee")
	bytes := sc.BytesToSequenceU8(key).Bytes()
	sessionKeys := sc.Sequence[primitives.SessionKey]{
		primitives.NewSessionKey(key, keyTypeId),
	}
	expect := sc.NewOption[sc.Sequence[primitives.SessionKey]](sessionKeys)

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(bytes)
	mockSessionKey.On("KeyTypeId").Return(keyTypeId)
	mockMemoryUtils.On("BytesToOffsetAndSize", expect.Bytes()).Return(ptrAndSize)

	result := target.DecodeSessionKeys(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockSessionKey.AssertCalled(t, "KeyTypeId")
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", expect.Bytes())
}

func setup() Module[primitives.Ed25519Signer] {
	mockCrypto = new(mocks.IoCrypto)
	mockMemoryUtils = new(mocks.MemoryTranslator)
	mockSessionKey = new(mocks.AuraModule)

	sessions := []primitives.Session{
		mockSessionKey,
	}

	target := New(sessions)
	target.crypto = mockCrypto
	target.memUtils = mockMemoryUtils

	return Module[primitives.Ed25519Signer](target)
}
