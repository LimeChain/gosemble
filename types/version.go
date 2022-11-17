package types

import (
	"bytes"

	"github.com/LimeChain/goscale"
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
	StateVersion       uint8
}

func (v *VersionData) Encode() ([]byte, error) {
	var buffer = bytes.Buffer{}
	var encoder = goscale.Encoder{Writer: &buffer}

	encoder.EncodeByteSlice(v.SpecName)
	encoder.EncodeByteSlice(v.ImplName)
	encoder.EncodeUint32(v.AuthoringVersion)
	encoder.EncodeUint32(v.SpecVersion)
	encoder.EncodeUint32(v.ImplVersion)

	encoder.EncodeUint8(uint8(len(v.Apis)))
	for _, apiItem := range v.Apis {
		encoder.EncodeByteSlice(apiItem.Name[:])
		encoder.EncodeUint32(apiItem.Version)
	}

	encoder.EncodeUint32(v.TransactionVersion)
	encoder.EncodeUint8(v.StateVersion)

	return buffer.Bytes(), nil
}

func (v *VersionData) Decode(enc []byte) error {
	var buffer = bytes.NewBuffer(enc)
	var decoder = goscale.Decoder{Reader: buffer}

	v.SpecName = decoder.DecodeByteSlice()
	v.ImplName = decoder.DecodeByteSlice()
	v.AuthoringVersion = decoder.DecodeUint32()
	v.SpecVersion = decoder.DecodeUint32()
	v.ImplVersion = decoder.DecodeUint32()

	apisLength := decoder.DecodeUint8()
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
	v.StateVersion = decoder.DecodeUint8()

	return nil
}

func decodeApiName(decoder goscale.Decoder) [8]byte {
	var result [8]byte
	length := decoder.DecodeUintCompact()

	for i := 0; i < int(length); i++ {
		result[i] = decoder.DecodeByte()
	}

	return result
}
