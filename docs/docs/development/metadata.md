---
layout: default
title: Metadata
permalink: /development/metadata
---

# Metadata

The [Runtime Metadata](https://spec.polkadot.network/sect-metadata) consists of all the necessary information on how to interact with the Runtime. Given that
runtimes are upgradeable, changes to runtimes are applied to their metadata as well.


## Supported versions

Gosemble supports the following versions:

* [14](https://github.com/LimeChain/gosemble/blob/develop/primitives/types/metadata_v14.go#L9) (default)
* [15](https://github.com/LimeChain/gosemble/blob/develop/primitives/types/metadata_v15.go#L9)

## Generation process

The original implementation of Gosemble Metadata was based on hard-coded definition types for all the necessary
runtime types.

After the codebase was refactored and modularised, we began generating the metadata definitions with the help
of _Go_'s `reflect` package. We introduced a 
[metadata generator](https://github.com/LimeChain/gosemble/blob/develop/primitives/types/metadata_generator.go),
which builds and recursively adds the metadata type definitions.

**Implementation is under development**.

Currently, the metadata generator can generate metadata type definitions from:
    
* _Go_ types
* A module's extrinsic calls
* A module's errors
* A module's constants
* Extrinsic signed extensions

The goal is to add support for a module's storage types and remove all hard-coded definition types.

## Future guidelines

We intend to refactor the metadata generation to be at **compile-time**.
