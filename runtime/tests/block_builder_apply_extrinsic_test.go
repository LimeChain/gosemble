package main

import (
	"bytes"
	"testing"
	"time"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/runtime/wasmer"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/aura"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	call := types.Call{
		CallIndex: types.CallIndex{
			ModuleIndex:   system.Module.Index,
			FunctionIndex: system.Module.Functions["remark"].Index,
		},
		Args: sc.Sequence[sc.U8]{0xab, 0xcd},
	}

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

	call := types.Call{
		CallIndex: types.CallIndex{
			ModuleIndex:   system.Module.Index,
			FunctionIndex: system.Module.Functions["remark"].Index,
		},
		Args: sc.Sequence[sc.U8]{0xab, 0xcd},
	}
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

	call := types.Call{
		CallIndex: types.CallIndex{
			ModuleIndex:   system.Module.Index,
			FunctionIndex: system.Module.Functions["remark"].Index,
		},
		Args: sc.Sequence[sc.U8]{0xab, 0xcd},
	}

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
