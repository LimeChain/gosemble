# Milestone Deliverables

# M0. Auxiliaries, Supporting Tooling/Infra

There are a set of processes that will be performed/established during the whole development of the framework. Those processes will ensure the quality of the framework, its maintainability and usability.

1. CLI
    1. During the development of the framework, we will introduce a CLI that will help us and developers to get started with the framework. The CLI will be enhanced with time to support more features as the framework matures. Features include, but are not limited to:
        1. `init` → Once executed, it will create the appropriate folder structure (runtime, tests). The core modules of the Runtime will be placed in the `/runtime` folder, along with the core functionality such as the Runtime API, serialisation and deserialisation. This command will be used to bootstrap new Runtimes.
        2. `compile` → Once executed, the command will compile the Go Runtime and produce the WASM binary. During the compilation, we might need to generate code (similar to Rust macros). This will be more evident once the framework supports Metadata API.
2. Integration Tests
    1. During the development of the framework, we will set up integration tests, which will require a set of WASM APIs to be introduced in various modules. The integration tests will run inside a test bed/execution environment which must provide an easy way to “mock” data in order to write the integration tests in a timely manner.
    2. A separate project can be established based on this work → introduction for execution specs for Runtime modules that can be executed against both Substrate based and Go-based modules
3. CICD Pipeline
    1. GitHub actions are to be used to set up a pipeline that will `compile`, run `unit-test` and `integration-test`s

# M1. Runtime APIs, Storage, Core Primitives & Block Production

1. Setup sandbox environment to aid the development and testing
    1. Compile the runtime to WASM with [tinygo](https://tinygo.org/).
    2. run the compiled WASM inside [gossamer](https://github.com/ChainSafe/gossamer) VM instance.
2. Specify the API interfaces that need to be exposed by the WASM runtime. Reuse the core types and serialization logic for exchanging data that is already implemented in [gossamer](https://github.com/ChainSafe/gossamer) packages.
    
    ```go
    import (
    	"github.com/ChainSafe/gossamer/dot/types"
    	"github.com/ChainSafe/gossamer/lib/scale"
    )
    
    /*
    	SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
    	dataPtr - i32 pointer to the memory location.
    	dataLen - i32 length (in bytes) of the encoded arguments.
    	returns a pointer-size to the SCALE-encoded (version types.VersionData) data.
    */
    //export Core_version
    func Core_version(dataPtr int32, dataLen int32) uint64 {
    	version := runtime.VersionData{ /* ... */ }
    	scaleEncVersion, _ := version.Encode()
    	return utils.BytesToPointerAndSize(scaleEncVersion)
    }
    
    /*
    	SCALE encoded arguments (block types.Block) allocated in the Wasm VM memory, passed as:
    	dataPtr - i32 pointer to the memory location.
    	dataLen - i32 length (in bytes) of the encoded arguments.
    */
    //export Core_execute_block
    func ExecuteBlock(dataPtr int32, dataLen int32) {
    	// ...
    }
    
    /*
    	SCALE encoded arguments (header *types.Header) allocated in the Wasm VM memory, passed as:
    	dataPtr - i32 pointer to the memory location.
    	dataLen - i32 length (in bytes) of the encoded arguments.
    */
    //export Core_initialize_block
    func InitializeBlock(dataPtr int32, dataLen int32) {
    	// ...
    }
    ```
    
3. Make the necessary abstractions to utilize the interfaces exposed by the host environment. Runtime should be able to manipulate the storage and memory provided by the host.
    1. Implement the ability to
        1. Set storage at a given key using the exposed `ext_storage_set` function
        2. Retrieve a storage value by a given key using the exposed `ext_storage_get` function
        3. Read a storage value by a given key and offset params using the exposed `ext_storage_read` function.
        4. Clear the storage of the given key and its value using the exposed `ext_storage_clear` function.
        5. Check whether a given key exists in the storage using the exposed `ext_storage_exists` function.
    2. Implement the Accounts/Balances module such that:
        1. We are able to read the balance of an account
        2. We are able to change the balance of an account
    3. Be able to read genesis account balances
    4. Develop Unit tests for `Storage Module` read & write operations
    5. Deliver inline documentation of the code, a README file describing how one can run the Storage Module tests
    6. Export the package to [https://pkg.go.dev/](https://pkg.go.dev/)
    
    ```go
    // Imported functions available to the runtime.
    
    // Storage
    func ext_storage_set_version_1(key int64, value int64)
    func ext_storage_get_version_1(key int64) int64
    func ext_storage_read_version_1(key int64, value_out int64, offset int32) int64
    func ext_storage_clear_version_1(key_data int64)
    func ext_storage_exists_version_1(key_data int64) int32
    func ext_storage_clear_prefix_version_1(prefix int64)
    func ext_storage_clear_prefix_version_2(prefix uint64, limit uint64) uint64
    func ext_storage_append_version_1(key uint64, value uint64)
    func ext_storage_root_version_1() int64
    func ext_storage_root_version_2(version uint32) uint32
    func ext_storage_changes_root_version_1(parent_hash int64) int64
    func ext_storage_next_key_version_1(key int64) int64
    func ext_storage_start_transaction_version_1()
    func ext_storage_rollback_transaction_version_1()
    func ext_storage_commit_transaction_version_1()
    
    // Memory
    func ext_allocator_malloc_version_1(size int32) int32
    func ext_allocator_free(ptr int32)
    ```
    
4. Implement block production logic
    1. Implementing `initialize_block`, `execute_block`, `finalise_block`, `apply_extrinsic`, `inherent-extrinsics`, `random_seed` and `validate_transaction`
    2. Develop Unit tests for the Block Production module
    3. Deliver inline documentation of the code, a README file describing how one can run the tests
    4. Export the package to [https://pkg.go.dev/](https://pkg.go.dev/)
    
    ```go
    import (
    	"github.com/ChainSafe/gossamer/dot/types"
    )
    
    type BlockBuilder interface {
    	/*
    		SCALE encoded arguments (extrinsic types.Extrinsic) allocated in the Wasm VM memory, passed as:
    		dataPtr - i32 pointer to the memory location.
    		dataLen - i32 length (in bytes) of the encoded arguments.
    		returns a pointer-size to the SCALE-encoded ([]byte) data.
    	*/
    	//export BlockBuilder_apply_extrinsic
    	ApplyExtrinsic(dataPtr int32, dataLen int32) int64
    
    	/*
    		SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
    		dataPtr - i32 pointer to the memory location.
    		dataLen - i32 length (in bytes) of the encoded arguments.
    		returns a pointer-size to the SCALE-encoded (types.Header) data.
    	*/
    	//export BlockBuilder_finalize_block
    	FinalizeBlock(dataPtr int32, dataLen int32) int64
    
    	/*
    		SCALE encoded arguments (data types.InherentsData) allocated in the Wasm VM memory, passed as:
    		dataPtr - i32 pointer to the memory location.
    		dataLen - i32 length (in bytes) of the encoded arguments.
    		returns a pointer-size to the SCALE-encoded ([]types.Extrinsic) data.
    	*/
    	//export BlockBuilder_inherent_extrinisics
    	InherentExtrinisics(dataPtr int32, dataLen int32) int64
    
    	/*
    		SCALE encoded arguments (block types.Block, data types.InherentsData) allocated in the Wasm VM memory, passed as:
    		dataPtr - i32 pointer to the memory location.
    		dataLen - i32 length (in bytes) of the encoded arguments.
    		returns a pointer-size to the SCALE-encoded ([]byte) data.
    	*/
    	//export BlockBuilder_check_inherents
    	CheckInherents(dataPtr int32, dataLen int32) int64
    
    	/*
    		SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
    		dataPtr - i32 pointer to the memory location.
    		dataLen - i32 length (in bytes) of the encoded arguments.
    		returns a pointer-size to the SCALE-encoded ([32]byte) data.
    	*/
    	//export BlockBuilder_random_seed
    	RandomSeed(dataPtr int32, dataLen int32) int64
    }
    ```
    
5. Implement the transaction validation logic
    
    ```go
    // TaggedTransactionQueue handles validating transactions in the transaction queue.
    type TaggedTransactionQueue interface {
      //export TaggeTransactionQueue_validate_transaction
      ValidateTransaction(tx []byte) TransactionValidity
    }
    
    // TransactionValidity stores information concerning the validity of a transaction.
    type TransactionValidity struct {
      // Priority of the transaction.
    	//
    	// Priority determines the ordering of two transactions that have all
    	// their dependencies (required tags) satisfied.
    	Priority uint64,
    
    	// Transaction dependencies
    	//
    	// A non-empty list signifies that some other transactions which provide
    	// given tags are required to be included before that one.
      Requires [][]byte,
    
    	// Provided tags
    	//
    	// A list of tags this transaction provides. Successfully importing the transaction
    	// will enable other transactions that depend on (require) those tags to be included as well.
    	// Provided and required tags allow Substrate to build a dependency graph of transactions
    	// and import them in the right (linear) order.
    	Provides: [][]byte,
    
    	// Transaction longevity
    	//
    	// Longevity describes minimum number of blocks the validity is correct.
    	// After this period transaction should be removed from the pool or revalidated.
    	Longevity uint64,
    
    	// A flag indicating if the transaction should be propagated to other peers.
    	//
    	// By setting `false` here the transaction will still be considered for
    	// including in blocks that are authored on the current node, but will
    	// never be sent to other peers.
    	Propagate: bool,
    }
    ```
    
6. Make the `runtime_version` and `block_time` configurable
7. Add Unit and Integration tests to the implemented modules

Related Substrate Modules:

- `sp-api` - Substrate runtime API. It is required that each runtime implements at least the Core runtime API. This runtime API provides all the core functions that Substrate expects from a runtime.
- `sp-runtime-interface` - Substrate runtime interface.
- `sp-wasm-interface` - Types and traits for interfacing between the host and the wasm runtime.
- `sp-runtime` - Runtime Modules shared primitive types.
- `sp-version` - Substrate runtime version.
- `sp-state-machine` - Substrate state machine implementation.
- `sp-tracing` - Substrate tracing primitives and macros.
- `sp-std` - Lowest-abstraction level for the Substrate runtime: just exports useful primitives from std or client/alloc to be used with any code that depends on the runtime.
- `sp-core` - Shareable Substrate types.
- `sp-arithmetic` - Minimal fixed-point arithmetic primitives and types for runtime.
- `sp-core-hashing` - Hashing Functions.
- `sp-io` - I/O host interface for substrate runtime. Substrate runtime standard library as compiled when linked with Rust’s standard library.
- `sp-externalities` - Substrate externalities abstraction.
- `sp-storage` - Primitive types for storage-related stuff.
- `sp-trie` - Utility functions to interact with Substrate’s Base-16 Modified Merkle Patricia tree (“trie”).
- `sp-keystore` - Keystore traits.
- `frame-system` - The System pallet provides low-level access to core types and cross-cutting utilities. It acts as the base layer for other pallets to interact with the Substrate framework components.
- `frame-support` - Support code for the runtime.
- `frame-metadata` - Decodable variant of the RuntimeMetadata.
- `frame-executive` - The Executive module acts as the orchestration layer for the runtime. It dispatches incoming extrinsic calls to the respective modules in the runtime.
    
    

# M2. Timestamp & Aura

Implement the Aura (Proof-of-Authority) consensus engine responsible for block authorship

1. Implement the functionality for loading and setting the `authority list` defined in the genesis state.
2. Implement Aura Runtime API `slot_duration` and `check_inherent` functions
3. Develop Unit tests for the Aura module API.
4. Deliver both inline documentation of the code and README tutorial describing how one can run the Aura Runtime tests.
5. Export the package to [https://pkg.go.dev/](https://pkg.go.dev/)

```go
type AuraApi interface {
		//export slot_duration
		func SlotDuration() SlotDuration

		//export authorities
		func Authorities() []AuthorityId
	}
```

Related Substrate Modules

- `pallet-aura`
- `sp-consensus-aura` - Primitives for Aura.
- `sp-consensus` - Common utilities for building and using consensus engines in Substrate.
- `sp-consensus-slots` - Primitives for slots-based consensus engines.
- `sp-timestamp` - Substrate core types and inherent for timestamps.
- `sp-inherents` - Inherent extrinsics are extrinsics that are inherently added to each block.
- `pallet-timestamp` - The Timestamp pallet provides functionality to get and set the on-chain time.

# M3. Transaction API, Fees & Balances

1. Implement the state transition logic for sending balances between accounts
2. Implement Transaction Payment module to enforce fees on transactions
3. Add unit and integration tests
4. Add inline documentation and README file describing how one can run the tests for the module
5. Export the packages to [https://pkg.go.dev/](https://pkg.go.dev/)

For the successful completion of this, milestone we must be able to compile and run PoC Runtime successfully. The runtime must be able to:

1. Define genesis state and account balances
2. Run its Aura consensus
3. Sync/Initialize/Execute blocks
4. Being able to process Account Balance transfers (Extrinsics)

# M4. Metadata API

The Metadata module enables anyone to easily interact with the Runtime using the [https://polkadot.js.org/apps/](https://polkadot.js.org/apps/) interface. For the successful completion of this milestone we can define the following criteria:

1. Implement Metadata and its API for the already built modules (System, Aura, Balances, Timestamp)
2. Add inline documentation
3. Implement Unit Tests and add README file describing how one can build and run the tests for the module.
4. Export the packages to [https://pkg.go.dev/](https://pkg.go.dev/)