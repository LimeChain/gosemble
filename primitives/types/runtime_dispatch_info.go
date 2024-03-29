package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type RuntimeDispatchInfo struct {
	Weight     Weight
	Class      DispatchClass
	PartialFee Balance
}

func (rdi RuntimeDispatchInfo) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		rdi.Weight,
		rdi.Class,
		rdi.PartialFee,
	)
}

func DecodeRuntimeDispatchInfo(buffer *bytes.Buffer) (RuntimeDispatchInfo, error) {
	rdi := RuntimeDispatchInfo{}
	weight, err := DecodeWeight(buffer)
	if err != nil {
		return RuntimeDispatchInfo{}, err
	}
	class, err := DecodeDispatchClass(buffer)
	if err != nil {
		return RuntimeDispatchInfo{}, err
	}
	partialFee, err := sc.DecodeU128(buffer)
	if err != nil {
		return RuntimeDispatchInfo{}, err
	}
	rdi.Weight = weight
	rdi.Class = class
	rdi.PartialFee = partialFee
	return rdi, nil
}

func (rdi RuntimeDispatchInfo) Bytes() []byte {
	return sc.EncodedBytes(rdi)
}
