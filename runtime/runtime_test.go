package main

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/gossamer/lib/runtime/wasmer"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/types"
	"github.com/stretchr/testify/assert"
)

const WASM_RUNTIME = "../build/runtime.wasm" // node_template_runtime.wasm

func Test_CoreVersion(t *testing.T) {
	rt := wasmer.NewLocalTestInstanceWithTrie(t, WASM_RUNTIME, trie.NewEmptyTrie())
	res, err := rt.Exec("Core_version", []byte{})
	assert.Nil(t, err)
	// t.Logf("%x", res)
	buffer := bytes.Buffer{}
	buffer.Write(res)
	resultVersion := types.DecodeVersionData(&buffer)
	// t.Log(resultVersion)
	assert.Equal(t, constants.RuntimeVersion, resultVersion)
}

func Test_CoreInitializeBlock(t *testing.T) {

}
