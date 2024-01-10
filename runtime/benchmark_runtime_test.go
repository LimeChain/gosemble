package main

import (
	"testing"

	"github.com/ChainSafe/gossamer/lib/runtime"
	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/ChainSafe/gossamer/pkg/scale"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/assert"
)

func newBenchmarkingRuntime(b *testing.B) (*wazero_runtime.Instance, *runtime.Storage) {
	runtime := wazero_runtime.NewBenchInstanceWithTrie(b, WASM_RUNTIME, trie.NewEmptyTrie())
	return runtime, &runtime.Context.Storage
}

func newBenchmarkingRuntimeMetadata(b *testing.B, instance *wazero_runtime.Instance) *ctypes.Metadata {
	bMetadata, err := instance.Metadata()
	assert.NoError(b, err)

	var decoded []byte
	err = scale.Unmarshal(bMetadata, &decoded)
	assert.NoError(b, err)

	metadata := &ctypes.Metadata{}
	err = codec.Decode(decoded, metadata)
	assert.NoError(b, err)

	return metadata
}
