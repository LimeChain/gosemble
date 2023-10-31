package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedPerDispatchClassBytes, _ = hex.DecodeString("010203")

	targetPerDispatchClass = PerDispatchClass[sc.U8]{
		Normal:      sc.U8(1),
		Operational: sc.U8(2),
		Mandatory:   sc.U8(3),
	}
)

func Test_PerDispatchClass_Encode(t *testing.T) {
	buf := &bytes.Buffer{}

	targetPerDispatchClass.Encode(buf)

	assert.Equal(t, expectedPerDispatchClassBytes, buf.Bytes())
}

func Test_DecodePerDispatchClass(t *testing.T) {
	buf := bytes.NewBuffer(expectedPerDispatchClassBytes)

	result, err := DecodePerDispatchClass(buf, sc.DecodeU8)
	assert.NoError(t, err)

	assert.Equal(t, targetPerDispatchClass, result)
}

func Test_PerDispatchClass_Bytes(t *testing.T) {
	assert.Equal(t, expectedPerDispatchClassBytes, targetPerDispatchClass.Bytes())
}

func Test_PerDispatchClass_Get(t *testing.T) {
	assert.Equal(t, sc.U8(1), *targetPerDispatchClass.Get(NewDispatchClassNormal()))
	assert.Equal(t, sc.U8(2), *targetPerDispatchClass.Get(NewDispatchClassOperational()))
	assert.Equal(t, sc.U8(3), *targetPerDispatchClass.Get(NewDispatchClassMandatory()))
}

func Test_PerDispatchClass_Get_Panic(t *testing.T) {
	unknownDispatchClass := DispatchClass{sc.NewVaryingData(sc.U8(3))}

	assert.PanicsWithValue(t, "invalid DispatchClass type", func() {
		targetPerDispatchClass.Get(unknownDispatchClass)
	})
}
