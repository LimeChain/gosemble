package main

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/benchmarking"
	"github.com/LimeChain/gosemble/primitives/types"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

var blockLength, _ = system.MaxWithNormalRatio(constants.FiveMbPerBlockPerExtrinsic, constants.NormalDispatchRatio)

// TODO: implement components instead
func BenchmarkSystemRemarkStep1(b *testing.B) {
	benchmarkSystemRemark(b, 0)
}

func BenchmarkSystemRemarkStep2(b *testing.B) {
	benchmarkSystemRemark(b, 1)
}

func BenchmarkSystemRemarkStep3(b *testing.B) {
	benchmarkSystemRemark(b, 8)
}

func BenchmarkSystemRemarkStep4(b *testing.B) {
	benchmarkSystemRemark(b, 32)
}

func BenchmarkSystemRemarkStep5(b *testing.B) {
	benchmarkSystemRemark(b, 64)
}

func BenchmarkSystemRemarkStep6(b *testing.B) {
	benchmarkSystemRemark(b, 128)
}

func BenchmarkSystemRemarkStep7(b *testing.B) {
	benchmarkSystemRemark(b, 256)
}

func BenchmarkSystemRemarkStep8(b *testing.B) {
	benchmarkSystemRemark(b, 512)
}

func BenchmarkSystemRemarkStep9(b *testing.B) {
	benchmarkSystemRemark(b, 1024)
}

func BenchmarkSystemRemarkStep10(b *testing.B) {
	benchmarkSystemRemark(b, 128*1024)
}

func BenchmarkSystemRemarkStep11(b *testing.B) {
	benchmarkSystemRemark(b, 256*1024)
}

func BenchmarkSystemRemarkStep12(b *testing.B) {
	benchmarkSystemRemark(b, 512*1024)
}

func BenchmarkSystemRemarkStep13(b *testing.B) {
	benchmarkSystemRemark(b, 1024*1024)
}

func BenchmarkSystemRemarkStep14(b *testing.B) {
	benchmarkSystemRemark(b, blockLength.Max.Normal) // 3932100
}

func benchmarkSystemRemark(b *testing.B, size sc.U32) {
	rt, _ := newBenchmarkingRuntime(b)

	runtimeVersion, err := rt.Version()
	assert.NoError(b, err)

	metadata := runtimeMetadata(b, rt)

	// Setup the input params
	message := make([]byte, size)

	// Create the call
	call, err := ctypes.NewCall(metadata, "System.remark", message)
	assert.NoError(b, err)

	extrinsic := ctypes.NewExtrinsic(call)

	o := ctypes.SignatureOptions{
		BlockHash:          ctypes.Hash(parentHash),
		Era:                ctypes.ExtrinsicEra{IsImmortalEra: true},
		GenesisHash:        ctypes.Hash(parentHash),
		Nonce:              ctypes.NewUCompactFromUInt(0),
		SpecVersion:        ctypes.U32(runtimeVersion.SpecVersion),
		Tip:                ctypes.NewUCompactFromUInt(0),
		TransactionVersion: ctypes.U32(runtimeVersion.TransactionVersion),
	}

	// Sign the transaction using Alice's default account
	err = extrinsic.Sign(signature.TestKeyringPairAlice, o)
	assert.NoError(b, err)

	encodedExtrinsic := bytes.Buffer{}
	encoder := cscale.NewEncoder(&encodedExtrinsic)
	err = extrinsic.Encode(*encoder)
	assert.NoError(b, err)

	benchmarkConfig := benchmarking.BenchmarkConfig{
		InternalRepeats: sc.U32(b.N),
		Extrinsic:       sc.BytesToSequenceU8(encodedExtrinsic.Bytes()),
		Origin:          sc.NewOption[types.RawOrigin](nil),
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
