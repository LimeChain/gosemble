/*
  Core API - Version 3.
*/

#ifndef CORE_H
#define CORE_H
/*
  SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
  dataPtr - i32 pointer to the memory location.
  dataLen - i32 length (in bytes) of the encoded arguments.
  returns a pointer-size to the SCALE-encoded (version types.VersionData) data.
*/
WASM_EXPORT uint64_t Core_version(uint32_t dataPtr, uint32_t dataLen);

/*
  SCALE encoded arguments (block types.Block) allocated in the Wasm VM memory, passed as:
  dataPtr - i32 pointer to the memory location.
  dataLen - i32 length (in bytes) of the encoded arguments.
*/
WASM_EXPORT void Core_execute_block(uint32_t dataPtr, uint32_t dataLen);

/*
  SCALE encoded arguments (header *types.Header) allocated in the Wasm VM memory, passed as:
  dataPtr - i32 pointer to the memory location.
  dataLen - i32 length (in bytes) of the encoded arguments.
*/
WASM_EXPORT void Core_initialize_block(uint32_t dataPtr, uint32_t dataLen);
#endif // CORE_H