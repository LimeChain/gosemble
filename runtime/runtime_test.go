package main

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/gossamer/lib/runtime/storage"
	"github.com/ChainSafe/gossamer/lib/runtime/wasmer"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/types"
	"github.com/stretchr/testify/assert"
)

const WASM_RUNTIME = "../build/runtime.wasm" // -v9160

func Test_CoreVersion(t *testing.T) {
	rt := wasmer.NewLocalTestInstanceWithTrie(t, WASM_RUNTIME, trie.NewEmptyTrie())
	res, err := rt.Exec("Core_version", []byte{})
	assert.Nil(t, err)

	resultVersion := types.VersionData{}
	t.Logf("%#x", res)
	err = resultVersion.Decode(res)
	assert.Nil(t, err)
	t.Logf("%#x", res)
	assert.Equal(t, constants.RuntimeVersion, resultVersion)
}

func Test_CoreInitializeBlock(t *testing.T) {
	t.Skip()
	tt := trie.NewEmptyTrie()
	rt := wasmer.NewLocalTestInstanceWithTrie(t, WASM_RUNTIME, tt)

	scaleEncHeader, err := (&types.Header{}).Encode()
	assert.Nil(t, err)

	_, err = rt.Exec("Core_initialize_block", scaleEncHeader)
	assert.Nil(t, err)

	ts := storage.NewTrieState(tt)

	var buffer = bytes.Buffer{}
	var encoder = goscale.Encoder{Writer: &buffer}
	var decoder = goscale.Decoder{Reader: &buffer}

	key := "k1"
	encoder.EncodeString(key)
	encKey := buffer.Bytes()
	buffer.Read(encKey)

	encValue := ts.Get(encKey)
	encoder.Write(encValue)
	value := decoder.DecodeString()

	t.Logf("%s", value)

	key = "k2"
	encoder.EncodeString(key)
	encKey = buffer.Bytes()
	buffer.Read(encKey)

	encValue = ts.Get(encKey)
	encoder.Write(encValue)
	value = decoder.DecodeString()

	t.Logf("%s", value)
}
