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

// Benchmark `transfer_all` with the worst possible condition:
// * The recipient account is created
// * The sender is killed
func BenchmarkBalancesTransferAllAllowDeath(b *testing.B) {
	benchmarkInstance(b, BalancesTransferAll)
}

func BalancesTransferAll(b *testing.B, rt *wazero_runtime.Instance, storage *runtime.Storage, metadata *ctypes.Metadata, args ...interface{}) (ctypes.Call, []byte) {
	// Setup the input params
	balance := existentialMultiplier * existentialAmount

	call, err := ctypes.NewCall(metadata, "Balances.transfer_all", bobAddress, ctypes.NewBool(false))
	assert.NoError(b, err)

	benchmarkConfig := newExtrinsicCall(b, types.NewRawOriginSigned(aliceAccountId), call)

	// Setup the state
	aliceAccountInfo := gossamertypes.AccountInfo{
		Nonce:       0,
		Consumers:   0,
		Producers:   1,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(big.NewInt(balance)),
			Reserved:   scale.MustNewUint128(big.NewInt(0)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}
	setAccountInfo(b, storage, signature.TestKeyringPairAlice.PublicKey, aliceAccountInfo)

	// Execute the call
	res, err := rt.Exec("Benchmark_run", benchmarkConfig.Bytes())
	assert.NoError(b, err)

	// Validate the result
	aliceInfo := getAccountInfo(b, storage, aliceAccountIdBytes)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(int64(0))), aliceInfo.Data.Free)
	bobInfo := getAccountInfo(b, storage, bobAccountIdBytes)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(int64(balance))), bobInfo.Data.Free)

	return call, res
}
