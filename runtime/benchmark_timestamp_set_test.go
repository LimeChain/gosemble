package main

import (
	"bytes"
	"testing"
	"time"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/benchmarking"
	"github.com/LimeChain/gosemble/primitives/types"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func BenchmarkTimestampSet(b *testing.B) {
	benchmarkTimestampSet(b)
}

func benchmarkTimestampSet(b *testing.B) {
	rt, storage := newBenchmarkingRuntime(b)

	metadata := runtimeMetadata(b, rt)

	// Setup the input params
	now := uint64(time.Now().UnixMilli())
	call, err := ctypes.NewCall(metadata, "Timestamp.set", ctypes.NewUCompactFromUInt(now))
	assert.NoError(b, err)
	extrinsic := ctypes.NewExtrinsic(call)
	encodedExtrinsic := bytes.Buffer{}
	encoder := cscale.NewEncoder(&encodedExtrinsic)
	err = extrinsic.Encode(*encoder)
	assert.NoError(b, err)

	benchmarkConfig := benchmarking.BenchmarkConfig{
		InternalRepeats: sc.U32(b.N),
		Extrinsic:       sc.BytesToSequenceU8(encodedExtrinsic.Bytes()),
		Origin:          sc.NewOption[types.RawOrigin](types.NewRawOriginNone()),
	}

	// Setup the state
	(*storage).Put(append(keyTimestampHash, keyTimestampNowHash...), sc.U64(0).Bytes())
	assert.NoError(b, err)

	(*storage).DbWhitelistKey(string(append(keyTimestampHash, keyTimestampDidUpdateHash...)))

	res, err := rt.Exec("Benchmark_run", benchmarkConfig.Bytes())

	assert.NoError(b, err)

	// Validate the result/state
	nowStorageValue, err := sc.DecodeU64(bytes.NewBuffer((*storage).Get(append(keyTimestampHash, keyTimestampNowHash...))))
	assert.NoError(b, err)
	assert.Equal(b, sc.U64(now), nowStorageValue)

	didUpdateStorageValue, err := sc.DecodeBool(bytes.NewBuffer((*storage).Get(append(keyTimestampHash, keyTimestampDidUpdateHash...))))
	assert.NoError(b, err)
	assert.Equal(b, sc.Bool(true), didUpdateStorageValue)

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
