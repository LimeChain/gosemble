package main

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/gossamer/lib/runtime"
	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/lib/trie"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/benchmarking"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/stretchr/testify/assert"
)

func BenchmarkHooksOnInitialize(b *testing.B) {
	auraCurrentSlot := sc.U64(1)
	auraNewSlot := sc.U64(2)
	auraAuthorityPubKey, _ := types.NewSr25519PublicKey(sc.BytesToSequenceU8(signature.TestKeyringPairAlice.PublicKey)...)
	auraAuthorities := sc.Sequence[types.Sr25519PublicKey]{auraAuthorityPubKey}
	digest := types.NewDigest(sc.Sequence[types.DigestItem]{
		types.NewDigestItemPreRuntime(
			sc.BytesToFixedSequenceU8([]byte{'a', 'u', 'r', 'a'}),
			sc.BytesToSequenceU8(auraNewSlot.Bytes()),
		),
	})

	setupFn := func(storage *runtime.Storage) {
		(*storage).Put(append(keySystemHash, keyDigestHash...), digest.Bytes())
		(*storage).Put(append(keyAuraHash, keyCurrentSlotHash...), auraCurrentSlot.Bytes())
		(*storage).Put(append(keyAuraHash, keyAuthoritiesHash...), auraAuthorities.Bytes())
	}

	validateFn := func(storage *runtime.Storage) {
		assert.Equal(b, sc.U64(2).Bytes(), (*storage).Get(append(keyAuraHash, keyCurrentSlotHash...)))
	}

	runBenchmark(b, "on_initialize", setupFn, validateFn)
}

func BenchmarkHooksOnRuntimeUpgrade(b *testing.B) {
	setupFn := func(storage *runtime.Storage) {}
	validateFn := func(storage *runtime.Storage) {}
	runBenchmark(b, "on_runtime_upgrade", setupFn, validateFn)
}

func BenchmarkHooksOnFinalize(b *testing.B) {
	key := append(keyTimestampHash, keyTimestampDidUpdateHash...)

	setupFn := func(storage *runtime.Storage) {
		(*storage).Put(key, sc.Bool(true).Bytes())
	}

	validateFn := func(storage *runtime.Storage) {
		assert.Equal(b, []byte(nil), (*storage).Get(key))
	}

	runBenchmark(b, "on_finalize", setupFn, validateFn)
}

func BenchmarkHooksOnIdle(b *testing.B) {
	setupFn := func(storage *runtime.Storage) {}
	validateFn := func(storage *runtime.Storage) {}
	runBenchmark(b, "on_idle", setupFn, validateFn)
}

func runBenchmark(b *testing.B, hook string, setupFn func(storage *runtime.Storage), validateFn func(storage *runtime.Storage)) {
	rt := wazero_runtime.NewBenchInstanceWithTrie(b, "../build/runtime.wasm", trie.NewEmptyTrie())
	defer rt.Stop()

	storage := &rt.Context.Storage

	hookBytes := append(sc.Str(hook).Bytes(), sc.U64(1).Bytes()...)
	hookBytes = append(hookBytes, types.WeightZero().Bytes()...)

	benchmarkConfig := benchmarking.BenchmarkConfig{
		InternalRepeats: sc.U32(b.N),
		Benchmark:       sc.BytesToSequenceU8(hookBytes),
		Origin:          types.NewRawOriginNone(),
	}

	setupFn(storage)

	res, err := rt.Exec("Benchmark_hook", benchmarkConfig.Bytes())
	assert.NoError(b, err)

	validateFn(storage)

	benchmarkResult, err := benchmarking.DecodeBenchmarkResult(bytes.NewBuffer(res))
	assert.NoError(b, err)

	b.ReportMetric(float64(benchmarkResult.Time.ToBigInt().Int64()), "time")
	b.ReportMetric(float64(benchmarkResult.Reads), "reads")
	b.ReportMetric(float64(benchmarkResult.Writes), "writes")
}
