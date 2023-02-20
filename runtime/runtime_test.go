package main

import (
	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/LimeChain/gosemble/constants"
)

var (
	keySystemHash, _           = common.Twox128Hash(constants.KeySystem)
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
	extrinsicsRoot = common.MustHexToHash("0xfbe77e9def055a8d31a21675651765b9438e338d7ff02760b91dcca8bd6ff0fe")
	blockNumber    = uint(1)
	sealDigest     = gossamertypes.SealDigest{
		ConsensusEngineID: gossamertypes.BabeEngineID,
		// bytes for SealDigest that was created in setupHeaderFile function
		Data: []byte{158, 127, 40, 221, 220, 242, 124, 30, 107, 50, 141, 86, 148, 195, 104, 213, 178, 236, 93, 190,
			14, 65, 42, 225, 201, 143, 136, 213, 59, 228, 216, 80, 47, 172, 87, 31, 63, 25, 201, 202, 175, 40, 26,
			103, 51, 25, 36, 30, 12, 80, 149, 166, 131, 173, 52, 49, 98, 4, 8, 138, 54, 164, 189, 134},
	}

	preRuntimeDigest = gossamertypes.PreRuntimeDigest{
		ConsensusEngineID: gossamertypes.BabeEngineID,
		// bytes for PreRuntimeDigest that was created in setupHeaderFile function
		Data: []byte{1, 60, 0, 0, 0, 150, 89, 189, 15, 0, 0, 0, 0, 112, 237, 173, 28, 144, 100, 255,
			247, 140, 177, 132, 53, 34, 61, 138, 218, 245, 234, 4, 194, 75, 26, 135, 102, 227, 220, 1, 235, 3, 204,
			106, 12, 17, 183, 151, 147, 212, 227, 28, 192, 153, 8, 56, 34, 156, 68, 254, 209, 102, 154, 124, 124,
			121, 225, 230, 208, 169, 99, 116, 214, 73, 103, 40, 6, 157, 30, 247, 57, 226, 144, 73, 122, 14, 59, 114,
			143, 168, 143, 203, 221, 58, 85, 4, 224, 239, 222, 2, 66, 231, 168, 6, 221, 79, 169, 38, 12},
	}
)

// const WASM_RUNTIME = "../build/polkadot_runtime-v9370.compact.compressed.wasm"
// const WASM_RUNTIME = "../build/westend_runtime-v9370.compact.compressed.wasm"
// const WASM_RUNTIME = "../build/node_template_runtime.wasm"
// const WASM_RUNTIME = "../build/runtime-optimized.wasm" // min memory: 257
const WASM_RUNTIME = "../build/runtime.wasm"
