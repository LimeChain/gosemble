package benchmarking

import (
	"testing"

	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/lib/trie"
)

func BenchmarkBlockExecutionWeight(t *testing.B) {
	config := OverheadConfig{
		Warmup: 10,
		Repeat: 100,
	}

	// todo set heapPages and dbCache when Gossamer starts supporting db caching
	runtime := wazero_runtime.NewBenchInstanceWithTrie(t, "../build/runtime.wasm", trie.NewEmptyTrie())
	defer runtime.Stop()

	instance, err := newBenchmarkingInstance(runtime, *repeat)
	if err != nil {
		t.Fatalf("failed to create benchmarking instance: [%v]", err)
	}

	benchBlock(t, instance, config)
	// TODO: Generate weight files
}

func BenchmarkBaseExtrinsicWeight(t *testing.B) {
	config := OverheadConfig{
		Warmup:         10,
		Repeat:         100,
		MaxExtPerBlock: 2,
	}

	// todo set heapPages and dbCache when Gossamer starts supporting db caching
	runtime := wazero_runtime.NewBenchInstanceWithTrie(t, "../build/runtime.wasm", trie.NewEmptyTrie())
	defer runtime.Stop()

	instance, err := newBenchmarkingInstance(runtime, *repeat)
	if err != nil {
		t.Fatalf("failed to create benchmarking instance: [%v]", err)
	}

	benchExtrinsic(t, instance, config)
	// TODO: Generate weight files
}
