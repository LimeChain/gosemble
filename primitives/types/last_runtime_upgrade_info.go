package types

import (
	"bytes"
	"errors"

	sc "github.com/LimeChain/goscale"
)

type LastRuntimeUpgradeInfo struct {
	SpecVersion sc.Compact
	SpecName    sc.Str
}

func (lrui LastRuntimeUpgradeInfo) Encode(buffer *bytes.Buffer) error {
	specVersion, ok := lrui.SpecVersion.Number.(sc.U32)
	if !ok {
		return errors.New("invalid SpecVersion of LastRuntimeUpgradeInfo")
	}
	return sc.EncodeEach(buffer,
		sc.Compact{Number: specVersion},
		lrui.SpecName,
	)
}

func DecodeLastRuntimeUpgradeInfo(buffer *bytes.Buffer) (LastRuntimeUpgradeInfo, error) {
	specVersion, err := sc.DecodeCompact[sc.U32](buffer)
	if err != nil {
		return LastRuntimeUpgradeInfo{}, err
	}
	_, ok := specVersion.Number.(sc.U32)
	if !ok {
		return LastRuntimeUpgradeInfo{}, errors.New("invalid Spec Version Number When Decoding LastRuntimeUpgradeInfo")
	}
	specName, err := sc.DecodeStr(buffer)
	if err != nil {
		return LastRuntimeUpgradeInfo{}, err
	}
	return LastRuntimeUpgradeInfo{
		SpecVersion: specVersion,
		SpecName:    specName,
	}, nil
}

func (lrui LastRuntimeUpgradeInfo) Bytes() []byte {
	return sc.EncodedBytes(lrui)
}
