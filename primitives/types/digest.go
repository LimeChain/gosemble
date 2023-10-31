package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	DigestTypeConsensusMessage           = 4
	DigestTypeSeal                       = 5
	DigestTypePreRuntime                 = 6
	DigestTypeRuntimeEnvironmentUpgraded = 8
)

type Digest = sc.Dictionary[sc.U8, sc.FixedSequence[DigestItem]]

func DecodeDigest(buffer *bytes.Buffer) (Digest, error) {
	compactSize, err := sc.DecodeCompact(buffer)
	if err != nil {
		return Digest{}, err
	}
	size := int(compactSize.ToBigInt().Int64())

	decoder := sc.Decoder{Reader: buffer}

	result := Digest{}
	for i := 0; i < size; i++ {
		digestType, err := decoder.DecodeByte()
		if err != nil {
			return Digest{}, err
		}

		switch digestType {
		case DigestTypeConsensusMessage:
			consensusDigest, err := DecodeDigestItem(buffer)
			if err != nil {
				return Digest{}, err
			}
			result[DigestTypeConsensusMessage] = append(result[DigestTypeConsensusMessage], consensusDigest)
		case DigestTypeSeal:
			seal, err := DecodeDigestItem(buffer)
			if err != nil {
				return Digest{}, err
			}
			result[DigestTypeSeal] = append(result[DigestTypeSeal], seal)
		case DigestTypePreRuntime:
			preRuntimeDigest, err := DecodeDigestItem(buffer)
			if err != nil {
				return Digest{}, err
			}
			result[DigestTypePreRuntime] = append(result[DigestTypePreRuntime], preRuntimeDigest)
		case DigestTypeRuntimeEnvironmentUpgraded:
			sc.DecodeU8(buffer)
			// TODO:
		}
	}

	return result, nil
}
