package benchmarking

import (
	"flag"
	"fmt"
	"testing"

	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/lib/trie"
)

var (
	overheadWarmup         = flag.Int("overhead.warmup", 10, "How many warmup rounds before measuring.")
	overheadRepeat         = flag.Int("overhead.repeat", 100, "How many times the benchmark test should be repeated.")
	overheadMaxExtPerBlock = flag.Int("overhead.maxExtPerBlock", 500, "Maximum number of extrinsics per block")
)

func BenchmarkOverheadBlockExecutionWeight(t *testing.B) {
	config := OverheadConfig{
		Warmup: *overheadWarmup,
		Repeat: *overheadRepeat,
	}

	// todo set heapPages and dbCache when Gossamer starts supporting db caching
	runtime := wazero_runtime.NewBenchInstanceWithTrie(t, "../build/runtime.wasm", trie.NewEmptyTrie())
	defer runtime.Stop()

	instance, err := newBenchmarkingInstance(runtime, config.Repeat)
	if err != nil {
		t.Fatalf("failed to create benchmarking instance: [%v]", err)
	}

	stats := benchBlock(t, instance, config)
	fmt.Println("result stats")
	fmt.Println(stats)
	// TODO: Generate weight files
}

func BenchmarkOverheadBaseExtrinsicWeight(t *testing.B) {
	config := OverheadConfig{
		Warmup:         *overheadWarmup,
		Repeat:         *overheadRepeat,
		MaxExtPerBlock: *overheadMaxExtPerBlock,
	}

	// todo set heapPages and dbCache when Gossamer starts supporting db caching
	runtime := wazero_runtime.NewBenchInstanceWithTrie(t, "../build/runtime.wasm", trie.NewEmptyTrie())
	defer runtime.Stop()

	instance, err := newBenchmarkingInstance(runtime, config.Repeat)
	if err != nil {
		t.Fatalf("failed to create benchmarking instance: [%v]", err)
	}

	stats := benchExtrinsic(t, instance, config)
	fmt.Println("result stats")
	fmt.Println(stats)
	// TODO: Generate weight files
}
