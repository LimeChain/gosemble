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

func BenchmarkBalancesForceFree(b *testing.B) {
	benchmarkInstance(b, BalancesForceFree)
}

func BalancesForceFree(b *testing.B, rt *wazero_runtime.Instance, storage *runtime.Storage, metadata *ctypes.Metadata, args ...interface{}) (ctypes.Call, []byte) {
	// Setup the input params
	call, err := ctypes.NewCall(metadata, "Balances.force_free", aliceAddress, ctypes.NewU128(*big.NewInt(2 * existentialAmount)))
	assert.NoError(b, err)

	benchmarkConfig := newExtrinsicCall(b, types.NewRawOriginRoot(), call)

	// Setup the state
	accountInfo := gossamertypes.AccountInfo{
		Nonce:       0,
		Consumers:   0,
		Producers:   0,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(big.NewInt(existentialAmount)),
			Reserved:   scale.MustNewUint128(big.NewInt(existentialAmount)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}
	setAccountInfo(b, storage, aliceAccountIdBytes, accountInfo)

	info := getAccountInfo(b, storage, aliceAccountIdBytes)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(existentialAmount)), info.Data.Reserved)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(existentialAmount)), info.Data.Free)

	// Whitelist the keys
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyNumberHash...)))         // 1 read/write
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyExecutionPhaseHash...))) // 1 read
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyEventCountHash...)))     // 1 read/write
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyEventsHash...)))         // 1 read/write

	// Execute the call
	res, err := rt.Exec("Benchmark_run", benchmarkConfig.Bytes())
	assert.NoError(b, err)

	// Validate the result
	info = getAccountInfo(b, storage, aliceAccountIdBytes)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(0)), info.Data.Reserved)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(2*existentialAmount)), info.Data.Free)

	return call, res
}
