package types

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	key0 = [8]byte{
		't', 'e', 's', 't', 'i', 'n', 'h', '0',
	}
	key1 = [8]byte{
		't', 'e', 's', 't', 'i', 'n', 'h', '1',
	}

	value0 = sc.Sequence[sc.I32]{1, 2, 3}
	value1 = sc.U32(7)

	expectEncoded = common.MustHexToBytes("0x0874657374696e6830340c01000000020000000300000074657374696e68311007000000")
)

func Test_InherentData_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}
	inherent := NewInherentData()
	assert.Nil(t, inherent.Put(key0, value0))
	assert.Nil(t, inherent.Put(key1, value1))

	inherent.Encode(buffer)

	assert.Equal(t, expectEncoded, buffer.Bytes())
}

func Test_InherentData_Bytes(t *testing.T) {
	inherent := NewInherentData()
	assert.Nil(t, inherent.Put(key0, value0))
	assert.Nil(t, inherent.Put(key1, value1))

	encoded := inherent.Bytes()

	assert.Equal(t, expectEncoded, encoded)
}

func Test_InherentData_Decode(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.Write(expectEncoded)

	result, err := DecodeInherentData(buffer)
	assert.Nil(t, err)

	buffer.Reset()
	buffer.Write(sc.SequenceU8ToBytes(result.data[key0]))

	decodedValue0, err := sc.DecodeSequence[sc.I32](buffer)
	assert.NoError(t, err)
	assert.Equal(t, value0, decodedValue0)

	buffer.Reset()
	buffer.Write(sc.SequenceU8ToBytes(result.data[key1]))

	decodedValue1, err := sc.DecodeU32(buffer)
	assert.NoError(t, err)
	assert.Equal(t, value1, decodedValue1)
}

func Test_InherentData_Decode_Empty(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.Write(sc.ToCompact(sc.U32(1)).Bytes())

	result, err := DecodeInherentData(buffer)
	assert.Nil(t, result)
	assert.Equal(t, io.EOF, err)
}

func Test_InherentData_Decode_InvalidLength(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{1, 2, 3, 4, 5})

	result, err := DecodeInherentData(buffer)
	assert.Nil(t, result)
	assert.Equal(t, errors.New("invalid length"), err)
}

func Test_InherentData_Put_Error(t *testing.T) {
	inherent := NewInherentData()
	assert.Nil(t, inherent.Put(key0, value0))

	result := inherent.Put(key0, value1)

	assert.Equal(t, NewInherentErrorInherentDataExists(sc.BytesToSequenceU8(key0[:])), result)
}

func Test_InherentData_Get(t *testing.T) {
	inherent := NewInherentData()
	assert.Nil(t, inherent.Put(key0, value0))
	assert.Nil(t, inherent.Put(key1, value1))

	assert.Equal(t, sc.BytesToSequenceU8(value0.Bytes()), inherent.Get(key0))
	assert.Equal(t, sc.BytesToSequenceU8(value1.Bytes()), inherent.Get(key1))
}

func Test_InherentData_Clear(t *testing.T) {
	inherent := NewInherentData()
	assert.Nil(t, inherent.Put(key0, value0))

	assert.Equal(t, 1, len(inherent.data))

	inherent.Clear()

	assert.Equal(t, 0, len(inherent.data))
	assert.Equal(t, map[[8]uint8]sc.Sequence[sc.U8]{}, inherent.data)
	assert.Nil(t, nil, inherent.Get(key0))
}
