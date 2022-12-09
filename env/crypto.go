package env

/*
	Crypto: Interfaces for working with crypto related types from within the runtime.
*/
//go:wasm-module env
//go:export ext_crypto_ed25519_generate_version_1
func extCryptoEd25519GenerateVersion1(key_type_id int32, seed int64) int32

func ExtCryptoEd25519GenerateVersion1(key_type_id int32, seed int64) int32 {
	return extCryptoEd25519GenerateVersion1(key_type_id, seed)
}

//go:wasm-module env
//go:export ext_crypto_ed25519_verify_version_1
func extCryptoEd25519VerifyVersion1(sig int32, msg int64, key int32) int32

func ExtCryptoEd25519VerifyVersion1(sig int32, msg int64, key int32) int32 {
	return extCryptoEd25519VerifyVersion1(sig, msg, key)
}

//go:wasm-module env
//go:export ext_crypto_secp256k1_ecdsa_recover_compressed_version_1
func extCryptoSecp256k1ScdsaRecoverCompressedVersion1(sig int32, msg int32) int64

func ExtCryptoSecp256k1ScdsaRecoverCompressedVersion1(sig int32, msg int32) int64 {
	return extCryptoSecp256k1ScdsaRecoverCompressedVersion1(sig, msg)
}

//go:wasm-module env
//go:export ext_crypto_secp256k1_ecdsa_recover_version_1
func extCryptoSecp256k1EcdsaRecoverVersion1(sig int32, msg int32) int64

func ExtCryptoSecp256k1EcdsaRecoverVersion1(sig int32, msg int32) int64 {
	return extCryptoSecp256k1EcdsaRecoverVersion1(sig, msg)
}

//go:wasm-module env
//go:export ext_crypto_sr25519_generate_version_1
func extCryptoSr25519GenerateVersion1(key_type_id int32, seed int64) int32

func ExtCryptoSr25519GenerateVersion1(key_type_id int32, seed int64) int32 {
	return extCryptoSr25519GenerateVersion1(key_type_id, seed)
}

//go:wasm-module env
//go:export ext_crypto_sr25519_public_keys_version_1
func extCryptoSr25519PublicKeysVersion1(key_type_id int64) int64

func ExtCryptoSr25519PublicKeysVersion1(key_type_id int64) int64 {
	return extCryptoSr25519PublicKeysVersion1(key_type_id)
}

//go:wasm-module env
//go:export ext_crypto_sr25519_sign_version_1
func extCryptoSr25519SignVersion1(key_type_id int32, key int32, msg int64) int64

func ExtCryptoSr25519SignVersion1(key_type_id int32, key int32, msg int64) int64 {
	return extCryptoSr25519SignVersion1(key_type_id, key, msg)
}

//go:wasm-module env
//go:export ext_crypto_sr25519_verify_version_1
func extCryptoSr25519VerifyVersion1(sig int32, msg int64, key int32) int32

func ExtCryptoSr25519VerifyVersion1(sig int32, msg int64, key int32) int32 {
	return extCryptoSr25519VerifyVersion1(sig, msg, key)
}

//go:wasm-module env
//go:export ext_crypto_sr25519_verify_version_2
func extCryptoSr25519VerifyVersion2(sig int32, msg int64, key int32) int32

func ExtCryptoSr25519VerifyVersion2(sig int32, msg int64, key int32) int32 {
	return extCryptoSr25519VerifyVersion2(sig, msg, key)
}
