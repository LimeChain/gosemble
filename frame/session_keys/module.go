package session_keys

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/crypto"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

type Module struct {
	sessions []types.Session
}

func New(sessions []types.Session) Module {
	return Module{sessions: sessions}
}

// GenerateSessionKeys generates a set of session keys with an optional seed.
// The keys should be stored within the keystore exposed by the Host Api.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded optional seed.
// Returns a pointer-size of the SCALE-encoded set of keys.
// [Specification](https://spec.polkadot.network/chap-runtime-api#id-sessionkeys_generate_session_keys)
func (m Module) GenerateSessionKeys(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	seed := sc.DecodeOptionWith(buffer, sc.DecodeSequence[sc.U8])

	var publicKeys []byte
	for _, session := range m.sessions {
		keyGenerationFunc := getKeyFunction(session.KeyType())
		keyTypeId := session.KeyTypeId()

		publicKey := keyGenerationFunc(keyTypeId[:], seed.Bytes())
		publicKeys = append(publicKeys, publicKey...)
	}

	res := sc.BytesToSequenceU8(publicKeys)

	return utils.BytesToOffsetAndSize(res.Bytes())
}

// DecodeSessionKeys decodes the given session keys.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded keys.
// Returns a pointer-size of the SCALE-encoded set of raw keys and their respective key type.
// [Specification](https://spec.polkadot.network/chap-runtime-api#id-sessionkeys_decode_session_keys)
func (m Module) DecodeSessionKeys(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)
	sequence := sc.DecodeSequenceWith(buffer, sc.DecodeU8)

	buffer = bytes.NewBuffer(sc.SequenceU8ToBytes(sequence))
	sessionKeys := sc.Sequence[types.SessionKey]{}
	for _, session := range m.sessions {
		sessionKey := types.NewSessionKey(sc.FixedSequenceU8ToBytes(types.DecodePublicKey(buffer)), session.KeyTypeId())
		sessionKeys = append(sessionKeys, sessionKey)
	}

	result := sc.NewOption[sc.Sequence[types.SessionKey]](sessionKeys)
	return utils.BytesToOffsetAndSize(result.Bytes())
}

func getKeyFunction(keyType types.PublicKeyType) func([]byte, []byte) []byte {
	switch keyType {
	case types.PublicKeyEd25519:
		return crypto.ExtCryptoEd25519GenerateVersion1
	case types.PublicKeySr25519:
		return crypto.ExtCryptoSr25519GenerateVersion1
	case types.PublicKeyEcdsa:
		return crypto.ExtCryptoEcdsaGenerateVersion1
	default:
		log.Critical("invalid public key type")
	}

	panic("unreachable")
}
