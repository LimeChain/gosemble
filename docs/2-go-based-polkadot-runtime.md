# Research Feasibility of a Go Runtime - Aug 2022

- [0. Abstract](#0-abstract)
- [1. Introduction](#1-introduction)
- [2. Background](#2-background)
    * [2.1. The design decisions behind Polkadot's architecture](#21-the-design-decisions-behind-polkadots-architecture)
        + [2.1.1. WebAssembly specification](#211-webassembly-specification)
        + [2.1.2. WASI requirements](#212-wasi-requirements)
        + [2.1.3. SCALE codec](#213-scale-codec)
        + [2.1.4. Runtime calls](#214-runtime-calls)
        + [2.1.5. Exported globals](#215-exported-globals)
        + [2.1.6. Imported vs exported memory](#216-imported-vs-exported-memory)
        + [2.1.7. External memory management](#217-external-memory-management)
    * [2.2. Translating Go's language capabilities to WebAssembly MVP](#22-translating-gos-language-capabilities-to-webassembly-mvp)
        + [2.2.1. Limitations](#221-limitations)
        + [2.2.2. Proposals](#222-proposals)
    * [2.3. Compiler backends targeting Wasm - LLVM, Binaryen](#23-compiler-backends-targeting-wasm---llvm-binaryen)
    * [2.4. Go](#24-go)
        + [2.4.1. Compiler](#241-compiler)
        + [2.4.2. Runtime](#242-runtime)
            - [2.4.2.1. Memory Management](#2421-memory-management)
            - [2.4.2.2. Concurrency](#2422-concurrency)
    * [2.5. TinyGo](#25-tinygo)
        + [2.5.1. Compiler](#251-compiler)
        + [2.5.2. Runtime](#252-runtime)
            - [2.5.2.1. Basic features](#2521-basic-features)
            - [2.5.2.2. Memory Management](#2522-memory-management)
            - [2.5.2.3. Concurrency](#2523-concurrency)
            - [2.5.2.4. Parallelism](#2524-parallelism)
- [3. Technical challenges](#3-technical-challenges)
    * [3.1. Toolchain support for Wasm for non-browser environments](#31-toolchain-support-for-wasm-for-non-browser-environments)
    * [3.2. Wasm features that are not part of the MVP](#32-wasm-features-that-are-not-part-of-the-mvp)
    * [3.3. Standard library support](#33-standard-library-support)
    * [3.4. External memory allocator and GC](#34-external-memory-allocator-and-gc)
    * [3.5. Linker globals](#35-linker-globals)
    * [3.6. Developer experience](#36-developer-experience)
- [4. Solution](#4-solution)
    * [4.1. Proof of concept (alternative compiler + runtime + GC with external allocator)](#41-proof-of-concept-alternative-compiler--runtime--gc-with-external-allocator)
    * [4.2. Testing Guide](#42-testing-guide)
- [5. Conclusion](#5-conclusion)
- [6. Discussion](#6-discussion)
- [7. References](#7-references)

## 0. Abstract

The lack of diversity and ease of use of Polkadot Runtimes is a barrier that stops Polkadot from living up to its full
promise. The Polkadot community should as a whole address this problem.

While there are several choices for implementing Polkadot Hosts, [Rust](https://github.com/paritytech/substrate)
, [C++](https://github.com/soramitsu/kagome), and [Go](https://github.com/ChainSafe/gossamer), the only option for
writing Polkadot Runtimes is [Rust](https://github.com/paritytech/substrate). There are too many good things to say
about Rust, but it is well-known that it has a steep learning curve. On the other hand, Go is a language focused on
simplicity that is gaining popularity among software developers nowadays. It is modern, powerful, and fast, backed by
Google and used in many of their software, thus making it an ideal candidate for implementing Polkadot Runtimes.
Arguably, other Blockchain networks (e.g. [Cosmos](https://github.com/cosmos)) have gained significant adoption due to
the lower barrier for entry (compared to Rust).

To be feasible to develop Polkadot Runtime in Go, there are technological challenges that need to be cleared out first.
This research is aimed at those challenges.

## 1. Introduction

Writing Polkadot Runtimes in Go is exciting, mainly because of Go's simplicity and automatic memory management. However,
we have to be sure that Polkadot's design decisions are well suited to the design and toolchain support of Go. Notably,
producing Wasm runtime compatible with Polkadot's specification is not supported in any Go toolchain and Polkadot design
decisions around memory management pose the risk of making Go's main selling points irrelevant. Although the language
specification doesn't mention how it should manage its memory, the Go community recognizes it as a language with GC.

So aren't we setting up ourselves for failure with Go right from the start? This research aims to provide conclusions if
Go is a suitable choice to write Polkadot Runtimes and further aid the development of a Go toolchain capable of
producing compatible Wasm.

Here is a list of questions we are going to try to give answers to:

* What are the design decisions behind Polkadot's architecture - WebAssembly specification, Host/Runtime interaction,
  Runtime memory management?
* How well the incorporated WebAssembly specification in Polkadot aligns with the Go language capabilities - current
  limitations and upcoming proposals that might be leveraged?
* Which compiler framework is best suited for targeting Wasm from a language with GC-managed runtime - LLVM, Binaryen?
* How Go and TinyGo internals work and differ - compiler, runtime, GC?
* What are the challenges of implementing a Polkadot Runtime with GC-managed language like Go?

Along with the research document, a Proof of concept toolchain was developed that is capable of building a Go-based
Runtime that conforms to the Polkadot Wasm Runtime specification.

But why is it so important to write in Go beside the above said? We believe that the Polkadot ecosystem does not have
diversity of Runtime implementations. Developing a toolchain and framework that allows Runtimes to be developed in Go
would resolve that.

## 2. Background

### 2.1. The design decisions behind Polkadot's architecture

In addition to the [Polkadot spec](https://github.com/w3f/polkadot-spec) here is a list of important points that deserve
to be addressed.

#### 2.1.1. WebAssembly specification

At the time of writing this research document, there is the WebAssembly 1.0 specification and a draft for spec 2.0.
Polkadot/Substrate Runtimes target [WebAssembly MVP](https://github.com/WebAssembly/design/blob/main/MVP.md) without any
extensions enabled, that offers limited set of features compared to WebAssembly 1.0. Adding on top of that,
Polkadot/Substrate specifications for the Runtime module define very domain-specific API that consist of:

* imported Host provided functions (`ext_allocator_malloc_version_1`, `ext_allocator_free_version_1`, etc).
* imported Host provided memory.
* exported linker specific globals (`__heap_base`, `__data_end`).
* exported `__indirect_function_table` (WIP and not enabled currently).
* exported business logic API functions (`Core_version`, `Core_execute_block`, `Core_initialize_block`, etc).

```
Type: wasm
...
Imports:
  Functions:
    "env"."ext_allocator_malloc_version_1": [I32] -> [I32]
    "env"."ext_allocator_free_version_1": [I32] -> []
    ...
  Memories:
    "env"."memory": not shared (20 pages..)
  Tables:
  Globals:
Exports:
  Functions:
    "Core_version": [I32, I32] -> [I64]
    "Core_execute_block": [I32, I32] -> [I64]
    "Core_initialize_block": [I32, I32] -> [I64]
    ...
  Memories:
  Tables:
    "__indirect_function_table": FuncRef (352..352)
  Globals:
    "__data_end": I32 (constant)
    "__heap_base": I32 (constant)
```

```
 _________________________________________________________________________________________
| HOST                                                                                    |
|                                       _________________________________________         |
|                                      |               ALLOCATOR                 |        |
|                                      | ext_allocator_malloc,ext_allocator_free |        |
|                                      |_________________________________________|        |
|      ________________________________________________________|_____________________     |
|     | WASM                                                   | (imported)          |    |
|     |                                                        |                     |    |
|     |             ___________________________________________▼___________________  |    |
|     | (imported) |          |                  |         [---]                   | |    |
|  MEMORY  -----►  | Data     |         ◄- Stack | Heap -►                         | |    |
|     |            |__________|__________________|_________________________________| |    |
|     |            0     __data_end         __heap_base                  max memory  |    |
|     |                                                                              |    |
|     |                                                                              |    |
|     |______________________________________________________________________________|    |
|                                                                                         |
|_________________________________________________________________________________________|

```

#### 2.1.2. WASI requirements

Polkadot is a non-browser environment, but it is not an OS. It doesn't seek to provide access to an operating-system API
like files, networking, or any other major part of the things provided by WASI (WebAssembly System Interface).

#### 2.1.3. SCALE codec

Runtime data, coming in the form of byte code, needs to be as light as possible. The SCALE codec provides the capability
of efficiently encoding and decoding it. Since it is built for little-endian systems, it is compatible with Wasm
environments.

#### 2.1.4. Runtime calls

Each function call into the Runtime is done with newly allocated memory (via the shared allocator), either for sharing
input data or results. Arguments are SCALE encoded into a byte array and copied into this section of the Wasm shared
memory. Allocations do not persist between calls. It is important to note that the Runtime uses the same Host provided
allocator for all heap allocations, so the Host is in charge of the Wasm heap memory management. Data passing to the
Runtime API is always SCALE encoded, Host API calls on the other hand try to avoid all encoding.

#### 2.1.5. Exported globals

It is expected from the Runtime to export `__heap_base` global indicating the beginning of the heap. It is used by the
Host allocator to prevent memory allocations below that address and avoid clashes with the stack and data sections.

#### 2.1.6. Imported vs exported memory

Imported memory works a little better than exported memory since it avoids some edge cases, although it also has some
downsides. Working with exported memory is almost certainly still supported and in fact, this is how it worked in the
beginning. However, the current spec describes that memory should be made available to the Polkadot Runtime for import
under the symbol name `memory`.

#### 2.1.7. External memory management

The design in which allocation functions are on the Host side is dictated by the fact that some Host functions might
return buffers of data of unknown size. That means that the Wasm code cannot efficiently provide buffers upfront.

For example, let's examine the Host function that returns a given storage value. The storage value's size is not known
upfront in the general case, so the Wasm caller cannot pre-allocate the buffer upfront. A potential solution is to first
call the Host function without a buffer, which will return the value's size, and then do the second call passing a
buffer of the required size. For some Host functions, caches could be put in place for mitigation, some other functions
cannot be implemented in such model at all. To solve this problem, it was chosen to place the allocator on the Host
side. However, this is not the only possible solution, as there is an ongoing discussion about moving the allocator into
the Wasm: [[1]](https://github.com/paritytech/substrate/issues/11883). Notably, the allocator maintains some of its data
structures inside the linear memory and some other structures outside.

### 2.2. Translating Go's language capabilities to WebAssembly MVP

It is important to see how well the language features of Go translate to WebAssembly, its limitations, and upcoming
proposals that might help to overcome them. More specifically, the capabilities offered by WebAssembly MVP, targeted by
Polkadot. WebAssembly is a low-level bytecode instruction format for a typed stack-based virtual machine, that uses
little-endian byte order when translating between values and bytes and has a structured control flow. Incorporates
Harvard architecture - the program state is separate from the instructions, with an implicit stack that can't be
accessed, only the untrusted memory. It has single, contiguous, byte-addressable default linear memory, accessed by all
memory operators, spanning from offset 0 and extending up to varying memory sizes. It can be shared with other module
instances by either importing or defining it internally inside the module, but either way, after import or definition,
there is no difference when accessing it.

#### 2.2.1. Limitations

* the linear memory is great for languages with manual, reference counting or ownership memory model, but GC-managed
  memory requires a bit more work to port it. Since WebAssembly has no stack introspection to scan the roots, it
  requires to use mirrored shadow stack in the linear memory, pushed/popped along with the machine stack, thus making it
  less efficient.
* there is only one memory associated with a module or instance in MVP, this memory at index zero is the default.
* the linear memory is resizable, but only upward (`grow_memory`, `memory.grow`).
* in contrast to Go, WebAssembly does not yet support true parallelism, it lacks support for multiple threads, atomics,
  and memory barriers. However, in Polkadot, parallelism is achieved through different mechanics (ParaChains)

#### 2.2.2. Proposals

* [GC proposal](https://github.com/WebAssembly/gc/blob/main/proposals/gc/Overview.md) might be handy to support an
  automatic memory management, but there are concerns how performant it will
  be [[1]](https://github.com/WebAssembly/gc/issues/59), [[2]](https://github.com/WebAssembly/gc/issues/36). In addition
  to that, the Polkadot Runtime supports only WebAssembly MVP currently. To add on top of that, the GC proposal is 
  under development, and
  it is not yet clear if Polkadot will be able to leverage the GC proposal. Potential problems include determinism (is
  there anything in GC that causes ND? Can it be tamed efficiently?) and safety (Is it possible for a host to limit the
  resource consumption reliably and deterministically?).
* [Stack introspection](https://github.com/WebAssembly/design/issues/1340)
* [Multiple memories](https://github.com/WebAssembly/multi-memory)

### 2.3. Compiler backends targeting Wasm - LLVM, Binaryen

LLVM and Binaryen are both compiler infrastructures that can be used to produce WebAssembly. However, there are some
differences that are important to be noted:

* LLVM currently does not support Wasm GC, and it is not the main focus for LLVM as most languages that utilize it uses
  linear memory - C, C++, Rust, Zig. Binaryen is a better choice for languages with GC, which intend to compile to Wasm
  GC. Additionally, it will support compilation from Wasm GC to Wasm MVP, as a polyfill, though the opposite will not be
  possible. For example, if Polkadot's Wasm target switches to Wasm GC, having Binaryen as a compiler backend will have
  support for that, as the developers behind it are from the WebAssembly organization.
* LLVM is much more powerful as an optimizer.
* LLVM takes much more memory and is much slower than Binaryen's optimizations.
* LLVM creates larger output compared to Binaryen.

### 2.4. Go

#### 2.4.1. Compiler

The default compiler is `gc`. There is also `gccgo`, which uses the GCC back-end, and `gollvm`, which uses the LLVM
infrastructure (somewhat less mature).

**Pipeline (Lowering)**

```
// Frontend Compiler (lexing, parsing, typechecking)

CODE -> tokenizer/lexer/scanner (lexical analysis) -> TOKENS STREAM
TOKENS STREAM -> go/parser (syntactic analysis) -> AST (annotated)
AST -> go/types (type check) -> AST (with types)
AST (with types) -> golang.org/x/tools/go/ssa (convert) -> SSA

// Optimization

* variable capturing, inlining, escape analysis, closure rewriting, walk

SSA (higher level with Go-specific constructs like interfaces and *goroutines*) -> convert (tinygo) -> LLVM IR
LLVM IR -> optimize (llvm) -> LLVM IR (optimize by a mixture of handpicked LLVM optimization passes, TinyGo-specific optimizations ()

// Backend Compiler

LLVM IR (optimized) -> convert (llvm) -> MACHINE CODE
MACHINE CODE -> convert (llvm) -> OBJECT FILE
OBJECT FILE -> link (llvm) -> EXECUTABLE
```

#### 2.4.2. Runtime

Go has an extensive library, called the `runtime`, which is used by Go programs. 
The library includes a scheduler, memory allocator, garbage collector, stack management, data structures and other critical features of the `Go` language.

##### 2.4.2.1. Memory Management

Memory Management uses virtual memory that abstracts the access to the physical memory.

* Process Memory Layout

```
 ______________________
|        STACK         | Function Stack Frames
|----------------------| ◄--- Stack Pointer
|          |           |
|          ▼           |
|                      |
|          ▲           |
|          |           |
|----------------------|
|         HEAP         | Dynamic Allocated Variables
|----------------------|
|         BSS          | Uninitialized Static Variables
|----------------------|
|         DATA         | Initialized Static Variables
|----------------------|
|         TEXT         | Code
|______________________|
```

The Go Compiler decides (via escape analysis) when a value should be allocated on heap memory. The objects in the heap
are allocated by memory allocators and collected by garbage collectors.

Most of the function arguments, return values, and local variables of function calls are allocated on the stack, which
is managed by the compiler. When it comes down to passing pointers (sharing down), typically allocations are on the
stack. On the other hand, returning pointers (sharing up) typically allocate on heap memory.

Allocations on heap memory also occur when the value is:

* returned as a result of a function execution and is referenced
* too large to find on the stack
* with unknown size at compile time and the compiler does not know what to do with it

Memory management generally consists of three components - program, allocator and collector. Whenever the program
requests memory, it requests it through the memory allocator, which initialises the corresponding memory to the heap.

```
   Mutator (Process)
      |
      | malloc/free
      |
      ▼
  Allocator
      |
      | syscall mmap/munmap (address, size, permissions, flags)
      |
 _____▼____
|          | 
|   Heap   | ◄------ Garbage Collector
|__________|
```

* **Allocator**

  The Go runtime memory allocator is based on thread-caching malloc (TCMalloc) allocation strategy that classifies
  objects according to their size, including multi-level caching to improve its performance.

* **Garbage Collector**

  Garbage collectors are responsible for tracking memory allocations in heap memory, keeping those allocations that are
  still in-use, and releasing allocations when they are no longer needed.

  Go uses a concurrent mark-and-sweep algorithm with a write-barrier for its garbage collector, running concurrently
  with mutator threads and allowing multiple GC threads to run in parallel. The algorithm is decomposed into several
  phases.

    * *Collection Phases*

      There are three collection phases of the garbage collector.

      *Start/Stop The World (STW)* - whenever STW phase is found, application business logic is not executed.

        1. Mark Setup (STW)

           This phase turns on the *write barrier*, which makes sure that all concurrent activity is completely safe.
           This will stop every goroutine from running.

        2. Marking (concurrent)

           The goal of this phase is to mark values in heap memory that are still in-use. The collector inspects all
           stacks to find root pointers to heap memory and traverses the heap graph based on them. If the collector sees
           that it might run out of memory, *Mark Assist* is triggered, which slows down allocations to speed up
           calculations.

        3. Mark Termination (STW)

           This phase turns off the write barrier, and executes various cleanup tasks (e.g. flushing mcaches).

    * *Concurrent sweep*

      The sweep phase runs concurrently with normal execution. The heap is swept span-by-span both when a *goroutine*
      needs another span and concurrently in a background *goroutine*. In order to not request additional OS memory
      while there are not swept spans, when goroutine needs another span, it first tries to reclaim that much memory by
      sweeping. The cost of the sweeping is not on the GC, but on the new allocation itself.

    * *GC rate*

      The next GC is after an allocation of an extra amount of memory, proportional to the amount already in use. Go has
      an environment variable called `GOGC` (GC rate), which represents a ratio of how much new heap memory can be
      allocated before the next collection has to start. Adjusting `GOGC` changes the linear constant and the amount of
      extra memory used.

##### 2.4.2.2. Concurrency

Every Go program has an initial main *goroutine*, which is very similar to an application-level thread, and can have multiple other
*goroutines*, allowing to do concurrent tasks. *Goroutines* are managed by the Go scheduler, which depending on their state,
makes scheduling decisions.

The high-level states of a *goroutine* are:
* Running - the *goroutine* is on an OS thread and executes its instructions.
* Runnable - the *goroutine* is ready to execute its instructions, but does not have allocated time on an OS thread. This can happen when
multiple *goroutines* want time on an OS thread, which will automatically make *goroutines* wait longer to get time. This might lead to bad performance.
* Waiting/Blocked - the *goroutine* waits due to a system call or synchronization calls (blocking channels, atomic or mutex operations).

The Go scheduler implements [cooperative](https://en.wikipedia.org/wiki/Cooperative_multitasking) scheduling, which depending on events makes
scheduling decisions and does context-switching. Usually, you cannot predict what the Go scheduler is going to do. 

There are four types of events in the Go program that might allow the scheduler to make a scheduling decision:
* System calls
  
  Whenever a *goroutine* makes a system call, which will block the OS thread, the scheduler can decide to do context-switching by switching the *goroutine*
  off the thread and put another *goroutine* on it.

* Synchronization

  If atomic, mutex, or channel operations block a *goroutine*, the scheduler can again context-switch.

* Garbage collection
  
  The GC uses its own set of *goroutines*, which automatically allows scheduling decision to be made.

* The use of `go` word
  
  This creates a new *goroutine*, which gives the scheduler a chance to make a scheduling decision.

### 2.5. TinyGo

It is a subset of Go with very different goals from the standard Go. It is an alternative compiler and runtime aimed to
support different small embedded devices with a single processor core that require certain optimizations mostly toward
size.

*Goals*

* Have very small binary sizes.
* Support for most common microcontroller boards.
* Be usable on the web using WebAssembly.
* Good CGo support, with no more overhead than a regular function call.
* Support most standard library packages and compile most Go code without modification.

*Non-goals*

* Using more than one core.
* Be efficient while using zillions of *goroutines*. However, good *goroutine* support is certainly a goal.
* Be as fast as `gc`. However, LLVM will probably be better at optimizing certain things so TinyGo might actually turn
  out to be faster for number crunching.
* Be able to compile every Go program out there.

#### 2.5.1. Compiler

The compiler uses (mostly) the standard library to parse Go programs and LLVM to optimize the code and generate machine
code for the target architecture.

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

#### 2.5.2. Runtime

The runtime is written from scratch, optimized for size instead of speed and re-implements some compiler intrinsics and
packages:

* Startup code (device specific runtime initialization of memory, timers)
* GC (copied from micropython, super simple, optimized for size, runs the GC once it runs out of memory)
* Memory Allocator
* *Goroutines* scheduler
* Channels
* Time handling (every chip has a clock)
* Hashmap implementation (optimized for size instead of speed)
* Runtime functions like maps, slice append, defer, ... (various things that are not implemented in the compiler)
* Operations on strings
* `sync` & `reflect` package (strongly connected to the runtime)

##### 2.5.2.1. Basic features

All basic types, slices, all regular control flow including switch, closures and bound methods are supported, `defer`
keyword is almost entirely supported, except deferring some builtin functions, interfaces are quite stable and should
work well in almost all cases. Type switches and type asserts are also supported, as well as calling methods on
interfaces. The only exception is comparing two interface values. Maps are usable but not complete. You can use any type
as a value, but only some types are acceptable as map keys (strings, integers, pointers, and structs/arrays that contain
only these types). Also, they have not been optimized for performance and will cause linear lookup times in some cases.

##### 2.5.2.2. Memory Management

Heap with a garbage collector. While not directly a language feature, garbage collection is important for most Go
programs to make sure their memory usage stays in reasonable bounds.

Garbage collection is currently supported on all platforms, although it works best on 32-bit chips. A simple
conservative mark-sweep collector is used that will trigger a collection cycle when the heap runs out (that is fixed at
compile time) or when requested manually using `runtime.GC()`. Some other collector designs are used for other targets,
TinyGo will automatically pick a good GC for a given target.

Careful design may avoid memory allocations in main loops. You may want to compile with `-print-allocs=.` to find out
where allocations happen and why they happen.

##### 2.5.2.3. Concurrency

*Goroutines* and channels work for the most part, for platforms such as WebAssembly the support is a bit more limited (
calling a blocking function may for example allocate heap memory).

##### 2.5.2.4. Parallelism

Single core only.

## 3. Technical challenges

TinyGo's design decisions are mostly based on optimizations around small embedded devices. Of course, this is good for
the blockchain's use case too, but not always required and as crucial as is for devices with very limited resources. The
necessity of rewriting a large part of Go's runtime to align with those optimizations contributes to the effort of
supporting Go's capabilities. Another point where TinyGo diverges from the toolchain requirements of Polkadot is that it
supports Wasm aiming at the most recent Webassembly features. Here is a detailed breakdown of most of the problematic
points:

### 3.1. Toolchain support for Wasm for non-browser environments

* The official Go compiler does not support Wasm for non-browser
  environments [[1]](https://github.com/golang/go/issues/31105), only Wasm with browser-specific API.
* The alternative TinyGo compiler and runtime, supports Wasm outside the browser, but it is still not capable of
  producing Wasm compatible with Polkadot's requirements.

### 3.2. Wasm features that are not part of the MVP

TinyGo makes use of some features that are not supported in the targeted Wasm MVP, such as bulk memory
operations (`memory.copy`, `memory.fill` used to reduce the code size) and other extensions.

### 3.3. Standard library support

* The standard library relies on the `reflect` package (most common types like numbers, strings, and structs are
  supported), which is not fully supported by TinyGo [[1]](https://github.com/tinygo-org/tinygo/pull/2640). The core
  primitives and SCALE serialization logic that we intended to reuse
  from [gossamer](https://github.com/ChainSafe/gossamer) also rely on the `reflect` package.
* Many features of Cgo are still unsupported (#cgo statements are only partially supported).

### 3.4. External memory allocator and GC

* According to the Polkadot specification, the Wasm module does not include a memory allocator. It imports memory from
  the Host and relies on Host imported functions for all heap allocations. TinyGo implements simple GC and manages its
  memory by itself, contrary to specification. So it can't work out of the box on systems where the Host wants to manage
  its memory.

### 3.5. Linker globals

* The Runtime is expected to expose `__heap_base` global [[1]](https://github.com/tinygo-org/tinygo/issues/2045), but
  TinyGo doesn't support that out of the box.

### 3.6. Developer experience

* Implementing Wasm functionality makes you go pretty low-level and use some "unsafe" language constructs.

## 4. Solution

Taking into consideration all the technical challenges, the timeframe, and the number of different technologies, we
propose a solution based on a modified version of TinyGo/LLVM, which is going to serve as a proof of concept. The goal
is to implement a solution that incorporates runtime with GC and external memory allocator targeting Wasm MVP.

### 4.1. Proof of concept (alternative compiler + runtime + GC with external allocator)

No changes are supposed to happen on the front-end compiler. Most changes are going to be in LLVM backend and as part of
the runtime implementation.

**Add Toolchain support for Wasm compatible with Polkadot**

* [x] Fork and add [tinygo](https://github.com/LimeChain/tinygo) as a submodule.

  ```
    * root: contains the command line interface for the tinygo command and all its subcommands
    * builder: orchestrates the build
    * loader: loads and typechecks the code, and produces an AST
    * compiler: the compiler itself, makes little attempt at optimizing code
    * interp: tries to run package initializers at compile time as far as possible
    * transform: implements various optimizations necessary to produce working and efficient code
    * src: runtime implementation
  ```

* [x] Add separate Dockerfile and build script for building TinyGo with prebuild LLVM and dependencies for faster
   builds.

* [x] Add new target similar to Rust's `wasm32-unknown-unknown`, but aimed to support Polkadot's Wasm MVP.

* [x] add new target `polkawasm.json` in `targets/`.
* [x] add new `polkawasm` implementation in `runtime_polkawasm.go`
* [x] use the `polkawasm` directive to separate the new target from the existing Wasm/WASI functionality.
* [x] use linker flag to declare memory as imported.
* [x] use linker flags to export `__heap_base`, `__data_end` globals.
* [x] use linker flags to export `__indirect_function_table`.
* [x] change the stack placement not to start from the beginning of the linear memory.
* [x] disable the scheduler to remove the support of *goroutines* and channels (and JS/WASI exports).
* [x] remove the unsupported features by Wasm MVP (bulk memory operations, lang. ext) and add implementation
  of `memmove`, `memset`, `memcpy`, use opt flag as part of the target.
* [x] increase the memory size to 20 pages.
* [x] use the conservative GC as a starting point (there is a chance for memory corruption).
* [x] add GC implementation that can work with external memory allocator (remove memory allocation exports).
* [x] override the allocation functions used in the GC with such provided by the host.
* [x] remove the exported allocation functions.
* [ ] the `_start` export func should be called to init the heap (the host is not expected to support that).
* [ ] better abstractions, the extalloc GC depends on third party allocation API that might change in the future.

We have forked TinyGo and have created the following pull requests, which include all the completed steps above:

* [Polkawasm target](https://github.com/LimeChain/tinygo/pull/1), using TinyGo's `gc_conservative` garbage collector.
* An [extended version](https://github.com/LimeChain/tinygo/pull/2/) of `polkawasm` target, which uses an `extalloc`
  garbage collector, which uses the Host imported allocator functions (`ext_allocator_malloc_version_1`
  , `ext_allocator_free_version_1`).

**Setup Polkadot Host**

1. [x] Fork and add [gossamer](https://github.com/LimeChain/gossamer) as a submodule.
2. [x] Make the necessary changes to run locally provided Runtime inside the Host.
3. [x] Setup test instance to run the compiled Wasm (target MVP, import host provided functions and memory, implement
   bump allocator).

**Implement Polkadot Runtime**

1. [x] Implement SCALE codec without reflection.
2. [ ] Implement the minimal Runtime API (core API).
    * [x] `Core_version`
    * [ ] `Core_execute_block`
    * [ ] `Core_initialize_block`
3. [x] Add Makefile steps

**Future toolchain improvements**

1. [ ] write performance tests.
2. [ ] complete the SCALE codec implementation.
3. [ ] complete the `reflect` packages support
4. [ ] extalloc GC might need more work.
5. [ ] fix output errors in the Wasm memory.

SCALE codec

* [ ] fix: `panic("Assertion error: n>4 needed to compact-encode uint64")`
* [ ] fix: `Could not write " + strconv.Itoa(len(bytes)) + " bytes to writer`

TinyGo + conservative GC

* [ ] `runtimePanic("out of memory")`
* [ ] `runtimePanic("nil pointer dereference")`
* [ ] `runtimePanic("slice out of range")`
* [ ] `runtimePanic("index out of range")`

TinyGo + extalloc GC

* [ ] `return "reflect: call of reflect.Type." + e.Method + " on invalid type"`
* [ ] `return "reflect: call of " + e.Method + " on zero Value"`

* [ ] Gossamer instance freezes with Wasm compiled with extalloc GC while running the Core_version test. Here is the
  fragment that seems to cause that:

```go
  instance.ctx.Version, err = instance.version()
  if err != nil {
    instance.close()
    return nil, fmt.Errorf("getting instance version: %w", err)
  }
```

Comparison of `Core_version`'s output (with Apis field left empty):

* **Substrate**

  `\x10node4gosemble-node\x05\x00\x00\x00\x04\x00\x00\x00\x03\x00\x00\x00\x00\x02\x00\x00\x00\x01`

* **Gosemble + conservative GC**

  `\x10node4gosemble-node\x05\x00\x00\x00\x04\x00\x00\x00\x03\x00\x00\x00\x00\x02\x00\x00\x00\x01`

* **Gosemble + extalloc GC**

  `\x10node4gosemble-node\x05\x00\x00\x00\x04\x00\x00\x00\x03\x00\x00\x00\x00\x02\x00\x00\x00\x01`

### 4.2. Testing Guide

Detailed steps on how to test the PoC Runtime can be found [here](../README.md#poc-of-a-polkadot-runtime-in-go).

## 5. Conclusion

**The resulting proof of concept could be used to develop Polkadot Runtimes**. Also, the current research document is useful
for providing context as a starting point to make further contributions to the project or to take new direction for
implementation of alternative toolchain.

## 6. Discussion

* How good is the performance of the produced Wasm?
* How much effort will be to add full support of the `reflect` package?
* Is it possible in Wasm, to use separate heap regions, one reserved for GC allocations and another reserved for host
  allocations, to prevent clashes?
* Is LLVM the right choice for the long run or Binaryen is more suitable to provide stable foundation and support of
  Webassembly?

## 7. References

* [1] https://docs.substrate.io/
* [2] https://github.com/w3f/polkadot-spec
* [3] https://github.com/paritytech/substrate
* [4] https://webassembly.org/
* [5] https://www.w3.org/TR/2019/REC-wasm-core-1-20191205
* [6] https://dl.acm.org/doi/pdf/10.1145/3062341.3062363
* [7] https://www.cl.cam.ac.uk/~caw77/papers/mechanising-and-verifying-the-webassembly-specification.pdf
* [8] https://github.com/WebAssembly/design/blob/main/MVP.md
* [9] https://github.com/WebAssembly/spec
* [10] https://github.com/WebAssembly/proposals
* [11] https://github.com/WebAssembly/gc
* [12] https://github.com/WebAssembly/binaryen
* [13] https://llvm.org/
* [14] https://llvm.org/docs/GarbageCollection.html
* [15] https://go.dev/ref/spec
* [16] https://go.dev/doc/
* [17] https://go.dev/blog/
* [18] https://research.swtch.com/
* [19] https://github.com/golang/go
* [20] https://tinygo.org/
* [21] https://github.com/tinygo-org/tinygo
* [22] https://github.com/tinygo-org/tinygo/wiki
* [23] https://aykevl.nl/archive/
