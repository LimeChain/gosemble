package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/utils"
)

type ApiItem struct {
	Name    sc.FixedSequence[sc.U8] // size 8
	Version sc.U32
}

func NewApiItem(name [8]byte, version sc.U32) ApiItem {
	return ApiItem{
		Name:    sc.BytesToFixedSequenceU8(name[:]),
		Version: version,
	}
}

func (ai ApiItem) Encode(buffer *bytes.Buffer) error {
	return utils.EncodeEach(buffer,
		ai.Name,
		ai.Version,
	)
}

func DecodeApiItem(buffer *bytes.Buffer) (ApiItem, error) {
	name, err := sc.DecodeFixedSequence[sc.U8](8, buffer)
	if err != nil {
		return ApiItem{}, err
	}
	version, err := sc.DecodeU32(buffer)
	if err != nil {
		return ApiItem{}, err
	}
	return ApiItem{
		Name:    name,
		Version: version,
	}, nil
}

func (ai ApiItem) Bytes() []byte {
	return sc.EncodedBytes(ai)
}
