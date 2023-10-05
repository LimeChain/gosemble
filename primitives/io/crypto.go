package io

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

type Crypto interface {
	EcdsaGenerate(keyTypeId []byte, seed []byte) []byte

	Ed25519Generate(keyTypeId []byte, seed []byte) []byte
	Ed25519Verify(signature []byte, message []byte, pubKey []byte) bool

	Sr25519Generate(keyTypeId []byte, seed []byte) []byte
	Sr25519Verify(signature []byte, message []byte, pubKey []byte) bool
}

type crypto struct{}

func NewCrypto() Crypto {
	return crypto{}
}

func (c crypto) Ed25519Generate(keyTypeId []byte, seed []byte) []byte {
	r := env.ExtCryptoEd25519GenerateVersion1(utils.Offset32(keyTypeId), utils.BytesToOffsetAndSize(seed))
	return utils.ToWasmMemorySlice(r, 32)
}

func (c crypto) Ed25519Verify(signature []byte, message []byte, pubKey []byte) bool {
	return env.ExtCryptoEd25519VerifyVersion1(
		argsSigMsgPubKeyAsWasmMemory(signature, message, pubKey),
	) == 1
}

func (c crypto) Sr25519Generate(keyTypeId []byte, seed []byte) []byte {
	r := env.ExtCryptoSr25519GenerateVersion1(utils.Offset32(keyTypeId), utils.BytesToOffsetAndSize(seed))
	return utils.ToWasmMemorySlice(r, 32)
}

func (c crypto) Sr25519Verify(signature []byte, message []byte, pubKey []byte) bool {
	return env.ExtCryptoSr25519VerifyVersion2(
		argsSigMsgPubKeyAsWasmMemory(signature, message, pubKey),
	) == 1
}

func (c crypto) EcdsaGenerate(keyTypeId []byte, seed []byte) []byte {
	// TODO: ext_crypto_ecdsa_generate_version_1 is not exported by Gossamer
	panic("not exported by Gossamer")
	//r := env.ExtCryptoEcdsaGenerateVersion1(utils.Offset32(keyTypeId), utils.BytesToOffsetAndSize(seed))
	//return utils.ToWasmMemorySlice(r, 32)
}

func argsSigMsgPubKeyAsWasmMemory(signature []byte, message []byte, pubKey []byte) (sigOffset int32, msgOffsetSize int64, pubKeyOffset int32) {
	sigOffsetSize := utils.BytesToOffsetAndSize(signature)
	sigOffset, _ = utils.Int64ToOffsetAndSize(sigOffsetSize) // signature: 64-byte

	msgOffsetSize = utils.BytesToOffsetAndSize(message)

	pubKeyOffsetSize := utils.BytesToOffsetAndSize(pubKey)
	pubKeyOffset, _ = utils.Int64ToOffsetAndSize(pubKeyOffsetSize) // public key: 256-bit

	return sigOffset, msgOffsetSize, pubKeyOffset
}
