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

type Digest sc.Dictionary[sc.U8, sc.FixedSequence[DigestItem]]

func (d Digest) Encode(buffer *bytes.Buffer) {
	sc.Compact(len(d)).Encode(buffer)
	for k, v := range d {
		k.Encode(buffer)
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

	result := Digest{}
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
			sc.DecodeU8(buffer)
			// TODO:
		}
	}

	return result
}
