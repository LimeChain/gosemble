CURRENT_DIR = $(shell pwd)
SRC_DIR = /src/examples/wasm/gosemble
BUILD_PATH = build/runtime.wasm
TARGET = polkawasm
GC = extallocleak
VERSION = 0.29.0
IMAGE = tinygo/${TARGET}

DOCKER_BUILD_TINYGO = docker build --tag $(IMAGE):$(VERSION)-$(GC) -f tinygo/Dockerfile.$(TARGET) tinygo
DOCKER_RUN_TINYGO = docker run --rm -v $(CURRENT_DIR):$(SRC_DIR) -w $(SRC_DIR) $(IMAGE):$(VERSION)-$(GC) /bin/bash -c

TINYGO_BUILD_COMMAND_NODEBUG = tinygo build --no-debug -target=$(TARGET)
TINYGO_BUILD_COMMAND = tinygo build -target=$(TARGET)

RUNTIME_BUILD_NODEBUG = "$(TINYGO_BUILD_COMMAND_NODEBUG) -o=$(SRC_DIR)/$(BUILD_PATH) $(SRC_DIR)/runtime/"
RUNTIME_BUILD = "$(TINYGO_BUILD_COMMAND) -o=$(SRC_DIR)/$(BUILD_PATH) $(SRC_DIR)/runtime/"

build-wasi-libc:
	@cd tinygo/lib/wasi-libc && \
	if [ ! -e Makefile ]; then \
		echo "Submodules have not been downloaded. Please download them using:\n git submodule update --init"; \
		exit 1; \
	fi && \
	echo "Building \"wasi-libc\""; \
	make clean && make -j4 EXTRA_CFLAGS="-O2 -g -DNDEBUG" MALLOC_IMPL=none; \

build-tinygo:
	@cd tinygo && \
	if [ -e lib/wasi-libc/sysroot ]; then \
		go install; \
		tinygo version; \
	else \
		echo "Need to build wasi-libc first. Please run: \"make build-wasi-libc\""; \
		exit 1; \
	fi

build-docker-release:
	@set -e; \
	$(DOCKER_BUILD_TINYGO); \
	$(DOCKER_RUN_TINYGO) $(RUNTIME_BUILD_NODEBUG); \
	echo "Build - tinygo version: ${VERSION}, gc: ${GC} (no debug)"
	
build-docker-dev:
	@set -e; \
	$(DOCKER_BUILD_TINYGO); \
	$(DOCKER_RUN_TINYGO) $(RUNTIME_BUILD); \
	echo "Build - tinygo version: ${VERSION}, gc: ${GC}"

build-release: build-tinygo
	@$(TINYGO_BUILD_COMMAND_NODEBUG) -o=$(BUILD_PATH) runtime/runtime.go

build-dev: build-tinygo
	@$(TINYGO_BUILD_COMMAND) -o=$(BUILD_PATH) runtime/runtime.go

start-network:
	cp build/runtime.wasm substrate/bin/node-template/runtime.wasm; \
	cd substrate/bin/node-template; \
	cargo build --release; \
	cd ../..; \
	WASMTIME_BACKTRACE_DETAILS=1 RUST_LOG=runtime=trace ./target/release/node-template --dev --execution Wasm

test: test-unit test-integration

test-unit:
	@go test --tags "nonwasmenv" -v `go list ./... | grep -v runtime`

test-integration:
	@go test --tags="nonwasmenv" -v ./runtime/...