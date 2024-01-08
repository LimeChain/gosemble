package main

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/runtime"
	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/assert"
)

func Benchmark_Balances_ForceFree_Step1(b *testing.B) {
	benchmarkCallBalancesForceFree(b, 0)
}

func Benchmark_Balances_ForceFree_Step2(b *testing.B) {
	benchmarkCallBalancesForceFree(b, 3000)
}

func Benchmark_Balances_ForceFree_Step3(b *testing.B) {
	benchmarkCallBalancesForceFree(b, 10000000000)
}

func Benchmark_Balances_ForceFree_Step4(b *testing.B) {
	benchmarkCallBalancesForceFree(b, 19462877928727)
}

func Benchmark_Balances_ForceFree_Step5(b *testing.B) {
	benchmarkCallBalancesForceFree(b, 999999999999999999)
}

func benchmarkCallBalancesForceFree(b *testing.B, n int64) {
	rt, storage := newBenchRuntime(b)

	metadata := newBenchRuntimeMetadata(b, rt)

	// BEFORE:
	// set up inputs, storage, etc., based on the
	// case benchmark case (should be the worst-case)
	balance, ok := big.NewInt(0).SetString("500000000000000", 10)
	if !ok {
		b.Fatal("could not set balance")
	}
	setStorageAccountInfoForBench(b, storage, signature.TestKeyringPairAlice.PublicKey, balance, 0)

	alice, err := ctypes.NewMultiAddressFromAccountID(signature.TestKeyringPairAlice.PublicKey)
	assert.NoError(b, err)

	call, err := ctypes.NewCall(metadata, "Balances.force_unreserve", alice, ctypes.NewU128(*big.NewInt(n)))
	assert.NoError(b, err)

	callEnc, err := codec.Encode(call)
	assert.NoError(b, err)

	repeatsEnc, err := codec.Encode(uint8(b.N))
	assert.NoError(b, err)

	callEnc = append(repeatsEnc, callEnc...)

	// EXECUTE:
	// the actual benchmark
	var res []byte
	res, err = rt.Exec("Benchmark_run", callEnc)

	// AFTER:
	// validate the result and abort if errors
	assert.NoError(b, err)

	buf := bytes.NewBuffer(res)
	repeats, _ := sc.DecodeU8(buf)
	module, _ := sc.DecodeU8(buf)
	function, _ := sc.DecodeU8(buf)
	time, _ := sc.DecodeU64(buf)
	dbReads, _ := sc.DecodeU32(buf)
	dbWrites, _ := sc.DecodeU32(buf)
	b.ReportMetric(float64(repeats), "repeats")
	b.ReportMetric(float64(module), "module")
	b.ReportMetric(float64(function), "function")
	b.ReportMetric(float64(time), "time")
	b.ReportMetric(float64(dbReads), "reads")
	b.ReportMetric(float64(dbWrites), "writes")

	b.Cleanup(func() {
		rt.Stop()
	})
}

func newBenchRuntime(b *testing.B) (*wazero_runtime.Instance, *runtime.Storage) {
	runtime := wazero_runtime.NewBenchInstanceWithTrie(b, WASM_RUNTIME, trie.NewEmptyTrie())
	return runtime, &runtime.Context.Storage
}

func newBenchRuntimeMetadata(b *testing.B, instance *wazero_runtime.Instance) *ctypes.Metadata {
	bMetadata, err := instance.Metadata()
	assert.NoError(b, err)

	var decoded []byte
	err = scale.Unmarshal(bMetadata, &decoded)
	assert.NoError(b, err)

	metadata := &ctypes.Metadata{}
	err = codec.Decode(decoded, metadata)
	assert.NoError(b, err)

	return metadata
}

func setStorageAccountInfoForBench(b *testing.B, storage *runtime.Storage, account []byte, freeBalance *big.Int, nonce uint32) (storageKey []byte, info gossamertypes.AccountInfo) {
	accountInfo := gossamertypes.AccountInfo{
		Nonce:       nonce,
		Consumers:   0,
		Producers:   0,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(freeBalance),
			Reserved:   scale.MustNewUint128(big.NewInt(0)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}

	aliceHash, _ := common.Blake2b128(account)
	keyStorageAccount := append(keySystemHash, keyAccountHash...)
	keyStorageAccount = append(keyStorageAccount, aliceHash...)
	keyStorageAccount = append(keyStorageAccount, account...)

	bytesStorage, err := scale.Marshal(accountInfo)
	assert.NoError(b, err)

	err = (*storage).Put(keyStorageAccount, bytesStorage)
	assert.NoError(b, err)

	return keyStorageAccount, accountInfo
}

func setup(b *testing.B, rt *wazero_runtime.Instance, storage *runtime.Storage, metadata *ctypes.Metadata, rtVersion runtime.Version, n int64) []byte {
	header := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, uint(blockNumber), gossamertypes.NewDigest())
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(b, err)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(b, err)

	alice, err := ctypes.NewMultiAddressFromAccountID(signature.TestKeyringPairAlice.PublicKey)
	assert.NoError(b, err)

	call, err := ctypes.NewCall(metadata, "Balances.force_unreserve", alice, ctypes.NewU128(*big.NewInt(n)))
	assert.NoError(b, err)

	// Create the extrinsic
	ext := ctypes.NewExtrinsic(call)
	o := ctypes.SignatureOptions{
		BlockHash:          ctypes.Hash(parentHash),
		Era:                ctypes.ExtrinsicEra{IsImmortalEra: true},
		GenesisHash:        ctypes.Hash(parentHash),
		Nonce:              ctypes.NewUCompactFromUInt(0),
		SpecVersion:        ctypes.U32(rtVersion.SpecVersion),
		Tip:                ctypes.NewUCompactFromUInt(0),
		TransactionVersion: ctypes.U32(rtVersion.TransactionVersion),
	}

	// Sign the transaction using Alice's default account
	err = ext.Sign(signature.TestKeyringPairAlice, o)
	assert.NoError(b, err)

	extEnc := bytes.Buffer{}
	encoder := cscale.NewEncoder(&extEnc)
	err = ext.Encode(*encoder)
	assert.NoError(b, err)

	return extEnc.Bytes()
}

func benchmarkExtrinsicBalancesForceFreeBadOrigin(b *testing.B, n int64) {
	rt, storage := newBenchRuntime(b)

	runtimeVersion, err := rt.Version()
	assert.NoError(b, err)

	metadata := newBenchRuntimeMetadata(b, rt)

	// Set up the inputs, params, storage, etc., based on
	// the case we want to benchmark, including the worst-case.
	encodedExtrinsicBytes := setup(b, rt, storage, metadata, runtimeVersion, n)

	var res []byte

	// Exclude the setup from the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Each iteration should benchmark the same state
		balance, ok := big.NewInt(0).SetString("500000000000000", 10)
		if !ok {
			b.Fatal("could not set balance")
		}
		setStorageAccountInfoForBench(b, storage, signature.TestKeyringPairAlice.PublicKey, balance, 0)

		b.StartTimer()

		(*storage).DbResetTracker()
		// Whitelist the account that will be used to sign the extrinsic
		accountKey, _ := hex.DecodeString("26aa394eea5630e07c48ae0c9558cef7b99d880ec681799c0cf30e8886371da9de1e86a9a8c739864cf3cc5ec2bea59fd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d")
		(*storage).DbWhitelistKey(string(accountKey))

		// Measure the DB r/w for single benchmark case
		(*storage).DbStartTracker()
		// The actual benchmarking of an extrinsic
		res, err = rt.Exec("BlockBuilder_apply_extrinsic", encodedExtrinsicBytes)
		(*storage).DbStopTracker()
	}
	// Exclude the teardown from the benchmark
	b.StopTimer()

	// Report - add the DB r/w metrics to the benchmark report
	b.ReportMetric(float64((*storage).DbReadCount()), "reads")
	b.ReportMetric(float64((*storage).DbWriteCount()), "writes")

	// Validate - assert results and abort and report any errors.
	assert.NoError(b, err)
	assert.Equal(b, applyExtrinsicResultBadOriginErr.Bytes(), res)

	// Cleanup - stop the runtime
	b.Cleanup(func() {
		rt.Stop()
	})
}
