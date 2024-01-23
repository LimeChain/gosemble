package main

import (
	"testing"

	"github.com/ChainSafe/gossamer/lib/runtime"
	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

var blockLength, _ = system.MaxWithNormalRatio(constants.FiveMbPerBlockPerExtrinsic, constants.NormalDispatchRatio)

// TODO: implement components instead
func BenchmarkSystemRemarkStep1(b *testing.B) {
	benchmarkInstance(b, SystemRemark, 0)
}

func BenchmarkSystemRemarkStep2(b *testing.B) {
	benchmarkInstance(b, SystemRemark, 1)
}

func BenchmarkSystemRemarkStep3(b *testing.B) {
	benchmarkInstance(b, SystemRemark, 8)
}

func BenchmarkSystemRemarkStep4(b *testing.B) {
	benchmarkInstance(b, SystemRemark, 32)
}

func BenchmarkSystemRemarkStep5(b *testing.B) {
	benchmarkInstance(b, SystemRemark, 64)
}

func BenchmarkSystemRemarkStep6(b *testing.B) {
	benchmarkInstance(b, SystemRemark, 128)
}

func BenchmarkSystemRemarkStep7(b *testing.B) {
	benchmarkInstance(b, SystemRemark, 256)
}

func BenchmarkSystemRemarkStep8(b *testing.B) {
	benchmarkInstance(b, SystemRemark, 512)
}

func BenchmarkSystemRemarkStep9(b *testing.B) {
	benchmarkInstance(b, SystemRemark, 1024)
}

func BenchmarkSystemRemarkStep10(b *testing.B) {
	benchmarkInstance(b, SystemRemark, 128*1024)
}

func BenchmarkSystemRemarkStep11(b *testing.B) {
	benchmarkInstance(b, SystemRemark, 256*1024)
}

func BenchmarkSystemRemarkStep12(b *testing.B) {
	benchmarkInstance(b, SystemRemark, 512*1024)
}

func BenchmarkSystemRemarkStep13(b *testing.B) {
	benchmarkInstance(b, SystemRemark, 1024*1024)
}

func BenchmarkSystemRemarkStep14(b *testing.B) {
	benchmarkInstance(b, SystemRemark, int(blockLength.Max.Normal)) // 3932100
}

func SystemRemark(b *testing.B, rt *wazero_runtime.Instance, storage *runtime.Storage, metadata *ctypes.Metadata, args ...interface{}) (ctypes.Call, []byte) {
	// Setup the input params
	size := args[0].(int)

	message := make([]byte, size)
	call, err := ctypes.NewCall(metadata, "System.remark", message)
	assert.NoError(b, err)

	benchmarkConfig := newExtrinsicCall(b, types.NewRawOriginSigned(aliceAccountId), call)

	// Execute the call
	res, err := rt.Exec("Benchmark_run", benchmarkConfig.Bytes())
	assert.NoError(b, err)

	return call, res
}
