package benchmarking

import (
	"flag"
	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/benchmarking"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/stretchr/testify/assert"
	"io"
	"math/big"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = append(os.Args, "-steps", "3")
	os.Args = append(os.Args, "-test.benchtime", "1x")
	flag.Parse()

	component, err := NewLinear(1, 100)
	assert.NoError(t, err)

	pubKey := signature.TestKeyringPairAlice.PublicKey
	componentValues := []uint32{}

	testing.Benchmark(func(b *testing.B) {
		Run(b, func(instance *Instance) {
			acc, err := instance.GetAccountInfo(pubKey)
			assert.Error(t, err)

			accInfo := gossamertypes.AccountInfo{
				Nonce:       component.Value(),
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

			componentValues = append(componentValues, component.Value())

			instance.benchmarkResult = &benchmarking.BenchmarkResult{Reads: sc.U32(component.Value())}
		}, component)
	})

	assert.Equal(t, []uint32{1, 50, 100}, componentValues)

	w.Close()

	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	assert.Equal(t, "median slope analysis: benchmarking.analysis{baseExtrinsicTime:0x0, baseReads:0x32, baseWrites:0x0, slopesExtrinsicTime:[]uint64(nil), slopesReads:[]uint64(nil), slopesWrites:[]uint64(nil), minimumExtrinsicTime:0x0, minimumReads:0x1, minimumWrites:0x0}\n", string(out))
}
