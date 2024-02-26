# Gosemble

[![Go Report Card](https://goreportcard.com/badge/github.com/LimeChain/gosemble)](https://goreportcard.com/report/github.com/LimeChain/gosemble)
[![codecov](https://codecov.io/github/LimeChain/gosemble/graph/badge.svg?token=48SIN10OBK)](https://codecov.io/github/LimeChain/gosemble)

> [!WARNING]
> Gosemble is in pre-production and the code is not yet audited. Use at your own risk.

Go implementation of Polkadot/Substrate compatible runtimes. For more details, check
the [Official Documentation](https://limechain.github.io/gosemble/)

### Quick Start

#### Prerequisites

- [Git](https://git-scm.com/downloads)
- [Go 1.21](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/install/)
- [Rust](https://docs.substrate.io/install/) (for building the Substrate node)

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
