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

func (sv StorageValue[T]) GetBytes() sc.Option[sc.Sequence[sc.U8]] {
	return storage.Get(sv.key())
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

func (sv StorageValue[T]) DecodeLen() sc.Option[sc.U64] {
	// `Compact<u32>` is 5 bytes in maximum.
	data := [5]byte{}
	option := storage.Read(sv.key(), data[:], 0)

	if !option.HasValue {
		return sc.NewOption[sc.U64](nil)
	}

	length := option.Value
	if length > sc.U32(len(data)) {
		length = sc.U32(len(data))
	}

	buffer := &bytes.Buffer{}
	buffer.Write(data[:length])

	compact := sc.DecodeCompact(buffer)
	toLen := sc.U64(compact.ToBigInt().Uint64())

	return sc.NewOption[sc.U64](toLen)
}

func (sv StorageValue[T]) key() []byte {
	prefixHash := hashing.Twox128(sv.prefix)
	nameHash := hashing.Twox128(sv.name)

	return append(prefixHash, nameHash...)
}
