package dev

/*
#include <stdint.h>

uint32_t ext_allocator_malloc_version_1(uint32_t size);
void ext_allocator_free_version_1(uint32_t ptr);
*/
import "C"

import (
	"fmt"
	"path/filepath"
	"unsafe"

	"github.com/radkomih/gosemble/utils"
	"github.com/wasmerio/go-ext-wasm/wasmer"
)

//export ext_allocator_malloc_version_1
func ext_allocator_malloc_version_1(size uint32) uint32 { return 0 }

//export ext_allocator_free_version_1
func ext_allocator_free_version_1(ptr uint32) {}

func RunInWazmer(wasmRuntimeFile string) {
	modulePath, err := filepath.Abs(wasmRuntimeFile)
	check(err)

	bytes, err := wasmer.ReadBytes(modulePath)
	check(err)

	// Compile bytes into wasm binary
	module, err := wasmer.Compile(bytes)
	check(err)

	// Get current wasi version and corresponded import objects
	wasiVersion := wasmer.WasiGetVersion(module)
	if wasiVersion == 0 {
		// wasiVersion is unknow, use Latest instead
		wasiVersion = wasmer.Latest
	}

	// Instantiate WebAssembly module using derived import objects.
	// importObject := wasmer.NewDefaultWasiImportObjectForVersion(wasiVersion)
	importObject := wasmer.NewDefaultWasiImportObject()

	// Allocate memory from the host (the Wasm module expects to import 20 pages)
	memory, err := wasmer.NewMemory(2, 10)
	check(err)

	// Import host provided memory and functions into the Wasm module
	imports := wasmer.NewImports()
	imports.Namespace("env").AppendMemory("memory", memory)
	imports.Namespace("env").AppendFunction("ext_allocator_malloc_version_1", ext_allocator_malloc_version_1, unsafe.Pointer(C.ext_allocator_malloc_version_1))
	imports.Namespace("env").AppendFunction("ext_allocator_free_version_1", ext_allocator_free_version_1, unsafe.Pointer(C.ext_allocator_free_version_1))
	importObject.Extend(*imports)

	// Instantiate new module
	instance, err := module.InstantiateWithImportObject(importObject)
	check(err)
	defer importObject.Close()
	defer instance.Close()

	mem := memory.Data()

	// Write some data into memory from the host
	data := []byte("Go to Wasm")
	dataPtr := int32(0)
	dataSize := int32(len(data))
	for i := int32(0); i < dataSize; i++ {
		mem[i+dataPtr] = data[i]
	}
	fmt.Printf("%s\n", mem[dataPtr:dataSize+dataPtr])

	// Call an exported function from the Wasm module
	// by passing a ptr and size to the datas
	coreVersion := instance.Exports["Core_version"]
	result, err := coreVersion(dataPtr, dataSize)
	check(err)
	_, resultSize := utils.Int64ToPointerAndSize(uint64(result.ToI64()))
	fmt.Printf("%s\n", mem[dataPtr:resultSize])
}
