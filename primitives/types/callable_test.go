package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedCallableBytes, _ = hex.DecodeString("01020304")
)
var (
	targetCallable = Callable{
		ModuleId:   sc.U8(1),
		FunctionId: sc.U8(2),
		Arguments:  sc.NewVaryingData(sc.U8(3), sc.U8(4)),
	}
)

func Test_Callable_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetCallable.Encode(buffer)

	assert.Equal(t, expectedCallableBytes, buffer.Bytes())
}

func Test_Callable_Bytes(t *testing.T) {
	assert.Equal(t, expectedCallableBytes, targetCallable.Bytes())
}

func Test_Callable_ModuleIndex(t *testing.T) {
	assert.Equal(t, sc.U8(1), targetCallable.ModuleIndex())
}

func Test_Callable_FunctionIndex(t *testing.T) {
	assert.Equal(t, sc.U8(2), targetCallable.FunctionIndex())
}

func Test_Callable_Args(t *testing.T) {
	assert.Equal(t, sc.NewVaryingData(sc.U8(3), sc.U8(4)), targetCallable.Args())
}
