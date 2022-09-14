# Substrate Wasm interface

Wasm module compiled from `rust/substrate`.

```
Type: wasm
Size: 2.5 MB
Imports:
  Functions:
    "env"."ext_logging_log_version_1": [I32, I64, I64] -> []
    "env"."ext_sandbox_instance_teardown_version_1": [I32] -> []
    "env"."ext_sandbox_instantiate_version_1": [I32, I64, I64, I32] -> [I32]
    "env"."ext_sandbox_invoke_version_1": [I32, I64, I64, I32, I32, I32] -> [I32]
    "env"."ext_sandbox_memory_get_version_1": [I32, I32, I32, I32] -> [I32]
    "env"."ext_sandbox_memory_new_version_1": [I32, I32] -> [I32]
    "env"."ext_sandbox_memory_set_version_1": [I32, I32, I32, I32] -> [I32]
    "env"."ext_sandbox_memory_teardown_version_1": [I32] -> []
    "env"."ext_crypto_ed25519_generate_version_1": [I32, I64] -> [I32]
    "env"."ext_crypto_ed25519_verify_version_1": [I32, I64, I32] -> [I32]
    "env"."ext_crypto_finish_batch_verify_version_1": [] -> [I32]
    "env"."ext_crypto_secp256k1_ecdsa_recover_compressed_version_1": [I32, I32] -> [I64]
    "env"."ext_crypto_sr25519_generate_version_1": [I32, I64] -> [I32]
    "env"."ext_crypto_sr25519_public_keys_version_1": [I32] -> [I64]
    "env"."ext_crypto_sr25519_sign_version_1": [I32, I32, I64] -> [I64]
    "env"."ext_crypto_sr25519_verify_version_2": [I32, I64, I32] -> [I32]
    "env"."ext_crypto_start_batch_verify_version_1": [] -> []
    "env"."ext_trie_blake2_256_ordered_root_version_1": [I64] -> [I32]
    "env"."ext_misc_print_hex_version_1": [I64] -> []
    "env"."ext_misc_print_num_version_1": [I64] -> []
    "env"."ext_misc_print_utf8_version_1": [I64] -> []
    "env"."ext_misc_runtime_version_version_1": [I64] -> [I64]
    "env"."ext_default_child_storage_clear_version_1": [I64, I64] -> []
    "env"."ext_default_child_storage_get_version_1": [I64, I64] -> [I64]
    "env"."ext_default_child_storage_root_version_1": [I64] -> [I64]
    "env"."ext_default_child_storage_set_version_1": [I64, I64, I64] -> []
    "env"."ext_default_child_storage_storage_kill_version_1": [I64] -> []
    "env"."ext_allocator_free_version_1": [I32] -> []
    "env"."ext_allocator_malloc_version_1": [I32] -> [I32]
    "env"."ext_hashing_blake2_128_version_1": [I64] -> [I32]
    "env"."ext_hashing_blake2_256_version_1": [I64] -> [I32]
    "env"."ext_hashing_keccak_256_version_1": [I64] -> [I32]
    "env"."ext_hashing_sha2_256_version_1": [I64] -> [I32]
    "env"."ext_hashing_twox_128_version_1": [I64] -> [I32]
    "env"."ext_hashing_twox_64_version_1": [I64] -> [I32]
    "env"."ext_offchain_is_validator_version_1": [] -> [I32]
    "env"."ext_offchain_local_storage_compare_and_set_version_1": [I32, I64, I64, I64] -> [I32]
    "env"."ext_offchain_local_storage_get_version_1": [I32, I64] -> [I64]
    "env"."ext_offchain_local_storage_set_version_1": [I32, I64, I64] -> []
    "env"."ext_offchain_network_state_version_1": [] -> [I64]
    "env"."ext_offchain_random_seed_version_1": [] -> [I32]
    "env"."ext_offchain_submit_transaction_version_1": [I64] -> [I64]
    "env"."ext_storage_append_version_1": [I64, I64] -> []
    "env"."ext_storage_changes_root_version_1": [I64] -> [I64]
    "env"."ext_storage_clear_version_1": [I64] -> []
    "env"."ext_storage_clear_prefix_version_1": [I64] -> []
    "env"."ext_storage_commit_transaction_version_1": [] -> []
    "env"."ext_storage_get_version_1": [I64] -> [I64]
    "env"."ext_storage_next_key_version_1": [I64] -> [I64]
    "env"."ext_storage_read_version_1": [I64, I64, I32] -> [I64]
    "env"."ext_storage_rollback_transaction_version_1": [] -> []
    "env"."ext_storage_root_version_1": [] -> [I64]
    "env"."ext_storage_set_version_1": [I64, I64] -> []
    "env"."ext_storage_start_transaction_version_1": [] -> []
    "env"."ext_offchain_index_set_version_1": [I64, I64] -> []
  Memories:
    "env"."memory": not shared (20 pages..)
  Tables:
  Globals:
Exports:
  Functions:
    "Core_version": [I32, I32] -> [I64]
    "Core_execute_block": [I32, I32] -> [I64]
    "Core_initialize_block": [I32, I32] -> [I64]
    "Metadata_metadata": [I32, I32] -> [I64]
    "BlockBuilder_apply_extrinsic": [I32, I32] -> [I64]
    "BlockBuilder_finalize_block": [I32, I32] -> [I64]
    "BlockBuilder_inherent_extrinsics": [I32, I32] -> [I64]
    "BlockBuilder_check_inherents": [I32, I32] -> [I64]
    "BlockBuilder_random_seed": [I32, I32] -> [I64]
    "TaggedTransactionQueue_validate_transaction": [I32, I32] -> [I64]
    "OffchainWorkerApi_offchain_worker": [I32, I32] -> [I64]
    "GrandpaApi_grandpa_authorities": [I32, I32] -> [I64]
    "GrandpaApi_submit_report_equivocation_unsigned_extrinsic": [I32, I32] -> [I64]
    "GrandpaApi_generate_key_ownership_proof": [I32, I32] -> [I64]
    "BabeApi_configuration": [I32, I32] -> [I64]
    "BabeApi_current_epoch_start": [I32, I32] -> [I64]
    "BabeApi_generate_key_ownership_proof": [I32, I32] -> [I64]
    "BabeApi_submit_report_equivocation_unsigned_extrinsic": [I32, I32] -> [I64]
    "AuthorityDiscoveryApi_authorities": [I32, I32] -> [I64]
    "AccountNonceApi_account_nonce": [I32, I32] -> [I64]
    "ContractsApi_call": [I32, I32] -> [I64]
    "ContractsApi_get_storage": [I32, I32] -> [I64]
    "ContractsApi_rent_projection": [I32, I32] -> [I64]
    "TransactionPaymentApi_query_info": [I32, I32] -> [I64]
    "SessionKeys_generate_session_keys": [I32, I32] -> [I64]
    "SessionKeys_decode_session_keys": [I32, I32] -> [I64]
  Memories:
  Tables:
    "__indirect_function_table": FuncRef (352..352)
  Globals:
    "__data_end": I32 (constant)
    "__heap_base": I32 (constant)
```
