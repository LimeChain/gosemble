package types

import (
	"bytes"
	"fmt"

	module3 "github.com/LimeChain/gosemble/frame/balances/module"

	module2 "github.com/LimeChain/gosemble/frame/system/module"

	"github.com/LimeChain/gosemble/frame/timestamp/module"

	"github.com/LimeChain/gosemble/primitives/types"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/support"
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
	switch c.CallIndex.ModuleIndex {
	case module2.Module.Index():
		switch c.CallIndex.FunctionIndex {
		case module2.Module.Remark.Index():
			c.Function = module2.Module.Remark
		default:
			log.Trace(fmt.Sprintf("function index [%d] not found", c.CallIndex.FunctionIndex))
		}
	case module.Module.Index():
		switch c.CallIndex.FunctionIndex {
		case module.Module.Set.Index():
			c.Function = module.Module.Set
			c.Args = []sc.Encodable{
				sc.DecodeU64(buffer),
			}
		default:
			log.Trace(fmt.Sprintf("function index [%d] not found", c.CallIndex.FunctionIndex))
		}
	case module3.Module.Index():
		switch c.CallIndex.FunctionIndex {
		case module3.Module.Transfer.Index():
			c.Function = module3.Module.Transfer
			c.Args = []sc.Encodable{
				types.DecodeMultiAddress(buffer),
				sc.U128(sc.DecodeCompact(buffer)),
			}
		case module3.Module.SetBalance.Index():
			c.Function = module3.Module.SetBalance
			c.Args = []sc.Encodable{
				types.DecodeMultiAddress(buffer),
				sc.U128(sc.DecodeCompact(buffer)),
				sc.U128(sc.DecodeCompact(buffer)),
			}
		case module3.Module.ForceTransfer.Index():
			c.Function = module3.Module.ForceTransfer
			c.Args = []sc.Encodable{
				types.DecodeMultiAddress(buffer),
				types.DecodeMultiAddress(buffer),
				sc.U128(sc.DecodeCompact(buffer)),
			}
		case module3.Module.TransferKeepAlive.Index():
			c.Function = module3.Module.TransferKeepAlive
			c.Args = []sc.Encodable{
				types.DecodeMultiAddress(buffer),
				sc.U128(sc.DecodeCompact(buffer)),
			}
		case module3.Module.TransferAll.Index():
			c.Function = module3.Module.TransferAll
			c.Args = []sc.Encodable{
				types.DecodeMultiAddress(buffer),
				sc.DecodeBool(buffer),
			}
		case module3.Module.ForceFree.Index():
			c.Function = module3.Module.ForceFree
			c.Args = []sc.Encodable{
				types.DecodeMultiAddress(buffer),
				sc.U128(sc.DecodeCompact(buffer)),
			}
		}
	default:
		log.Trace(fmt.Sprintf("module with index [%d] not found", c.CallIndex.ModuleIndex))
	}
	return c
}

func (c Call) Bytes() []byte {
	return sc.EncodedBytes(c)
}
