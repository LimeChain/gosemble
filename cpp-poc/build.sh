#!/bin/bash
set -e

clang++ \
	--target=wasm32 \
	-nostdlibinc \
	-nostdlib \
	-O3 \
	-flto \
  -c \
	-o runtime.o \
	runtime/runtime.cpp

wasm-ld \
	--no-entry \
	--export="__heap_base" \
  --export="__data_end" \
	--export-dynamic \
	--export-table \
	--import-memory \
	--allow-undefined \
	--lto-O3 \
	runtime.o \
	-o runtime.wasm

rm -rf runtime.o

mv runtime.wasm build/dev_runtime.wasm

# --include-directory="/opt/homebrew/Cellar/llvm/14.0.6_1/include/c++/v1" \
# --library-directory="/opt/homebrew/Cellar/llvm/14.0.6_1/lib" \

# docker build -t radkomih/cpp-dev:1.0 .
# docker run -it radkomih/cpp-dev:1.0