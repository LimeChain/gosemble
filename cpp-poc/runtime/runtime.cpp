/*
  Targets WebAssembly MVP.
*/

// #include <cmath>
#include "runtime.hpp"

const char *SPEC_NAME = "gosemble";
const char *IMPL_NAME = "C++";
const int AUTHORING_VERSION = 1;
const int SPEC_VERSION = 1;
const int IMPL_VERSION = 1;
const int TRANSACTION_VERSION = 1;
const int STATE_VERSION = 1;

uint64_t Core_version(uint32_t dataPtr, uint32_t dataSize)
{
  // TODO use the external memory management functions to allocate memory
  // uint8_t *mem = &memory;
  // if (mem != nullptr) {}
  ext_allocator_malloc_version_1(0);
  ext_allocator_free_version_1(0);

  return uint64_t(dataPtr + dataSize) | (uint64_t(dataSize + dataSize) << 32);
}
