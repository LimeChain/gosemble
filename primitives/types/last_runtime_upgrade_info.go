package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type LastRuntimeUpgradeInfo struct {
	SpecVersion sc.U32
	SpecName    sc.Str
}

func (lrui LastRuntimeUpgradeInfo) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		sc.ToCompact(lrui.SpecVersion),
		lrui.SpecName,
	)
}

func DecodeLastRuntimeUpgradeInfo(buffer *bytes.Buffer) (LastRuntimeUpgradeInfo, error) {
	specVersion, err := sc.DecodeCompact(buffer)
	if err != nil {
		return LastRuntimeUpgradeInfo{}, err
	}
	specName, err := sc.DecodeStr(buffer)
	if err != nil {
		return LastRuntimeUpgradeInfo{}, err
	}
	return LastRuntimeUpgradeInfo{
		SpecVersion: sc.U32(specVersion.ToBigInt().Uint64()),
		SpecName:    specName,
	}, nil
}

func (lrui LastRuntimeUpgradeInfo) Bytes() []byte {
	return sc.EncodedBytes(lrui)
}
