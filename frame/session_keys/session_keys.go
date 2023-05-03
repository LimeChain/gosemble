package session_keys

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/aura"
	"github.com/LimeChain/gosemble/constants/grandpa"
	"github.com/LimeChain/gosemble/primitives/crypto"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

func GenerateSessionKeys(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	seed := sc.DecodeOptionWith(buffer, sc.DecodeSequence[sc.U8])

	auraPubKey := crypto.ExtCryptoSr25519GenerateVersion1(aura.KeyTypeId[:], seed.Bytes())
	grandpaPubKey := crypto.ExtCryptoEd25519GenerateVersion1(grandpa.KeyTypeId[:], seed.Bytes())

	res := sc.BytesToSequenceU8(append(auraPubKey, grandpaPubKey...))

	return utils.BytesToOffsetAndSize(res.Bytes())
}

func DecodeSessionKeys(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)
	sequence := sc.DecodeSequenceWith(buffer, sc.DecodeU8)

	buffer = bytes.NewBuffer(sc.SequenceU8ToBytes(sequence))
	sessionKeys := sc.Sequence[types.SessionKey]{
		types.NewSessionKey(sc.FixedSequenceU8ToBytes(types.DecodePublicKey(buffer)), aura.KeyTypeId),
		types.NewSessionKey(sc.FixedSequenceU8ToBytes(types.DecodePublicKey(buffer)), grandpa.KeyTypeId),
	}

	result := sc.NewOption[sc.Sequence[types.SessionKey]](sessionKeys)
	return utils.BytesToOffsetAndSize(result.Bytes())
}
