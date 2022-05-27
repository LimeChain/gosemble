#!/usr/bin/env bash

# -scheduler=none
tinygo build \
  -wasm-abi=generic \
  -target=wasi \
  -o=build/dev_runtime.wasm \
  runtime/runtime.go