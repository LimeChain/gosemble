package main

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/gossamer/lib/runtime/wasmer"
	"github.com/ChainSafe/gossamer/lib/trie"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_Metadata_Encoding_Success(t *testing.T) {
	runtime := wasmer.NewTestInstanceWithTrie(t, POLKADOT_RUNTIME, trie.NewEmptyTrie())
	bMetadata, err := runtime.Metadata()
	assert.NoError(t, err)

	buffer := bytes.NewBuffer(bMetadata)

	// Decode Compact Length
	_ = sc.DecodeCompact(buffer)

	// Copy bytes for assertion after re-encode.
	bMetadataCopy := make([]byte, buffer.Len())
	copy(bMetadataCopy, buffer.Bytes())

	metadata, err := types.DecodeMetadata(buffer)
	assert.NoError(t, err)

	// Assert encoding of previously decoded
	assert.Equal(t, bMetadataCopy, metadata.Bytes())
}
