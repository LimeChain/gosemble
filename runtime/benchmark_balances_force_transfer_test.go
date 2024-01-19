package main

import (
	"bytes"
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/benchmarking"
	"github.com/LimeChain/gosemble/primitives/types"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func BenchmarkBalancesForceTransferStep1(b *testing.B) {
	benchmarkBalancesForceTransfer(b)
}

func benchmarkBalancesForceTransfer(b *testing.B) {
	rt, storage := newBenchmarkingRuntime(b)

	metadata := runtimeMetadata(b, rt)

	// Setup the input params
	alice, err := ctypes.NewMultiAddressFromAccountID(signature.TestKeyringPairAlice.PublicKey)
	assert.NoError(b, err)
	aliceAccountId := alice.AsID.ToBytes()

	bob, err := ctypes.NewMultiAddressFromHexAccountID("0x90b5ab205c6974c9ea841be688864633dc9ca8a357843eeacf2314649965fe22")
	assert.NoError(b, err)
	bobAccountId := bob.AsID.ToBytes()

	balance := existentialMultiplier * existentialAmount
	transferAmount := uint64(existentialAmount*(existentialMultiplier-1) + 1)

	// Create the call
	call, err := ctypes.NewCall(metadata, "Balances.force_transfer", alice, bob, ctypes.NewUCompactFromUInt(transferAmount))
	assert.NoError(b, err)

	// Create the extrinsic
	extrinsic := ctypes.NewExtrinsic(call)

	encodedExtrinsic := bytes.Buffer{}
	encoder := cscale.NewEncoder(&encodedExtrinsic)
	err = extrinsic.Encode(*encoder)
	assert.NoError(b, err)

	benchmarkConfig := benchmarking.BenchmarkConfig{
		InternalRepeats: sc.U32(b.N),
		Extrinsic:       sc.BytesToSequenceU8(encodedExtrinsic.Bytes()),
		Origin:          sc.NewOption[types.RawOrigin](types.NewRawOriginRoot()),
	}

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
	setAccountInfo(b, storage, aliceAccountId, aliceAccountInfo)

	aliceInfo := getAccountInfo(b, storage, aliceAccountId)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(balance)), aliceInfo.Data.Free)

	(*storage).DbWhitelistKey(string(append(keySystemHash, keyNumberHash...)))         // 1 read/write
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyExecutionPhaseHash...))) // 1 read
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyEventCountHash...)))     // 1 read/write
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyEventsHash...)))         // 1 read/write

	res, err := rt.Exec("Benchmark_run", benchmarkConfig.Bytes())

	assert.NoError(b, err)

	// Validate the result/state
	aliceInfo = getAccountInfo(b, storage, aliceAccountId)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(balance-int64(transferAmount))), aliceInfo.Data.Free)
	bobInfo := getAccountInfo(b, storage, bobAccountId)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(int64(transferAmount))), bobInfo.Data.Free)

	benchmarkResult, err := benchmarking.DecodeBenchmarkResult(bytes.NewBuffer(res))
	assert.NoError(b, err)

	// Report the results
	b.ReportMetric(float64(call.CallIndex.SectionIndex), "module")
	b.ReportMetric(float64(call.CallIndex.MethodIndex), "function")
	b.ReportMetric(float64(b.N), "repeats")
	b.ReportMetric(float64(benchmarkResult.ExtrinsicTime.ToBigInt().Int64()), "time")
	b.ReportMetric(float64(benchmarkResult.Reads), "reads")
	b.ReportMetric(float64(benchmarkResult.Writes), "writes")

	b.Cleanup(func() {
		rt.Stop()
	})
}
