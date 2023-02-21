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

// const WASM_RUNTIME = "../build/polkadot_runtime-v9370.compact.compressed.wasm"
// const WASM_RUNTIME = "../build/westend_runtime-v9370.compact.compressed.wasm"
// const WASM_RUNTIME = "../build/node_template_runtime.wasm"
// const WASM_RUNTIME = "../build/runtime-optimized.wasm" // min memory: 257
const WASM_RUNTIME = "../build/runtime.wasm"
