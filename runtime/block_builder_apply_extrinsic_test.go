package main

import (
	"bytes"
	"math/big"
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

	extra := newTestExtra(types.NewImmortalEra(), 0, 0)
	call := newTestCall(0, 0, 0xab, 0xcd)
	signer := newTestSigner()
	signature := newTestSignature("3594650a752f5b8a4853b66390ea2078b44725246d654ec94509ee2fa89916a74fa8d770f34b81b16adb08899ff223e4f61178e98c42c6a9bfefb987a4176208")
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

	extra := newTestExtra(types.NewImmortalEra(), 1, 0) // instead of 0 to make the signature invalid
	call := newTestCall(0, 0, 0xab, 0xcd)
	signer := newTestSigner()
	invalidSignature := newTestSignature("b51bc6f6e9f7079003adc98ed21622d0d887d68e8645b5cb99e3d6080aa8c9dd408ca39c91d9c70a49a377bc2b5504e36427e1845b3820c58c95f146f0cec203")

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

func Test_ApplyExtrinsic_FutureError(t *testing.T) {
	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	pubKey1 := []byte{0x15, 0xb0, 0x7f, 0xe2, 0xe7, 0x81, 0x87, 0x4a, 0xd9, 0x7f, 0xbe, 0x3f, 0xcb, 0xf9, 0xab, 0xaf, 0x8e, 0x96, 0x5d, 0x2d, 0xb5, 0x30, 0xba, 0xb0, 0x89, 0xc1, 0xf3, 0xaa, 0x21, 0xf4, 0x20, 0x63}

	accountInfo := gossamertypes.AccountInfo{
		Nonce:       3,
		Consumers:   2,
		Producers:   3,
		Sufficients: 4,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(big.NewInt(5)),
			Reserved:   scale.MustNewUint128(big.NewInt(6)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(7)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(8)),
		},
	}

	hash, _ := common.Blake2b128(pubKey1)
	key := append(keySystemHash, keyAccountHash...)
	key = append(key, hash...)
	key = append(key, pubKey1...)

	bytesStorage, err := scale.Marshal(accountInfo)
	assert.NoError(t, err)

	err = storage.Put(key, bytesStorage)
	assert.NoError(t, err)

	storageRoot := common.MustHexToHash("0x733cbee365f04eb93cd369eeaaf47bb94c1c98603944ba43c39b33070ae90880") // Depends on timestamp
	digest := gossamertypes.NewDigest()

	header := gossamertypes.NewHeader(parentHash, storageRoot, extrinsicsRoot, blockNumber, digest)
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

	extra := newTestExtra(types.NewImmortalEra(), 5, 0)
	call := newTestCall(0, 0, 0xab, 0xcd)
	signer := newTestSigner()
	signature := newTestSignature("64278529b1d03ca91fd4435d75899e145328d1af0dfb7ab2e324e49196dbcc22963ec1fa6257bc7476be532b8f91bbd0ec46ccac704517cb39a64103929e7800")
	tx := types.NewSignedUncheckedExtrinsic(call, signer, signature, extra)

	encTransactionValidityResult, err := rt.Exec("BlockBuilder_apply_extrinsic", tx.Bytes())
	assert.NoError(t, err)

	buffer := &bytes.Buffer{}
	buffer.Write(encTransactionValidityResult)
	transactionValidityResult := types.DecodeTransactionValidityResult(buffer)

	assert.Equal(t,
		types.NewTransactionValidityResult(
			types.NewTransactionValidityError(
				types.NewInvalidTransaction(types.FutureError),
			),
		),
		transactionValidityResult,
	)
}
