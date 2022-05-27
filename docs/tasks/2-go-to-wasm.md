# Go to Wasm

## Notes
*Which WebAssembly specification (1.0/2.0) is the Polkadot runtime targeted at?*
* The Wasm runtime targets [WebAssembly MVP](https://github.com/WebAssembly/design/blob/main/MVP.md), i.e. without any extensions enabled.

*Does it target WASI, as it should be run outside the browser?*
* No, it does not. WASI is a standard that seeks to provide a system-level API comparable to an OS. Substrate/Polkadot is no OS and it does not support files, networking, or a major part of other things provided by WASI. Instead, a more domain-specific API is provided.

*Is it possible to use language with automatic memory management?*
* Theoretically, it might be possible, but the support would be limited, performance might be unsatisfactory and the toolchain would need to polyfill the GC. To support an automatic memory management, the [GC proposal](https://github.com/WebAssembly/gc/blob/main/proposals/gc/Overview.md) would be handy. But the Wasm runtime supports only WebAssembly MVP currently, also the GC proposal is under development. It is not yet clear if Polkadot will be able to leverage the GC proposal. Potential problems include determinism (is there anything in GC that causes ND? Can it be tamed efficiently?) and safety (Is it possible for a host to limit the resource consumption reliably and deterministically?).

The [automatic reference counting](https://docs.swift.org/swift-book/LanguageGuide/AutomaticReferenceCounting.html) used by Swift can work just well as of now.

*External memory management*
* The beginning of this heap is marked by the `__heap_base` symbol exported by the Wasm module.
No memory should be allocated below that address, to avoid clashes with the stack and data section.
```
 -------------------------------------------------------------------
| Data     |          <- Stack | Heap ->                            |
 -------------------------------------------------------------------
0     __data_end          __heap_base                           max memory
```

## Tasks

1. Utilise `TinyGo` to compile the `Go` runtime as Wasm module matching the Polkadot spec.

* [x] Create build script `/scripts/build.sh`.

2. Expose the expected API from the Wasm module `/runtime/runtime.go`.

* [x] Exported API functions.
* [x] Imported host functions.
* [ ] Exported linker specific globals.
* [ ] Imported host memory.

1. Setup dev environment `/dev`.

* [x] Setup host runtime to execute the compiled Wasm module inside.
* [x] Import the host provided functions used inside the Wasm module.
* [x] Read/write from/to the host/Wasm's shared memory by using a pointer-size.
* [ ] Import the host provided memory inside the Wasm module.

## Issues

### Toolchain support for Wasm outside the browser

* The official Go compiler does not support Wasm for non-browser environments [read more](https://github.com/golang/go/issues/31105) [read more](https://substrate.stackexchange.com/questions/60/what-is-gossamer-and-how-does-it-compare-to-substrate/89#89). The only options is to use the `TinyGo` compiler.

### Standard library support

* The `reflect` package is not fully supported by the `TinyGo` compiler [read more](https://github.com/tinygo-org/tinygo/pull/2640). The core primitives and SCALE serialization logic that we intended to reuse from [gossamer](https://github.com/ChainSafe/gossamer) all rely on the `reflect` package.

### External memory allocator

* By specification, the Wasm module does not include a memory allocator, it should import memory from the host and rely on host imported functions for all heap allocations. `TinyGo` has GC and manages its memory by itself. So it can't work directly on systems where the host wants to manage the memory. It exports allocation functions while `Rust`'s toolchain does not. Also `Rust` eagerly collects memory before returning from a Wasm function while `TinyGo` does not. Wasm compiled from Rust and Tinygo differ in both terms of exports and runtime behavior around allocation, because there is no WebAssembly specification for it.
```sh
# We currently don't support external memory. We used to have an "extalloc" memory allocator, but it was rather complicated and had some bugs.

acdaa723 runtime: fix compilation errors when using gc.extalloc
959442dc unix: use conservative GC by default
747336f0 runtime: remove extalloc
38b14706 internal/task: fix two missed instances of extalloc
```
[read more](https://github.com/golang/go/issues/13761)

* The runtime is expected to expose `__heap_base` global [read more](https://github.com/tinygo-org/tinygo/issues/2045), but `TinyGo` doesn't support that out of the box.

**Hacks:**
* Hardcode the same value for `__heap_base` inside the host and the runtime module.
* Inside the Wasm module, the GC could just allocate a large amount of memory and work with that.

**Solutions:**
1. Implement external memory allocator and add support for it in `TinyGo`. Extend `TinyGo` to support importing memory and exporting linker specific globals `__heap_base`, `__data_end`.

2. Use different language that supports WebAssembly MVP (no GC support) and LLVM as compilation toolchain (C, C++, Zig, Swift. Rust, AssemblyScript are already used in implementations).
[read more](https://www.fermyon.com/wasm-languages/webassembly-language-support)
[read more](https://github.com/appcypher/awesome-wasm-langs).


# Further Research
1. Go internals, runtime, memory allocation, GC
* get a deep understanding of how internals, runtime, GC and memory allocation in Go works.

2. WebAssembly GC proposal
* thoroughly research the GC proposal for WebAssembly, such as its design and progress by other entities so far.

3. Build a PoC manual memory allocator via FFI
* `syscall.Mmap` from `syscall` package.
* `malloc/calloc`, `free` via `cgo`.
* `jemalloc` via `cgo`.

4. Research through `TinyGo` or alternative compiler toolchain for the following addition of:
* compilation from Go to LLVM IR
* compilation from  LLVM IR to Wasm
* toolchain support for Wasm MVP

5. Propose an extension logic for `TinyGo` or build alternative compiler toolchain, based on 4.

We expect the research to give us vast knowledge of the missing pieces that Go needs, so that it becomes compatible for runtime implementation.
