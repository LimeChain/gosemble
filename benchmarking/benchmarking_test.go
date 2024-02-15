package benchmarking

import (
	"flag"
	"path/filepath"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/benchmarking"

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
