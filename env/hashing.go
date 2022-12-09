package env

/*
	Hashing: Interface that provides functions for hashing with diﬀerent algorithms.
*/
//go:wasm-module env
//go:export ext_hashing_blake2_128_version_1
func extHashingBlake2128Version1(data int64) int32

func ExtHashingBlake2128Version1(data int64) int32 {
	return extHashingBlake2128Version1(data)
}

//go:wasm-module env
//go:export ext_hashing_blake2_256_version_1
func extHashingBlake2256Version1(data int64) int32

func ExtHashingBlake2256Version1(data int64) int32 {
	return extHashingBlake2256Version1(data)
}

//go:wasm-module env
//go:export ext_hashing_keccak_256_version_1
func extHashingKeccak256Version1(data int64) int32

func ExtHashingKeccak256Version1(data int64) int32 {
	return extHashingKeccak256Version1(data)
}

//go:wasm-module env
//go:export ext_hashing_twox_128_version_1
func extHashingTwox128Version1(data int64) int32

func ExtHashingTwox128Version1(data int64) int32 {
	return extHashingTwox128Version1(data)
}

//go:wasm-module env
//go:export ext_hashing_twox_64_version_1
func extHashingTwox64Version1(data int64) int32

func ExtHashingTwox64Version1(data int64) int32 {
	return extHashingTwox64Version1(data)
}
