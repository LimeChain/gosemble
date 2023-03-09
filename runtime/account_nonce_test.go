package main

import (
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_AccountNonceApi_account_nonce_Empty(t *testing.T) {
	pubKey1 := common.MustHexToBytes("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee")

	rt, _ := newTestRuntime(t)

	result, err := rt.Exec("AccountNonceApi_account_nonce", pubKey1)
	assert.NoError(t, err)

	assert.Equal(t, sc.U32(0).Bytes(), result)
}

func Test_AccountNonceApi_account_nonce(t *testing.T) {
	pubKey1 := common.MustHexToBytes("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee")

	accountInfo := gossamertypes.AccountInfo{
		Nonce:       1,
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

	rt, storage := newTestRuntime(t)

	hash, _ := common.Blake2b128(pubKey1)
	key := append(keySystemHash, keyAccountHash...)
	key = append(key, hash...)
	key = append(key, pubKey1...)

	bytesStorage, err := scale.Marshal(accountInfo)
	assert.NoError(t, err)

	err = storage.Put(key, bytesStorage)
	assert.NoError(t, err)

	result, err := rt.Exec("AccountNonceApi_account_nonce", pubKey1)
	assert.NoError(t, err)

	assert.Equal(t, sc.U32(accountInfo.Nonce).Bytes(), result)
}
