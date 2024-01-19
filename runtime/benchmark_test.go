package main

import (
	"bytes"
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/runtime"
	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/benchmarking"
	"github.com/LimeChain/gosemble/primitives/types"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

// TODO: switch to Gosemble types

var (
	aliceAddress, _     = ctypes.NewMultiAddressFromHexAccountID("0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d")
	aliceAccountIdBytes = aliceAddress.AsID.ToBytes()
	aliceAccountId, _   = types.NewAccountId(sc.BytesToSequenceU8(aliceAccountIdBytes)...)

	bobAddress, _     = ctypes.NewMultiAddressFromHexAccountID("0x90b5ab205c6974c9ea841be688864633dc9ca8a357843eeacf2314649965fe22")
	bobAccountIdBytes = bobAddress.AsID.ToBytes()
)

var (
	existentialAmount     = int64(BalancesExistentialDeposit.ToBigInt().Int64())
	existentialMultiplier = int64(10)
)

func benchmarkInstance(
	b *testing.B,
	fn func(*testing.B, *wazero_runtime.Instance, *runtime.Storage, *ctypes.Metadata, ...interface{}) (ctypes.Call, []byte),
	args ...interface{},
) {
	rt, storage := newBenchmarkingRuntime(b)
	metadata := runtimeMetadata(b, rt)

	call, res := fn(b, rt, storage, metadata, args...)

	benchmarkResult, err := benchmarking.DecodeBenchmarkResult(bytes.NewBuffer(res))
	assert.NoError(b, err)

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

func newExtrinsicCall(b *testing.B, origin types.RawOrigin, call ctypes.Call) *benchmarking.BenchmarkConfig {
	extrinsic := ctypes.NewExtrinsic(call)

	encodedExtrinsic := bytes.Buffer{}
	encoder := cscale.NewEncoder(&encodedExtrinsic)
	err := extrinsic.Encode(*encoder)
	assert.NoError(b, err)

	return &benchmarking.BenchmarkConfig{
		InternalRepeats: sc.U32(b.N),
		Extrinsic:       sc.BytesToSequenceU8(encodedExtrinsic.Bytes()),
		Origin:          origin,
	}
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
