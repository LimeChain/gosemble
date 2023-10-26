package session_keys

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/io"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "SessionKeys"
	apiVersion    = 1
)

type Module struct {
	sessions []types.Session
	crypto   io.Crypto
	memUtils utils.WasmMemoryTranslator
}

func New(sessions []types.Session) Module {
	return Module{
		sessions: sessions,
		crypto:   io.NewCrypto(),
		memUtils: utils.NewMemoryTranslator(),
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() types.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return types.NewApiItem(hash, apiVersion)
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
	b := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	seed, err := sc.DecodeOptionWith(buffer, sc.DecodeSequence[sc.U8])
	if err != nil {
		return m.memUtils.BytesToOffsetAndSize([]byte(err.Error()))
	}

	var publicKeys []byte
	for _, session := range m.sessions {
		keyGenerationFunc := getKeyFunction(m, session.KeyType())
		keyTypeId := session.KeyTypeId()

		publicKey := keyGenerationFunc(keyTypeId[:], seed.Bytes())
		publicKeys = append(publicKeys, publicKey...)
	}

	res := sc.BytesToSequenceU8(publicKeys)

	return m.memUtils.BytesToOffsetAndSize(res.Bytes())
}

// DecodeSessionKeys decodes the given session keys.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded keys.
// Returns a pointer-size of the SCALE-encoded set of raw keys and their respective key type.
// [Specification](https://spec.polkadot.network/chap-runtime-api#id-sessionkeys_decode_session_keys)
func (m Module) DecodeSessionKeys(dataPtr int32, dataLen int32) int64 {
	b := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)
	sequence, err := sc.DecodeSequenceWith(buffer, sc.DecodeU8)
	if err != nil {
		log.Critical(err.Error())
		return 0
	}

	buffer = bytes.NewBuffer(sc.SequenceU8ToBytes(sequence))
	sessionKeys := sc.Sequence[types.SessionKey]{}
	for _, session := range m.sessions {
		pk, err := types.DecodePublicKey(buffer)
		if err != nil {
			log.Critical(err.Error())
			return 0
		}
		sessionKey := types.NewSessionKey(sc.FixedSequenceU8ToBytes(pk), session.KeyTypeId())
		sessionKeys = append(sessionKeys, sessionKey)
	}

	result := sc.NewOption[sc.Sequence[types.SessionKey]](sessionKeys)
	return m.memUtils.BytesToOffsetAndSize(result.Bytes())
}

func getKeyFunction(m Module, keyType types.PublicKeyType) func([]byte, []byte) []byte {
	switch keyType {
	case types.PublicKeyEd25519:
		return m.crypto.Ed25519Generate
	case types.PublicKeySr25519:
		return m.crypto.Sr25519Generate
	case types.PublicKeyEcdsa:
		return m.crypto.EcdsaGenerate
	default:
		log.Critical("invalid public key type")
	}

	panic("unreachable")
}
