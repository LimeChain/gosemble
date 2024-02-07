---
layout: default
title: Toolchain
permalink: /overview/toolchain
---

# Toolchain 🛠️

Since we use toolchain based on [modified version of TinyGo](https://github.com/LimeChain/tinygo) for compiling Go to Wasm, some of Go's
language capabilities cannot be applied due to the limited support in TinyGo which also affects some of the design decisions.

## 1. Tinygo modifications

There are several changes made to TinyGo in order to make it compatible with Polkadot's Wasm target:

### 1.1. New target 🎯

* `polkawasm` - targeting standalone **Wasm MVP**, similar to Rust's `wasm32-unknown-unknown`.

### 1.2. Custom garbage collector 🗑️

* `extalloc` - custom GC that utilizes an external allocator (as per Polkadot specification). It is conservative, tracing (mark and sweep) garbage collector that relies on an external memory allocator (via`ext_allocator_malloc`, `ext_allocator_free`) for the WebAssembly (`polkawasm`) target.
* `extalloc_leaking` - leaking GC implementation that only allocates memory through the external allocator, but never frees it (not a real GC), however, it is useful for testing purposes and performance comparisons.

### 1.3. Wasi-libc 🔩

* fork with disabled **bulk memory** operations.

### 1.4. Binaryen 🔩

* version that lowers away the **sign extension** operations in `wasm-opt`.