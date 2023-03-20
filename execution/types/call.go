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
	Args      sc.VaryingData
	Function  support.FunctionMetadata
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

	c.Function = function
	c.Args = function.Decode(buffer)

	return c
}

func (c Call) Bytes() []byte {
	return sc.EncodedBytes(c)
}
