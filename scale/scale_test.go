package scale

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if reflect.DeepEqual(a, b) {
		return
	}
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func assertPanics(code func()) {
	panicked := false
	assertPanicsInner(code, &panicked)
	if !panicked {
		panic("Panic was expected, but code executed successfully")
	}
}

func assertPanicsInner(code func(), panicked *bool) {
	defer func() {
		if r := recover(); r != nil {
			*panicked = true
		}
	}()
	code()
}

func hexify(bytes []byte) string {
	res := make([]string, len(bytes))
	for i, b := range bytes {
		res[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(res, " ")
}

func encodeToBytes(value Encodeable) []byte {
	var buffer = bytes.Buffer{}
	value.ParityEncode(Encoder{&buffer})
	return buffer.Bytes()
}

func verifyEncodingReturnDecoder(t *testing.T, encode func(pe Encoder), expectedHex string) Decoder {
	var buffer = bytes.Buffer{}
	encode(Encoder{&buffer})
	assertEqual(t, hexify(buffer.Bytes()), expectedHex)
	return Decoder{&buffer}
}

func TestSliceOfBytesEncodedAsExpected(t *testing.T) {
	value := []byte{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
	pd := verifyEncodingReturnDecoder(t,
		func(pe Encoder) { pe.EncodeByteSlice(value) },
		"28 00 01 01 02 03 05 08 0d 15 22",
	)
	assertEqual(t, pd.DecodeByteSlice(), value)
}

// You can either define a type for your slice:
// type Int16Slice []int16
// with methods ParityEncode and ParityDecode

func TestSliceOfInt16EncodedAsExpected(t *testing.T) {
	value := []int16{0, 1, -1, 2, -2, 3, -3}

	pd := verifyEncodingReturnDecoder(t,
		func(pe Encoder) {
			pe.EncodeCollection(
				len(value),
				func(i int) { pe.EncodeInt16(value[i]) },
			)
		},
		"1c 00 00 01 00 ff ff 02 00 fe ff 03 00 fd ff",
	)
	var v2 []int16
	pd.DecodeCollection(
		func(n int) { v2 = make([]int16, n) },
		func(i int) { v2[i] = pd.DecodeInt16() },
	)
	assertEqual(t, v2, value)
}

// OptionInt8 is an example implementation of an "Option" type, mirroring Option<u8> in Rust version.
// Since Go does not support generics, one has to define such types manually.
// See below for ParityEncode / ParityDecode implementations.
type OptionInt8 struct {
	hasValue bool
	value    int8
}

func (o OptionInt8) ParityEncode(pe Encoder) {
	pe.EncodeBool(o.hasValue)
	if o.hasValue {
		pe.EncodeInt8(o.value)
	}
}

func (o *OptionInt8) ParityDecode(pd Decoder) {
	o.hasValue = pd.DecodeBool()
	if o.hasValue {
		o.value = pd.DecodeInt8()
	}
}

func TestOptionInt8EncodedAsExpected(t *testing.T) {
	tests := map[OptionInt8]string{
		OptionInt8{true, 1}:  "01 01",
		OptionInt8{true, -1}: "01 ff",
		OptionInt8{false, 0}: "00",
	}

	for value, hex := range tests {
		pd := verifyEncodingReturnDecoder(t,
			func(pe Encoder) { value.ParityEncode(pe) },
			hex,
		)
		v2 := OptionInt8{}
		v2.ParityDecode(pd)
		assertEqual(t, v2, value)
	}
}

func TestOptionBoolEncodedAsExpected(t *testing.T) {
	tests := map[OptionBool]string{
		NewOptionBool(true):  "01",
		NewOptionBool(false): "02",
		NewOptionBoolEmpty(): "00",
	}

	for value, hex := range tests {
		pd := verifyEncodingReturnDecoder(t,
			func(pe Encoder) { value.ParityEncode(pe) },
			hex,
		)
		v2 := OptionBool{}
		v2.ParityDecode(pd)
		assertEqual(t, v2, value)
	}
}

func TestSliceOfStringEncodedAsExpected(t *testing.T) {
	value := []string{
		"Hamlet",
		"Война и мир",
		"三国演义",
		"أَلْف لَيْلَة وَلَيْلَة‎"}

	pd := verifyEncodingReturnDecoder(t,
		func(pe Encoder) {
			pe.EncodeCollection(
				len(value),
				func(i int) { pe.EncodeString(value[i]) },
			)
		},
		"10 18 48 61 6d 6c 65 74 50 d0 92 d0 be d0 b9 d0 bd d0 b0 20 d0 "+
			"b8 20 d0 bc d0 b8 d1 80 30 e4 b8 89 e5 9b bd e6 bc 94 e4 b9 89 bc d8 a3 d9 8e d9 84 d9 92 "+
			"d9 81 20 d9 84 d9 8e d9 8a d9 92 d9 84 d9 8e d8 a9 20 d9 88 d9 8e d9 84 d9 8e d9 8a d9 92 "+
			"d9 84 d9 8e d8 a9 e2 80 8e",
	)
	var v2 []string
	pd.DecodeCollection(
		func(n int) { v2 = make([]string, n) },
		func(i int) { v2[i] = pd.DecodeString() },
	)
	assertEqual(t, v2, value)
}

func TestCompactIntegersEncodedAsExpected(t *testing.T) {
	tests := map[uint64]string{
		0:              "00",
		63:             "fc",
		64:             "01 01",
		16383:          "fd ff",
		16384:          "02 00 01 00",
		1073741823:     "fe ff ff ff",
		1073741824:     "03 00 00 00 40",
		1<<32 - 1:      "03 ff ff ff ff",
		1 << 32:        "07 00 00 00 00 01",
		1 << 40:        "0b 00 00 00 00 00 01",
		1 << 48:        "0f 00 00 00 00 00 00 01",
		1<<56 - 1:      "0f ff ff ff ff ff ff ff",
		1 << 56:        "13 00 00 00 00 00 00 00 01",
		math.MaxUint64: "13 ff ff ff ff ff ff ff ff"}

	for value, hex := range tests {
		pd := verifyEncodingReturnDecoder(t,
			func(pe Encoder) { pe.EncodeUintCompact(value) },
			hex,
		)
		v2 := pd.DecodeUintCompact()
		assertEqual(t, v2, value)
	}
}
