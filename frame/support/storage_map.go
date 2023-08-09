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

func (sm StorageMap[K, V]) Get(k K) V {
	return storage.GetDecode(sm.key(k), sm.decodeFunc)
}

func (sm StorageMap[K, V]) Exists(k K) bool {
	exists := storage.Exists(sm.key(k))

	return exists != 0
}

func (sm StorageMap[K, V]) Put(k K, value V) {
	storage.Set(sm.key(k), value.Bytes())
}

func (sm StorageMap[K, V]) Append(k K, value V) {
	storage.Append(sm.key(k), value.Bytes())
}

func (sm StorageMap[K, V]) TakeBytes(k K) []byte {
	return storage.TakeBytes(sm.key(k))
}

func (sm StorageMap[K, V]) Remove(k K) {
	storage.Clear(sm.key(k))
}

func (sm StorageMap[K, V]) Clear(limit sc.U32) {
	prefixHash := hashing.Twox128(sm.prefix)
	nameHash := hashing.Twox128(sm.name)

	storage.ClearPrefix(append(prefixHash, nameHash...), sc.NewOption[sc.U32](limit).Bytes())
}

func (sm StorageMap[K, V]) key(key K) []byte {
	prefixHash := hashing.Twox128(sm.prefix)
	nameHash := hashing.Twox128(sm.name)

	keyBytes := key.Bytes()
	keyHash := sm.keyHashFunc(keyBytes)

	concatKey := append(prefixHash, nameHash...)
	concatKey = append(concatKey, keyHash...)
	concatKey = append(concatKey, keyBytes...)

	return concatKey
}
