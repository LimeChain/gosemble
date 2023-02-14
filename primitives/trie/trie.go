package trie

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

func Blake2256OrderedRoot(key []byte) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	// TODO: switch to v2 once gossamer supports it or whenever runtime is imported in Substrate
	r := env.ExtTrieBlake2256OrderedRootVersion1(keyOffsetSize)
	return utils.ToWasmMemorySlice(r, 32)
}
