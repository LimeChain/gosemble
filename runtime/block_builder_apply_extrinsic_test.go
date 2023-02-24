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
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
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

	extra := newTestExtra(types.NewEra(types.ImmortalEra{}), 0, 0)
	call := newTestCall(0, 0, 0xab, 0xcd)
	signer := newTestSigner()
	signature := newTestSignature(0x35, 0x94, 0x65, 0x0a, 0x75, 0x2f, 0x5b, 0x8a, 0x48, 0x53, 0xb6, 0x63, 0x90, 0xea, 0x20, 0x78, 0xb4, 0x47, 0x25, 0x24, 0x6d, 0x65, 0x4e, 0xc9, 0x45, 0x09, 0xee, 0x2f, 0xa8, 0x99, 0x16, 0xa7, 0x4f, 0xa8, 0xd7, 0x70, 0xf3, 0x4b, 0x81, 0xb1, 0x6a, 0xdb, 0x08, 0x89, 0x9f, 0xf2, 0x23, 0xe4, 0xf6, 0x11, 0x78, 0xe9, 0x8c, 0x42, 0xc6, 0xa9, 0xbf, 0xef, 0xb9, 0x87, 0xa4, 0x17, 0x62, 0x08)
	uxt := types.NewSignedUncheckedExtrinsic(call, signer, signature, extra)

	res, err := rt.Exec("BlockBuilder_apply_extrinsic", uxt.Bytes())

	currentExtrinsicIndex := sc.U32(1)
	extrinsicIndexValue := rt.GetContext().Storage.Get(constants.KeyExtrinsicIndex)
	assert.Equal(t, currentExtrinsicIndex.Bytes(), extrinsicIndexValue)

	keyExtrinsicDataPrefixHash := append(keySystemHash, keyExtrinsicDataHash...)

	prevExtrinsic := currentExtrinsicIndex - 1
	hashIndex, err := common.Twox64(prevExtrinsic.Bytes())
	assert.NoError(t, err)

	keyExtrinsic := append(keyExtrinsicDataPrefixHash, hashIndex...)
	storageUxt := rt.GetContext().Storage.Get(append(keyExtrinsic, prevExtrinsic.Bytes()...))

	assert.Equal(t, uxt.Bytes(), storageUxt)

	assert.NoError(t, err)

	assert.Equal(t,
		types.NewApplyExtrinsicResult(types.NewDispatchOutcome(nil)).Bytes(),
		res,
	)
}

func Test_ApplyExtrinsic_Unsigned_DispatchOutcome(t *testing.T) {
	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	call := newTestCall(0, 0, 0xab, 0xcd)
	uxt := types.NewUnsignedUncheckedExtrinsic(call)

	res, err := rt.Exec("BlockBuilder_apply_extrinsic", uxt.Bytes())

	assert.NoError(t, err)

	assert.Equal(t,
		types.NewApplyExtrinsicResult(types.NewDispatchOutcome(nil)).Bytes(),
		res,
	)
}

func Test_ApplyExtrinsic_DispatchError_BadProofError(t *testing.T) {
	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	extra := newTestExtra(types.NewEra(types.ImmortalEra{}), 1, 0) // instead of 0 to make the signature invalid
	call := newTestCall(0, 0, 0xab, 0xcd)
	signer := newTestSigner()
	invalidSignature := newTestSignature(0xb5, 0x1b, 0xc6, 0xf6, 0xe9, 0xf7, 0x07, 0x90, 0x03, 0xad, 0xc9, 0x8e, 0xd2, 0x16, 0x22, 0xd0, 0xd8, 0x87, 0xd6, 0x8e, 0x86, 0x45, 0xb5, 0xcb, 0x99, 0xe3, 0xd6, 0x08, 0x0a, 0xa8, 0xc9, 0xdd, 0x40, 0x8c, 0xa3, 0x9c, 0x91, 0xd9, 0xc7, 0x0a, 0x49, 0xa3, 0x77, 0xbc, 0x2b, 0x55, 0x04, 0xe3, 0x64, 0x27, 0xe1, 0x84, 0x5b, 0x38, 0x20, 0xc5, 0x8c, 0x95, 0xf1, 0x46, 0xf0, 0xce, 0xc2, 0x03)
	uxt := types.NewSignedUncheckedExtrinsic(call, signer, invalidSignature, extra)

	res, err := rt.Exec("BlockBuilder_apply_extrinsic", uxt.Bytes())

	extrinsicIndex := sc.U32(0)
	extrinsicIndexValue := rt.GetContext().Storage.Get(append(keySystemHash, sc.NewOption[sc.U32](extrinsicIndex).Bytes()...))
	assert.Equal(t, []byte(nil), extrinsicIndexValue)

	assert.NoError(t, err)

	assert.Equal(t,
		types.NewApplyExtrinsicResult(
			types.NewTransactionValidityError(types.NewInvalidTransaction(types.BadProofError)),
		).Bytes(),
		res,
	)
}

func Test_ApplyExtrinsic_InherentsFails(t *testing.T) {
	t.Skip()
}
