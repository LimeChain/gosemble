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

type crypto struct {
	memoryTranslator utils.WasmMemoryTranslator
}

func NewCrypto() Crypto {
	return crypto{
		memoryTranslator: utils.NewMemoryTranslator(),
	}
}

func (c crypto) EcdsaGenerate(keyTypeId []byte, seed []byte) []byte {
	// TODO: ext_crypto_ecdsa_generate_version_1 is not exported by Gossamer
	panic("not exported by Gossamer")
	//r := env.ExtCryptoEcdsaGenerateVersion1(c.memoryTranslator.Offset32(keyTypeId), c.memoryTranslator.BytesToOffsetAndSize(seed))
	//return c.memoryTranslator.ToWasmMemorySlice(r, 32)
}

func (c crypto) Ed25519Generate(keyTypeId []byte, seed []byte) []byte {
	r := env.ExtCryptoEd25519GenerateVersion1(
		c.memoryTranslator.Offset32(keyTypeId),
		c.memoryTranslator.BytesToOffsetAndSize(seed),
	)
	return c.memoryTranslator.GetWasmMemorySlice(r, 32)
}

func (c crypto) Ed25519Verify(signature []byte, message []byte, pubKey []byte) bool {
	return env.ExtCryptoEd25519VerifyVersion1(
		argsSigMsgPubKeyAsWasmMemory(c.memoryTranslator, signature, message, pubKey),
	) == 1
}

func (c crypto) Sr25519Generate(keyTypeId []byte, seed []byte) []byte {
	r := env.ExtCryptoSr25519GenerateVersion1(
		c.memoryTranslator.Offset32(keyTypeId),
		c.memoryTranslator.BytesToOffsetAndSize(seed),
	)
	return c.memoryTranslator.GetWasmMemorySlice(r, 32)
}

func (c crypto) Sr25519Verify(signature []byte, message []byte, pubKey []byte) bool {
	return env.ExtCryptoSr25519VerifyVersion2(
		argsSigMsgPubKeyAsWasmMemory(c.memoryTranslator, signature, message, pubKey),
	) == 1
}

func argsSigMsgPubKeyAsWasmMemory(mem utils.WasmMemoryTranslator, signature []byte, message []byte, pubKey []byte) (sigOffset int32, msgOffsetSize int64, pubKeyOffset int32) {
	sigOffsetSize := mem.BytesToOffsetAndSize(signature)
	sigOffset, _ = mem.Int64ToOffsetAndSize(sigOffsetSize) // signature: 64-byte

	msgOffsetSize = mem.BytesToOffsetAndSize(message)

	pubKeyOffsetSize := mem.BytesToOffsetAndSize(pubKey)
	pubKeyOffset, _ = mem.Int64ToOffsetAndSize(pubKeyOffsetSize) // public key: 256-bit

	return sigOffset, msgOffsetSize, pubKeyOffset
}
