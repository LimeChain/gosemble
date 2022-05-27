/*
  BlockBuilder API - Version 4.
*/

#ifndef BLOCK_BUILDER_H
#define BLOCK_BUILDER_H
/*
SCALE encoded arguments (extrinsic types.Extrinsic) allocated in the Wasm VM memory, passed as:
dataPtr - i32 pointer to the memory location.
dataLen - i32 length (in bytes) of the encoded arguments.
returns a pointer-size to the SCALE-encoded ([]byte) data.
*/
WASM_EXPORT uint64_t BlockBuilder_apply_extrinsic(uint32_t dataPtr, uint32_t dataLen);

/*
  SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
  dataPtr - i32 pointer to the memory location.
  dataLen - i32 length (in bytes) of the encoded arguments.
  returns a pointer-size to the SCALE-encoded (types.Header) data.
*/
WASM_EXPORT uint64_t BlockBuilder_finalize_block(uint32_t dataPtr, uint32_t dataLen);

/*
  SCALE encoded arguments (data types.InherentsData) allocated in the Wasm VM memory, passed as:
  dataPtr - i32 pointer to the memory location.
  dataLen - i32 length (in bytes) of the encoded arguments.
  returns a pointer-size to the SCALE-encoded ([]types.Extrinsic) data.
*/
WASM_EXPORT uint64_t BlockBuilder_inherent_extrinisics(uint32_t dataPtr, uint32_t dataLen);

/*
  SCALE encoded arguments (block types.Block, data types.InherentsData) allocated in the Wasm VM memory, passed as:
  dataPtr - i32 pointer to the memory location.
  dataLen - i32 length (in bytes) of the encoded arguments.
  returns a pointer-size to the SCALE-encoded ([]byte) data.
*/
WASM_EXPORT uint64_t BlockBuilder_check_inherents(uint32_t dataPtr, uint32_t dataLen);

/*
  SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
  dataPtr - i32 pointer to the memory location.
  dataLen - i32 length (in bytes) of the encoded arguments.
  returns a pointer-size to the SCALE-encoded ([32]byte) data.
*/
WASM_EXPORT uint64_t BlockBuilder_random_seed(uint32_t dataPtr, uint32_t dataLen);
#endif // BLOCK_BUILDER_H