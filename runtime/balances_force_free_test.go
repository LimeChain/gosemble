package main

import (
	"bytes"
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/pkg/scale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func Test_Balances_ForceFree_BadOrigin(t *testing.T) {
	rt, storage := newTestRuntime(t)
	runtimeVersion := rt.Version()

	metadata := runtimeMetadata(t)

	alice, err := ctypes.NewMultiAddressFromAccountID(signature.TestKeyringPairAlice.PublicKey)

	call, err := ctypes.NewCall(metadata, "Balances.force_unreserve", alice, ctypes.NewUCompactFromUInt(10000000000))
	assert.NoError(t, err)

	// Create the extrinsic
	ext := ctypes.NewExtrinsic(call)
	o := ctypes.SignatureOptions{
		BlockHash:          ctypes.Hash(parentHash),
		Era:                ctypes.ExtrinsicEra{IsImmortalEra: true},
		GenesisHash:        ctypes.Hash(parentHash),
		Nonce:              ctypes.NewUCompactFromUInt(0),
		SpecVersion:        ctypes.U32(runtimeVersion.SpecVersion),
		Tip:                ctypes.NewUCompactFromUInt(0),
		TransactionVersion: ctypes.U32(runtimeVersion.TransactionVersion),
	}

	mockBalance, ok := big.NewInt(0).SetString("500000000000000", 10)
	assert.True(t, ok)

	accountInfo := gossamertypes.AccountInfo{
		Nonce:       0,
		Consumers:   0,
		Producers:   0,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(mockBalance),
			Reserved:   scale.MustNewUint128(big.NewInt(0)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}

	aliceHash, _ := common.Blake2b128(signature.TestKeyringPairAlice.PublicKey)
	keyStorageAccountAlice := append(keySystemHash, keyAccountHash...)
	keyStorageAccountAlice = append(keyStorageAccountAlice, aliceHash...)
	keyStorageAccountAlice = append(keyStorageAccountAlice, signature.TestKeyringPairAlice.PublicKey...)

	bytesStorage, err := scale.Marshal(accountInfo)
	assert.NoError(t, err)

	err = storage.Put(keyStorageAccountAlice, bytesStorage)
	assert.NoError(t, err)

	// Sign the transaction using Alice's default account
	err = ext.Sign(signature.TestKeyringPairAlice, o)
	assert.NoError(t, err)

	extEnc := bytes.Buffer{}
	encoder := cscale.NewEncoder(&extEnc)
	err = ext.Encode(*encoder)
	assert.NoError(t, err)

	header := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, blockNumber, gossamertypes.NewDigest())
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

	res, err := rt.Exec("BlockBuilder_apply_extrinsic", extEnc.Bytes())
	assert.NoError(t, err)

	expectedResult :=
		primitives.NewApplyExtrinsicResult(
			primitives.NewDispatchOutcome(
				primitives.NewDispatchErrorBadOrigin()))

	assert.Equal(t, expectedResult.Bytes(), res)
}
