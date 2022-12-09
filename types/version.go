package types

import (
	"bytes"
	"fmt"

	sc "github.com/LimeChain/goscale"
)

type ApiItem struct {
	Name    sc.FixedSequence[sc.U8] // TODO: https://github.com/LimeChain/goscale/issues/37
	Version sc.U32
}

func (api ApiItem) Encode(buffer *bytes.Buffer) {
	api.Name.Encode(buffer)
	api.Version.Encode(buffer)
}

func DecodeApiItem(buffer *bytes.Buffer) ApiItem {
	return ApiItem{
		Name:    sc.DecodeFixedSequence[sc.U8](8, buffer),
		Version: sc.DecodeU32(buffer),
	}
}

func (api ApiItem) String() string {
	return fmt.Sprintf("ApiItem { Name: %#x, Version: %d}", api.Name, api.Version)
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

func (v VersionData) Encode(buffer *bytes.Buffer) {
	v.SpecName.Encode(buffer)
	v.ImplName.Encode(buffer)
	v.AuthoringVersion.Encode(buffer)
	v.SpecVersion.Encode(buffer)
	v.ImplVersion.Encode(buffer)
	v.Apis.Encode(buffer)
	v.TransactionVersion.Encode(buffer)
	v.StateVersion.Encode(buffer)
}

func DecodeVersionData(buffer *bytes.Buffer) VersionData {
	var v VersionData

	v.SpecName = sc.DecodeStr(buffer)
	v.ImplName = sc.DecodeStr(buffer)
	v.AuthoringVersion = sc.DecodeU32(buffer)
	v.SpecVersion = sc.DecodeU32(buffer)
	v.ImplVersion = sc.DecodeU32(buffer)

	apisLength := sc.DecodeCompact(buffer)
	if apisLength != 0 {
		var apis []ApiItem
		for i := 0; i < int(apisLength); i++ {
			apis = append(apis, DecodeApiItem(buffer))
		}
		v.Apis = apis
	}

	v.TransactionVersion = sc.DecodeU32(buffer)
	v.StateVersion = sc.DecodeU8(buffer)

	return v
}

func (v VersionData) String() string {
	var result string

	result = "VersionData {\n"
	result += fmt.Sprintf("SpecName: %s\n", v.SpecName)
	result += fmt.Sprintf("ImplName: %s\n", v.ImplName)
	result += fmt.Sprintf("AuthoringVersion: %d\n", v.AuthoringVersion)
	result += fmt.Sprintf("SpecVersion: %d\n", v.SpecVersion)
	result += fmt.Sprintf("ImplVersion: %d\n", v.ImplVersion)
	result += "Apis: ["
	for _, v := range v.Apis {
		result += v.String()
	}
	result += "]\n"
	result += fmt.Sprintf("TransactionVersion: %d\n", v.TransactionVersion)
	result += fmt.Sprintf("StateVersion: %d\n", v.StateVersion)
	result += "}"

	return result
}
