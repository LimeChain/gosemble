# Gosemble - Go-based Polkadot/Substrate Runtimes 

The following project is based on a previous [research](https://github.com/LimeChain/gosemble-research) funded by [**Web3 Foundation**](https://web3.foundation) **Grants** and developed, maintained by
[**LimeChain**](https://limechain.tech).

[Research Results](./docs/2-go-based-polkadot-runtime.md)

## PoC of a Polkadot Runtime in Go

The steps below will showcase testing a PoC Polkadot Runtime implementation in Go.

**Prerequisites**

- [git](https://git-scm.com/downloads)
- [Go 1.19+](https://golang.org/doc/install)
- [docker](https://docs.docker.com/install/)

**Cloning the repository**

```bash
git clone https://github.com/LimeChain/gosemble.git
cd gosemble
```

**Pull all necessary git submodules**

```bash
git submodule update --init --recursive
```

**Build the Runtime**

Using a [forked version of TinyGo](https://github.com/LimeChain/tinygo), we build the Runtime with target `polkawasm`,
exported in `build/runtime.wasm`.

There are currently two options to choose from for the GC implementation that can be switched by setting the `GC` environment variable.

* simple GC implementation that works with the external host's allocator as per specification.

```bash
GC="" make build
```

* conservative GC, which works by using a different heap base, offset from the allocator's one (as a workaround), 
so that the GC uses a separate heap region for its allocations and does not interfere with the host's allocator region (only for testing).

```bash
GC="conservative" make build
```

**Run Tests**

After the runtime has been built, we execute standard Go tests with the help of
[Gossamer](https://github.com/LimeChain/gossamer), which we use to import necessary Polkadot Host
functionality and interact with the Runtime.

```bash
make test_unit
make test_integration
```

**Run a Network**

```bash
make start-network
```

**Optional steps**

* Inspect WASM Runtime - [wasmer](https://wasmer.io/)

```bash
wasmer inspect build/runtime.wasm
```

* To inspect the WASM Runtime in more detail, and view the actual memory segments

```bash
wasm-objdump -x build/runtime.wasm
```

## Architecture/Development Notes

The proposed solution is based on an alternative Go compiler that aims at supporting Wasm runtimes compatible with [Polkadot spec](https://spec.polkadot.network/id-polkadot-protocol)/[Substrate](https://docs.substrate.io/main-docs/) and conforming to the decisions behind Polkadot's architecture.

#### Toolchain

Since we use modified Tinygo for compiling Go to Wasm, some of Go's language capabilities can not be applied due to the limited support in Tinygo which also affects some of the design decisions.

#### WebAssembly specification

It targets [WebAssembly MVP](https://github.com/WebAssembly/design/blob/main/MVP.md) without any
extensions enabled, that offers limited set of features compared to WebAssembly 1.0. Adding on top of that,
Polkadot/Substrate specifications for the Runtime module define very domain-specific API that consist of:

* imported Host provided functions for dealing with memory, storage, crypto, logging, etc.
* imported Host provided memory.
* exported linker specific globals (`__heap_base`, `__data_end`).
* exported `__indirect_function_table` (WIP and not enabled currently).
* exported business logic API functions (`Core_version`, `Core_execute_block`, `Core_initialize_block`, etc).

Polkadot is a non-browser environment, but it is not an OS. It doesn't seek to provide access to an operating-system API
like files, networking, or any other major part of the things provided by WASI (WebAssembly System Interface).

#### SCALE codec

Runtime data, coming in the form of byte code, needs to be as light as possible. The SCALE codec provides the capability
of efficiently encoding and decoding it. Since it is built for little-endian systems, it is compatible with Wasm
environments.
The runtime works with custom-defined SCALE types compatible with Tinygo. At the boundary where it interacts with the host (memory, storage), those are converted to ordinary Go types. 

#### Host/Runtime interaction

Each function call into the Runtime is done with newly allocated memory (via the shared allocator), either for sharing
input data or results. Arguments are SCALE encoded into a byte array and copied into this section of the Wasm shared
memory. Allocations do not persist between calls. It is important to note that the Runtime uses the same Host provided
allocator for all heap allocations, so the Host is in charge of the Wasm heap memory management. Data passing to the
Runtime API is always SCALE encoded, Host API calls on the other hand try to avoid all encoding.

#### GC with external memory allocator

Since Go is a language that uses GC, such with external memory allocator is implemented in our Tinygo fork to meet the requirements of the Polkadot specification.

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


#### Stack placement

The stack placement differs from the one compiled from Substrate/Rust. The stack is placed before the data section.

#### Exported globals

It is expected from the Runtime to export `__heap_base` global indicating the beginning of the heap. It is used by the
Host allocator to prevent memory allocations below that address and avoid clashes with the stack and data sections.

#### Imported vs exported memory

Imported memory works a little better than exported memory since it avoids some edge cases, although it also has some
downsides. Working with exported memory is almost certainly still supported and in fact, this is how it worked in the
beginning. However, the current spec describes that memory should be made available to the Polkadot Runtime for import
under the symbol name `memory`.

#### No concurrency

Parallelism is achieved through Parachains 

```
 ________________________________________________________________________________________
| HOST                                                                                   |
|                                       _________________________________________        |
|                                      |               ALLOCATOR                 |       |
|                                      | ext_allocator_malloc,ext_allocator_free |       |
|                                      |_________________________________________|       |
|     ________________________________________________________|_____________________     |
|    | WASM                                                   | (imported)          |    |
|    |                                                        |                     |    |
|    |             ___________________________________________▼___________________  |    |
|    | (imported) |          |             |            |         [---]           | |    |
| MEMORY  -----►  | ◄- Stack | .Data .BSS  |            | Heap -►                 | |    |
|    |            |__________|_____________|____________|_________________________| |    |
|    |            0     stack base    __data_end    __heap_base         max memory  |    |
|    |                                                                              |    |
|    |______________________________________________________________________________|    |
|                                                                                        |
|________________________________________________________________________________________|

```

### Developer experience

Implementing Wasm functionality makes you go pretty low-level and use some "unsafe" language constructs.

#### Package Structure

* `build` - the output for the the runtime Wasm file.
* `config` - configuration of the used runtime modules (pallets).
* `constants` - constants used in the runtime.
* `env` - stubs for the host provided functions.
* `execution` - runtime execution logic.
* `frame` - runtime modules.
* `primitives` - runtime primitives.
* `runtime` - runtime entry point and tests.
* `utils` - utility functions.
* `tinygo` - submodule for the Tinygo compiler.
* `goscale` - submodule for the SCALE codec.
* `gossamer` - submodule for the Gossamer host, used during development and for running tests.
* `substrate` - submodule for the Substrate host, used for running a network.

