package main

import (
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/runtime"
	"github.com/ChainSafe/gossamer/lib/runtime/wasmer"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/assert"
)

const POLKADOT_RUNTIME = "../build/polkadot_runtime-v9370.compact.compressed.wasm"

// const WASM_RUNTIME = "../build/polkadot_runtime-v9370.compact.compressed.wasm"
// const WASM_RUNTIME = "../build/westend_runtime-v9370.compact.compressed.wasm"
// const WASM_RUNTIME = "../build/node_template_runtime.wasm"
// const WASM_RUNTIME = "../build/runtime-optimized.wasm" // min memory: 257
const WASM_RUNTIME = "../build/runtime.wasm"

var (
	keySystemHash, _           = common.Twox128Hash(constants.KeySystem)
	keyAccountHash, _          = common.Twox128Hash(constants.KeyAccount)
	keyAllExtrinsicsLenHash, _ = common.Twox128Hash(constants.KeyAllExtrinsicsLen)
	keyAuraHash, _             = common.Twox128Hash(constants.KeyAura)
	keyAuthoritiesHash, _      = common.Twox128Hash(constants.KeyAuthorities)
	keyBlockHash, _            = common.Twox128Hash(constants.KeyBlockHash)
	keyCurrentSlotHash, _      = common.Twox128Hash(constants.KeyCurrentSlot)
	keyDigestHash, _           = common.Twox128Hash(constants.KeyDigest)
	keyExecutionPhaseHash, _   = common.Twox128Hash(constants.KeyExecutionPhase)
	keyExtrinsicCountHash, _   = common.Twox128Hash(constants.KeyExtrinsicCount)
	keyExtrinsicDataHash, _    = common.Twox128Hash(constants.KeyExtrinsicData)
	keyLastRuntime, _          = common.Twox128Hash(constants.KeyLastRuntimeUpgrade)
	keyNumberHash, _           = common.Twox128Hash(constants.KeyNumber)
	keyParentHash, _           = common.Twox128Hash(constants.KeyParentHash)
	keyTimestampHash, _        = common.Twox128Hash(constants.KeyTimestamp)
	keyTimestampNowHash, _     = common.Twox128Hash(constants.KeyNow)
	keyTimestampDidUpdate, _   = common.Twox128Hash(constants.KeyDidUpdate)
	keyBlockWeight, _          = common.Twox128Hash(constants.KeyBlockWeight)
)

var (
	parentHash     = common.MustHexToHash("0x0f6d3477739f8a65886135f58c83ff7c2d4a8300a010dfc8b4c5d65ba37920bb")
	stateRoot      = common.MustHexToHash("0xd9e8bf89bda43fb46914321c371add19b81ff92ad6923e8f189b52578074b073")
	extrinsicsRoot = common.MustHexToHash("0x041a2ad6a911e57dd3135882449ac1cee00981b111b8b782ddffc93206fbe86b")
	blockNumber    = uint(1)
	sealDigest     = gossamertypes.SealDigest{
		ConsensusEngineID: gossamertypes.BabeEngineID,
		// bytes for SealDigest that was created in setupHeaderFile function
		Data: []byte{158, 127, 40, 221, 220, 242, 124, 30, 107, 50, 141, 86, 148, 195, 104, 213, 178, 236, 93, 190,
			14, 65, 42, 225, 201, 143, 136, 213, 59, 228, 216, 80, 47, 172, 87, 31, 63, 25, 201, 202, 175, 40, 26,
			103, 51, 25, 36, 30, 12, 80, 149, 166, 131, 173, 52, 49, 98, 4, 8, 138, 54, 164, 189, 134},
	}
)

func newTestRuntime(t *testing.T) (*wasmer.Instance, *runtime.Storage) {
	runtime := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, trie.NewEmptyTrie())
	storage := &runtime.GetContext().Storage
	return runtime, storage
}

// TODO: Remove once metadata() is implemented
func runtimeMetadata(t *testing.T) *ctypes.Metadata {
	storage := trie.NewEmptyTrie()
	polkadotRuntime := wasmer.NewTestInstanceWithTrie(t, POLKADOT_RUNTIME, storage) // Used only for call encodings
	bMetadata, err := polkadotRuntime.Metadata()
	assert.NoError(t, err)

	var decoded []byte
	err = scale.Unmarshal(bMetadata, &decoded)
	assert.NoError(t, err)

	metadata := &ctypes.Metadata{}
	err = codec.Decode(decoded, metadata)
	assert.NoError(t, err)

	return metadata
}

func setBlockNumber(t *testing.T, storage *runtime.Storage, blockNumber sc.U64) {
	blockNumberBytes, err := scale.Marshal(uint64(blockNumber))
	assert.NoError(t, err)

	systemHash := hashing.Twox128(constants.KeySystem)
	numberHash := hashing.Twox128(constants.KeyNumber)
	err = (*storage).Put(append(systemHash, numberHash...), blockNumberBytes)
	assert.NoError(t, err)
}

func setStorageAccountInfo(t *testing.T, storage *runtime.Storage, account []byte, freeBalance *big.Int, nonce uint32) (storageKey []byte, info gossamertypes.AccountInfo) {
	accountInfo := gossamertypes.AccountInfo{
		Nonce:       nonce,
		Consumers:   0,
		Producers:   0,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(freeBalance),
			Reserved:   scale.MustNewUint128(big.NewInt(0)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}

	aliceHash, _ := common.Blake2b128(account)
	keyStorageAccount := append(keySystemHash, keyAccountHash...)
	keyStorageAccount = append(keyStorageAccount, aliceHash...)
	keyStorageAccount = append(keyStorageAccount, account...)

	bytesStorage, err := scale.Marshal(accountInfo)
	assert.NoError(t, err)

	err = (*storage).Put(keyStorageAccount, bytesStorage)
	assert.NoError(t, err)

	return keyStorageAccount, accountInfo
}
