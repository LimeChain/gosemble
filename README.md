# Go-based Polkadot/Substrate Runtimes 

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

```bash
make build
```

**Run Tests**

After the runtime has been built, we execute standard Go tests with the help of
a [forked version of Gossamer](https://github.com/LimeChain/gossamer), which we use to import necessary Polkadot Host
functionality and interact with the Runtime.

```bash
make test_unit
make test_integration
```

**Optional steps**

* Inspect WASM Runtime - [wasmer](https://wasmer.io/)

```bash
wasmer inspect build/runtime.wasm
```

**Architecture/Development Notes**

* Toolchain - uses a modified Tinygo fork for compiling Go to Wasm compatible with Substrate. Thus a lot of Go's language capabilities can not be applied due to the limited support in Tinygo which affects some of the design decisions.

* SCALE codec - the runtime works only with custom-defined SCALE types compatible with Tinygo. At the boundary where it interacts with the host (memory,storage), those are converted to ordinary Go types.

* Host - uses Gossamer implementation for running integration tests.

* Go pragmas:
  * `//go:build` - provides a separate implementation of the host imported functions in a nonwasm environment (for running tests).
  * `//go:export` - before a function declaration, acts as a function import. The function needs to be referenced somewhere to be actually exported.