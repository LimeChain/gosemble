package main

import (
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/runtime"
	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/pkg/scale"
	"github.com/LimeChain/gosemble/primitives/types"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

// Benchmark extrinsic with the worst possible conditions:
// * Transfer will kill the sender account.
// * Transfer will create the recipient account.
func BenchmarkBalancesForceTransfer(b *testing.B) {
	benchmarkInstance(b, BalancesForceTransfer, nil)
}

func BalancesForceTransfer(b *testing.B, rt *wazero_runtime.Instance, storage *runtime.Storage, metadata *ctypes.Metadata, args ...interface{}) (ctypes.Call, []byte) {
	// Setup the input params
	balance := existentialMultiplier * existentialAmount
	transferAmount := uint64(existentialAmount*(existentialMultiplier-1) + 1)

	call, err := ctypes.NewCall(metadata, "Balances.force_transfer", aliceAddress, bobAddress, ctypes.NewUCompactFromUInt(transferAmount))
	assert.NoError(b, err)

	benchmarkConfig := newExtrinsicCall(b, types.NewRawOriginRoot(), call)

	// Setup the state
	aliceAccountInfo := gossamertypes.AccountInfo{
		Nonce:       0,
		Consumers:   0,
		Producers:   0,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(big.NewInt(balance)),
			Reserved:   scale.MustNewUint128(big.NewInt(0)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}
	setAccountInfo(b, storage, aliceAccountIdBytes, aliceAccountInfo)

	aliceInfo := getAccountInfo(b, storage, aliceAccountIdBytes)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(balance)), aliceInfo.Data.Free)

	// Whitelist the keys
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyNumberHash...)))         // 1 read/write
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyExecutionPhaseHash...))) // 1 read
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyEventCountHash...)))     // 1 read/write
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyEventsHash...)))         // 1 read/write

	// Execute the call
	res, err := rt.Exec("Benchmark_run", benchmarkConfig.Bytes())
	assert.NoError(b, err)

	// Validate the result
	aliceInfo = getAccountInfo(b, storage, aliceAccountIdBytes)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(balance-int64(transferAmount))), aliceInfo.Data.Free)
	bobInfo := getAccountInfo(b, storage, bobAccountIdBytes)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(int64(transferAmount))), bobInfo.Data.Free)

	return call, res
}
