# Gosemble

> **Warning**
> The Gosemble is in pre-production

Go implementation of Polkadot/Substrate compatible runtimes. For more details, check
the [Official Documentation](https://limechain.github.io/gosemble/)

### Quick Start

#### Prerequisites

- [Git](https://git-scm.com/downloads)
- [Go 1.19+](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/install/)
- [Rust](https://docs.substrate.io/install/)

#### Clone the repository

```bash
git clone https://github.com/LimeChain/gosemble.git
cd gosemble
```

#### Pull all necessary git submodules

```bash
git submodule update --init --recursive
```

#### Build a Runtime

Using our [fork of TinyGo](https://github.com/LimeChain/tinygo), there are currently two options to choose from for the
GC implementation. Modify the `GC` environment variable to switch between them.

##### Extalloc GC

It works with the host's external allocator as per specification.

```bash
make build
```

##### Conservative GC

It is used only for **development** and **testing** and works by using a different heap base offset from the allocator's
one (as a workaround), so the GC can use a separate heap region for its allocations and not interfere with the
allocator's region.

```bash
GC="conservative" make build
```

#### Run Tests

After the Runtime is built, execute the tests with the help of [Gossamer](https://github.com/LimeChain/gossamer), which
is used to import necessary Polkadot Host functionality and interact with the Runtime.

```bash
make test_unit
make test_integration
```

#### Start a local network

Once you build the runtime wasm blob, you can start a local network using Substrate as a host.

```bash
make start-network
```
