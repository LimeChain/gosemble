package main

import (
	"bytes"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/runtime/wasmer"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateTransaction(t *testing.T) {
	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	storageRoot := common.MustHexToHash("0x733cbee365f04eb93cd369eeaaf47bb94c1c98603944ba43c39b33070ae90880") // Depends on timestamp
	digest := gossamertypes.NewDigest()

	header := gossamertypes.NewHeader(parentHash, storageRoot, extrinsicsRoot, blockNumber, digest)
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

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

	tx := types.NewSignedUncheckedExtrinsic(call, signer, signature, extra)
	txSource := types.NewTransactionSource(types.External)
	blockHash := types.Blake2bHash{FixedSequence: sc.BytesToFixedSequenceU8(parentHash.ToBytes())}

	buffer := &bytes.Buffer{}
	txSource.Encode(buffer)
	tx.Encode(buffer)
	blockHash.Encode(buffer)

	// TODO: substrate (source, exttrinsic, hash)
	// buffer.Reset()
	// buffer.Write([]byte{2})
	// buffer.Write([]byte{49, 2, 132, 0, 6, 196, 28, 36, 60, 116, 41, 76, 197, 21, 40, 124, 17, 142, 128, 189, 115, 168, 219, 199, 151, 158, 208, 8, 177, 131, 105, 116, 42, 17, 129, 26, 1, 60, 192, 208, 181, 87, 44, 143, 114, 72, 255, 50, 152, 244, 18, 67, 236, 14, 64, 195, 182, 131, 122, 108, 125, 212, 102, 99, 219, 120, 29, 8, 33, 198, 29, 127, 51, 5, 172, 155, 106, 153, 41, 233, 18, 4, 113, 57, 91, 211, 184, 81, 203, 162, 155, 162, 73, 255, 43, 179, 45, 160, 43, 252, 132, 10, 0, 0, 0, 0, 6, 0, 0, 212, 53, 147, 199, 21, 253, 211, 28, 97, 20, 26, 189, 4, 169, 159, 214, 130, 44, 133, 88, 133, 76, 205, 227, 154, 86, 132, 231, 165, 109, 162, 125, 0})
	// buffer.Write([]byte{69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69})

	encTransactionValidityResult, err := rt.Exec("TaggedTransactionQueue_validate_transaction", buffer.Bytes())
	assert.NoError(t, err)

	buffer.Reset()
	buffer.Write(encTransactionValidityResult)
	transactionValidityResult := types.DecodeTransactionValidityResult(buffer)

	assert.Equal(t,
		types.NewTransactionValidityResult(types.ValidTransaction{}),
		transactionValidityResult,
	)
}
