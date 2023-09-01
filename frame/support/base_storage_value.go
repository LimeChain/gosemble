package support

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/storage"
)

type baseStorageValue[T sc.Encodable] struct {
	decodeFunc   func(buffer *bytes.Buffer) T
	defaultValue *T
}

func (bsv baseStorageValue[T]) get(key []byte) T {
	if bsv.defaultValue == nil {
		return storage.GetDecode(key, bsv.decodeFunc)
	}

	return storage.GetDecodeOnEmpty(key, bsv.decodeFunc, *bsv.defaultValue)
}

func (bsv baseStorageValue[T]) getBytes(key []byte) sc.Option[sc.Sequence[sc.U8]] {
	return storage.Get(key)
}

func (bsv baseStorageValue[T]) exists(key []byte) bool {
	exists := storage.Exists(key)

	return exists != 0
}

func (bsv baseStorageValue[T]) put(key []byte, value T) {
	storage.Set(key, value.Bytes())
}

func (bsv baseStorageValue[T]) clear(key []byte) {
	storage.Clear(key)
}

func (bsv baseStorageValue[T]) append(key []byte, value T) {
	storage.Append(key, value.Bytes())
}

func (bsv baseStorageValue[T]) takeBytes(key []byte) []byte {
	return storage.TakeBytes(key)
}

func (bsv baseStorageValue[T]) take(key []byte) T {
	return storage.TakeDecode(key, bsv.decodeFunc)
}

func (bsv baseStorageValue[T]) decodeLen(key []byte) sc.Option[sc.U64] {
	// `Compact<u32>` is 5 bytes in maximum.
	data := [5]byte{}
	option := storage.Read(key, data[:], 0)

	if !option.HasValue {
		return sc.NewOption[sc.U64](nil)
	}

	length := option.Value
	if length.Gt(sc.U32(len(data))) {
		length = sc.U32(len(data))
	}

	buffer := &bytes.Buffer{}
	buffer.Write(data[:length])

	compact := sc.DecodeCompact(buffer)
	toLen := sc.U64(compact.ToBigInt().Uint64())

	return sc.NewOption[sc.U64](toLen)
}
