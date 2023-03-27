package types

import (
	"bytes"
	"fmt"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/support"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Call struct {
	CallIndex types.CallIndex
	function  support.FunctionMetadata
	Args      sc.VaryingData
}

func NewCall(index types.CallIndex, args sc.VaryingData, function support.FunctionMetadata) Call {
	return Call{
		CallIndex: index,
		Args:      args,
		function:  function,
	}
}

func (c Call) Encode(buffer *bytes.Buffer) {
	c.CallIndex.Encode(buffer)
	c.Args.Encode(buffer)
}

func DecodeCall(buffer *bytes.Buffer) Call {
	c := Call{}
	c.CallIndex = types.DecodeCallIndex(buffer)

	module, ok := Modules[c.CallIndex.ModuleIndex]
	if !ok {
		log.Critical(fmt.Sprintf("module with index [%d] not found", c.CallIndex.ModuleIndex))
	}

	function, ok := module.Functions()[c.CallIndex.FunctionIndex]
	if !ok {
		log.Critical(fmt.Sprintf("function index [%d] for module [%d] not found", c.CallIndex.FunctionIndex, c.CallIndex.ModuleIndex))
	}

	c.function = function
	c.Args = function.Decode(buffer)

	return c
}

func (c Call) Bytes() []byte {
	return sc.EncodedBytes(c)
}

func (c Call) Function() support.FunctionMetadata {
	return c.function
}
