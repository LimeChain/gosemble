package main

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/assert"
)

func Test_Metadata_Encoding_Success(t *testing.T) {
	runtime, _ := newTestRuntime(t)
	gossamerMetadata := runtimeMetadata(t, runtime)

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

	// Encode gossamer Metadata
	bGossamerMetadata, err := codec.Encode(gossamerMetadata)
	assert.NoError(t, err)

	assert.Equal(t, metadata.Bytes(), bGossamerMetadata)
}

func Test_Metadata_Versions_Correct_Versions(t *testing.T) {
	runtime, _ := newTestRuntime(t)

	metadataVersions, err := runtime.Exec("Metadata_versions", []byte{})
	assert.NoError(t, err)

	buffer := &bytes.Buffer{}
	buffer.Write(metadataVersions)

	versions := sc.DecodeSequence[sc.U32](buffer)
	buffer.Reset()

	expectedVersions := sc.Sequence[sc.U32]{
		sc.U32(14), sc.U32(15),
	}

	assert.Equal(t, versions, expectedVersions)
}

func Test_Metadata_At_Version_14(t *testing.T) {
	runtime, _ := newTestRuntime(t)
	gossamerMetadata := runtimeMetadata(t, runtime)

	version14 := sc.U32(14)

	bMetadata, err := runtime.Exec("Metadata_at_version", version14.Bytes())
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

	// Encode gossamer Metadata
	bGossamerMetadata, err := codec.Encode(gossamerMetadata)
	assert.NoError(t, err)

	assert.Equal(t, metadata.Bytes(), bGossamerMetadata)
}

func Test_Metadata_At_Version_15(t *testing.T) {
	runtime, _ := newTestRuntime(t)

	version15 := sc.U32(15)

	bMetadata, err := runtime.Exec("Metadata_at_version", version15.Bytes())
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
