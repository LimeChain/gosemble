---
layout: default
permalink: /development/onboarding
---
# Onboarding 🎓

💬 *Please provide feedback on our onboarding guide. Let us know if any parts are unclear, confusing, or if you have other suggestions for improvement.*

The goal of this project is to develop a framework in Go that can be used by a Polkadot node.

The Polkadot node is divided into two components, the Polkadot **Runtime** and the Polkadot **Host**. The Runtime handles the state transition logic for the Polkadot protocol and is designed to be upgradeable without the need of a fork. The Polkadot Host provides necessary functionality required for the Runtime to execute its state transition logic.

Currently, Host implementations can be developed in Rust (Substrate), C++ (Kagome) and Go (Gossamer). On the other hand, Runtime can be developed only in Rust (Substrate).

In this project, we focus only on the Runtime development in Go. The output of the Runtime code is Wasm bytecode, which can be plugged into any Host (you can consider this similar as how Solidity smart contracts are executed in Ethereum). The difference is that this bytecode takes care of the core functionality of a network.

# Agenda

[1. Tech Stack 💘](https://www.notion.so/1-Tech-Stack-21e0196a5a8a4515a7bef4b004a31049?pvs=21) 

[1.1. **WebAssembly 🏗️**](https://www.notion.so/1-1-WebAssembly-d173f14596e94a60bbb3c94f8a66479f?pvs=21) 

[1.2. Go & **TinyGo 🦖 🐣**](https://www.notion.so/1-2-Go-TinyGo-3b41e9fc4d68469fa500e6b98573a332?pvs=21) 

[1.3. **TinyGo fork 🧪**](https://www.notion.so/1-3-TinyGo-fork-a89033f83732407baee991669ecc30be?pvs=21) 

[2. Architecture of a Polkadot Node 🏛️](https://www.notion.so/2-Architecture-of-a-Polkadot-Node-8eaca4ff0b474b8fabaaee01b6f7c9d9?pvs=21) 

[2.1. **Node** 💻](https://www.notion.so/2-1-Node-b81ba7fef13b47fc8b2dd96614d5bb96?pvs=21) 

[2.2. **Host-Runtime Interaction 🤝**](https://www.notion.so/2-2-Host-Runtime-Interaction-c4c0eca9f7334fa09a2a6c7291f7ccff?pvs=21) 

[2.3. **Runtime Internals ⚙️**](https://www.notion.so/2-3-Runtime-Internals-1df0dd180d144e569b13d642f93b3bda?pvs=21) 

[3. Specifics of implementing a Runtime in Go 🔍🐛 ](https://www.notion.so/3-Specifics-of-implementing-a-Runtime-in-Go-5605505e54fe4874a8583f29dc4b9574?pvs=21) 

[4. Tasks 📝](https://www.notion.so/4-Tasks-c893c18882f047e1950a8169ea1a62fe?pvs=21) 

[4.1. Compile a Runtime from Gosemble and run it in a Substrate node 🛠️](https://www.notion.so/4-1-Compile-a-Runtime-from-Gosemble-and-run-it-in-a-Substrate-node-a39a322c65e048f6810808324b5b5fc8?pvs=21) 

[4.2. Implement simple Runtime function and add tests 🛠️](https://www.notion.so/4-2-Implement-simple-Runtime-function-and-add-tests-6ce7220b0eb34f57b43e59265eddad29?pvs=21) 

[](https://www.notion.so/9471abd55d9f45b4ae737ea010777f47?pvs=21) 

# 1. Tech Stack 💘

## 1.1. **WebAssembly 🏗️**

WebAssembly (abbreviated as Wasm) is a binary instruction format for a stack-based virtual machine. It is designed as a compact, portable, and fast compilation target for high-level languages like C, C++, and Rust and many others that are being adapted. It enables execution of code at near-native speeds directly in web browsers and various host environments, as it resembles low-level machine code that modern CPUs understand.

### Binary Format

WebAssembly code is delivered in a low-level binary format, which is more compact than its textual representation. This binary format is designed to be fast to decode and execute. There's also a textual representation of this binary format, which is useful for debugging and testing.

### Low-Level Virtual Machine

WebAssembly provides a set of low-level virtual instructions that are closer to machine code than high-level programming languages. This enhances its execution efficiency and also makes it a suitable compilation target for other languages.

### Typed Instructions

WebAssembly instructions are strongly typed. It supports several numerical types like i32, i64, f32, and f64 and a few others for handling memory and tables.

### Stack-Based Architecture

Its computational model is designed around a stack-based architecture. Operations are performed by pushing and popping values from an implicit stack (the stack is inaccessible and distinct from the untrusted linear memory). Though it is not a pure stack machine, as it accommodates features like unlimited virtual registers (local variables).

### Modules & Sections

In the WebAssembly binary format, code and data are organized into modules. Each module consists of various sections arranged in a specific order, though some sections are optional. The module structure is defined by the WebAssembly specification and is validated before execution.

- **Type Section** - declares all function signatures used within the module.
- **Import Section** - specifies all module imports, such as functions, memories, globals, and tables.
- **Function Section** - contains a list of function declarations, each referencing the type section for its signature.
- **Table Section** - declares the table of function references used by the module.
- **Global Section** - declares global variables.
- **Export Section** - specifies all module exports, such as functions, memories, globals, and tables.
- **Elem Section** - initializes the table with function references for indirect function calls.
- **Code Section** - contains the binary code for the module's functions.
- **Data Section** - contains initial values for the module's memory.
- **Custom Section** - contains custom data for the host environment, possibly containing toolchain-specific information. This section can appear multiple times and is not restricted to a specific position in the list of sections.

Modules can be dynamically loaded and combined, making it possible to build and manage larger applications effectively.

### Memory Model

WebAssembly uses a single contiguous byte buffer as its memory model, referred to as linear memory. This memory is resizable, byte-addressable and is accessible by all memory operations.

**Imported/Exported Memory**

It is isolated from the host system, thus providing a safe environment for the execution of untrusted code. It can be imported or exported, thus facilitates reading and writing operations by both WebAssembly and the Host.

### Functions

Every module may contain functions which can be either exported (made accessible outside the module) or imported (indicating dependency on an external function).

**Exported/Imported Functions**

The WebAssembly module and its Host communicate using host-imported and exported functions.

Host-imported functions act as an additional bridge between WebAssembly and its host, enabling the module to access resources, input-output operations, or system-specific functionality. An example would be to get the time of the machine the Host is currently running on.

On the other hand, exported functions enables customization, allowing developers to expose specific functionalities to the host environment. An example would be to export a runtime function, which takes care of the execution of blockchain transactions.

### WebAssembly Extensions

WebAssembly has a minimal core specification, but it's designed to be extensible. Proposals like threads, garbage collection, and SIMD (single instruction, multiple data) operations are being worked on or have been added to provide more capabilities over time.

### WebAssembly & JavaScript

A significant feature of WebAssembly is its seamless interaction with JavaScript. The two can work in tandem within web applications.

### WebAssembly & WASI

Beyond browser capabilities, with the introduction of the WebAssembly System Interface (WASI), it can be integrated into a wide range of environments (Host), including web applications, desktop software, and more. In our case, we are going to embed it into another Rust/Go/C++ application. This is achieved by the WebAssembly module exposing a well-defined interface, facilitating communication with the host.

### WebAssembly MVP & Polkadot

Polkadot uses a version of WebAssembly (Wasm MVP) that does not support reference types or multiple return types. Therefore, non-numeric values are exchanged through shared memory using pointer-sized allocations. This mechanism allows the WebAssembly module to interact manipulate data within the host's memory space, facilitating data exchange between the WebAssembly module and its host.

In the case of Polkadot, the WebAssembly bytecode takes care of the state transition and block execution of each Polkadot Node, which is the most critical part. It is plugged into the Polkadot Node and it is called Runtime. In case of bugs, upgrades or updates, the logic can just be replaced with a new WebAssembly bytecode. This allows the Runtime to be updated on-chain, without the need of a network fork.

1. [Documentation](https://webassembly.org/)
2. [Intro to WebAssembly](https://hacks.mozilla.org/2017/02/a-cartoon-intro-to-webassembly/)
3. Install latest Go (1.21) - https://go.dev/dl/ or brew
4. Guides
    1. [Go to Wasm with JS apis - executed it in the browser](https://github.com/golang/go/wiki/WebAssembly#getting-started)
    2. [Go to WASI](https://pkg.go.dev/github.com/stealthrocket/wasi-go#readme-with-go)

WebAssembly has different [platform targets](https://snarky.ca/webassembly-and-its-platform-targets/) and extensions. Go supports [Wasm depending on JavaScript supported APIs](https://webassembly.org/getting-started/js-api/) and with the release of Go 1.21, they’ve added support for [WASI](https://wasi.dev/)**.**

**Unfortunately, Polkadot targets an old version of WebAssembly, called WebAssembly MVP, before [spec version 1](https://www.w3.org/TR/wasm-core-1/).** This is why we will not use the Go toolchain for building wasm blobs, but [TinyGo](https://github.com/tinygo-org/tinygo).

## 1.2. Go & **TinyGo 🦖 🐣**

TinyGo is a subset of Go with different goals from the standard Go. It is an alternative compiler and runtime aimed to support different small embedded devices and WebAssembly with a single processor core, emphasizing size optimizations.

1. [Documentation](https://github.com/tinygo-org/tinygo/blob/release/README.md)
2. [Install](https://tinygo.org/getting-started/install/macos/)
3. Guides
    1. [TinyGo to Wasm with JS apis - compile a Wasm module and execute it inside a JS environment (browser)](https://tinygo.org/docs/guides/webassembly/wasm/)
    2. [TinyGo to WASI](https://wasmbyexample.dev/home.en-us.html) - compile a Wasm blob and execute it inside another Go host application ([Wazero VM](https://wazero.io/))
        1. https://github.com/tetratelabs/wazero/tree/main/examples/allocation - check README and `tinygo` folder
        2. https://github.com/tetratelabs/wazero/tree/main/examples/import-go - check README and `tinygo` folder

## 1.3. **TinyGo fork 🧪**

We have [forked TinyGo](https://github.com/LimeChain/tinygo/) as we need to add a new target for the Polkadot-specific wasm blob, targeting standalone **Wasm MVP**, similar to Rust's `wasm32-unknown-unknown`, **without bulk memory operations and other extensions,** also incorporating **custom GC** that utilizes an external allocator. In the [polkawasm-target-dev branch](https://github.com/LimeChain/tinygo/tree/polkawasm-target-dev), you can see the changes created by Rado for the specific TinyGo releases.

- Example: https://github.com/LimeChain/tinygo/pull/16 shows the changes added to TinyGo `v0.29`.

We use a local build of TinyGo and do not depend on the already-built brew dependency.

[Here are the steps](https://github.com/LimeChain/gosemble/blob/develop/docs/docs/development/toolchain-setup.md) how to install and build it locally.

After you have built TinyGo, execute the following:

```bash
tinygo version
```

The output should be similar to:

```bash
tinygo version 0.30.0 darwin/arm64 (using go version go1.21.0 and LLVM version 16.0.6)
```

# 2. Architecture of a Polkadot Node 🏛️

Now that you have learned about WebAssembly, shared memory, runtime imported/exported functions and the TinyGo toolchain, let’s look at the Polkadot specification.

Polkadot node architecture and protocol specification is heavily influenced by the tech stack: WebAssembly MVP and Rust. Some implementation details, like the memory management, are not well abstracted and tightly coupled with the Rust implementation and even included as part of the protocol specification.

## 2.1. **Node** 💻

[Polkadot protocol](https://spec.polkadot.network/id-polkadot-protocol) has been divided into two parts, the [Polkadot Runtime](https://spec.polkadot.network/part-polkadot-runtime) and the [Polkadot Host](https://spec.polkadot.network/part-polkadot-host).

## 2.2. **Host-Runtime Interaction 🤝**

### I**mported Functions 📥**

External functions provided by the **Host** environment (Substrate/Kagome/Gossamer host) that **Runtime** (WebAssembly module) can invoke when needed, for more details check the [Host API](https://spec.polkadot.network/chap-host-api). The Host API provides access to memory, storage, crypto, hashing, logging and misc functionality.

Example (Storage):

- Rust implementation using Substrate
- Go [implementation](https://github.com/LimeChain/gosemble/blob/master/env/storage.go) using  Gosemble

### E**xported Functions** 📤

Defined within the **Runtime** (WebAssembly module) and can be invoked by the **Host** application, for more details check the [Runtime API](https://spec.polkadot.network/chap-runtime-api).

The Runtime API provides core and chains specific functionality.

Example (Core API):

- Rust [implementation](https://github.com/paritytech/polkadot-sdk/blob/master/substrate/bin/node-template/runtime/src/lib.rs#L345) using Substrate
- Go [implementation](https://github.com/LimeChain/gosemble/blob/develop/runtime/runtime.go#L214) using Gosemble

### **Memory 🧠**

Shared between the Host and the Runtime and is [managed by the Host allocator](https://spec.polkadot.network/chap-state#sect-memory-management) for all heap allocations. All data passed between Host and Runtime, like arguments to exported or imported functions or returned results, is encoded using SCALE encoding. Non numeric types, like byte buffers, are [shared](https://spec.polkadot.network/chap-state#sect-runtime-return-value) using a [pointer-size](https://spec.polkadot.network/chap-host-api#defn-runtime-pointer-size) to the allocation in the heap.

- [SCALE encoding (Spec)](https://spec.polkadot.network/id-cryptography-encoding#sect-scale-codec)
- [SCALE Codec (Substrate)](https://docs.substrate.io/reference/scale-codec/)

## 2.3. **Runtime Internals ⚙️**

### Extrinsics (Transactions) 💳

- [Transaction Types](https://docs.substrate.io/learn/transaction-types/)
- [Transaction Lifecycle](https://docs.substrate.io/learn/transaction-lifecycle/)

### Weights & Fees ⚖️ 💸

- [Transaction Weights (Gas) and Fees](https://docs.substrate.io/build/tx-weights-fees/)

### Accounts, addresses, and keys 👤 🔑

- [Accounts Addresses Keys](https://docs.substrate.io/learn/accounts-addresses-keys/)

### Storage 💾

Storing and retrieving data, key/value generation and types of storage values.

- [State Transitions and Storage](https://docs.substrate.io/learn/state-transitions-and-storage/#querying-storage)
- [Transactional Storage](https://docs.substrate.io/build/runtime-storage/#transactional-storage)

### Pallets (Modules) 🧱

Pallets communicate and interact with each other via events, storage, calls, hooks, etc.

- [Pallet Coupling](https://docs.substrate.io/build/pallet-coupling/)
- [Events and Errors](https://docs.substrate.io/build/events-and-errors/)

# 3. Specifics of implementing a Runtime in Go 🔍🐛

Developing a framework for writing Polkadot runtimes in Go is not a straight forward process, accompanied with many blockers and issues that need to be resolved. Most of the issues are related to the incompatibilities between the design decisions around the Polkadot protocol and the Go language. Here are some of the major challenges that we faced while working on the project:

- Missing support for standalone Wasm (**MVP**) 🕳️
- GC that is required to work with an external allocator, provided by the Host 💣
- Immature toolchain based on [custom Tinygo](https://github.com/LimeChain/tinygo) 🐣
- [SCALE codec](https://github.com/LimeChain/goscale) implementation with minimal reflection 🪩
- Writing mostly low-level and unsafe Go (none of the concurrency capabilities are utilized) ⚠️
- The spec lacks details regarding the Runtime; thus, you should be able to read Rust code, which is the actual source of truth. 🦀

Most of the things are documented here [Gosemble Runtime Architecture](https://limechain.github.io/gosemble/overview/runtime-architecture/) , but more up to date version of the docs can be found [Gosemble docs](https://github.com/LimeChain/gosemble/tree/develop/docs/docs) 🚧.

# 4. Tasks 📝

### 4.1. Compile a Runtime from Gosemble and run it in a Substrate node 🛠️

1. Install `git` and `docker`
2. Clone the Gosemble repo - `git clone https://github.com/LimeChain/gosemble.git`
3. Checkout the development branch - `git checkout develop`
4. Pull all necessary git submodules - `git submodule update --init --recursive`
5. Build the runtime - `make build-docker-dev`
6. Run the tests - `make test`
7. Start a local network - https://limechain.github.io/gosemble/tutorials/start-a-network/
8. Connect to Polkadot.js and do a simple transfer - https://limechain.github.io/gosemble/tutorials/transfer-funds/

### 4.2. Implement simple Runtime function and add tests 🛠️

1. Declare a runtime exported function, [example](https://github.com/LimeChain/gosemble/blob/8d77db41b91b51984769c9d68b0b347ed29f1c32/runtime/runtime.go#L240C1-L241C70).
2. Read a byte buffer passed as an argument, [example](https://github.com/LimeChain/gosemble/blob/8d77db41b91b51984769c9d68b0b347ed29f1c32/api/block_builder/module.go#L58C1-L60C30).
3. SCALE decode the byte buffer, [example](https://github.com/LimeChain/gosemble/blob/8d77db41b91b51984769c9d68b0b347ed29f1c32/api/block_builder/module.go#L62).
4. Call 2-3 host imported functions, [example](https://github.com/LimeChain/gosemble/blob/8d77db41b91b51984769c9d68b0b347ed29f1c32/primitives/log/log.go#L41C1-L45C2), [example](https://github.com/LimeChain/gosemble/blob/8d77db41b91b51984769c9d68b0b347ed29f1c32/primitives/storage/storage.go#L135-L139C2).
5. Return a byte buffer as a result, [example](https://github.com/LimeChain/gosemble/blob/8d77db41b91b51984769c9d68b0b347ed29f1c32/api/block_builder/module.go#L75).
6. Add unit & integration tests, [example](https://github.com/LimeChain/gosemble/blob/8d77db41b91b51984769c9d68b0b347ed29f1c32/runtime/block_builder_apply_extrinsic_test.go#L76).

# 5. Add GPG Key to your Github account

https://docs.github.com/en/authentication/managing-commit-signature-verification/adding-a-gpg-key-to-your-github-account