package main

import (
	"bytes"
	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/types"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/ChainSafe/gossamer/lib/runtime/wasmer"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/stretchr/testify/assert"
)

var (
	keySystemHash, _         = common.Twox128Hash(constants.KeySystem)
	keyBlockHash, _          = common.Twox128Hash(constants.KeyBlockHash)
	keyExecutionPhaseHash, _ = common.Twox128Hash(constants.KeyExecutionPhase)
	keyLastRuntime, _        = common.Twox128Hash(constants.KeyLastRuntimeUpgrade)
	keyNumberHash, _         = common.Twox128Hash(constants.KeyNumber)
	keyParentHash, _         = common.Twox128Hash(constants.KeyParentHash)
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
	parentHash := common.Hash{}
	stateRoot := common.Hash{}
	extrinsicsRoot := common.Hash{}
	blockNumber := uint(1)

	digest := gossamertypes.NewDigest()

	header, err := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, blockNumber, digest)
	require.NoError(t, err)

	encodedHeader, err := scale.Marshal(*header)
	if err != nil {
		panic(err)
	}

	storage := trie.NewEmptyTrie()

	rt := wasmer.NewLocalTestInstanceWithTrie(t, WASM_RUNTIME, storage)
	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.Nil(t, err)

	buffer := &bytes.Buffer{}

	lrui := types.LastRuntimeUpgradeInfo{
		SpecVersion: constants.SPEC_VERSION,
		SpecName:    constants.SPEC_NAME,
	}
	lrui.Encode(buffer)
	assert.Equal(t, buffer.Bytes(), storage.Get(append(keySystemHash, keyLastRuntime...)))
	buffer.Reset()

	sc.U32(0).Encode(buffer)
	assert.Equal(t, buffer.Bytes(), storage.Get(constants.KeyExtrinsicIndex))
	buffer.Reset()

	sc.U32(constants.ExecutionPhaseApplyExtrinsic).Encode(buffer)
	assert.Equal(t, buffer.Bytes(), storage.Get(append(keySystemHash, keyExecutionPhaseHash...)))
	buffer.Reset()

	sc.U32(blockNumber).Encode(buffer)
	assert.Equal(t, buffer.Bytes(), storage.Get(append(keySystemHash, keyNumberHash...)))
	buffer.Reset()
	assert.Equal(t, parentHash.ToBytes(), storage.Get(append(keySystemHash, keyParentHash...)))

	blockHashKey := append(keySystemHash, keyBlockHash...)
	prevBlock := sc.U32(blockNumber - 1)
	prevBlock.Encode(buffer)
	numHash, err := common.Twox64(buffer.Bytes())
	if err != nil {
		panic(err)
	}
	assert.Equal(t, parentHash.ToBytes(), storage.Get(append(blockHashKey, numHash...)))
	buffer.Reset()
}
