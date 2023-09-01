package support

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
)

// HashStorageMap is a key-value storage map, which takes `prefix` and `name` that are hashed using hashing.Twox128 and appended before each key value.
type HashStorageMap[K, V sc.Encodable] struct {
	prefix      []byte
	name        []byte
	keyHashFunc func([]byte) []byte
	decodeFunc  func(buffer *bytes.Buffer) V
}

func NewHashStorageMap[K, V sc.Encodable](prefix []byte, name []byte, keyHashFunc func([]byte) []byte, decodeFunc func(buffer *bytes.Buffer) V) StorageMap[K, V] {
	return HashStorageMap[K, V]{
		prefix,
		name,
		keyHashFunc,
		decodeFunc,
	}
}

func (hsm HashStorageMap[K, V]) Get(k K) V {
	return storage.GetDecode(hsm.key(k), hsm.decodeFunc)
}

func (hsm HashStorageMap[K, V]) Exists(k K) bool {
	exists := storage.Exists(hsm.key(k))

	return exists != 0
}

func (hsm HashStorageMap[K, V]) Put(k K, value V) {
	storage.Set(hsm.key(k), value.Bytes())
}

func (hsm HashStorageMap[K, V]) Append(k K, value V) {
	storage.Append(hsm.key(k), value.Bytes())
}

func (hsm HashStorageMap[K, V]) TakeBytes(k K) []byte {
	return storage.TakeBytes(hsm.key(k))
}

func (hsm HashStorageMap[K, V]) Remove(k K) {
	storage.Clear(hsm.key(k))
}

func (hsm HashStorageMap[K, V]) Clear(limit sc.U32) {
	prefixHash := hashing.Twox128(hsm.prefix)
	nameHash := hashing.Twox128(hsm.name)

	storage.ClearPrefix(append(prefixHash, nameHash...), sc.NewOption[sc.U32](limit).Bytes())
}

func (hsm HashStorageMap[K, V]) key(key K) []byte {
	prefixHash := hashing.Twox128(hsm.prefix)
	nameHash := hashing.Twox128(hsm.name)

	keyBytes := key.Bytes()
	keyHash := hsm.keyHashFunc(keyBytes)

	concatKey := append(prefixHash, nameHash...)
	concatKey = append(concatKey, keyHash...)
	concatKey = append(concatKey, keyBytes...)

	return concatKey
}
