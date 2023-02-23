package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// Weight information that is only available post dispatch.
// NOTE: This can only be used to reduce the weight or fee, not increase it.
type PostDispatchInfo struct {
	// Actual weight consumed by a call or `None` which stands for the worst case static weight.
	ActualWeight sc.Option[Weight]

	// Whether this transaction should pay fees when all is said and done.
	PaysFee sc.U8
}

func (pdi PostDispatchInfo) Encode(buffer *bytes.Buffer) {
	pdi.ActualWeight.Encode(buffer)
	pdi.PaysFee.Encode(buffer)
}

func DecodePostDispatchInfo(buffer *bytes.Buffer) PostDispatchInfo {
	pdi := PostDispatchInfo{}
	pdi.ActualWeight = sc.DecodeOptionWith(buffer, DecodeWeight)
	pdi.PaysFee = sc.DecodeU8(buffer)
	return pdi
}

func (pdi PostDispatchInfo) Bytes() []byte {
	return sc.EncodedBytes(pdi)
}
