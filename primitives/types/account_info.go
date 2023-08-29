package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type RefCount = sc.U32

type AccountInfo struct {
	Nonce       AccountIndex
	Consumers   RefCount
	Providers   RefCount
	Sufficients RefCount
	Data        AccountData
}

func (ai AccountInfo) Encode(buffer *bytes.Buffer) {
	ai.Nonce.Encode(buffer)
	ai.Consumers.Encode(buffer)
	ai.Providers.Encode(buffer)
	ai.Sufficients.Encode(buffer)
	ai.Data.Encode(buffer)
}

func (ai AccountInfo) Bytes() []byte {
	return sc.EncodedBytes(ai)
}

func DecodeAccountInfo(buffer *bytes.Buffer) AccountInfo {
	return AccountInfo{
		Nonce:       sc.DecodeU32(buffer),
		Consumers:   sc.DecodeU32(buffer),
		Providers:   sc.DecodeU32(buffer),
		Sufficients: sc.DecodeU32(buffer),
		Data:        DecodeAccountData(buffer),
	}
}

func (ai AccountInfo) Frozen(reasons Reasons) sc.U128 {
	switch reasons {
	case ReasonsAll:
		if ai.Data.MiscFrozen.Gt(ai.Data.FeeFrozen) {
			return ai.Data.MiscFrozen
		}
		return ai.Data.FeeFrozen
	case ReasonsMisc:
		return ai.Data.MiscFrozen
	case ReasonsFee:
		return ai.Data.MiscFrozen
	}

	return sc.NewU128(0)
}
