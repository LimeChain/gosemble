//go:build !nonwasmenv

package env

/*
	Allocator: Provides functionality for calling into the memory allocator.
	Actual GC functions are in TinyGo runtime_polkawasm.go
*/

//go:wasmimports env ext_allocator_free_version_1
func ExtAllocatorFreeVersion1(ptr int32)

//go:wasmimports env ext_allocator_malloc_version_1
func ExtAllocatorMallocVersion1(size int32) int32
