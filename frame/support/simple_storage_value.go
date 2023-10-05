package support

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// SimpleStorageValue takes a key upon initialisation and uses it in raw format to get storage values.
type SimpleStorageValue[T sc.Encodable] struct {
	baseStorage[T]
	key []byte
}

func NewSimpleStorageValue[T sc.Encodable](key []byte, decodeFunc func(buffer *bytes.Buffer) T) StorageValue[T] {
	return SimpleStorageValue[T]{
		baseStorage: newBaseStorage[T](decodeFunc, nil),
		key:         key,
	}
}

func (ssv SimpleStorageValue[T]) Get() T {
	return ssv.baseStorage.get(ssv.key)
}

func (ssv SimpleStorageValue[T]) GetBytes() sc.Option[sc.Sequence[sc.U8]] {
	return ssv.baseStorage.getBytes(ssv.key)
}

func (ssv SimpleStorageValue[T]) Exists() bool {
	return ssv.baseStorage.exists(ssv.key)
}

func (ssv SimpleStorageValue[T]) Put(value T) {
	ssv.baseStorage.put(ssv.key, value)
}

func (ssv SimpleStorageValue[T]) Clear() {
	ssv.baseStorage.clear(ssv.key)
}

func (ssv SimpleStorageValue[T]) Append(value T) {
	ssv.baseStorage.append(ssv.key, value)
}

func (ssv SimpleStorageValue[T]) TakeBytes() []byte {
	return ssv.baseStorage.takeBytes(ssv.key)
}

func (ssv SimpleStorageValue[T]) Take() T {
	return ssv.baseStorage.take(ssv.key)
}

func (ssv SimpleStorageValue[T]) DecodeLen() sc.Option[sc.U64] {
	return ssv.baseStorage.decodeLen(ssv.key)
}
