package io

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

type Trie interface {
	Blake2256OrderedRoot(key []byte, version int32) []byte
}

type trie struct {
	memoryTranslator utils.WasmMemoryTranslator
}

func NewTrie() Trie {
	return trie{
		memoryTranslator: utils.NewMemoryTranslator(),
	}
}

func (t trie) Blake2256OrderedRoot(key []byte, version int32) []byte {
	keyOffsetSize := t.memoryTranslator.BytesToOffsetAndSize(key)
	r := env.ExtTrieBlake2256OrderedRootVersion2(keyOffsetSize, version)
	return t.memoryTranslator.GetWasmMemorySlice(r, 32)
}
