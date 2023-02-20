package main

import (
	"bytes"
	"github.com/LimeChain/gosemble/constants/aura"
	"github.com/LimeChain/gosemble/constants/timestamp"
	primitivestrie "github.com/LimeChain/gosemble/primitives/trie"
	"testing"
	"time"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	runtimetypes "github.com/ChainSafe/gossamer/lib/runtime"
	"github.com/ChainSafe/gossamer/lib/runtime/wasmer"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
)

// const WASM_RUNTIME = "../../build/polkadot_runtime-v9370.compact.compressed.wasm"
// const WASM_RUNTIME = "../../build/westend_runtime-v9370.compact.compressed.wasm"
// const WASM_RUNTIME = "../../build/node_template_runtime.wasm"
// const WASM_RUNTIME = "../../build/runtime-optimized.wasm" // min memory: 257
const WASM_RUNTIME = "../../build/runtime.wasm"

func Test_CoreVersion(t *testing.T) {
	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	res, err := rt.Exec("Core_version", []byte{})
	assert.NoError(t, err)

	buffer := bytes.Buffer{}
	buffer.Write(res)
	dec := scale.NewDecoder(&buffer)
	runtimeVersion := runtimetypes.Version{}
	err = dec.Decode(&runtimeVersion)
	assert.NoError(t, err)
	assert.Equal(t, "node-template", string(runtimeVersion.SpecName))
	assert.Equal(t, "node-template", string(runtimeVersion.ImplName))
	assert.Equal(t, uint32(1), runtimeVersion.AuthoringVersion)
	assert.Equal(t, uint32(100), runtimeVersion.SpecVersion)
	assert.Equal(t, uint32(1), runtimeVersion.ImplVersion)
	assert.Equal(t, uint32(1), runtimeVersion.TransactionVersion)
	assert.Equal(t, uint32(1), runtimeVersion.StateVersion)
}

func Test_CoreInitializeBlock(t *testing.T) {
	preRuntimeDigest := gossamertypes.PreRuntimeDigest{
		ConsensusEngineID: gossamertypes.BabeEngineID,
		// bytes for PreRuntimeDigest that was created in setupHeaderFile function
		Data: []byte{1, 60, 0, 0, 0, 150, 89, 189, 15, 0, 0, 0, 0, 112, 237, 173, 28, 144, 100, 255,
			247, 140, 177, 132, 53, 34, 61, 138, 218, 245, 234, 4, 194, 75, 26, 135, 102, 227, 220, 1, 235, 3, 204,
			106, 12, 17, 183, 151, 147, 212, 227, 28, 192, 153, 8, 56, 34, 156, 68, 254, 209, 102, 154, 124, 124,
			121, 225, 230, 208, 169, 99, 116, 214, 73, 103, 40, 6, 157, 30, 247, 57, 226, 144, 73, 122, 14, 59, 114,
			143, 168, 143, 203, 221, 58, 85, 4, 224, 239, 222, 2, 66, 231, 168, 6, 221, 79, 169, 38, 12},
	}

	expectedStorageDigest := gossamertypes.NewDigest()

	digest := gossamertypes.NewDigest()

	preRuntimeDigestItem := gossamertypes.NewDigestItem()
	assert.NoError(t, preRuntimeDigestItem.Set(preRuntimeDigest))

	sealDigestItem := gossamertypes.NewDigestItem()
	assert.NoError(t, sealDigestItem.Set(sealDigest))

	prdi, err := preRuntimeDigestItem.Value()
	assert.NoError(t, err)
	assert.NoError(t, digest.Add(prdi))

	sdi, err := sealDigestItem.Value()
	assert.NoError(t, err)
	assert.NoError(t, digest.Add(sdi))
	assert.NoError(t, expectedStorageDigest.Add(prdi))

	header := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, blockNumber, digest)
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

	lrui := types.LastRuntimeUpgradeInfo{
		SpecVersion: sc.ToCompact(constants.SPEC_VERSION),
		SpecName:    constants.SPEC_NAME,
	}
	assert.Equal(t, lrui.Bytes(), storage.Get(append(keySystemHash, keyLastRuntime...)))

	encExtrinsicIndex0, _ := scale.Marshal(uint32(0))
	assert.Equal(t, encExtrinsicIndex0, storage.Get(constants.KeyExtrinsicIndex))

	encExecutionPhaseApplyExtrinsic, _ := scale.Marshal(uint32(0))
	assert.Equal(t, encExecutionPhaseApplyExtrinsic, storage.Get(append(keySystemHash, keyExecutionPhaseHash...)))

	encBlockNumber, _ := scale.Marshal(uint32(blockNumber))
	assert.Equal(t, encBlockNumber, storage.Get(append(keySystemHash, keyNumberHash...)))

	encExpectedDigest, err := scale.Marshal(expectedStorageDigest)
	assert.NoError(t, err)
	assert.Equal(t, encExpectedDigest, storage.Get(append(keySystemHash, keyDigestHash...)))
	assert.Equal(t, parentHash.ToBytes(), storage.Get(append(keySystemHash, keyParentHash...)))

	blockHashKey := append(keySystemHash, keyBlockHash...)
	encPrevBlock, _ := scale.Marshal(uint32(blockNumber - 1))
	numHash, err := common.Twox64(encPrevBlock)
	assert.NoError(t, err)

	blockHashKey = append(blockHashKey, numHash...)
	blockHashKey = append(blockHashKey, encPrevBlock...)
	assert.Equal(t, parentHash.ToBytes(), storage.Get(blockHashKey))
}

func Test_BlockBuilder_Inherent_Extrinsics(t *testing.T) {
	idata := gossamertypes.NewInherentData()
	time := time.Now().UnixMilli()
	err := idata.SetInherent(gossamertypes.Timstap0, uint64(time))

	assert.NoError(t, err)

	expectedExtrinsic := types.UncheckedExtrinsic{
		Version: types.ExtrinsicFormatVersion,
		Function: types.Call{
			CallIndex: types.CallIndex{
				ModuleIndex:   timestamp.ModuleIndex,
				FunctionIndex: timestamp.FunctionIndex,
			},
			Args: sc.BytesToSequenceU8(sc.U64(time).Bytes()),
		},
	}

	ienc, err := idata.Encode()
	assert.NoError(t, err)

	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	inherentExt, err := rt.Exec("BlockBuilder_inherent_extrinsics", ienc)
	assert.NoError(t, err)

	assert.NotNil(t, inherentExt)

	buffer := &bytes.Buffer{}
	buffer.Write([]byte{inherentExt[0]})

	totalInherents := sc.DecodeCompact(buffer)
	assert.Equal(t, int64(1), totalInherents.ToBigInt().Int64())
	buffer.Reset()

	buffer.Write(inherentExt[1:])
	extrinsic := types.DecodeUncheckedExtrinsic(buffer)

	assert.Equal(t, expectedExtrinsic, extrinsic)
}

func Test_ApplyExtrinsic_Timestamp(t *testing.T) {
	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	bytesSlotDuration, err := rt.Exec("AuraApi_slot_duration", []byte{})
	assert.NoError(t, err)

	idata := gossamertypes.NewInherentData()
	time := time.Now().UnixMilli()

	buffer := &bytes.Buffer{}
	buffer.Write(bytesSlotDuration)

	slotDuration := sc.DecodeU64(buffer)
	buffer.Reset()

	slot := sc.U64(time) / slotDuration

	preRuntimeDigest := gossamertypes.PreRuntimeDigest{
		ConsensusEngineID: aura.EngineId,
		Data:              slot.Bytes(),
	}

	digest := gossamertypes.NewDigest()
	assert.NoError(t, digest.Add(preRuntimeDigest))

	header := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, blockNumber, digest)
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

	err = idata.SetInherent(gossamertypes.Timstap0, uint64(time))
	assert.NoError(t, err)

	ienc, err := idata.Encode()
	assert.NoError(t, err)
	inherentExt, err := rt.Exec("BlockBuilder_inherent_extrinsics", ienc)
	assert.NoError(t, err)

	applyResult, err := rt.Exec("BlockBuilder_apply_extrinsic", inherentExt[1:])
	assert.NoError(t, err)

	assert.Equal(t,
		types.NewApplyExtrinsicResult(types.NewDispatchOutcome(nil)).Bytes(),
		applyResult,
	)

	assert.Equal(t, []byte{1}, storage.Get(append(keyTimestampHash, keyTimestampDidUpdate...)))
	assert.Equal(t, sc.U64(time).Bytes(), storage.Get(append(keyTimestampHash, keyTimestampNowHash...)))

	assert.Equal(t, slot.Bytes(), storage.Get(append(keyAuraHash, keyCurrentSlotHash...)))
}

func Test_ApplyExtrinsic_DispatchOutcome(t *testing.T) {
	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	call := types.NewCall("System", "remark", sc.Sequence[sc.U8]{0xab, 0xcd})

	// privKey - 0x11, 0xb2, 0x1e, 0x9d, 0xd8, 0xd9, 0x22, 0x61, 0xe2, 0xf5, 0xa4, 0xa5, 0x93, 0xf5, 0x7a, 0xd1, 0xce, 0xd5, 0xbf, 0x0d, 0x94, 0xb8, 0xdc, 0x06, 0x2d, 0xb1, 0x11, 0x42, 0x7d, 0x3b, 0xf6, 0x35, 0x15, 0xb0, 0x7f, 0xe2, 0xe7, 0x81, 0x87, 0x4a, 0xd9, 0x7f, 0xbe, 0x3f, 0xcb, 0xf9, 0xab, 0xaf, 0x8e, 0x96, 0x5d, 0x2d, 0xb5, 0x30, 0xba, 0xb0, 0x89, 0xc1, 0xf3, 0xaa, 0x21, 0xf4, 0x20, 0x63
	signer := types.NewMultiAddress(
		types.NewAddress32(
			0x15, 0xb0, 0x7f, 0xe2, 0xe7, 0x81, 0x87, 0x4a, 0xd9, 0x7f, 0xbe, 0x3f, 0xcb, 0xf9, 0xab, 0xaf, 0x8e, 0x96, 0x5d, 0x2d, 0xb5, 0x30, 0xba, 0xb0, 0x89, 0xc1, 0xf3, 0xaa, 0x21, 0xf4, 0x20, 0x63,
		),
	)

	signature := types.NewMultiSignature(
		types.NewEd25519(
			0xb5, 0x1b, 0xc6, 0xf6, 0xe9, 0xf7, 0x07, 0x90, 0x03, 0xad, 0xc9, 0x8e, 0xd2, 0x16, 0x22, 0xd0, 0xd8, 0x87, 0xd6, 0x8e, 0x86, 0x45, 0xb5, 0xcb, 0x99, 0xe3, 0xd6, 0x08, 0x0a, 0xa8, 0xc9, 0xdd, 0x40, 0x8c, 0xa3, 0x9c, 0x91, 0xd9, 0xc7, 0x0a, 0x49, 0xa3, 0x77, 0xbc, 0x2b, 0x55, 0x04, 0xe3, 0x64, 0x27, 0xe1, 0x84, 0x5b, 0x38, 0x20, 0xc5, 0x8c, 0x95, 0xf1, 0x46, 0xf0, 0xce, 0xc2, 0x03,
		),
	)

	extra := types.Extra{
		Era:   types.ExtrinsicEra{},
		Nonce: sc.ToCompact(0),
		Fee:   sc.ToCompact(0),
	}

	uxt := types.NewSignedUncheckedExtrinsic(call, signer, signature, extra)

	res, err := rt.Exec("BlockBuilder_apply_extrinsic", uxt.Bytes())

	currentExtrinsicIndex := sc.U32(1)
	extrinsicIndexValue := rt.GetContext().Storage.Get(constants.KeyExtrinsicIndex)
	require.Equal(t, currentExtrinsicIndex.Bytes(), extrinsicIndexValue)

	keyExtrinsicDataPrefixHash := append(keySystemHash, keyExtrinsicDataHash...)

	prevExtrinsic := currentExtrinsicIndex - 1
	hashIndex, err := common.Twox64(prevExtrinsic.Bytes())
	assert.NoError(t, err)

	keyExtrinsic := append(keyExtrinsicDataPrefixHash, hashIndex...)
	storageUxt := rt.GetContext().Storage.Get(append(keyExtrinsic, prevExtrinsic.Bytes()...))

	require.Equal(t, uxt.Bytes(), storageUxt)

	require.NoError(t, err)

	require.Equal(t,
		types.NewApplyExtrinsicResult(types.NewDispatchOutcome(nil)).Bytes(),
		res,
	)
}

func Test_ApplyExtrinsic_Unsigned_DispatchOutcome(t *testing.T) {
	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	call := types.NewCall("System", "remark", sc.Sequence[sc.U8]{0xab, 0xcd})
	// substrate - [0, 0, 8, 171, 205]
	// gosemble -  [0, 0, 8, 171, 205]

	uxt := types.NewUnsignedUncheckedExtrinsic(call)
	// substrate - [1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 171, 205]
	// gosemble -  [18 04 00 00 08 ab cd]

	res, err := rt.Exec("BlockBuilder_apply_extrinsic", uxt.Bytes())

	require.NoError(t, err)

	require.Equal(t,
		types.NewApplyExtrinsicResult(types.NewDispatchOutcome(nil)).Bytes(),
		res,
	)
}

func Test_ApplyExtrinsic_DispatchError_BadProofError(t *testing.T) {
	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	call := types.NewCall("System", "remark", sc.Sequence[sc.U8]{0xab, 0xcd})

	// privKey - 0x11, 0xb2, 0x1e, 0x9d, 0xd8, 0xd9, 0x22, 0x61, 0xe2, 0xf5, 0xa4, 0xa5, 0x93, 0xf5, 0x7a, 0xd1, 0xce, 0xd5, 0xbf, 0x0d, 0x94, 0xb8, 0xdc, 0x06, 0x2d, 0xb1, 0x11, 0x42, 0x7d, 0x3b, 0xf6, 0x35, 0x15, 0xb0, 0x7f, 0xe2, 0xe7, 0x81, 0x87, 0x4a, 0xd9, 0x7f, 0xbe, 0x3f, 0xcb, 0xf9, 0xab, 0xaf, 0x8e, 0x96, 0x5d, 0x2d, 0xb5, 0x30, 0xba, 0xb0, 0x89, 0xc1, 0xf3, 0xaa, 0x21, 0xf4, 0x20, 0x63
	signer := types.NewMultiAddress(
		types.NewAddress32(
			0x15, 0xb0, 0x7f, 0xe2, 0xe7, 0x81, 0x87, 0x4a, 0xd9, 0x7f, 0xbe, 0x3f, 0xcb, 0xf9, 0xab, 0xaf, 0x8e, 0x96, 0x5d, 0x2d, 0xb5, 0x30, 0xba, 0xb0, 0x89, 0xc1, 0xf3, 0xaa, 0x21, 0xf4, 0x20, 0x63,
		),
	)

	extra := types.Extra{
		Era:   types.ExtrinsicEra{},
		Nonce: sc.ToCompact(1), // instead of 0 to make the signature invalid
		Fee:   sc.ToCompact(0),
	}

	invalidSignature := types.NewMultiSignature(
		types.NewEd25519(
			0xb5, 0x1b, 0xc6, 0xf6, 0xe9, 0xf7, 0x07, 0x90, 0x03, 0xad, 0xc9, 0x8e, 0xd2, 0x16, 0x22, 0xd0, 0xd8, 0x87, 0xd6, 0x8e, 0x86, 0x45, 0xb5, 0xcb, 0x99, 0xe3, 0xd6, 0x08, 0x0a, 0xa8, 0xc9, 0xdd, 0x40, 0x8c, 0xa3, 0x9c, 0x91, 0xd9, 0xc7, 0x0a, 0x49, 0xa3, 0x77, 0xbc, 0x2b, 0x55, 0x04, 0xe3, 0x64, 0x27, 0xe1, 0x84, 0x5b, 0x38, 0x20, 0xc5, 0x8c, 0x95, 0xf1, 0x46, 0xf0, 0xce, 0xc2, 0x03,
		),
	)

	uxt := types.NewSignedUncheckedExtrinsic(call, signer, invalidSignature, extra)

	res, err := rt.Exec("BlockBuilder_apply_extrinsic", uxt.Bytes())

	extrinsicIndex := sc.U32(0)
	extrinsicIndexValue := rt.GetContext().Storage.Get(append(keySystemHash, sc.NewOption[sc.U32](extrinsicIndex).Bytes()...))
	require.Equal(t, []byte(nil), extrinsicIndexValue)

	require.NoError(t, err)

	require.Equal(t,
		types.NewApplyExtrinsicResult(
			types.NewTransactionValidityError(types.NewInvalidTransaction(types.BadProofError)),
		).Bytes(),
		res,
	)
}

func Test_ApplyExtrinsic_InherentsFails(t *testing.T) {
	t.Skip()
}

func Test_CheckInherents(t *testing.T) {
	expectedCheckInherentsResult := types.NewCheckInherentsResult()

	idata := gossamertypes.NewInherentData()
	time := time.Now().UnixMilli()
	err := idata.SetInherent(gossamertypes.Timstap0, uint64(time))

	assert.NoError(t, err)

	ienc, err := idata.Encode()
	assert.NoError(t, err)

	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	inherentExt, err := rt.Exec("BlockBuilder_inherent_extrinsics", ienc)
	assert.NoError(t, err)
	assert.NotNil(t, inherentExt)

	header := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, blockNumber, gossamertypes.NewDigest())

	var exts [][]byte
	err = scale.Unmarshal(inherentExt, &exts)
	assert.NoError(t, err)

	block := gossamertypes.Block{
		Header: *header,
		Body:   gossamertypes.BytesArrayToExtrinsics(exts),
	}

	encodedBlock, err := scale.Marshal(block)
	assert.NoError(t, err)

	inputData := append(encodedBlock, ienc...)
	bytesCheckInherentsResult, err := rt.Exec("BlockBuilder_check_inherents", inputData)
	assert.NoError(t, err)

	buffer := &bytes.Buffer{}
	buffer.Write(bytesCheckInherentsResult)
	checkInherentsResult := types.DecodeCheckInherentsResult(buffer)

	assert.Equal(t, expectedCheckInherentsResult, checkInherentsResult)
}

func Test_BlockExecution(t *testing.T) {
	// core.InitializeBlock
	// blockBuilder.InherentExtrinsics
	// blockBuilder.ApplyExtrinsics
	// blockBuilder.FinalizeBlock

	storageRoot := common.MustHexToHash("0xd9e8bf89bda43fb46914321c371add19b81ff92ad6923e8f189b52578074b073") // Depends on timestamp
	time := time.Date(2023, time.January, 2, 3, 4, 5, 6, time.UTC)

	expectedStorageDigest := gossamertypes.NewDigest()
	digest := gossamertypes.NewDigest()

	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	bytesSlotDuration, err := rt.Exec("AuraApi_slot_duration", []byte{})
	assert.NoError(t, err)

	buffer := &bytes.Buffer{}
	buffer.Write(bytesSlotDuration)

	slotDuration := sc.DecodeU64(buffer)
	buffer.Reset()

	slot := sc.U64(time.UnixMilli()) / slotDuration

	preRuntimeDigest := gossamertypes.PreRuntimeDigest{
		ConsensusEngineID: aura.EngineId,
		Data:              slot.Bytes(),
	}
	assert.NoError(t, digest.Add(preRuntimeDigest))
	assert.NoError(t, expectedStorageDigest.Add(preRuntimeDigest))

	header := gossamertypes.NewHeader(parentHash, storageRoot, extrinsicsRoot, blockNumber, digest)
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

	lrui := types.LastRuntimeUpgradeInfo{
		SpecVersion: sc.ToCompact(constants.SPEC_VERSION),
		SpecName:    constants.SPEC_NAME,
	}
	assert.Equal(t, lrui.Bytes(), storage.Get(append(keySystemHash, keyLastRuntime...)))

	encExtrinsicIndex0, _ := scale.Marshal(uint32(0))
	assert.Equal(t, encExtrinsicIndex0, storage.Get(constants.KeyExtrinsicIndex))

	encExecutionPhaseApplyExtrinsic, _ := scale.Marshal(uint32(0))
	assert.Equal(t, encExecutionPhaseApplyExtrinsic, storage.Get(append(keySystemHash, keyExecutionPhaseHash...)))

	encBlockNumber, _ := scale.Marshal(uint32(blockNumber))
	assert.Equal(t, encBlockNumber, storage.Get(append(keySystemHash, keyNumberHash...)))

	encExpectedDigest, err := scale.Marshal(expectedStorageDigest)
	assert.NoError(t, err)

	assert.Equal(t, encExpectedDigest, storage.Get(append(keySystemHash, keyDigestHash...)))
	assert.Equal(t, parentHash.ToBytes(), storage.Get(append(keySystemHash, keyParentHash...)))

	blockHashKey := append(keySystemHash, keyBlockHash...)
	encPrevBlock, _ := scale.Marshal(uint32(blockNumber - 1))
	numHash, err := common.Twox64(encPrevBlock)
	assert.NoError(t, err)

	blockHashKey = append(blockHashKey, numHash...)
	blockHashKey = append(blockHashKey, encPrevBlock...)
	assert.Equal(t, parentHash.ToBytes(), storage.Get(blockHashKey))

	idata := gossamertypes.NewInherentData()
	err = idata.SetInherent(gossamertypes.Timstap0, uint64(time.UnixMilli()))

	assert.NoError(t, err)

	expectedExtrinsic := types.UncheckedExtrinsic{
		Version: types.ExtrinsicFormatVersion,
		Function: types.Call{
			CallIndex: types.CallIndex{
				ModuleIndex:   timestamp.ModuleIndex,
				FunctionIndex: timestamp.FunctionIndex,
			},
			Args: sc.BytesToSequenceU8(sc.U64(time.UnixMilli()).Bytes()),
		},
	}

	ienc, err := idata.Encode()
	assert.NoError(t, err)

	inherentExt, err := rt.Exec("BlockBuilder_inherent_extrinsics", ienc)
	assert.NoError(t, err)
	assert.NotNil(t, inherentExt)

	buffer.Write([]byte{inherentExt[0]})

	totalInherents := sc.DecodeCompact(buffer)
	assert.Equal(t, int64(1), totalInherents.ToBigInt().Int64())
	buffer.Reset()

	buffer.Write(inherentExt[1:])
	extrinsic := types.DecodeUncheckedExtrinsic(buffer)
	buffer.Reset()

	assert.Equal(t, expectedExtrinsic, extrinsic)

	applyResult, err := rt.Exec("BlockBuilder_apply_extrinsic", inherentExt[1:])
	assert.NoError(t, err)

	assert.Equal(t,
		types.NewApplyExtrinsicResult(types.NewDispatchOutcome(nil)).Bytes(),
		applyResult,
	)

	bytesResult, err := rt.Exec("BlockBuilder_finalize_block", []byte{})
	assert.NoError(t, err)

	resultHeader := gossamertypes.NewEmptyHeader()
	assert.NoError(t, scale.Unmarshal(bytesResult, resultHeader))
	resultHeader.Hash() // Call this to be set, otherwise structs do not match...

	assert.Equal(t, header, resultHeader)

	assert.Equal(t, []byte(nil), storage.Get(append(keyTimestampHash, keyTimestampDidUpdate...)))
	assert.Equal(t, sc.U64(time.UnixMilli()).Bytes(), storage.Get(append(keyTimestampHash, keyTimestampNowHash...)))

	assert.Equal(t, []byte(nil), storage.Get(constants.KeyExtrinsicIndex))
	assert.Equal(t, []byte(nil), storage.Get(append(keySystemHash, keyExecutionPhaseHash...)))
	assert.Equal(t, []byte(nil), storage.Get(append(keySystemHash, keyAllExtrinsicsLenHash...)))
	assert.Equal(t, []byte(nil), storage.Get(append(keySystemHash, keyExtrinsicCountHash...)))

	assert.Equal(t, parentHash.ToBytes(), storage.Get(append(keySystemHash, keyParentHash...)))
	assert.Equal(t, encExpectedDigest, storage.Get(append(keySystemHash, keyDigestHash...)))
	assert.Equal(t, encBlockNumber, storage.Get(append(keySystemHash, keyNumberHash...)))

	assert.Equal(t, slot.Bytes(), storage.Get(append(keyAuraHash, keyCurrentSlotHash...)))
}

func Test_ExecuteBlock(t *testing.T) {
	// blockBuilder.Inherent_Extrinsics
	// blockBuilder.ExecuteBlock

	storageRoot := common.MustHexToHash("0xd9e8bf89bda43fb46914321c371add19b81ff92ad6923e8f189b52578074b073") // Depends on timestamp
	time := time.Date(2023, time.January, 2, 3, 4, 5, 6, time.UTC)

	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	bytesSlotDuration, err := rt.Exec("AuraApi_slot_duration", []byte{})
	assert.NoError(t, err)

	buffer := &bytes.Buffer{}
	buffer.Write(bytesSlotDuration)

	slotDuration := sc.DecodeU64(buffer)
	buffer.Reset()

	slot := sc.U64(time.UnixMilli()) / slotDuration

	preRuntimeDigest := gossamertypes.PreRuntimeDigest{
		ConsensusEngineID: aura.EngineId,
		Data:              slot.Bytes(),
	}

	idata := gossamertypes.NewInherentData()
	err = idata.SetInherent(gossamertypes.Timstap0, uint64(time.UnixMilli()))

	assert.NoError(t, err)

	ienc, err := idata.Encode()
	assert.NoError(t, err)

	expectedExtrinsic := types.UncheckedExtrinsic{
		Version: types.ExtrinsicFormatVersion,
		Function: types.Call{
			CallIndex: types.CallIndex{
				ModuleIndex:   timestamp.ModuleIndex,
				FunctionIndex: timestamp.FunctionIndex,
			},
			Args: sc.BytesToSequenceU8(sc.U64(time.UnixMilli()).Bytes()),
		},
	}

	inherentExt, err := rt.Exec("BlockBuilder_inherent_extrinsics", ienc)
	assert.NoError(t, err)
	assert.NotNil(t, inherentExt)

	buffer.Write([]byte{inherentExt[0]})

	totalInherents := sc.DecodeCompact(buffer)
	assert.Equal(t, int64(1), totalInherents.ToBigInt().Int64())
	buffer.Reset()

	buffer.Write(inherentExt[1:])
	extrinsic := types.DecodeUncheckedExtrinsic(buffer)

	assert.Equal(t, expectedExtrinsic, extrinsic)

	var exts [][]byte
	err = scale.Unmarshal(inherentExt, &exts)
	assert.Nil(t, err)

	extrinsicsRoot := primitivestrie.Blake2256OrderedRoot(inherentExt, constants.StorageVersion)

	digest := gossamertypes.NewDigest()

	assert.NoError(t, err)
	assert.NoError(t, digest.Add(preRuntimeDigest))

	expectedStorageDigest, err := scale.Marshal(digest)
	assert.NoError(t, err)
	encBlockNumber, _ := scale.Marshal(uint32(blockNumber))

	header := gossamertypes.NewHeader(parentHash, storageRoot, common.BytesToHash(extrinsicsRoot), blockNumber, digest)

	block := gossamertypes.Block{
		Header: *header,
		Body:   gossamertypes.BytesArrayToExtrinsics(exts),
	}

	encodedBlock, err := scale.Marshal(block)
	assert.Nil(t, err)

	_, err = rt.Exec("Core_execute_block", encodedBlock)
	assert.NoError(t, err)

	assert.Equal(t, []byte(nil), storage.Get(append(keyTimestampHash, keyTimestampDidUpdate...)))
	assert.Equal(t, sc.U64(time.UnixMilli()).Bytes(), storage.Get(append(keyTimestampHash, keyTimestampNowHash...)))

	assert.Equal(t, []byte(nil), storage.Get(constants.KeyExtrinsicIndex))
	assert.Equal(t, []byte(nil), storage.Get(append(keySystemHash, keyExecutionPhaseHash...)))
	assert.Equal(t, []byte(nil), storage.Get(append(keySystemHash, keyAllExtrinsicsLenHash...)))
	assert.Equal(t, []byte(nil), storage.Get(append(keySystemHash, keyExtrinsicCountHash...)))

	assert.Equal(t, parentHash.ToBytes(), storage.Get(append(keySystemHash, keyParentHash...)))
	assert.Equal(t, expectedStorageDigest, storage.Get(append(keySystemHash, keyDigestHash...)))
	assert.Equal(t, encBlockNumber, storage.Get(append(keySystemHash, keyNumberHash...)))

	assert.Equal(t, slot.Bytes(), storage.Get(append(keyAuraHash, keyCurrentSlotHash...)))
}

func Test_Aura_Authorities_Empty(t *testing.T) {
	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	result, err := rt.Exec("AuraApi_authorities", []byte{})
	assert.NoError(t, err)

	assert.Equal(t, []byte{0}, result)
}

func Test_Aura_Authorities(t *testing.T) {
	pubKey1 := common.MustHexToBytes("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee")
	pubKey2 := common.MustHexToBytes("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ef")

	buffer := &bytes.Buffer{}
	buffer.Write(pubKey1)

	bytesPubKey1, err := common.Read32Bytes(buffer)
	assert.NoError(t, err)

	buffer.Write(pubKey2)
	bytesPubKey2, err := common.Read32Bytes(buffer)
	assert.NoError(t, err)

	authorities := [][32]byte{
		bytesPubKey1,
		bytesPubKey2,
	}

	bytesAuthorities, err := scale.Marshal(authorities)
	assert.NoError(t, err)

	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	err = storage.Put(append(keyAuraHash, keyAuthoritiesHash...), bytesAuthorities)
	assert.NoError(t, err)

	result, err := rt.Exec("AuraApi_authorities", []byte{})
	assert.NoError(t, err)

	assert.Equal(t, bytesAuthorities, result)

	var resultAuthorities [][32]byte
	err = scale.Unmarshal(result, &resultAuthorities)
	assert.NoError(t, err)

	assert.Equal(t, authorities, resultAuthorities)
}
