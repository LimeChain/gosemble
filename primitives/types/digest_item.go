package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	DigestItemOther                      sc.U8 = 0
	DigestItemConsensusMessage           sc.U8 = 4
	DigestItemSeal                       sc.U8 = 5
	DigestItemPreRuntime                 sc.U8 = 6
	DigestItemRuntimeEnvironmentUpgraded sc.U8 = 8
)

type DigestItem struct {
	sc.VaryingData
}

func NewDigestItemOther(message sc.Sequence[sc.U8]) DigestItem {
	return DigestItem{sc.NewVaryingData(DigestItemOther, message)}
}

func NewDigestItemConsensusMessage(consensusEngineId sc.FixedSequence[sc.U8], message sc.Sequence[sc.U8]) DigestItem {
	// TODO: type safety consensusEngineId must be [4]byte
	return DigestItem{sc.NewVaryingData(DigestItemConsensusMessage, consensusEngineId, message)}
}

func NewDigestItemSeal(consensusEngineId sc.FixedSequence[sc.U8], message sc.Sequence[sc.U8]) DigestItem {
	// TODO: type safety consensusEngineId must be [4]byte
	return DigestItem{sc.NewVaryingData(DigestItemSeal, consensusEngineId, message)}
}

func NewDigestItemPreRuntime(consensusEngineId sc.FixedSequence[sc.U8], message sc.Sequence[sc.U8]) DigestItem {
	// TODO: type safety consensusEngineId must be [4]byte
	return DigestItem{sc.NewVaryingData(DigestItemPreRuntime, consensusEngineId, message)}
}

func NewDigestItemRuntimeEnvironmentUpgrade() DigestItem {
	return DigestItem{sc.NewVaryingData(DigestItemRuntimeEnvironmentUpgraded)}
}

func DecodeDigestItem(buffer *bytes.Buffer) (DigestItem, error) {
	itemType, err := sc.DecodeU8(buffer)
	if err != nil {
		return DigestItem{}, err
	}

	switch itemType {
	case DigestItemOther:
		message, err := sc.DecodeSequence[sc.U8](buffer)
		if err != nil {
			return DigestItem{}, err
		}
		return NewDigestItemOther(message), nil
	case DigestItemConsensusMessage:
		engine, err := sc.DecodeFixedSequence[sc.U8](4, buffer)
		if err != nil {
			return DigestItem{}, err
		}
		message, err := sc.DecodeSequence[sc.U8](buffer)
		if err != nil {
			return DigestItem{}, err
		}

		return NewDigestItemConsensusMessage(engine, message), nil
	case DigestItemSeal:
		engine, err := sc.DecodeFixedSequence[sc.U8](4, buffer)
		if err != nil {
			return DigestItem{}, err
		}
		message, err := sc.DecodeSequence[sc.U8](buffer)
		if err != nil {
			return DigestItem{}, err
		}

		return NewDigestItemSeal(engine, message), nil
	case DigestItemPreRuntime:
		engine, err := sc.DecodeFixedSequence[sc.U8](4, buffer)
		if err != nil {
			return DigestItem{}, err
		}
		message, err := sc.DecodeSequence[sc.U8](buffer)
		if err != nil {
			return DigestItem{}, err
		}

		return NewDigestItemPreRuntime(engine, message), nil
	case DigestItemRuntimeEnvironmentUpgraded:
		return NewDigestItemRuntimeEnvironmentUpgrade(), nil
	default:
		return DigestItem{}, newTypeError("DigestItem")
	}
}

func (di DigestItem) IsSeal() bool {
	// TODO: sanity check
	return di.VaryingData[0] == DigestItemSeal
}

func (di DigestItem) IsPreRuntime() bool {
	return di.VaryingData[0] == DigestItemPreRuntime
}

func (di DigestItem) AsPreRuntime() (DigestPreRuntime, error) {
	if di.IsPreRuntime() {
		return NewDigestPreRuntime(di.VaryingData[1].(sc.FixedSequence[sc.U8]), di.VaryingData[2].(sc.Sequence[sc.U8])), nil
	}
	return DigestPreRuntime{}, newTypeError("DigestPreRuntime")
}

// TODO: has the same fields as DigestPreRuntime, merge at some point
func (di DigestItem) AsSeal() (DigestSeal, error) {
	if di.IsSeal() {
		return NewDigestSeal(di.VaryingData[1].(sc.FixedSequence[sc.U8]), di.VaryingData[2].(sc.Sequence[sc.U8])), nil
	}
	return DigestSeal{}, newTypeError("DigestSeal")
}
