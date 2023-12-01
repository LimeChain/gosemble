package types

import (
	"bytes"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewDispatchClassNormal(t *testing.T) {
	assert.Equal(t,
		DispatchClass{sc.NewVaryingData(DispatchClassNormal)},
		NewDispatchClassNormal(),
	)
}

func Test_NewDispatchClassOperational(t *testing.T) {
	assert.Equal(t,
		DispatchClass{sc.NewVaryingData(DispatchClassOperational)},
		NewDispatchClassOperational(),
	)
}

func Test_NewDispatchClassMandatory(t *testing.T) {
	assert.Equal(t,
		DispatchClass{sc.NewVaryingData(DispatchClassMandatory)},
		NewDispatchClassMandatory(),
	)
}

func Test_DecodeDispatchClass_Normal(t *testing.T) {
	targetDispatchClass := NewDispatchClassNormal()

	buf := &bytes.Buffer{}
	err := targetDispatchClass.Encode(buf)
	assert.NoError(t, err)

	dispatchClass, err := DecodeDispatchClass(buf)
	assert.NoError(t, err)
	assert.Equal(t, targetDispatchClass, dispatchClass)
}

func Test_DecodeDispatchClass_Operational(t *testing.T) {
	targetDispatchClass := NewDispatchClassOperational()

	buf := &bytes.Buffer{}
	err := targetDispatchClass.Encode(buf)
	assert.NoError(t, err)

	dispatchClass, err := DecodeDispatchClass(buf)
	assert.NoError(t, err)
	assert.Equal(t, targetDispatchClass, dispatchClass)
}

func Test_DecodeDispatchClass_Mandatory(t *testing.T) {
	targetDispatchClass := NewDispatchClassMandatory()

	buf := &bytes.Buffer{}
	err := targetDispatchClass.Encode(buf)
	assert.NoError(t, err)

	dispatchClass, err := DecodeDispatchClass(buf)
	assert.NoError(t, err)
	assert.Equal(t, targetDispatchClass, dispatchClass)
}

func Test_DecodeDispatchClass_TypeError(t *testing.T) {
	buf := &bytes.Buffer{}
	targetDispatchClass := DispatchClass{sc.NewVaryingData(sc.U8(3))}
	err := targetDispatchClass.Encode(buf)
	assert.NoError(t, err)

	result, err := DecodeDispatchClass(buf)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'DispatchClass' type", err.Error())
	assert.Equal(t, DispatchClass{}, result)
}

func Test_DecodeDispatchClass_Empty(t *testing.T) {
	buffer := &bytes.Buffer{}

	dispatchClass, err := DecodeDispatchClass(buffer)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, DispatchClass{}, dispatchClass)
}

func Test_DispatchClass_Is_Empty_TypeError(t *testing.T) {
	result, err := DispatchClass{}.Is(DispatchClassNormal)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'DispatchClass' type", err.Error())
	assert.Equal(t, sc.Bool(false), result)
}

func Test_DispatchClass_Is(t *testing.T) {
	result, err := NewDispatchClassNormal().Is(DispatchClassNormal)

	assert.NoError(t, err)
	assert.Equal(t, sc.Bool(true), result)
}

func Test_DispatchClass_Is_TypeError(t *testing.T) {
	unknownDispatchClass := DispatchClassType(3)

	result, err := NewDispatchClassNormal().Is(unknownDispatchClass)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'DispatchClass' type", err.Error())
	assert.Equal(t, sc.Bool(false), result)
}

func Test_DispatchClassAll(t *testing.T) {
	assert.Equal(t,
		[]DispatchClass{NewDispatchClassNormal(), NewDispatchClassOperational(), NewDispatchClassMandatory()},
		DispatchClassAll(),
	)
}

func Test_DispatchClassType_Bytes(t *testing.T) {
	result := DispatchClassNormal.Bytes()

	assert.Equal(t, []byte{0}, result)
}
