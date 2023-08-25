package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

type Perbill struct {
	Percentage sc.U32
}

func (p Perbill) Encode(buffer *bytes.Buffer) {
	p.Percentage.Encode(buffer)
}

func DecodePerbill(buffer *bytes.Buffer) Perbill {
	p := Perbill{}
	p.Percentage = sc.DecodeU32(buffer)
	return p
}

func (p Perbill) Bytes() []byte {
	return sc.EncodedBytes(p)
}

func (p Perbill) Mul(v sc.Encodable) sc.Encodable {
	switch v := v.(type) {
	case sc.U32:
		return ((v.Div(sc.U32(100))).Mul(p.Percentage))
	case Weight:
		return Weight{
			RefTime:   (v.RefTime.Div(sc.U64(100))).Mul(sc.U64(p.Percentage)).(sc.U64),
			ProofSize: (v.ProofSize.Div(sc.U64(100))).Mul(sc.U64(p.Percentage)).(sc.U64),
		}
	default:
		log.Critical("unsupported type")
	}

	panic("unreachable")
}
