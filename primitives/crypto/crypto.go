//go:build !nonwasmenv

package crypto

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

type Ed25519 interface {
	GenerateVersion1(keyTypeId []byte, seed []byte) []byte
	VerifyVersion1(signature []byte, message []byte, pubKey []byte) bool
}

type ed25519 struct{}

func NewEd25519() ed25519 {
	return ed25519{}
}

func (cr ed25519) GenerateVersion1(keyTypeId []byte, seed []byte) []byte {
	r := env.ExtCryptoEd25519GenerateVersion1(utils.Offset32(keyTypeId), utils.BytesToOffsetAndSize(seed))
	return utils.ToWasmMemorySlice(r, 32)
}

func (cr ed25519) VerifyVersion1(signature []byte, message []byte, pubKey []byte) bool {
	return env.ExtCryptoEd25519VerifyVersion1(
		argsSigMsgPubKeyAsWasmMemory(signature, message, pubKey),
	) == 1
}

type Sr25519 interface {
	GenerateVersion1(keyTypeId []byte, seed []byte) []byte
	VerifyVersion2(signature []byte, message []byte, pubKey []byte) bool
}

type sr25519 struct{}

func NewSr25519() sr25519 {
	return sr25519{}
}

func (sr sr25519) GenerateVersion1(keyTypeId []byte, seed []byte) []byte {
	r := env.ExtCryptoSr25519GenerateVersion1(utils.Offset32(keyTypeId), utils.BytesToOffsetAndSize(seed))
	return utils.ToWasmMemorySlice(r, 32)
}

func (sr sr25519) VerifyVersion2(signature []byte, message []byte, pubKey []byte) bool {
	return env.ExtCryptoSr25519VerifyVersion2(
		argsSigMsgPubKeyAsWasmMemory(signature, message, pubKey),
	) == 1
}

type Ecdsa interface {
	GenerateVersion1(keyTypeId []byte, seed []byte) []byte
}

type ecdsa struct{}

func NewEcdsa() ecdsa {
	return ecdsa{}
}

func (ed ecdsa) GenerateVersion1(keyTypeId []byte, seed []byte) []byte {
	// TODO: ext_crypto_ecdsa_generate_version_1 is not exported by Gossamer
	panic("not exported by Gossamer")
	//r := env.ExtCryptoEcdsaGenerateVersion1(utils.Offset32(keyTypeId), utils.BytesToOffsetAndSize(seed))
	//return utils.ToWasmMemorySlice(r, 32)
}

type SignatureBatcher interface {
	StartBatchVerify()
	FinishBatchVerify() int32
}

type signatureBatcher struct{}

func NewSignatureBatcher() signatureBatcher {
	return signatureBatcher{}
}

func (sb signatureBatcher) StartBatchVerify() {
	env.ExtCryptoStartBatchVerifyVersion1()
}

func (sb signatureBatcher) FinishBatchVerify() int32 {
	return env.ExtCryptoFinishBatchVerifyVersion1()
}

func argsSigMsgPubKeyAsWasmMemory(signature []byte, message []byte, pubKey []byte) (sigOffset int32, msgOffsetSize int64, pubKeyOffset int32) {
	sigOffsetSize := utils.BytesToOffsetAndSize(signature)
	sigOffset, _ = utils.Int64ToOffsetAndSize(sigOffsetSize) // signature: 64-byte

	msgOffsetSize = utils.BytesToOffsetAndSize(message)

	pubKeyOffsetSize := utils.BytesToOffsetAndSize(pubKey)
	pubKeyOffset, _ = utils.Int64ToOffsetAndSize(pubKeyOffsetSize) // public key: 256-bit

	return sigOffset, msgOffsetSize, pubKeyOffset
}
