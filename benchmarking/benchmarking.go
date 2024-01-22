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
	steps     *int = flag.Int("steps", 50, "Select how many samples we should take across the variable components.")
	heapPages *int = flag.Int("heap-pages", 2048, "Cache heap allocation pages.")
	dbCache   *int = flag.Int("db-cache", 1024, "Limit the memory the database cache can use.")
)

// Executes a benchmark test.
// b is a go benchmark instance which must be provided by the calling test.
// testFn is a closure function that defines the test. It accepts a benchmarking instance param which is used to setup storage and run extrinsics.
// components is a registry for linear components variables which testFn may use.
func Run(b *testing.B, name string, testFn func(i *Instance) *benchmarkingtypes.BenchmarkResult, components ...*linear) {
	if *steps < 2 {
		b.Fatal("`steps` must be at least 2.")
	}

	results := []benchmarkResult{}

	if len(components) == 0 {
		res := runTestFn(b, testFn)
		results = append(results, newBenchmarkResult(res, []uint32{}))
	}

	// iterate components for each step
	for i, component := range components {
		cValues, err := component.values(*steps)
		if err != nil {
			b.Fatal(err)
		}

		for y, v := range cValues {
			component.setValue(v)

			componentValues := componentValues(components)

			testName := fmt.Sprintf("Step %d/ComponentIndex %d/ComponentValues %d", y+1, i, componentValues)

			b.Run(testName, func(b *testing.B) {
				res := runTestFn(b, testFn)
				results = append(results, newBenchmarkResult(res, componentValues))
			})
		}
	}

	// todo weights
	fmt.Printf("min squares analysis: %v", minSquareAnalysis(results))
}

func runTestFn(b *testing.B, testFn func(i *Instance) *benchmarkingtypes.BenchmarkResult) benchmarkingtypes.BenchmarkResult {
	// todo set heapPages and dbCache when gosammer starts supporting db caching
	runtime := wazero_runtime.NewBenchInstanceWithTrie(b, "../build/runtime.wasm", trie.NewEmptyTrie())
	defer runtime.Stop()

	instance, err := newBenchmarkingInstance(runtime, b.N)
	if err != nil {
		b.Fatalf("failed to create benchmarking instance: %v", err)
	}

	benchmarkResult := testFn(instance)
	if benchmarkResult == nil {
		b.Fatal("testFn must return non-nil excecution result, unless it returns non-nil error")
	}

	b.ReportMetric(float64(benchmarkResult.ExtrinsicTime.ToBigInt().Int64()), "ns/op")
	b.ReportMetric(float64(benchmarkResult.Reads), "reads")
	b.ReportMetric(float64(benchmarkResult.Writes), "writes")

	return *benchmarkResult
}
