package benchmarking

import (
	"flag"
	"fmt"
	"testing"

	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/lib/trie"
	benchmarkingtypes "github.com/LimeChain/gosemble/primitives/benchmarking"
)

var (
	stepsFlag     *int    = flag.Int("steps", 50, "Select how many samples we should take across the variable components")
	buildPathFlag *string = flag.String("BUILD_PATH", "build/runtime-benchmarking.wasm", "Path to compiled runtime wasm file. Path must be relative to project root.")
)

// todo analysis
// type benchmarkResult struct {
// 	extrinsicTime uint64
// 	reads, writes uint32
// 	components    []uint32
// }

// func newBenchmarkResult(benchmarkRes benchmarkingtypes.BenchmarkResult, componentValues []uint32) benchmarkResult {
// 	return benchmarkResult{
// 		extrinsicTime: benchmarkRes.ExtrinsicTime.ToBigInt().Uint64(),
// 		reads:         uint32(benchmarkRes.Reads),
// 		writes:        uint32(benchmarkRes.Writes),
// 		components:    componentValues,
// 	}
// }

// Executes a benchmark test.
// b is a go benchmark instance which must be provided by the calling test.
// testFn is a closure function that defines the test. It accepts a benchmarking instance param which is used to setup storage and run extrinsics.
// components is a registry for linear components variables which testFn may use.
func Run(b *testing.B, name string, testFn func(i *Instance) *benchmarkingtypes.BenchmarkResult, components ...*linear) {
	if *stepsFlag < 2 {
		b.Fatal("`steps` must be at least 2.")
	}

	runtimePath, ok := fullPath(*buildPathFlag)
	if !ok {
		b.Fatalf("`BUILD_PATH` must be a valid path to runtime wasm file: %s", *buildPathFlag)
	}

	// todo analysis
	// results := []benchmarkResult{}
	err := forEachStep(*stepsFlag, &components, func(currentStep int, currentComponentIndex int) {
		testDetails := "run"
		if len(components) > 0 {
			testDetails = fmt.Sprintf("run(step:%d,componentIndex:%d,componentValues:%d)", currentStep, currentComponentIndex, componentValues(components))
		}

		b.Run(testDetails, func(b *testing.B) {
			runtime := wazero_runtime.NewBenchInstanceWithTrie(b, runtimePath, trie.NewEmptyTrie())
			defer runtime.Stop()

			instance, err := newBenchmarkingInstance(runtime, 5)
			if err != nil {
				b.Fatalf("failed to create benchmarking instance: %v", err)
			}

			benchmarkResult := testFn(instance)
			if benchmarkResult == nil {
				b.Fatal("testFn must return non-nil excecution result, unless it returns non-nil error")
			}

			b.ReportMetric(float64(benchmarkResult.ExtrinsicTime.ToBigInt().Int64()), "time")
			b.ReportMetric(float64(benchmarkResult.Reads), "reads")
			b.ReportMetric(float64(benchmarkResult.Writes), "writes")

			// todo analysis
			// results = append(results, newBenchmarkResult(*benchmarkResult, componentValues(components)))
		})
	})

	if err != nil {
		b.Fatal(err)
	}

	// todo analysis
	// if *analysisFlag {
	// 	analyse(results)
	// }
}
