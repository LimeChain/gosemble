package scale

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"
)

// Implementation for Parity codec in Tinygo.
const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// Encoder is a wrapper around a Writer that allows encoding data items to a stream.
type Encoder struct {
	Writer io.Writer
}

// Write several bytes to the encoder.
func (pe Encoder) Write(bytes []byte) {
	c, err := pe.Writer.Write(bytes)
	check(err)
	if c < len(bytes) {
		panic("Could not write " + strconv.Itoa(len(bytes)) + " bytes to writer")
	}
}

// EncodeByte writes a single byte to an encoder.
func (pe Encoder) EncodeByte(b byte) {
	intBuf[0] = b
	pe.Write(intBuf[:1])
}

// Reusable buffer to encode/decode integers. Should be safe, since WASM is
// single-threaded by default.
var intBuf [8]byte

// EncodeUintCompact writes an unsigned integer to the stream using the compact encoding.
// A typical usage is storing the length of a collection.
// Definition of compact encoding:
// 0b00 00 00 00 / 00 00 00 00 / 00 00 00 00 / 00 00 00 00
//
//	xx xx xx 00															(0 ... 2**6 - 1)		(u8)
//	yL yL yL 01 / yH yH yH yL												(2**6 ... 2**14 - 1)	(u8, u16)  low LH high
//	zL zL zL 10 / zM zM zM zL / zM zM zM zM / zH zH zH zM					(2**14 ... 2**30 - 1)	(u16, u32)  low LMMH high
//	nn nn nn 11 [ / zz zz zz zz ]{4 + n}									(2**30 ... 2**536 - 1)	(u32, u64, u128, U256, U512, U520) straight LE-encoded
//
// Rust implementation: see impl<'a> Encode for CompactRef<'a, u64>
func (pe Encoder) EncodeUintCompact(v uint64) {

	// TODO: handle numbers wide than 64 bits (byte slices?)
	// Currently, Rust implementation only seems to support u128

	if v < 1<<30 {
		if v < 1<<6 {
			pe.EncodeByte(byte(v) << 2)
		} else if v < 1<<14 {
			buf := intBuf[:2]
			binary.LittleEndian.PutUint16(buf, uint16(v<<2)+1)
			pe.Write(buf)
		} else {
			buf := intBuf[:4]
			binary.LittleEndian.PutUint32(buf, uint32(v<<2)+2)
			pe.Write(buf)
		}
		return
	}

	n := byte(0)
	limit := uint64(1 << 32)
	for v >= limit && limit > 256 { // when overflows, limit will be < 256
		n++
		limit <<= 8
	}
	if n > 4 {
		panic("Assertion error: n>4 needed to compact-encode uint64")
	}
	pe.EncodeByte((n << 2) + 3)
	binary.LittleEndian.PutUint64(intBuf[:8], v)
	pe.Write(intBuf[:4+n])
}

// EncodeUint64 writes a uint64 to the encoder
func (pe Encoder) EncodeUint64(value uint64) {
	buf := intBuf[:8]
	binary.LittleEndian.PutUint64(buf, value)
	pe.Write(buf)
}

// EncodeUint32 writes a uint32 to the encoder
func (pe Encoder) EncodeUint32(value uint32) {
	buf := intBuf[:4]
	binary.LittleEndian.PutUint32(buf, value)
	pe.Write(buf)
}

// EncodeUint16 writes a uint16 to the encoder
func (pe Encoder) EncodeUint16(value uint16) {
	buf := intBuf[:2]
	binary.LittleEndian.PutUint16(buf, value)
	pe.Write(buf)
}

// EncodeInt64 writes a int64 to the encoder
func (pe Encoder) EncodeInt64(value int64) {
	pe.EncodeUint64(uint64(value))
}

// EncodeInt32 writes a int32 to the encoder
func (pe Encoder) EncodeInt32(value int32) {
	pe.EncodeUint32(uint32(value))
}

// EncodeInt16 writes a int16 to the encoder
func (pe Encoder) EncodeInt16(value int16) {
	pe.EncodeUint16(uint16(value))
}

// EncodeInt8 writes a int8 to the encoder
func (pe Encoder) EncodeInt8(value int8) {
	pe.EncodeByte(byte(value))
}

// EncodeBool writes a boolean to the encoder
func (pe Encoder) EncodeBool(b bool) {
	if b {
		pe.EncodeByte(1)
	} else {
		pe.EncodeByte(0)
	}
}

// EncodeByteSlice writes a slice of bytes to the encoder
func (pe Encoder) EncodeByteSlice(value []byte) {
	pe.EncodeUintCompact(uint64(len(value)))
	pe.Write(value)
}

// EncodeString writes a UTF-8 string to the encoder
func (pe Encoder) EncodeString(value string) {
	pe.EncodeByteSlice([]byte(value))
}

// EncodeOption stores optionally present value to the stream.
func (pe Encoder) EncodeOption(hasValue bool, value Encodeable) {
	if !hasValue {
		pe.EncodeByte(0)
	} else {
		pe.EncodeByte(1)
		value.ParityEncode(pe)
	}
}

// Decoder - a wraper around a Reader that allows decoding data items from a stream.
// Unlike Rust implementations, decoder methods do not return success state, but just
// panic on error. Since decoding failue is an "unexpected" error, this approach should
// be justified.
type Decoder struct {
	Reader io.Reader
}

// Read reads bytes from a stream into a buffer and panics if cannot read the required
// number of bytes.
func (pd Decoder) Read(bytes []byte) {
	c, err := pd.Reader.Read(bytes)
	check(err)
	if c < len(bytes) {
		panic("Cannot read the required number of bytes " + strconv.Itoa(len(bytes)) + ", only " + strconv.Itoa(c) + " available")
	}
}

// DecodeByte reads a next byte from the stream.
func (pd Decoder) DecodeByte() byte {
	pd.Read(intBuf[:1])
	return intBuf[0]
}

// DecodeUint64 reads a uint64 from the stream.
func (pd Decoder) DecodeUint64() uint64 {
	buf := intBuf[:8]
	pd.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

// DecodeUint32 reads a uint32 from the stream.
func (pd Decoder) DecodeUint32() uint32 {
	buf := intBuf[:4]
	pd.Read(buf)
	return binary.LittleEndian.Uint32(buf)
}

// DecodeUint16 reads a uint16 from the stream.
func (pd Decoder) DecodeUint16() uint16 {
	buf := intBuf[:2]
	pd.Read(buf)
	return binary.LittleEndian.Uint16(buf)
}

// DecodeInt64 reads a int64 from the stream.
func (pd Decoder) DecodeInt64() int64 {
	return int64(pd.DecodeUint64())
}

// DecodeInt32 reads a int32 from the stream.
func (pd Decoder) DecodeInt32() int32 {
	return int32(pd.DecodeUint32())
}

// DecodeInt16 reads a int16 from the stream.
func (pd Decoder) DecodeInt16() int16 {
	return int16(pd.DecodeUint16())
}

// DecodeInt8 reads a int8 from the stream.
func (pd Decoder) DecodeInt8() int8 {
	return int8(pd.DecodeByte())
}

// DecodeBool reads a bool from the stream.
func (pd Decoder) DecodeBool() bool {
	return pd.DecodeByte() > 0
}

// DecodeUintCompact decodes a compact-encoded integer. See EncodeUintCompact method.
func (pd Decoder) DecodeUintCompact() uint64 {
	b := pd.DecodeByte()
	mode := b & 3
	switch mode {
	case 0:
		return uint64(b >> 2)
	case 1:
		r := uint64(pd.DecodeByte())
		r <<= 6
		r += uint64(b >> 2)
		return r
	case 2:
		buf := intBuf[:4]
		buf[0] = b
		pd.Read(intBuf[1:4])
		r := binary.LittleEndian.Uint32(buf)
		r >>= 2
		return uint64(r)
	case 3:
		n := b >> 2
		if n > 4 {
			panic("Not supported: n>4 encountered when decoding a compact-encoded uint")
		}
		pd.Read(intBuf[:n+4])
		for i := n + 4; i < 8; i++ {
			intBuf[i] = 0
		}
		return binary.LittleEndian.Uint64(intBuf[:8])
	default:
		panic("Code should be unreachable")
	}
}

// DecodeByteSlice reads a slice of bytes from the stream
func (pd Decoder) DecodeByteSlice() []byte {
	value := make([]byte, uintptr(pd.DecodeUintCompact()))
	pd.Read(value)
	return value
}

// DecodeString reads a UTF-8 string from the stream
func (pd Decoder) DecodeString() string {
	return string(pd.DecodeByteSlice())
}

// DecodeOption decodes a optionally available value into a boolean presence field and a value.
func (pd Decoder) DecodeOption(hasValue *bool, valuePointer Decodeable) {
	b := pd.DecodeByte()
	switch b {
	case 0:
		*hasValue = false
	case 1:
		*hasValue = true
		valuePointer.ParityDecode(pd)
	default:
		panic("Unknown byte prefix for encoded Option: " + strconv.Itoa(int(b)))
	}
}

// Encodeable is an interface that defines a custom encoding rules for a data type.
// Should be defined for structs (not pointers to them).
// See OptionBool for an example implementation.
type Encodeable interface {
	// ParityEncode encodes and write this structure into a stream
	ParityEncode(encoder Encoder)
}

// EncodeCollection encodes a collection using a per-element callback.
// See []int16 in tests as an example
func (pe Encoder) EncodeCollection(length int, encodeElem func(int)) {
	pe.EncodeUintCompact(uint64(length))
	for i := 0; i < length; i++ {
		encodeElem(i)
	}
}

// Decodeable is an interface that defines a custom encoding rules for a data type.
// Should be defined for pointers to structs.
// See OptionBool for an example implementation.
type Decodeable interface {
	// ParityDecode populates this structure from a stream (overwriting the current contents), return false on failure
	ParityDecode(decoder Decoder)
}

// DecodeCollection encodes a collection using a maker callback and a per-element callback.
// See []int16 in tests as an example
func (pd Decoder) DecodeCollection(setSize func(int), decodeElem func(int)) {
	l := int(pd.DecodeUintCompact())
	setSize(l)
	for i := 0; i < l; i++ {
		decodeElem(i)
	}
}

// OptionBool is a structure that can store a boolean or a missing value.
// Note that encoding rules are slightly different from other "Option" fields.
type OptionBool struct {
	hasValue bool
	value    bool
}

// NewOptionBoolEmpty creates an OptionBool without a value.
func NewOptionBoolEmpty() OptionBool {
	return OptionBool{false, false}
}

// NewOptionBool creates an OptionBool with a value.
func NewOptionBool(value bool) OptionBool {
	return OptionBool{true, value}
}

// ParityEncode implements encoding for OptionBool as per Rust implementation.
func (o OptionBool) ParityEncode(encoder Encoder) {
	if !o.hasValue {
		encoder.EncodeByte(0)
	} else {
		if o.value {
			encoder.EncodeByte(1)
		} else {
			encoder.EncodeByte(2)
		}
	}
}

// ParityDecode implements decoding for OptionBool as per Rust implementation.
func (o *OptionBool) ParityDecode(decoder Decoder) {
	b := decoder.DecodeByte()
	switch b {
	case 0:
		o.hasValue = false
		o.value = false
	case 1:
		o.hasValue = true
		o.value = true
	case 2:
		o.hasValue = true
		o.value = false
	default:
		panic("Unknown byte prefix for encoded OptionBool: " + strconv.Itoa(int(b)))
	}
}

// ToBytes is a helper method to encode an encodeable value as a byte slice
func ToBytes(value Encodeable) []byte {
	var buffer = bytes.Buffer{}
	value.ParityEncode(Encoder{&buffer})
	return buffer.Bytes()
}

// ToBytesCustom is a helper method to run a custom encoding sequence and return result as a byte slice
func ToBytesCustom(encode func(Encoder)) []byte {
	var buffer = bytes.Buffer{}
	encode(Encoder{&buffer})
	return buffer.Bytes()
}

// FromBytes is a method to decode a decodeable value from a byte slice
func FromBytes(value Decodeable, encoded []byte) {
	var buffer = bytes.NewBuffer(encoded)
	value.ParityDecode(Decoder{buffer})
}
