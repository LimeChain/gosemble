package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/utils"
)

type Callable struct {
	ModuleId   sc.U8
	FunctionId sc.U8
	Arguments  sc.VaryingData
}

func (c Callable) Encode(buffer *bytes.Buffer) error {
	return utils.EncodeEach(buffer,
		c.ModuleId,
		c.FunctionId,
		c.Arguments,
	)
}

func (c Callable) Bytes() []byte {
	return sc.EncodedBytes(c)
}

func (c Callable) ModuleIndex() sc.U8 {
	return c.ModuleId
}

func (c Callable) FunctionIndex() sc.U8 {
	return c.FunctionId
}

func (c Callable) Args() sc.VaryingData {
	return c.Arguments
}
