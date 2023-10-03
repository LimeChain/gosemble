//go:build nonwasmenv

package crypto

import (
	crypto_ed25519 "crypto/ed25519"
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
	panic("not implemented")
}

func (cr ed25519) VerifyVersion1(signature []byte, message []byte, pubKey []byte) bool {
	return crypto_ed25519.Verify(pubKey, message, signature)
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
	panic("not implemented")
}

func (sr sr25519) VerifyVersion2(signature []byte, message []byte, pubKey []byte) bool {
	panic("not implemented")
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
	panic("not implemented")
}

func (sb signatureBatcher) FinishBatchVerify() int32 {
	panic("not implemented")
}
