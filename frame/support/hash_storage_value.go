package support

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/hashing"
)

// HashStorageValue is a storage value, which takes `prefix` and `name` that are hashed using hashing.Twox128 and appended before each key value.
type HashStorageValue[T sc.Encodable] struct {
	baseStorageValue[T]
	prefix []byte
	name   []byte
}

func NewHashStorageValue[T sc.Encodable](prefix []byte, name []byte, decodeFunc func(buffer *bytes.Buffer) T) StorageValue[T] {
	return HashStorageValue[T]{
		baseStorageValue: baseStorageValue[T]{
			decodeFunc: decodeFunc,
		},
		prefix: prefix,
		name:   name,
	}
}

func NewHashStorageValueWithDefault[T sc.Encodable](prefix []byte, name []byte, decodeFunc func(buffer *bytes.Buffer) T, defaultValue *T) StorageValue[T] {
	return HashStorageValue[T]{
		baseStorageValue: baseStorageValue[T]{
			decodeFunc:   decodeFunc,
			defaultValue: defaultValue,
		},
		prefix: prefix,
		name:   name,
	}
}

func (hsv HashStorageValue[T]) Get() T {
	return hsv.baseStorageValue.get(hsv.key())
}

func (hsv HashStorageValue[T]) GetBytes() sc.Option[sc.Sequence[sc.U8]] {
	return hsv.baseStorageValue.getBytes(hsv.key())
}

func (hsv HashStorageValue[T]) Exists() bool {
	return hsv.baseStorageValue.exists(hsv.key())
}

func (hsv HashStorageValue[T]) Put(value T) {
	hsv.baseStorageValue.put(hsv.key(), value)
}

func (hsv HashStorageValue[T]) Clear() {
	hsv.baseStorageValue.clear(hsv.key())
}

func (hsv HashStorageValue[T]) Append(value T) {
	hsv.baseStorageValue.append(hsv.key(), value)
}

func (hsv HashStorageValue[T]) Take() T {
	return hsv.baseStorageValue.take(hsv.key())
}

func (hsv HashStorageValue[T]) TakeBytes() []byte {
	return hsv.baseStorageValue.takeBytes(hsv.key())
}

func (hsv HashStorageValue[T]) DecodeLen() sc.Option[sc.U64] {
	return hsv.baseStorageValue.decodeLen(hsv.key())
}

func (hsv HashStorageValue[T]) key() []byte {
	prefixHash := hashing.Twox128(hsv.prefix)
	nameHash := hashing.Twox128(hsv.name)

	return append(prefixHash, nameHash...)
}
