# Research on Go to Wasm Toolchain

The idea of writing Polkadot Runtimes in Go is exciting, mainly because of Go's simplicity and automatic memory management. But Polkadot's spec describes that memory should be managed by the Host, which defeats one of Go's main selling points. Although the language specification doesn't mention how it should manage its memory, the Go community recognizes it as a language with GC, and anything else would be more like a different language with a similar syntax. So aren't we setting up ourselves for failure with Go right from the beginning?

Questions we are gonna try to answer:
* How Go internals work - compiler, runtime, GC
* Steps required to start implementing new runtime with GC and external allocator

## Polkadot's Wasm
**WebAssembly specification version**

The Wasm runtime module targets [WebAssembly MVP](https://github.com/WebAssembly/design/blob/main/MVP.md) without any extensions enabled and with domain-specific API.
* [x] module exports API functions (the business logic).
* [x] module imports Host provided functions (`ext_allocator_free/ext_allocator_malloc`).
* [ ] module imports Host provided memory.
* [ ] module exports linker specific globals (`__heap_base`).
* [ ] ??? `__indirect_function_table` and `__data_end` are not in the spec

**WASI interface requirements**

Polkadot is non-browser environment, but it is not an OS and doesn't seeks to provide a system-level API comparable to an OS like files, networking, or a major part of other things provided by WASI.

**Language with an automatic memory management**

Theoretically, it might be possible, but the support would be limited, performance might be unsatisfactory and the toolchain would need to polyfill the GC. To support an automatic memory management, the [GC proposal](https://github.com/WebAssembly/gc/blob/main/proposals/gc/Overview.md) might be handy. But the Wasm runtime supports only WebAssembly MVP currently, also the GC proposal is under development and it is not yet clear if Polkadot will be able to leverage the GC proposal. Potential problems include determinism (is there anything in GC that causes ND? Can it be tamed efficiently?) and safety (Is it possible for a host to limit the resource consumption reliably and deterministically?).

**External memory management**

The beginning of this heap is marked by the `__heap_base` symbol exported by the Wasm module.
No memory should be allocated below that address, to avoid clashes with the stack and data section.
```
 _________________________________________________________________________________________
| HOST                                                                                    |                             
|                                                     ext_malloc/ext_free                 |
|      ____________________________________________________|_________________________     |
|     | WASM                                               | (imported)              |    |
|     |                                                    |                         |    |
|     |             ______________________________________\|/______________________  |    |
|     | (imported) |          |                  |                                 | |    |
|  MEMORY  ----->  | Data     |         <- Stack | Heap ->                         | |    |
|     |            |__________|__________________|_________________________________| |    |
|     |            0     __data_end         __heap_base         /|\      max memory  |    |
|     |                                                          |                   |    |
|     |                                                          GC                  |    |
|     |______________________________________________________________________________|    |
|                                                                                         |
|_________________________________________________________________________________________|

```

**Support of concurrency**

The runtime executes in serial, the parallelism is accomplished through a network of parallel running chains (Parachains). It is primarily because creating a semantic for deterministic parallel executing is really difficult in general.

### Technical challenges

**Toolchain support for Wasm outside the browser**
* The official Go compiler does not support Wasm for non-browser environments [1](https://github.com/golang/go/issues/31105), [2](https://substrate.stackexchange.com/questions/60/what-is-gossamer-and-how-does-it-compare-to-substrate/89#89), only Wasm with browser specific API. The only options that supports that is TinyGo.

**Standard library support**
* The `reflect` package is not fully supported by TinyGo [1](https://github.com/tinygo-org/tinygo/pull/2640). The core primitives and SCALE serialization logic that we intended to reuse from [gossamer](https://github.com/ChainSafe/gossamer) all rely on the `reflect` package.

**External memory allocator**
* According to the Polkadot specification, the Wasm module does not include a memory allocator, it imports memory from the Host and relies on Host imported functions for all heap allocations (ext_allocator_malloc/ext_allocator_free). TinyGo implements simple GC and manages its memory by itself, contrary to specification. So it can't work directly on systems where the host wants to manage the memory. It used to have an `extalloc` external memory allocator, but has been rather complicated and buggy [1](https://github.com/golang/go/issues/13761).

**Linker globals**
* The runtime is expected to expose `__heap_base` global [1](https://github.com/tinygo-org/tinygo/issues/2045), but TinyGo doesn't support that out of the box.

## Wasm/Go/TinyGo compiler, internals, runtime, GC, memory allocation

## WebAssembly 
---

WebAssembly/gc proposal most likely won't help [1](https://github.com/WebAssembly/gc/issues/59).

## Go
---

### Compiler

The default compiler is `gc`. There are also `gccgo` which uses the GCC back-end and `gollvm` which uses the LLVM infrastructure (somewhat less mature).

```
// Frontend Compiler (lexing, parsing, typechecking)

CODE -> tokenizer/lexer/scanner (lexical analysis) -> TOKENS STREAM
TOKENS STREAM -> go/parser (syntactic analysis) -> AST (annotated)
AST -> go/types (type check) -> AST (with types)
AST (with types) -> golang.org/x/tools/go/ssa (convert) -> SSA

// Optimization

* variable capturing, inlining, escape analysis, closure rewriting, walk

SSA (higher level with Go-specific constructs like interfaces and goroutines) -> convert (tinygo) -> LLVM IR
LLVM IR -> optimize (llvm) -> LLVM IR (optimize by a mixture of handpicked LLVM optimization passes, TinyGo-specific optimizations ()

// Backend Compiler

LLVM IR (optimized) -> convert (llvm) -> MACHINE CODE
MACHINE CODE -> convert (llvm) -> OBJECT FILE
OBJECT FILE -> link (llvm) -> EXECUTABLE
```

### Runtime

Implements GC, scheduler etc. included in every Go program. Contains a lot of type information at runtime.

*Memory Management*

* tracing garbage collector (not reference counting)
* hybrid stop-the-world/concurrent collector (stop-the-world part limited by a 10ms deadline)
* CPU cores dedicated to running the concurrent collector
* tri-color mark-and-sweep algorithm
* non-generational
* non-compacting
* fully precise
* incurs a small cost if the program is moving pointers around
* lower latency, but most likely also lower throughput

*Concurrency*

*Parallelism*

*Packages*

## TinyGo (took ~1.5-2 years to get something working)
---

It is a subset of Go with very different goals from the standard Go. It is a new compiler and runtime aimed to support many different small embedded devices with a single processor core that require certain optimizations mostly toward size.

*Goals*
* Have very small binary sizes.
* Support for most common microcontroller boards.
* Be usable on the web using WebAssembly.
* Good CGo support, with no more overhead than a regular function call.
* Support most standard library packages and compile most Go code without modification.

*Non-goals*
* Using more than one core.
* Be efficient while using zillions of goroutines. However, good goroutine support is certainly a goal.
* Be as fast as `gc`. However, LLVM will probably be better at optimizing certain things so TinyGo might actually turn out to be faster for number crunching.
* Be able to compile every Go program out there.

### Compiler

The compiler uses (mostly) the standard library to parse Go programs and LLVM to optimize the code and generate machine code for the target architecture.

**Pipeline (Lowering)**
```
// Frontend Compiler (lexing, parsing, typechecking)

CODE -> parse (go/parser) -> AST
AST -> type check (go/types) -> AST (with types)
AST (with types) -> convert (golang.org/x/tools/go/ssa) -> SSA

// Optimization

SSA (higher level with Go-specific constructs like interfaces and goroutines) -> convert (tinygo) -> LLVM IR
LLVM IR -> optimize (llvm) -> LLVM IR (optimize by a mixture of handpicked LLVM optimization passes, TinyGo-specific optimizations (escape analysis, string-to-[]byte optimizations, etc.) and custom lowering.)

// Backend Compiler

LLVM IR (optimized) -> convert (llvm) -> MACHINE CODE
MACHINE CODE -> convert (llvm) -> OBJECT FILE
OBJECT FILE -> link (llvm) -> EXECUTABLE
```

* root: contains the command line interface for the tinygo command and all its subcommands
* builder: orchestrates the build
* loader: loads and typechecks the code, and produces an AST
* compiler: the compiler itself, makes little attempt at optimizing code
* interp: tries to run package initializers at compile time as far as possible
* transform: implements various optimizations necessary to produce working and efficient code

### Runtime

The runtime is written from scratch, optimized for size instead of speed and re-implements some compiler intrinsics and packages:

* Startup code (device specific runtime initialization of memory, timers)
* GC (copied from micropython, super simple, optimized for size, runs the GC once it runs out of memory)
* Memory Allocator
* Goroutines scheduler
* Channels
* Time handling (every chip has a clock)
* Hashmap implementation (optimized for size instead of speed)
* Runtime functions like maps, slice append, defer, ... (various things that are not implemented in the compiler)
* Operations on strings
* `sync` & `reflect` package (strongly connected to the runtime)

*Basic features*

All basic types, slices, all regular control flow including switch, closures and bound methods are supported, `defer` keyword is almost entirely supported, with the exception of deferring some builtin functions, interfaces are quite stable and should work well in almost all cases. Type switches and type asserts are also supported, as well as calling methods on interfaces. The only exception is comparing two interface values.
Maps are usable but not complete. You can use any type as a value, but only some types are acceptable as map keys (strings, integers, pointers, and structs/arrays that contain only these types). Also, they have not been optimized for performance and will cause linear lookup times in some cases.

*Memory Management*

Heap with a garbage collector. While not directly a language feature, garbage collection is important for most Go programs to make sure their memory usage stays in reasonable bounds.

Garbage collection is currently supported on all platforms, although it works best on 32-bit chips. A simple conservative mark-sweep collector is used that will trigger a collection cycle when the heap runs out (that is fixed at compile time) or when requested manually using `runtime.GC()`. Some other collector designs are used for other targets, TinyGo will automatically pick a good GC for a given target.

Careful design may avoid memory allocations in main loops. You may want to compile with `-print-allocs=.` to find out where allocations happen and why they happen.

*Concurrency*

Goroutines and channels work for the most part, for platforms such as WebAssembly the support is a bit more limited (calling a blocking function may for example allocate heap memory).

*Parallelism*

Single core only.

*Packages*

Only partial support of the `reflect` package (most common types like numbers, strings, and structs are supported). Standard library relies on reflection.
Many features of Cgo are still unsupported (#cgo statements are only partially supported).

---

## PoC of an alternative compiler + runtime

TinyGo design decisions are based on optimizations around small embedded devices, which might not be always suitable or required in Wasm. Also, it will be difficult to support custom and non-standard APIs (import memory + export linker-specific globals) in the long run, the TinyGo core contributors are a bit against supporting things that are not standardized as is the case with Polkadot Wasm's custom API.
It will be best to implement a solution similar to TinyGo (compiler + runtime + external memory allocator) only targeting Wasm.

1. Use TinyGo as a starting point setup + Docker.
2. Frontend-compiler should be mostly the same as in TinyGo.
3. Remove everything related to micro devices, only the Wasm related stuff.
4. Runtime does not need to be super size-optimized.
5. Implement custom GC that can work with external memory allocators (`extalloc`) via FFI.
6. Implement the export of `__data_end`, `__heap_base` globals.
7. Remove exported allocation functions in TinyGo.
8. Guide used to test the PoC.

*Hacks*
* Hardcode the same value for `__heap_base` inside the host and the runtime module.
* Inside the Wasm module, the GC could just allocate a large amount of memory and work with that.
