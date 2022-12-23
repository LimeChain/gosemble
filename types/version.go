package types

import (
	"bytes"
	"fmt"

	sc "github.com/LimeChain/goscale"
)

type ApiItem struct {
	Name    sc.FixedSequence[sc.U8] // size 8
	Version sc.U32
}

func (ai ApiItem) Encode(buffer *bytes.Buffer) {
	ai.Name.Encode(buffer)
	ai.Version.Encode(buffer)
	// sc.Tuple[ApiItem]{Data: ai}.Encode(buffer)
}

func (ai ApiItem) Bytes() []byte {
	buffer := &bytes.Buffer{}
	ai.Encode(buffer)

	return buffer.Bytes()
}

func DecodeApiItem(buffer *bytes.Buffer) ApiItem {
	return ApiItem{
		Name:    sc.DecodeFixedSequence[sc.U8](8, buffer),
		Version: sc.DecodeU32(buffer),
	}
}

func (ai ApiItem) String() string {
	return fmt.Sprintf("ApiItem { Name: %#x, Version: %d}", ai.Name, ai.Version)
}

type VersionData struct {
	SpecName           sc.Str
	ImplName           sc.Str
	AuthoringVersion   sc.U32
	SpecVersion        sc.U32
	ImplVersion        sc.U32
	Apis               sc.Sequence[ApiItem]
	TransactionVersion sc.U32
	StateVersion       sc.U8
}

func (vd VersionData) Encode(buffer *bytes.Buffer) {
	vd.SpecName.Encode(buffer)
	vd.ImplName.Encode(buffer)
	vd.AuthoringVersion.Encode(buffer)
	vd.SpecVersion.Encode(buffer)
	vd.ImplVersion.Encode(buffer)
	vd.Apis.Encode(buffer)
	vd.TransactionVersion.Encode(buffer)
	vd.StateVersion.Encode(buffer)
	// sc.Tuple[VersionData]{Data: vd}.Encode(buffer)
}

func DecodeVersionData(buffer *bytes.Buffer) VersionData {
	var vd VersionData

	vd.SpecName = sc.DecodeStr(buffer)
	vd.ImplName = sc.DecodeStr(buffer)
	vd.AuthoringVersion = sc.DecodeU32(buffer)
	vd.SpecVersion = sc.DecodeU32(buffer)
	vd.ImplVersion = sc.DecodeU32(buffer)

	apisLength := sc.DecodeCompact(buffer).ToBigInt().Int64()
	if apisLength != 0 {
		var apis []ApiItem
		for i := 0; i < int(apisLength); i++ {
			apis = append(apis, DecodeApiItem(buffer))
		}
		vd.Apis = apis
	}

	vd.TransactionVersion = sc.DecodeU32(buffer)
	vd.StateVersion = sc.DecodeU8(buffer)

	return vd
}

func (vd VersionData) String() string {
	var result string

	result = "VersionData {\n"
	result += fmt.Sprintf("SpecName: %s\n", vd.SpecName)
	result += fmt.Sprintf("ImplName: %s\n", vd.ImplName)
	result += fmt.Sprintf("AuthoringVersion: %d\n", vd.AuthoringVersion)
	result += fmt.Sprintf("SpecVersion: %d\n", vd.SpecVersion)
	result += fmt.Sprintf("ImplVersion: %d\n", vd.ImplVersion)
	result += "Apis: ["
	for _, v := range vd.Apis {
		result += v.String()
	}
	result += "]\n"
	result += fmt.Sprintf("TransactionVersion: %d\n", vd.TransactionVersion)
	result += fmt.Sprintf("StateVersion: %d\n", vd.StateVersion)
	result += "}"

	return result
}
