---
layout: default
title: Inspect
permalink: /build/inspect
---

Install [wasmer](https://wasmer.io/) to get a simple view of the compiled WASM.

```bash
wasmer inspect build/runtime.wasm
```

To inspect the WASM in more detail, and view the actual memory segments, you can install [wabt](https://github.com/WebAssembly/wabt)

```bash
wasm-objdump -x build/runtime.wasm
```