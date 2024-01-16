package main

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/benchmarking"
	"github.com/LimeChain/gosemble/primitives/types"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func BenchmarkSystemRemarkStep1(b *testing.B) {
	benchmarkSystemRemark(b)
}

func benchmarkSystemRemark(b *testing.B) {
	rt, _ := newBenchmarkingRuntime(b)

	metadata := runtimeMetadata(b, rt)

	// Setup the input params
	call, err := ctypes.NewCall(metadata, "System.remark", []byte{})
	assert.NoError(b, err)

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

	res, err := rt.Exec("Benchmark_run", benchmarkConfig.Bytes())

	assert.NoError(b, err)

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
