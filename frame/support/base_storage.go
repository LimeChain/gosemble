package support

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/io"
)

type baseStorage[T sc.Encodable] struct {
	storage      io.Storage
	decodeFunc   func(buffer *bytes.Buffer) T
	defaultValue *T
}

func newBaseStorage[T sc.Encodable](decodeFunc func(buffer *bytes.Buffer) T, defaultValue *T) baseStorage[T] {
	return baseStorage[T]{
		storage:      io.NewStorage(),
		decodeFunc:   decodeFunc,
		defaultValue: defaultValue,
	}
}

func (bs baseStorage[T]) get(key []byte) T {
	if bs.defaultValue == nil {
		return bs.getDecode(key, bs.decodeFunc)
	}

	return bs.getDecodeOnEmpty(key, bs.decodeFunc, *bs.defaultValue)
}

func (bs baseStorage[T]) getBytes(key []byte) sc.Option[sc.Sequence[sc.U8]] {
	return bs.storage.Get(key)
}

func (bs baseStorage[T]) exists(key []byte) bool {
	return bs.storage.Exists(key)
}

func (bs baseStorage[T]) put(key []byte, value T) {
	bs.storage.Set(key, value.Bytes())
}

func (bs baseStorage[T]) clear(key []byte) {
	bs.storage.Clear(key)
}

func (bs baseStorage[T]) append(key []byte, value T) {
	bs.storage.Append(key, value.Bytes())
}

func (bs baseStorage[T]) take(key []byte) T {
	return bs.takeDecode(key, bs.decodeFunc)
}

func (bs baseStorage[T]) decodeLen(key []byte) sc.Option[sc.U64] {
	// `Compact<u32>` is 5 bytes in maximum.
	data := [5]byte{}
	option := bs.storage.Read(key, data[:], 0)

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

// getDecode gets the storage value and returns it decoded. The result from Get is Option<sc.Sequence[sc.U8]>.
// If the option is empty, it returns the default value T.
// If the option is not empty, it decodes it using decodeFunc and returns it.
func (bs baseStorage[T]) getDecode(key []byte, decodeFunc func(buffer *bytes.Buffer) T) T {
	option := bs.storage.Get(key)

	if !option.HasValue {
		return *new(T)
	}

	buffer := &bytes.Buffer{}
	buffer.Write(sc.SequenceU8ToBytes(option.Value))

	return decodeFunc(buffer)
}

func (bs baseStorage[T]) getDecodeOnEmpty(key []byte, decodeFunc func(buffer *bytes.Buffer) T, onEmpty T) T {
	option := bs.storage.Get(key)

	if !option.HasValue {
		return onEmpty
	}

	buffer := &bytes.Buffer{}
	buffer.Write(sc.SequenceU8ToBytes(option.Value))

	return decodeFunc(buffer)
}

// takeBytes gets the storage value. The result from Get is Option<sc.Sequence[sc.U8]>.
// If the option is empty, it returns nil.
// If the option is not empty, it clears it and returns the sequence as bytes.
func (bs baseStorage[T]) takeBytes(key []byte) []byte {
	option := bs.storage.Get(key)

	if !option.HasValue {
		return nil
	}

	bs.storage.Clear(key)

	return sc.SequenceU8ToBytes(option.Value)
}

// TakeDecode gets the storage value and returns it decoded. The result from Get is Option<sc.Sequence[sc.U8]>.
// If the option is empty, it returns default value T.
// If the option is not empty, it clears it and returns decodeFunc(value).
func (bs baseStorage[T]) takeDecode(key []byte, decodeFunc func(buffer *bytes.Buffer) T) T {
	option := bs.storage.Get(key)

	if !option.HasValue {
		return *new(T)
	}

	bs.storage.Clear(key)

	buffer := &bytes.Buffer{}
	buffer.Write(sc.SequenceU8ToBytes(option.Value))

	return decodeFunc(buffer)
}
