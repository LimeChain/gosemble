package main

import (
	"bytes"
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/runtime"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/benchmarking"
	"github.com/LimeChain/gosemble/primitives/types"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func BenchmarkBalancesForceFreeStep1(b *testing.B) {
	benchmarkBalancesForceFree(b)
}

func benchmarkBalancesForceFree(b *testing.B) {
	rt, storage := newBenchmarkingRuntime(b)

	metadata := runtimeMetadata(b, rt)

	// Setup the input params
	alice, err := ctypes.NewMultiAddressFromAccountID(signature.TestKeyringPairAlice.PublicKey)
	assert.NoError(b, err)

	// Create the call
	call, err := ctypes.NewCall(metadata, "Balances.force_free", alice, ctypes.NewU128(*big.NewInt(600000000)))
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
	// This setup differs from the one in the Substrate, and the
	// total (free + reserved) is less than the existential deposit,
	// resulting in additional read/write for "Balances TotalIssuance" key
	pubKey := signature.TestKeyringPairAlice.PublicKey
	accountInfo := gossamertypes.AccountInfo{
		Nonce:       0,
		Consumers:   0,
		Producers:   0,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(big.NewInt(500000000)),
			Reserved:   scale.MustNewUint128(big.NewInt(500000000)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}
	setAccountInfo(b, storage, pubKey, accountInfo)

	info := getAccountInfo(b, storage, pubKey)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(500000000)), info.Data.Reserved)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(500000000)), info.Data.Free)

	(*storage).DbWhitelistKey(string(append(keySystemHash, keyBlockNumberHash...)))    // 1 read/write
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyExecutionPhaseHash...))) // 1 read
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyEventCountHash...)))     // 1 read/write
	(*storage).DbWhitelistKey(string(append(keySystemHash, keyEventsHash...)))         // 1 read/write

	res, err := rt.Exec("Benchmark_run", benchmarkConfig.Bytes())

	assert.NoError(b, err)

	// Validate the result/state
	info = getAccountInfo(b, storage, pubKey)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(0)), info.Data.Reserved)
	assert.Equal(b, scale.MustNewUint128(big.NewInt(2*500000000)), info.Data.Free)

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
