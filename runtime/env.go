package main

// "//export" on a func is actually an import in TinyGo.
// The function needs to be referenced somewhere to be actualy exported.

//go:wasm-module env
//export ext_allocator_malloc_version_1
func extAllocatorMallocVersion1(size uint32) uint32

//go:wasm-module env
//export ext_allocator_free_version_1
func extAllocatorFreeVersion1(ptr uint32)
