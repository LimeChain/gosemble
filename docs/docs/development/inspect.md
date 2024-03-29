---
layout: default
title: Inspect
permalink: /development/inspect
---

# Inspect 🔬

Install [wasmer](https://wasmer.io/) to get a simple view of the compiled WASM.

```bash
wasmer inspect build/runtime.wasm
```

```bash
Type: wasm
Size: 873.2 KB
Imports:
  Functions:
    "env"."ext_allocator_malloc_version_1": [I32] -> [I32]
    "env"."ext_allocator_free_version_1": [I32] -> []
    "env"."ext_crypto_secp256k1_ecdsa_recover_compressed_version_2": [I32, I32] -> [I64]
    "env"."ext_crypto_ed25519_generate_version_1": [I32, I64] -> [I32]
    "env"."ext_crypto_ed25519_verify_version_1": [I32, I64, I32] -> [I32]
    "env"."ext_crypto_sr25519_generate_version_1": [I32, I64] -> [I32]
    "env"."ext_crypto_sr25519_verify_version_2": [I32, I64, I32] -> [I32]
    "env"."ext_hashing_blake2_128_version_1": [I64] -> [I32]
    "env"."ext_hashing_blake2_256_version_1": [I64] -> [I32]
    "env"."ext_hashing_twox_128_version_1": [I64] -> [I32]
    "env"."ext_hashing_twox_64_version_1": [I64] -> [I32]
    "env"."ext_storage_append_version_1": [I64, I64] -> []
    "env"."ext_storage_clear_version_1": [I64] -> []
    "env"."ext_storage_clear_prefix_version_2": [I64, I64] -> [I64]
    "env"."ext_storage_exists_version_1": [I64] -> [I32]
    "env"."ext_storage_get_version_1": [I64] -> [I64]
    "env"."ext_storage_read_version_1": [I64, I64, I32] -> [I64]
    "env"."ext_storage_root_version_2": [I32] -> [I64]
    "env"."ext_storage_set_version_1": [I64, I64] -> []
    "env"."ext_storage_commit_transaction_version_1": [] -> []
    "env"."ext_storage_rollback_transaction_version_1": [] -> []
    "env"."ext_storage_start_transaction_version_1": [] -> []
    "env"."ext_trie_blake2_256_ordered_root_version_2": [I64, I32] -> [I32]
    "env"."ext_logging_log_version_1": [I32, I64, I64] -> []
  Memories:
    "env"."memory": not shared (20 pages..)
  Tables:
  Globals:
Exports:
  Functions:
    "Core_version": [I32, I32] -> [I64]
    "Core_initialize_block": [I32, I32] -> [I64]
    "Core_execute_block": [I32, I32] -> [I64]
    "BlockBuilder_apply_extrinsic": [I32, I32] -> [I64]
    "BlockBuilder_finalize_block": [I32, I32] -> [I64]
    "BlockBuilder_inherent_extrinsics": [I32, I32] -> [I64]
    "BlockBuilder_check_inherents": [I32, I32] -> [I64]
    "TaggedTransactionQueue_validate_transaction": [I32, I32] -> [I64]
    "AuraApi_slot_duration": [I32, I32] -> [I64]
    "AuraApi_authorities": [I32, I32] -> [I64]
    "AccountNonceApi_account_nonce": [I32, I32] -> [I64]
    "TransactionPaymentApi_query_info": [I32, I32] -> [I64]
    "TransactionPaymentApi_query_fee_details": [I32, I32] -> [I64]
    "TransactionPaymentCallApi_query_call_info": [I32, I32] -> [I64]
    "TransactionPaymentCallApi_query_call_fee_details": [I32, I32] -> [I64]
    "Metadata_metadata": [I32, I32] -> [I64]
    "Metadata_metadata_at_version": [I32, I32] -> [I64]
    "Metadata_metadata_versions": [I32, I32] -> [I64]
    "SessionKeys_generate_session_keys": [I32, I32] -> [I64]
    "SessionKeys_decode_session_keys": [I32, I32] -> [I64]
    "GrandpaApi_grandpa_authorities": [I32, I32] -> [I64]
    "OffchainWorkerApi_offchain_worker": [I32, I32] -> [I64]
    "GenesisBuilder_create_default_config": [I32, I32] -> [I64]
    "GenesisBuilder_build_config": [I32, I32] -> [I64]
  Memories:
  Tables:
    "__indirect_function_table": FuncRef (105..105)
  Globals:
    "__heap_base": I32 (constant)
    "__data_end": I32 (constant)
```

To inspect the WASM in more detail, and view the actual memory segments, you can
install [wabt](https://github.com/WebAssembly/wabt).

```bash
wasm-objdump -x build/runtime.wasm
```