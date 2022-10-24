package types

import (
	"bytes"

	"github.com/LimeChain/gosemble/scale"
)

type ApiItem struct {
	Name    [8]byte
	Version uint32
}

type VersionData struct {
	SpecName           []byte
	ImplName           []byte
	AuthoringVersion   uint32
	SpecVersion        uint32
	ImplVersion        uint32
	Apis               []ApiItem
	TransactionVersion uint32
	StateVersion       uint32
}

func (v *VersionData) Encode() ([]byte, error) {
	var buffer = bytes.Buffer{}
	var encoder = scale.Encoder{Writer: &buffer}

	encoder.EncodeByteSlice(v.SpecName)
	encoder.EncodeByteSlice(v.ImplName)
	encoder.EncodeUint32(v.AuthoringVersion)
	encoder.EncodeUint32(v.SpecVersion)
	encoder.EncodeUint32(v.ImplVersion)
	encoder.EncodeUint32(uint32(len(v.Apis)))

	for _, apiItem := range v.Apis {
		encoder.EncodeByteSlice(apiItem.Name[:])
		encoder.EncodeUint32(apiItem.Version)
	}

	encoder.EncodeUint32(v.TransactionVersion)
	encoder.EncodeUint32(v.StateVersion)

	return buffer.Bytes(), nil
}

func (v *VersionData) Decode(enc []byte) error {
	var buffer = bytes.NewBuffer(enc)
	var decoder = scale.Decoder{Reader: buffer}

	v.SpecName = decoder.DecodeByteSlice()
	v.ImplName = decoder.DecodeByteSlice()
	v.AuthoringVersion = decoder.DecodeUint32()
	v.SpecVersion = decoder.DecodeUint32()
	v.ImplVersion = decoder.DecodeUint32()

	apisLength := decoder.DecodeUint32()
	if apisLength != 0 {
		var apis []ApiItem

		for i := 0; i < int(apisLength); i++ {
			apis = append(apis, ApiItem{
				Name:    decodeApiName(decoder),
				Version: decoder.DecodeUint32(),
			})
		}
		v.Apis = apis
	}

	v.TransactionVersion = decoder.DecodeUint32()
	v.StateVersion = decoder.DecodeUint32()

	return nil
}

func decodeApiName(decoder scale.Decoder) [8]byte {
	var result [8]byte
	length := decoder.DecodeUintCompact()

	for i := 0; i < int(length); i++ {
		result[i] = decoder.DecodeByte()
	}

	return result
}
