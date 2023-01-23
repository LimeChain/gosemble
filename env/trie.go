package env

/*
	Trie: Interface that provides trie related functionality
*/
//go:wasm-module env
//go:export ext_trie_blake2_256_ordered_root_version_2
func extTrieBlake2256OrderedRootVersion2(input int64, version int32) int32

func ExtTrieBlake2256OrderedRootVersion2(input int64, version int32) int32 {
	return extTrieBlake2256OrderedRootVersion2(input, version)
}
