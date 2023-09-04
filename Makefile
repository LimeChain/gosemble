SHELL := /bin/bash
CURRENT_DIR = $(shell pwd)
SRC_DIR = /src/examples/wasm/gosemble
BUILD_PATH = build/runtime.wasm
IMAGE = polkawasm/tinygo
TAG = 0.29.0

build-docker:
	docker build --tag $(IMAGE):$(TAG)-extallocleak -f tinygo/Dockerfile.polkawasm tinygo; \
	docker run --rm -v $(CURRENT_DIR):$(SRC_DIR) -w $(SRC_DIR) $(IMAGE):$(TAG)-extallocleak /bin/bash -c "tinygo build -target=polkawasm -o=$(SRC_DIR)/$(BUILD_PATH) $(SRC_DIR)/runtime/"; \
	echo "build - tinygo version: ${TAG}, gc: extallocleak"; \
	
build-wasi-libc:
	@cd tinygo && \
	if [ ! -e lib/wasi-libc/Makefile ]; then \
		echo "Submodules have not been downloaded. Please download them using:\n git submodule update --init"; \
		exit 1; \
	fi && \
	if [ -e lib/wasi-libc/sysroot ]; then \
		echo "wasi-libc is already built"; \
	else \
		echo "building wasi-libc"; \
		cd lib/wasi-libc && make clean && make -j4 EXTRA_CFLAGS="-O2 -g -DNDEBUG" MALLOC_IMPL=none; \
	fi

build-tinygo: build-wasi-libc
	@cd tinygo; \
		go install;
	@tinygo version

build-release: build-tinygo
	@tinygo build --no-debug -target=polkawasm -o=$(BUILD_PATH) runtime/runtime.go

build-dev: build-tinygo
	@tinygo build -target=polkawasm -o=$(BUILD_PATH) runtime/runtime.go

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