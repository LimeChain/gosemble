package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewTimestampErrorTooEarly(t *testing.T) {
	expect := TimestampError{sc.NewVaryingData(TimestampErrorTooEarly)}

	assert.Equal(t, expect, NewTimestampErrorTooEarly())
}

func Test_NewTimestampErrorTooFarInFuture(t *testing.T) {
	expect := TimestampError{sc.NewVaryingData(TimestampErrorTooFarInFuture)}
	assert.Equal(t, expect, NewTimestampErrorTooFarInFuture())
}

func Test_NewTimestampErrorInvalid(t *testing.T) {
	expect := TimestampError{sc.NewVaryingData(TimestampErrorInvalid)}

	assert.Equal(t, expect, NewTimestampErrorInvalid())
}

func Test_TimestampError_IsFatal_True(t *testing.T) {
	assert.Equal(t, sc.Bool(true), NewTimestampErrorTooEarly().IsFatal())
	assert.Equal(t, sc.Bool(true), NewTimestampErrorTooFarInFuture().IsFatal())
	assert.Equal(t, sc.Bool(true), NewTimestampErrorInvalid().IsFatal())
}

func Test_TimestampError_IsFatal_False(t *testing.T) {
	unknownTimestampError := TimestampError{sc.NewVaryingData(sc.U8(3))}
	assert.Equal(t, sc.Bool(false), unknownTimestampError.IsFatal())
}

func Test_TimestampError_Error_TooEarly(t *testing.T) {
	target := NewTimestampErrorTooEarly()
	assert.Equal(t, "The time since the last timestamp is lower than the minimum period.", target.Error())
}

func Test_TimestampError_Error_TooFarInFuture(t *testing.T) {
	target := NewTimestampErrorTooFarInFuture()
	assert.Equal(t, "The timestamp of the block is too far in the future.", target.Error())
}

func Test_TimestampError_Error_Invalid(t *testing.T) {
	target := NewTimestampErrorInvalid()
	assert.Equal(t, "invalid inherent check for timestamp module", target.Error())
}

func Test_TimestampError_Error_Panics(t *testing.T) {
	target := TimestampError{sc.NewVaryingData(sc.U8(3))}

	assert.PanicsWithValue(t, "invalid TimestampError", func() {
		target.Error()
	})
}
