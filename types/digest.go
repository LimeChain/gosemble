package types

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
)

const (
	DigestTypeNonSystem                  = 0
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

func DecodeDigest(buffer *bytes.Buffer) Digest {
	length := sc.DecodeCompact(buffer)

	decoder := sc.Decoder{buffer}

	result := map[uint8]sc.FixedSequence[DigestItem]{}
	for i := 0; i < int(length); i++ {
		digestType := decoder.DecodeByte()

		switch digestType {
		case DigestTypeNonSystem:
			// TODO: DecodeSequence[Byte] == ByteArray
			//result = append(result, sc.DecodeSliceU8(buffer))
		case DigestTypeConsensusMessage:
			consensusDigest := DigestItem{
				Engine:  sc.DecodeFixedSequence[sc.U8](4, buffer),
				Payload: sc.DecodeSequence[sc.U8](buffer),
			}
			result[DigestTypeConsensusMessage] = append(result[DigestTypeConsensusMessage], consensusDigest)
		case DigestTypeSeal:
			seal := DigestItem{
				Engine:  sc.DecodeFixedSequence[sc.U8](4, buffer),
				Payload: sc.DecodeSequence[sc.U8](buffer),
			}
			result[DigestTypeSeal] = append(result[DigestTypeSeal], seal)
		case DigestTypePreRuntime:
			preRuntimeDigest := DigestItem{
				Engine:  sc.DecodeFixedSequence[sc.U8](4, buffer),
				Payload: sc.DecodeSequence[sc.U8](buffer),
			}
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
