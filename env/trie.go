package env

/*
	Trie: Interface that provides trie related functionality
*/
//go:wasm-module env
//go:export ext_trie_blake2_256_ordered_root_version_1
func extTrieBlake2256OrderedRootVersion1(data int64) int32

func ExtTrieBlake2256OrderedRootVersion1(data int64) int32 {
	return extTrieBlake2256OrderedRootVersion1(data)
}
