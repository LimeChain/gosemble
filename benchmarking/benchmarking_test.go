package benchmarking

import (
	"flag"
	"fmt"
	"path/filepath"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/benchmarking"

	"io"
	"math/big"
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/stretchr/testify/assert"
)

var (
	pubKey = signature.TestKeyringPairAlice.PublicKey
)

func TestRun(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "output.go")

	// redirect os.Stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// set benchmarking flags
	os.Args = append(os.Args, "-steps", "3")
	os.Args = append(os.Args, "-test.benchtime", "1x")
	flag.Parse()

	// run with components
	testing.Benchmark(func(b *testing.B) {
		component, err := NewLinear(1, 100)
		assert.NoError(t, err)

		componentValues := []uint32{}

		RunDispatchCall(b, outputPath, func(instance *Instance) {
			testFn(t, instance, component.Value())
			componentValues = append(componentValues, component.Value())
		}, component)

		assert.Equal(t, []uint32{1, 50, 100}, componentValues)
	})
	expectedMedianSlopesAnalysis := "median slope analysis: benchmarking.analysis{baseExtrinsicTime:0x0, baseReads:0x0, baseWrites:0x0, slopesExtrinsicTime:[]uint64{0x0}, slopesReads:[]uint64{0x1}, slopesWrites:[]uint64{0x0}, minimumExtrinsicTime:0x0, minimumReads:0x1, minimumWrites:0x0}\n"

	// run with no components
	testing.Benchmark(func(b *testing.B) {
		value := uint32(100)
		RunDispatchCall(b, outputPath, func(instance *Instance) {
			testFn(t, instance, value)
		})
	})
	expectedMedianValuesAnalysis := "median slope analysis: benchmarking.analysis{baseExtrinsicTime:0x0, baseReads:0x64, baseWrites:0x0, slopesExtrinsicTime:[]uint64(nil), slopesReads:[]uint64(nil), slopesWrites:[]uint64(nil), minimumExtrinsicTime:0x0, minimumReads:0x64, minimumWrites:0x0}\n"

	// stop redirecting os.Stdout
	w.Close()
	os.Stdout = rescueStdout

	// assert output
	out, _ := io.ReadAll(r)
	assert.Equal(t, fmt.Sprintf("%s%s", expectedMedianSlopesAnalysis, expectedMedianValuesAnalysis), string(out))
}

func testFn(t *testing.T, instance *Instance, value uint32) {
	acc, err := instance.GetAccountInfo(pubKey)
	assert.Error(t, err)

	accInfo := gossamertypes.AccountInfo{
		Nonce:       value,
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

	err = instance.SetAccountInfo(pubKey, accInfo)
	assert.NoError(t, err)

	acc, err = instance.GetAccountInfo(pubKey)
	assert.NoError(t, err)
	assert.Equal(t, accInfo, acc)

	instance.benchmarkResult = &benchmarking.BenchmarkResult{Reads: sc.U32(value)}
}
