CURRENT_DIR = $(shell pwd)
SRC_DIR = /src/examples/wasm/gosemble
BUILD_PATH = build/runtime.wasm
TARGET = polkawasm
GC = custom-extalloc-leak
VERSION = 0.30.0
IMAGE = tinygo/${TARGET}

WASMOPT_PATH = /tinygo/lib/binaryen/bin/wasm-opt

DOCKER_BUILD_TINYGO = docker build --tag $(IMAGE):$(VERSION)-$(GC) -f tinygo/Dockerfile.$(TARGET) tinygo
DOCKER_RUN_TINYGO = docker run --rm -v $(CURRENT_DIR):$(SRC_DIR) -w $(SRC_DIR) $(IMAGE):$(VERSION)-$(GC) /bin/bash -c

TINYGO_BUILD_COMMAND_NODEBUG = tinygo build --no-debug -target=$(TARGET)
TINYGO_BUILD_COMMAND = tinygo build -target=$(TARGET)

RUNTIME_BUILD_NODEBUG = "WASMOPT="$(WASMOPT_PATH)" $(TINYGO_BUILD_COMMAND_NODEBUG) -o=$(SRC_DIR)/$(BUILD_PATH) $(SRC_DIR)/runtime/"
RUNTIME_BUILD = "WASMOPT="$(WASMOPT_PATH)" $(TINYGO_BUILD_COMMAND) -o=$(SRC_DIR)/$(BUILD_PATH) $(SRC_DIR)/runtime/"

clear-wasi-libc:
	@cd tinygo/lib/wasi-libc && \
	make clean

clear-binaryen:
	@cd tinygo/lib/binaryen && \
	rm -rf CMakeCache.txt

build-docker-release: clear-binaryen
	@set -e; \
	$(DOCKER_BUILD_TINYGO);
	$(DOCKER_RUN_TINYGO) $(RUNTIME_BUILD_NODEBUG); \
	echo "Build - tinygo version: ${VERSION}, gc: ${GC} (no debug)"
	
build-docker-dev: clear-binaryen
	@set -e; \
	$(DOCKER_BUILD_TINYGO);
	$(DOCKER_RUN_TINYGO) $(RUNTIME_BUILD); \
	echo "Build - tinygo version: ${VERSION}, gc: ${GC}"

build-wasi-libc: clear-wasi-libc
	@cd tinygo/lib/wasi-libc && \
	if [ ! -e Makefile ]; then \
		echo "Submodules have not been downloaded. Please download them using:\n git submodule update --init"; \
		exit 1; \
	fi && \
	echo "Building \"wasi-libc\""; \
	make -j4 EXTRA_CFLAGS="-O2 -g -DNDEBUG" MALLOC_IMPL=none; \

build-binaryen: clear-binaryen
	@cd tinygo/lib/binaryen && \
	if [ ! -e Makefile ]; then \
		echo "Submodules have not been downloaded. Please download them using:\n git submodule update --init"; \
		exit 1; \
	fi && \
	echo "Building \"binaryen\""; \
	cmake . && make; \

build-tinygo:
	@cd tinygo && \
	if [ ! -e lib/wasi-libc/sysroot ]; then \
		echo "Need to build wasi-libc. Please run: \"make build-wasi-libc\""; \
		exit 1; \
	fi; \
	if [ ! -e lib/binaryen/bin/wasm-opt ]; then \
		echo "Need to build binaryen. Please run: \"make build-binaryen\""; \
		exit 1; \
	fi; \
	echo "Building \"tinygo\""; \
	go install; \
	tinygo version; \

build-release: build-tinygo
	@echo "Building \"runtime.wasm\" (no-debug)"; \
	WASMOPT="$(CURRENT_DIR)/$(WASMOPT_PATH)" $(TINYGO_BUILD_COMMAND_NODEBUG) -o=$(BUILD_PATH) runtime/runtime.go

build-dev: build-tinygo
	@echo "Building \"runtime.wasm\""; \
	WASMOPT="$(CURRENT_DIR)/$(WASMOPT_PATH)" $(TINYGO_BUILD_COMMAND) -o=$(BUILD_PATH) runtime/runtime.go

start-network:
	cp build/runtime.wasm polkadot-sdk/substrate/bin/node-template/runtime.wasm; \
	cd polkadot-sdk/substrate/bin/node-template/node; \
	cargo build --release; \
	cd ../../../..; \
	WASMTIME_BACKTRACE_DETAILS=1 RUST_LOG=runtime=trace ./target/release/node-template --dev --execution=wasm

test: test-unit test-integration

test-unit:
	@go test --tags "nonwasmenv" -cover -v `go list ./... | grep -v runtime`

test-integration:
	@go test --tags="nonwasmenv" -v ./runtime/...