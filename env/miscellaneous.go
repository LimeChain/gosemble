package env

/*
	Miscellaneous: Interface that provides miscellaneous functions for communicating between the runtime and the node.
*/
//go:wasm-module env
//go:export ext_misc_print_hex_version_1
func extMiscPrintHexVersion1(data int64)

func ExtMiscPrintHexVersion1(data int64) {
	extMiscPrintHexVersion1(data)
}

//go:wasm-module env
//go:export ext_misc_print_num_version_1
func extMiscPrintNumVersion1(value int64)

func ExtMiscPrintNumVersion1(value int64) {
	extMiscPrintNumVersion1(value)
}

//go:wasm-module env
//go:export ext_misc_print_utf8_version_1
func extMiscPrintUtf8Version1(data int64)

func ExtMiscPrintUtf8Version1(data int64) {
	extMiscPrintUtf8Version1(data)
}

//go:wasm-module env
//go:export ext_misc_runtime_version_version_1
func extMiscRuntimeVersionVersion1(data int64) int64

func ExtMiscRuntimeVersionVersion1(data int64) int64 {
	return extMiscRuntimeVersionVersion1(data)
}
