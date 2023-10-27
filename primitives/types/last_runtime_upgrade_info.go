package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type LastRuntimeUpgradeInfo struct {
	SpecVersion sc.U32
	SpecName    sc.Str
}

func (lrui LastRuntimeUpgradeInfo) Encode(buffer *bytes.Buffer) {
	sc.ToCompact(lrui.SpecVersion).Encode(buffer)
	lrui.SpecName.Encode(buffer)
}

func DecodeLastRuntimeUpgradeInfo(buffer *bytes.Buffer) LastRuntimeUpgradeInfo {
	return LastRuntimeUpgradeInfo{
		SpecVersion: sc.U32(sc.DecodeCompact(buffer).ToBigInt().Uint64()),
		SpecName:    sc.DecodeStr(buffer),
	}
}

func (lrui LastRuntimeUpgradeInfo) Bytes() []byte {
	return sc.EncodedBytes(lrui)
}
