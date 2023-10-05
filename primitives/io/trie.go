package io

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

type Trie interface {
	Blake2256OrderedRoot(key []byte, version int32) []byte
}

func NewTrie() Trie {
	return trie{}
}

type trie struct {
}

func (t trie) Blake2256OrderedRoot(key []byte, version int32) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	r := env.ExtTrieBlake2256OrderedRootVersion2(keyOffsetSize, version)
	return utils.ToWasmMemorySlice(r, 32)
}
