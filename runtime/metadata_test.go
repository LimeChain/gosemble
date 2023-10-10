package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_Metadata_Encoding_Success(t *testing.T) {
	runtime, _ := newTestRuntime(t)

	bMetadata, err := runtime.Metadata()
	assert.NoError(t, err)

	buffer := bytes.NewBuffer(bMetadata)

	// Decode Compact Length
	_ = sc.DecodeCompact(buffer)

	// Copy bytes for assertion after re-encode.
	bMetadataCopy := make([]byte, buffer.Len())
	copy(bMetadataCopy, buffer.Bytes())

	fmt.Println(hex.EncodeToString(buffer.Bytes()))

	metadata, err := types.DecodeMetadata(buffer)
	assert.NoError(t, err)

	// Assert encoding of previously decoded
	assert.Equal(t, bMetadataCopy, metadata.Bytes())
}
