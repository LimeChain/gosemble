package main

import (
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/pkg/scale"
	"github.com/LimeChain/gosemble/benchmarking"
	"github.com/LimeChain/gosemble/primitives/types"

	// "github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func BenchmarkBalancesForceFreeNew(b *testing.B) {
	benchmarking.Run(b, func(i *benchmarking.Instance) {
		// arrange
		// pubKey := signature.TestKeyringPairAlice.PublicKey

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

		err := i.SetAccountInfo(aliceAccountIdBytes, accountInfo)
		assert.NoError(b, err)

		// (*i.Storage()).DbWhitelistKey(string(append(keySystemHash, keyNumberHash...)))         // 1 read/write
		// (*i.Storage()).DbWhitelistKey(string(append(keySystemHash, keyExecutionPhaseHash...))) // 1 read
		// (*i.Storage()).DbWhitelistKey(string(append(keySystemHash, keyEventCountHash...)))     // 1 read/write
		// (*i.Storage()).DbWhitelistKey(string(append(keySystemHash, keyEventsHash...)))         // 1 read/write

		// act
		err = i.ExecuteExtrinsic(
			"Balances.force_free",
			types.NewRawOriginRoot(),
			aliceAddress,
			ctypes.NewU128(*big.NewInt(2 * existentialAmount)),
		)

		// assert
		assert.NoError(b, err)

		existingAccountInfo, err := i.GetAccountInfo(aliceAccountIdBytes)
		assert.NoError(b, err)
		// assert.Nil(b, existingAccountInfo)
		// assert.NoError(b, err)
		assert.Equal(b, scale.MustNewUint128(big.NewInt(0)), existingAccountInfo.Data.Reserved)
		assert.Equal(b, scale.MustNewUint128(big.NewInt(2*existentialAmount)), existingAccountInfo.Data.Free)
	})
}

// func benchmarkBalancesForceFree(b *testing.B) {
// 	rt, storage := newBenchmarkingRuntime(b)

// 	metadata := runtimeMetadata(b, rt)

// 	// Setup the input params
// 	alice, err := ctypes.NewMultiAddressFromAccountID(signature.TestKeyringPairAlice.PublicKey)
// 	assert.NoError(b, err)

// 	// Create the call
// 	call, err := ctypes.NewCall(metadata, "Balances.force_free", alice, ctypes.NewU128(*big.NewInt(2 * existentialAmount)))
// 	assert.NoError(b, err)

// 	// Create the extrinsic
// 	extrinsic := ctypes.NewExtrinsic(call)

// 	encodedExtrinsic := bytes.Buffer{}
// 	encoder := cscale.NewEncoder(&encodedExtrinsic)
// 	err = extrinsic.Encode(*encoder)
// 	assert.NoError(b, err)

// 	benchmarkConfig := benchmarkingtypes.BenchmarkConfig{
// 		InternalRepeats: sc.U32(b.N),
// 		Extrinsic:       sc.BytesToSequenceU8(encodedExtrinsic.Bytes()),
// 		Origin:          sc.NewOption[types.RawOrigin](types.NewRawOriginRoot()),
// 	}

// 	// Setup the state
// 	pubKey := signature.TestKeyringPairAlice.PublicKey
// 	accountInfo := gossamertypes.AccountInfo{
// 		Nonce:       0,
// 		Consumers:   0,
// 		Producers:   0,
// 		Sufficients: 0,
// 		Data: gossamertypes.AccountData{
// 			Free:       scale.MustNewUint128(big.NewInt(existentialAmount)),
// 			Reserved:   scale.MustNewUint128(big.NewInt(existentialAmount)),
// 			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
// 			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
// 		},
// 	}
// 	setAccountInfo(b, storage, pubKey, accountInfo)

// 	info := getAccountInfo(b, storage, pubKey)
// 	assert.Equal(b, scale.MustNewUint128(big.NewInt(existentialAmount)), info.Data.Reserved)
// 	assert.Equal(b, scale.MustNewUint128(big.NewInt(existentialAmount)), info.Data.Free)

// 	(*storage).DbWhitelistKey(string(append(keySystemHash, keyNumberHash...)))         // 1 read/write
// 	(*storage).DbWhitelistKey(string(append(keySystemHash, keyExecutionPhaseHash...))) // 1 read
// 	(*storage).DbWhitelistKey(string(append(keySystemHash, keyEventCountHash...)))     // 1 read/write
// 	(*storage).DbWhitelistKey(string(append(keySystemHash, keyEventsHash...)))         // 1 read/write

// 	res, err := rt.Exec("Benchmark_run", benchmarkConfig.Bytes())

// 	assert.NoError(b, err)

// 	// Validate the result/state
// 	info = getAccountInfo(b, storage, pubKey)
// 	assert.Equal(b, scale.MustNewUint128(big.NewInt(0)), info.Data.Reserved)
// 	assert.Equal(b, scale.MustNewUint128(big.NewInt(2*existentialAmount)), info.Data.Free)

// 	benchmarkResult, err := benchmarking.DecodeBenchmarkResult(bytes.NewBuffer(res))
// 	assert.NoError(b, err)

// 	// Report the results
// 	b.ReportMetric(float64(call.CallIndex.SectionIndex), "module")
// 	b.ReportMetric(float64(call.CallIndex.MethodIndex), "function")
// 	b.ReportMetric(float64(b.N), "repeats")
// 	b.ReportMetric(float64(benchmarkResult.ExtrinsicTime.ToBigInt().Int64()), "time")
// 	b.ReportMetric(float64(benchmarkResult.Reads), "reads")
// 	b.ReportMetric(float64(benchmarkResult.Writes), "writes")

// 	b.Cleanup(func() {
// 		rt.Stop()
// 	})
// }
