package benchmarking

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/ChainSafe/gossamer/lib/runtime"
	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/lib/trie"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/benchmarking"
	benchmarkingtypes "github.com/LimeChain/gosemble/primitives/benchmarking"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

const WASM_RUNTIME = "../build/runtime.wasm"

var (
	steps         *int    = flag.Int("steps", 50, "Select how many samples we should take across the variable components.")
	repeat        *int    = flag.Int("repeat", 20, "Select how many repetitions of this benchmark should run from within the wasm.")
	heapPages     *int    = flag.Int("heap-pages", 4096, "Cache heap allocation pages.")
	dbCache       *int    = flag.Int("db-cache", 1024, "Limit the memory the database cache can use.")
	gc            *string = flag.String("gc", "", "GC flag used for building the runtime")
	tinyGoVersion *string = flag.String("tinygoversion", "", "TinyGO version used for building the runtime")
	target        *string = flag.String("target", "", "Target used for building the runtime")
)

// Executes a benchmark test.
// b is a go benchmark instance which must be provided by the calling test.
// outputPath is the path where weight file will be generated.
// testFn is a closure function that defines the test. It accepts a benchmarking instance param which is used to setup storage and run extrinsics.
// components is a registry for linear components variables which testFn may use.
func RunDispatchCall(b *testing.B, outputPath string, testFn func(i *Instance), components ...*linear) {
	if *steps < 2 {
		b.Fatal("`steps` must be at least 2.")
	}

	results := []benchmarkResult{}

	if len(components) == 0 {
		b.Run("Step 1", func(b *testing.B) {
			res := runTestFn(b, testFn)
			results = append(results, newBenchmarkResult(res, []uint32{}))
		})
	}

	// iterate components for each step
	for i, component := range components {
		for y, v := range component.values(*steps) {
			component.setValue(v)

			componentValues := componentValues(components)

			testName := fmt.Sprintf("Step %d/ComponentIndex %d/ComponentValues %d", y+1, i, componentValues)

			b.Run(testName, func(b *testing.B) {
				res := runTestFn(b, testFn)
				results = append(results, newBenchmarkResult(res, componentValues))
			})
		}
	}

	// generate weight file
	analysis := medianSlopesAnalysis(results)
	fmt.Println(analysis.String())

	if err := generateWeightFile(extrinsicTemplate, outputPath, analysis.String(), analysis.baseExtrinsicTime, analysis.baseReads, analysis.baseWrites); err != nil {
		b.Fatalf("failed to generate weight file: %v", err)
	}
}

func RunHook(b *testing.B,
	hookName string,
	setupFn func(storage *runtime.Storage),
	validateFn func(storage *runtime.Storage),
) benchmarkingtypes.BenchmarkResult {
	// todo set heapPages and dbCache when gosammer starts supporting db caching
	runtime := wazero_runtime.NewBenchInstanceWithTrie(b, WASM_RUNTIME, trie.NewEmptyTrie())
	defer runtime.Stop()

	instance, err := newBenchmarkingInstance(runtime, *repeat)
	if err != nil {
		b.Fatalf("failed to create benchmarking instance: %v", err)
	}

	hookBytes := append(sc.Str(hookName).Bytes(), sc.U64(1).Bytes()...)
	hookBytes = append(hookBytes, types.WeightZero().Bytes()...)

	benchmarkConfig := benchmarking.BenchmarkConfig{
		InternalRepeats: sc.U32(b.N),
		Benchmark:       sc.BytesToSequenceU8(hookBytes),
		Origin:          types.NewRawOriginNone(),
	}

	setupFn(instance.Storage())

	res, err := runtime.Exec("Benchmark_hook", benchmarkConfig.Bytes())
	assert.NoError(b, err)

	validateFn(instance.Storage())

	benchmarkResult, err := benchmarking.DecodeBenchmarkResult(bytes.NewBuffer(res))
	assert.NoError(b, err)

	b.ReportMetric(float64(benchmarkResult.Time.ToBigInt().Int64()), "time")
	b.ReportMetric(float64(benchmarkResult.Reads), "reads")
	b.ReportMetric(float64(benchmarkResult.Writes), "writes")

	return benchmarkResult
}

func runTestFn(b *testing.B, testFn func(i *Instance)) benchmarkingtypes.BenchmarkResult {
	// todo set heapPages and dbCache when gosammer starts supporting db caching
	runtime := wazero_runtime.NewBenchInstanceWithTrie(b, WASM_RUNTIME, trie.NewEmptyTrie())
	defer runtime.Stop()

	instance, err := newBenchmarkingInstance(runtime, *repeat)
	if err != nil {
		b.Fatalf("failed to create benchmarking instance: %v", err)
	}

	err = instance.InitializeBlock(blockNumber, dateTime)
	if err != nil {
		b.Fatalf("failed to initialize block: number [%d], date time [%d], err [%v]", blockNumber, dateTime, err)
	}

	testFn(instance)

	benchmarkResult := instance.benchmarkResult
	if benchmarkResult == nil {
		b.Fatal("No valid extrinsic or block call could be found in testFn")
	}

	b.ReportMetric(float64(benchmarkResult.Time.ToBigInt().Int64()), "time")
	b.ReportMetric(float64(benchmarkResult.Reads), "reads")
	b.ReportMetric(float64(benchmarkResult.Writes), "writes")

	return *benchmarkResult
}
