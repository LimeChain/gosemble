package main

import (
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/runtime/wasmer"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_CoreInitializeBlock(t *testing.T) {
	expectedStorageDigest := gossamertypes.NewDigest()

	digest := gossamertypes.NewDigest()

	preRuntimeDigestItem := gossamertypes.NewDigestItem()
	assert.NoError(t, preRuntimeDigestItem.Set(preRuntimeDigest))

	sealDigestItem := gossamertypes.NewDigestItem()
	assert.NoError(t, sealDigestItem.Set(sealDigest))

	prdi, err := preRuntimeDigestItem.Value()
	assert.NoError(t, err)
	assert.NoError(t, digest.Add(prdi))

	sdi, err := sealDigestItem.Value()
	assert.NoError(t, err)
	assert.NoError(t, digest.Add(sdi))
	assert.NoError(t, expectedStorageDigest.Add(prdi))

	header := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, blockNumber, digest)
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

	lrui := types.LastRuntimeUpgradeInfo{
		SpecVersion: sc.ToCompact(constants.SPEC_VERSION),
		SpecName:    constants.SPEC_NAME,
	}
	assert.Equal(t, lrui.Bytes(), storage.Get(append(keySystemHash, keyLastRuntime...)))

	encExtrinsicIndex0, _ := scale.Marshal(uint32(0))
	assert.Equal(t, encExtrinsicIndex0, storage.Get(constants.KeyExtrinsicIndex))

	encExecutionPhaseApplyExtrinsic, _ := scale.Marshal(uint32(0))
	assert.Equal(t, encExecutionPhaseApplyExtrinsic, storage.Get(append(keySystemHash, keyExecutionPhaseHash...)))

	encBlockNumber, _ := scale.Marshal(uint32(blockNumber))
	assert.Equal(t, encBlockNumber, storage.Get(append(keySystemHash, keyNumberHash...)))

	encExpectedDigest, err := scale.Marshal(expectedStorageDigest)
	assert.NoError(t, err)
	assert.Equal(t, encExpectedDigest, storage.Get(append(keySystemHash, keyDigestHash...)))
	assert.Equal(t, parentHash.ToBytes(), storage.Get(append(keySystemHash, keyParentHash...)))

	blockHashKey := append(keySystemHash, keyBlockHash...)
	encPrevBlock, _ := scale.Marshal(uint32(blockNumber - 1))
	numHash, err := common.Twox64(encPrevBlock)
	assert.NoError(t, err)

	blockHashKey = append(blockHashKey, numHash...)
	blockHashKey = append(blockHashKey, encPrevBlock...)
	assert.Equal(t, parentHash.ToBytes(), storage.Get(blockHashKey))
}
