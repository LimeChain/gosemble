package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewPaysYes(t *testing.T) {
	assert.Equal(t, sc.NewVaryingData(PaysYes), NewPaysYes())
}

func Test_NewPaysNo(t *testing.T) {
	assert.Equal(t, sc.NewVaryingData(PaysNo), NewPaysNo())
}

func Test_DecodePays_Yes(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)

	result := DecodePays(buffer)

	assert.Equal(t, NewPaysYes(), result)
}

func Test_DecodePays_No(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(1)

	result := DecodePays(buffer)

	assert.Equal(t, NewPaysNo(), result)
}

func Test_DecodePays_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(2)

	assert.PanicsWithValue(t, "invalid Pays type", func() {
		DecodePays(buffer)
	})
}
