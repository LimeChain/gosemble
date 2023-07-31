package support

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/storage"
)

// SimpleStorageValue takes a key upon initialisation and uses it in raw format to get storage values.
type SimpleStorageValue[T sc.Encodable] struct {
	key        []byte
	decodeFunc func(buffer *bytes.Buffer) T
}

func NewSimpleStorageValue[T sc.Encodable](key []byte, decodeFunc func(buffer *bytes.Buffer) T) *SimpleStorageValue[T] {
	return &SimpleStorageValue[T]{
		key,
		decodeFunc,
	}
}

func (ssv SimpleStorageValue[T]) Get() T {
	return storage.GetDecode(ssv.key, ssv.decodeFunc)
}

func (ssv SimpleStorageValue[T]) Exists() bool {
	exists := storage.Exists(ssv.key)

	return exists != 0
}

func (ssv SimpleStorageValue[T]) Put(value T) {
	storage.Set(ssv.key, value.Bytes())
}

func (ssv SimpleStorageValue[T]) Clear() {
	storage.Clear(ssv.key)
}

func (ssv SimpleStorageValue[T]) Append(value T) {
	storage.Append(ssv.key, value.Bytes())
}

func (ssv SimpleStorageValue[T]) TakeBytes() []byte {
	return storage.TakeBytes(ssv.key)
}

func (ssv SimpleStorageValue[T]) Take() T {
	return storage.TakeDecode(ssv.key, ssv.decodeFunc)
}
