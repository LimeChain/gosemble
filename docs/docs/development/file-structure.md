---
layout: default
title: File structure
permalink: /development/file-structure
---

# File structure 📁

* `api` - api interface providing access to the modules (pallets) functionality.
* `benchmarking` - cli and utilities for running benchmark tests.
* `build` - the output directory for the compiled Wasm file.
* `constants` - constants used in the runtime.
* `docs` - project documentation.
* `env` - stubs for the host-provided functions.
* `execution` - runtime execution logic.
* `frame` - runtime modules (pallets).
* `hooks` - hooks implemented by the modules.
* `mocks` - mock implementations for testing.
* `primitives` - runtime primitive types and host functions.
* `runtime` - runtime entry point and integration tests.
* `utils` - utility functions.
* `scripts` - scripts used during deployment.
* `tinygo` - submodule for the TinyGo compiler, used for WASM compilation.
* `goscale` - submodule for the SCALE codec.
* `gossamer` - submodule for the Gossamer host, used during development and for running tests.
* `polkadot-sdk` - submodule for the Substrate host, used for running a network.
