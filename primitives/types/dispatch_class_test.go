package types

import (
	"bytes"
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
	targetDispatchClass.Encode(buf)

	assert.Equal(t, targetDispatchClass, DecodeDispatchClass(buf))
}

func Test_DecodeDispatchClass_Operational(t *testing.T) {
	targetDispatchClass := NewDispatchClassOperational()

	buf := &bytes.Buffer{}
	targetDispatchClass.Encode(buf)

	assert.Equal(t, targetDispatchClass, DecodeDispatchClass(buf))
}

func Test_DecodeDispatchClass_Mandatory(t *testing.T) {
	targetDispatchClass := NewDispatchClassMandatory()

	buf := &bytes.Buffer{}
	targetDispatchClass.Encode(buf)

	assert.Equal(t, targetDispatchClass, DecodeDispatchClass(buf))
}

func Test_DecodeDispatchClass_Panic(t *testing.T) {
	targetDispatchClass := DispatchClass{sc.NewVaryingData(sc.U8(3))}

	buf := &bytes.Buffer{}
	targetDispatchClass.Encode(buf)

	assert.PanicsWithValue(t, "invalid DispatchClass type", func() {
		DecodeDispatchClass(buf)
	})
}

func Test_Is(t *testing.T) {
	assert.Equal(t,
		sc.Bool(true),
		NewDispatchClassNormal().Is(DispatchClassNormal),
	)
}

func Test_Is_Panic(t *testing.T) {
	unknownDispatchClass := sc.U8(3)

	assert.PanicsWithValue(t, "invalid DispatchClass value", func() {
		NewDispatchClassNormal().Is(unknownDispatchClass)
	})
}

func Test_DispatchClassAll(t *testing.T) {
	assert.Equal(t,
		[]DispatchClass{NewDispatchClassNormal(), NewDispatchClassOperational(), NewDispatchClassMandatory()},
		DispatchClassAll(),
	)
}
