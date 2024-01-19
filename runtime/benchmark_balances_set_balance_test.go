package main

import (
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/runtime"
	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/pkg/scale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

var value = uint64(existentialMultiplier * existentialAmount)

// Coming from ROOT account. This always creates an account.
func BenchmarkBalancesSetBalanceCreating(b *testing.B) {
	benchmarkInstance(b, BalancesSetBalance, value, value)
}

// Coming from ROOT account. This always kills an account.
func BenchmarkBalancesSetBalanceKilling(b *testing.B) {
	benchmarkInstance(b, BalancesSetBalance, value, uint64(0))
}

func BalancesSetBalance(b *testing.B, rt *wazero_runtime.Instance, storage *runtime.Storage, metadata *ctypes.Metadata, args ...interface{}) (ctypes.Call, []byte) {
	// Setup the input params
	balance := args[0].(uint64)
	amount := args[1].(uint64)

	call, err := ctypes.NewCall(metadata, "Balances.set_balance", aliceAddress, ctypes.NewUCompactFromUInt(amount), ctypes.NewUCompactFromUInt(amount))
	assert.NoError(b, err)

	benchmarkConfig := newExtrinsicCall(b, types.NewRawOriginRoot(), call)

	// Setup the state
	aliceAccountInfo := gossamertypes.AccountInfo{
		Nonce:       0,
		Consumers:   0,
		Producers:   1,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(big.NewInt(int64(balance))),
			Reserved:   scale.MustNewUint128(big.NewInt(0)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}
	setAccountInfo(b, storage, signature.TestKeyringPairAlice.PublicKey, aliceAccountInfo)

	aliceInfo := getAccountInfo(b, storage, aliceAccountIdBytes)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(int64(balance))), aliceInfo.Data.Free)

	// Execute the call
	res, err := rt.Exec("Benchmark_run", benchmarkConfig.Bytes())
	assert.NoError(b, err)

	// Validate the result
	aliceInfo = getAccountInfo(b, storage, aliceAccountIdBytes)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(int64(amount))), aliceInfo.Data.Free)

	return call, res
}
