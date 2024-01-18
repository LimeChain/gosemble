package main

import (
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/runtime"
	"github.com/ChainSafe/gossamer/pkg/scale"
	"github.com/stretchr/testify/assert"
)

var (
	existentialAmount = int64(BalancesExistentialDeposit.ToBigInt().Int64())
)

func setAccountInfo(b *testing.B, storage *runtime.Storage, account []byte, info gossamertypes.AccountInfo) {
	bytesStorage, err := scale.Marshal(info)
	assert.NoError(b, err)

	err = (*storage).Put(accountStorageKey(account), bytesStorage)
	assert.NoError(b, err)
}

func getAccountInfo(b *testing.B, storage *runtime.Storage, account []byte) *gossamertypes.AccountInfo {
	accountInfo := gossamertypes.AccountInfo{
		Nonce:       0,
		Consumers:   0,
		Producers:   0,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(big.NewInt(0)),
			Reserved:   scale.MustNewUint128(big.NewInt(0)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}

	bytesStorage := (*storage).Get(accountStorageKey(account))

	err := scale.Unmarshal(bytesStorage, &accountInfo)
	assert.NoError(b, err)

	return &accountInfo
}

func accountStorageKey(account []byte) []byte {
	aliceHash, _ := common.Blake2b128(account)
	keyStorageAccount := append(keySystemHash, keyAccountHash...)
	keyStorageAccount = append(keyStorageAccount, aliceHash...)
	keyStorageAccount = append(keyStorageAccount, account...)
	return keyStorageAccount
}
