---
layout: default
permalink: /
---

# Overview 👀 
 
## Framework for building Parachains 🎨

This is an alternative solution for building Polkadot/Substrate-compatible runtimes in *Go*.
It aims to streamline the process of building a parachain with emphasis on simplicity and ease of use over configurability and feature richness. 
It is designed to be straightforward to understand and use, while still providing the necessary tools to build a parachain.

## Why choose an alternative implementation in Go 🌱

While there are several implementations of Polkadot Hosts in [Rust](https://github.com/paritytech/substrate), 
[C++](https://github.com/soramitsu/kagome), and [Go](https://github.com/ChainSafe/gossamer), the only option for writing
Polkadot Runtimes is in [Rust](https://github.com/paritytech/substrate). Writing Polkadot Runtimes in *Go* is exciting,
mainly because of *Go*'s simplicity and automatic memory management. It is a modern, powerful, and fast language, backed
by Google and used in many of their software, thus making it an ideal candidate for implementing Polkadot Runtimes.


## How it started 💡

This project is a result of [previous research](https://github.com/LimeChain/gosemble-research), funded by [Web3 Foundation](https://web3.foundation) **Grants**, and developed and maintained by [LimeChain](https://limechain.tech). The research provides conclusions if *Go* is a suitable choice for writing Polkadot Runtimes and further aids the development of a *Go* [**toolchain**](./overview/toolchain), capable of producing compatible runtimes.

If you are new to Polkadot, please check our [**onboarding guide**](./development/onboarding) on how to get started with the project.
