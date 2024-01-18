# Gosemble

[![codecov](https://codecov.io/github/LimeChain/gosemble/graph/badge.svg?token=48SIN10OBK)](https://codecov.io/github/LimeChain/gosemble)

> **Warning**
> The Gosemble is in pre-production

Go implementation of Polkadot/Substrate compatible runtimes. For more details, check
the [Official Documentation](https://limechain.github.io/gosemble/)

### Quick Start

#### Prerequisites

- [Git](https://git-scm.com/downloads)
- [Go 1.21](https://golang.org/doc/install)
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

#### Build

To build a runtime, execute: 

```bash
make build-docker-release
```

#### Start a local network

After the runtime is built, start a local network using Substrate host:

```bash
make start-network
```

#### Run Tests

After the Runtime is built, execute the tests with the help of [Gossamer](https://github.com/LimeChain/gossamer), which
is used to import necessary Polkadot Host functionality and interact with the Runtime.

```bash
make test-unit
make test-integration
```

#### Benchmarking

// todo quick intro to benchmarking in polkadot.

// todo purpose of benchmarking - weights, which are essential part of the consensus

// todo specifics of polkadot testing - linear components, instance, whitelisted keys, db reads/writes, etc

// todo quick overview of example benchmark test in gosemble, point out important details - benchmarking tag, build_config.json which is generated when building with benchmark-test command through BUILD_RUNTIME flag, etc.

Write benchmarks
// todo example benchmark test

Build benchmarking runtime

```bash
# build with local tinygo binary
make build-benchmarking

# build with docker
make build-docker-benchmarking
```

Run benchmarks

```bash
make benchmark steps=5 
```
