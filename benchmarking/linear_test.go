package benchmarking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinearComponent(t *testing.T) {
	l, err := NewLinear(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, uint32(10), l.Value())

	values, err := l.values(5)
	assert.NoError(t, err)
	assert.Equal(t, []uint32{0, 2, 5, 7, 10}, values)

	l.setValue(99)
	assert.Equal(t, uint32(99), l.Value())

	componentValues := componentValues([]*linear{l})
	assert.Equal(t, []uint32{99}, componentValues)

	l, err = NewLinear(1, 0)
	assert.Equal(t, errInvalidValues, err)
}
