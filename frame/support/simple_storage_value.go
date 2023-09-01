package support

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// SimpleStorageValue takes a key upon initialisation and uses it in raw format to get storage values.
type SimpleStorageValue[T sc.Encodable] struct {
	baseStorageValue[T]
	key []byte
}

func NewSimpleStorageValue[T sc.Encodable](key []byte, decodeFunc func(buffer *bytes.Buffer) T) StorageValue[T] {
	return SimpleStorageValue[T]{
		baseStorageValue: baseStorageValue[T]{
			decodeFunc: decodeFunc,
		},
		key: key,
	}
}

func (ssv SimpleStorageValue[T]) Get() T {
	return ssv.baseStorageValue.get(ssv.key)
}

func (ssv SimpleStorageValue[T]) GetBytes() sc.Option[sc.Sequence[sc.U8]] {
	return ssv.baseStorageValue.getBytes(ssv.key)
}

func (ssv SimpleStorageValue[T]) Exists() bool {
	return ssv.baseStorageValue.exists(ssv.key)
}

func (ssv SimpleStorageValue[T]) Put(value T) {
	ssv.baseStorageValue.put(ssv.key, value)
}

func (ssv SimpleStorageValue[T]) Clear() {
	ssv.baseStorageValue.clear(ssv.key)
}

func (ssv SimpleStorageValue[T]) Append(value T) {
	ssv.baseStorageValue.append(ssv.key, value)
}

func (ssv SimpleStorageValue[T]) TakeBytes() []byte {
	return ssv.baseStorageValue.takeBytes(ssv.key)
}

func (ssv SimpleStorageValue[T]) Take() T {
	return ssv.baseStorageValue.take(ssv.key)
}

func (ssv SimpleStorageValue[T]) DecodeLen() sc.Option[sc.U64] {
	return ssv.baseStorageValue.decodeLen(ssv.key)
}
