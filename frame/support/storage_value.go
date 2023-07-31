package support

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
)

type StorageValue[T sc.Encodable] struct {
	prefix     []byte
	name       []byte
	decodeFunc func(buffer *bytes.Buffer) T
}

func NewStorageValue[T sc.Encodable](prefix []byte, name []byte, decodeFunc func(buffer *bytes.Buffer) T) *StorageValue[T] {
	return &StorageValue[T]{
		prefix,
		name,
		decodeFunc,
	}
}

func (sv StorageValue[T]) Get() T {
	return storage.GetDecode(sv.key(), sv.decodeFunc)
}

func (sv StorageValue[T]) Exists() bool {
	exists := storage.Exists(sv.key())

	return exists != 0
}

func (sv StorageValue[T]) Put(value T) {
	storage.Set(sv.key(), value.Bytes())
}

func (sv StorageValue[T]) Clear() {
	storage.Clear(sv.key())
}

func (sv StorageValue[T]) Append(value T) {
	storage.Append(sv.key(), value.Bytes())
}

func (sv StorageValue[T]) Take() T {
	return storage.TakeDecode(sv.key(), sv.decodeFunc)
}

func (sv StorageValue[T]) TakeBytes() []byte {
	return storage.TakeBytes(sv.key())
}

func (sv StorageValue[T]) key() []byte {
	prefixHash := hashing.Twox128(sv.prefix)
	nameHash := hashing.Twox128(sv.name)

	return append(prefixHash, nameHash...)
}
