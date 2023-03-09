package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

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
)

var (
	parentHash     = common.MustHexToHash("0x0f6d3477739f8a65886135f58c83ff7c2d4a8300a010dfc8b4c5d65ba37920bb")
	stateRoot      = common.MustHexToHash("0xd9e8bf89bda43fb46914321c371add19b81ff92ad6923e8f189b52578074b073")
	extrinsicsRoot = common.MustHexToHash("0x311f481f0ad8739cc513de030c2b99cb6539438560f282ca6fba6e44e8a68120")
	blockNumber    = uint(1)
	sealDigest     = gossamertypes.SealDigest{
		ConsensusEngineID: gossamertypes.BabeEngineID,
		// bytes for SealDigest that was created in setupHeaderFile function
		Data: []byte{158, 127, 40, 221, 220, 242, 124, 30, 107, 50, 141, 86, 148, 195, 104, 213, 178, 236, 93, 190,
			14, 65, 42, 225, 201, 143, 136, 213, 59, 228, 216, 80, 47, 172, 87, 31, 63, 25, 201, 202, 175, 40, 26,
			103, 51, 25, 36, 30, 12, 80, 149, 166, 131, 173, 52, 49, 98, 4, 8, 138, 54, 164, 189, 134},
	}
)

func newTestKeyPair() ([]byte, []byte) {
	privKey := []byte{
		0x11, 0xb2, 0x1e, 0x9d, 0xd8, 0xd9, 0x22, 0x61,
		0xe2, 0xf5, 0xa4, 0xa5, 0x93, 0xf5, 0x7a, 0xd1,
		0xce, 0xd5, 0xbf, 0x0d, 0x94, 0xb8, 0xdc, 0x06,
		0x2d, 0xb1, 0x11, 0x42, 0x7d, 0x3b, 0xf6, 0x35,
		0x15, 0xb0, 0x7f, 0xe2, 0xe7, 0x81, 0x87, 0x4a,
		0xd9, 0x7f, 0xbe, 0x3f, 0xcb, 0xf9, 0xab, 0xaf,
		0x8e, 0x96, 0x5d, 0x2d, 0xb5, 0x30, 0xba, 0xb0,
		0x89, 0xc1, 0xf3, 0xaa, 0x21, 0xf4, 0x20, 0x63,
	}

	pubKey := []byte{
		0x15, 0xb0, 0x7f, 0xe2, 0xe7, 0x81, 0x87, 0x4a,
		0xd9, 0x7f, 0xbe, 0x3f, 0xcb, 0xf9, 0xab, 0xaf,
		0x8e, 0x96, 0x5d, 0x2d, 0xb5, 0x30, 0xba, 0xb0,
		0x89, 0xc1, 0xf3, 0xaa, 0x21, 0xf4, 0x20, 0x63,
	}

	return privKey, pubKey
}

func newTestSigner() types.MultiAddress {
	_, pubKey := newTestKeyPair()
	return types.NewMultiAddress32(types.NewAddress32(sc.BytesToFixedSequenceU8(pubKey)...))
}

func newTestSignature(hexSig string) types.MultiSignature {
	bytes, err := hex.DecodeString(hexSig)
	if err != nil {
		panic(err)
	}
	res := []sc.U8{}
	for _, b := range bytes {
		res = append(res, sc.U8(b))
	}

	return types.NewMultiSignatureEd25519(types.NewEd25519(res...))
}

func signEd25519(digest []byte, privKey []byte) []byte {
	return ed25519.Sign(privKey, digest[:])
}

func newTestExtra(era types.Era, nonce sc.U32, fee sc.U64) types.SignedExtra {
	return types.SignedExtra{
		Era:   era,
		Nonce: nonce,
		Fee:   fee,
	}
}

func newTestCall(moduleIndex sc.U8, functionIndex sc.U8, args ...byte) types.Call {
	return types.Call{
		CallIndex: types.CallIndex{
			ModuleIndex:   sc.U8(moduleIndex),
			FunctionIndex: sc.U8(functionIndex),
		},
		Args: sc.BytesToSequenceU8(args),
	}
}

func setBlockNumber(t *testing.T, storage *trie.Trie, blockNumber sc.U64) {
	blockNumberBytes, err := scale.Marshal(uint64(blockNumber))
	assert.NoError(t, err)

	systemHash := hashing.Twox128(constants.KeySystem)
	numberHash := hashing.Twox128(constants.KeyNumber)
	err = storage.Put(append(systemHash, numberHash...), blockNumberBytes)
	assert.NoError(t, err)
}
