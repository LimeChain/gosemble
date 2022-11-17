package env

/*
	Allocator: Provides functionality for calling into the memory allocator.
*/
//go:wasm-module env
//go:export ext_allocator_malloc_version_1
func extAllocatorMallocVersion1(size int32) int32

func ExtAllocatorMallocVersion1(size int32) int32 {
	return extAllocatorMallocVersion1(size)
}

//go:wasm-module env
//go:export ext_allocator_free_version_1
func extAllocatorFreeVersion1(ptr int32)

func ExtAllocatorFreeVersion1(ptr int32) {
	extAllocatorFreeVersion1(ptr)
}
