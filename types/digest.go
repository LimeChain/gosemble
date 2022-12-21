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

type Digest struct {
	Values map[uint8]sc.FixedSequence[DigestItem]
}

func (d Digest) Encode(buffer *bytes.Buffer) {
	sc.Compact(len(d.Values)).Encode(buffer)
	for k, v := range d.Values {
		sc.U8(k).Encode(buffer)
		v.Encode(buffer)
	}
}

func (d Digest) Bytes() []byte {
	buffer := &bytes.Buffer{}
	d.Encode(buffer)

	return buffer.Bytes()
}

func DecodeDigest(buffer *bytes.Buffer) Digest {
	length := sc.DecodeCompact(buffer)

	decoder := sc.Decoder{buffer}

	result := map[uint8]sc.FixedSequence[DigestItem]{}
	for i := 0; i < int(length); i++ {
		digestType := decoder.DecodeByte()

		switch digestType {
		case DigestTypeConsensusMessage:
			consensusDigest := DecodeDigestItem(buffer)
			result[DigestTypeConsensusMessage] = append(result[DigestTypeConsensusMessage], consensusDigest)
		case DigestTypeSeal:
			seal := DecodeDigestItem(buffer)
			result[DigestTypeSeal] = append(result[DigestTypeSeal], seal)
		case DigestTypePreRuntime:
			preRuntimeDigest := DecodeDigestItem(buffer)
			result[DigestTypePreRuntime] = append(result[DigestTypePreRuntime], preRuntimeDigest)
		case DigestTypeRuntimeEnvironmentUpgraded:
			// TODO:
		default:
			panic("invalid digest type")
		}
	}

	return Digest{
		result,
	}
}
