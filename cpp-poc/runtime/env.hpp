/*
  Externaly imported functions providing access to memory and storage.
  Need to be referenced somewhere to be actualy imported.
*/

#ifndef ENV_H
#define ENV_H
// __attribute__((__import_name__("ext_allocator_malloc_version_1")))
extern "C" uint32_t ext_allocator_malloc_version_1(uint32_t size);

// __attribute__((__import_name__("ext_allocator_free_version_1")))
extern "C" void ext_allocator_free_version_1(uint32_t ptr);

extern uint8_t memory;
#endif // ENV_H