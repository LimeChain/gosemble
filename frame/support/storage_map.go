package support

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
)

type StorageMap[K, V sc.Encodable] struct {
	prefix      []byte
	name        []byte
	keyHashFunc func([]byte) []byte
	decodeFunc  func(buffer *bytes.Buffer) V
}

func NewStorageMap[K, V sc.Encodable](prefix []byte, name []byte, keyHashFunc func([]byte) []byte, decodeFunc func(buffer *bytes.Buffer) V) *StorageMap[K, V] {
	return &StorageMap[K, V]{
		prefix,
		name,
		keyHashFunc,
		decodeFunc,
	}
}

func (sm StorageMap[K, V]) Get(key K) V {
	prefixHash := hashing.Twox128(sm.prefix)
	nameHash := hashing.Twox128(sm.name)

	keyBytes := key.Bytes()
	keyHash := sm.keyHashFunc(keyBytes)

	concatKey := append(prefixHash, nameHash...)
	concatKey = append(concatKey, keyHash...)
	concatKey = append(concatKey, keyBytes...)

	return storage.GetDecode(concatKey, sm.decodeFunc)
}

func (sm StorageMap[K, V]) Put(key K, value V) {
	prefixHash := hashing.Twox128(sm.prefix)
	nameHash := hashing.Twox128(sm.name)

	keyBytes := key.Bytes()
	keyHash := sm.keyHashFunc(keyBytes)

	concatKey := append(prefixHash, nameHash...)
	concatKey = append(concatKey, keyHash...)
	concatKey = append(concatKey, keyBytes...)

	storage.Set(concatKey, value.Bytes())
}

func (sm StorageMap[K, V]) Append(key K, value V) {
	prefixHash := hashing.Twox128(sm.prefix)
	nameHash := hashing.Twox128(sm.name)

	keyBytes := key.Bytes()
	keyHash := sm.keyHashFunc(keyBytes)

	concatKey := append(prefixHash, nameHash...)
	concatKey = append(concatKey, keyHash...)
	concatKey = append(concatKey, keyBytes...)

	storage.Append(concatKey, value.Bytes())
}
