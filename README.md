# Go-based Polkadot/Substrate Runtimes 

The following repository contains research and PoC on toolchains for building Go-based Polkadot/Substrate runtimes.
The research is funded by [**Web3 Foundation**](https://web3.foundation) **Grants** and developed, maintained by
[**LimeChain**](https://limechain.tech).

[Research Results](./docs/2-go-based-polkadot-runtime.md)

## PoC of a Polkadot Runtime in Go

The steps below will showcase testing a PoC Polkadot Runtime implementation in Go.

**Prerequisites**

- [git](https://git-scm.com/downloads)
- [Go 1.18+](https://golang.org/doc/install)
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
make test
```

**Optional steps**

* Inspect WASM Runtime - [wasmer](https://wasmer.io/)

```bash
wasmer inspect build/runtime.wasm
```

* Convert WASM from binary to text format - [wasm2wat](https://command-not-found.com/wasm2wat)

```bash
wasm2wat build/runtime.wasm -o build/runtime.wat
cat build/runtime.wat
```

