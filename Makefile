SHELL := /bin/bash
IMAGE := polkawasm/tinygo
TAG := 0.25.0
BRANCH_CONSERVATIVE_GC := new-polkawasm-target-release-$(TAG)
BRANCH_EXTALLOC_GC := new-polkawasm-target-extalloc-gc-release-$(TAG)
CURRENT_DIR := $(shell pwd)
SRC_DIR := /src/examples/wasm/gosemble
BUILD_PATH := build/runtime.wasm

# Build with the standard system installed TinyGo.
sys_tinygo_build:
	@tinygo build -target=./target.json -o=$(BUILD_PATH) runtime/runtime.go

# Build with our forked TinyGo.
.PHONY: build
build:
	@if [[ "$(GC)" == "extalloc" ]]; then \
		cd tinygo; \
		git checkout $(BRANCH_EXTALLOC_GC); \
		cd ..; \
		docker build --tag $(IMAGE):$(TAG)-ext -f tinygo/Dockerfile.polkawasm tinygo; \
		docker run --rm -v $(CURRENT_DIR):$(SRC_DIR) -w $(SRC_DIR) $(IMAGE):$(TAG)-ext /bin/bash -c "tinygo build -target=polkawasm -o=$(SRC_DIR)/$(BUILD_PATH) $(SRC_DIR)/runtime/"; \
		echo "Compiled with extalloc GC..."; \
	else \
		cd tinygo; \
		git checkout $(BRANCH_CONSERVATIVE_GC); \
		cd ..; \
		docker build --tag $(IMAGE):$(TAG) -f tinygo/Dockerfile.polkawasm tinygo; \
		docker run --rm -v $(CURRENT_DIR):$(SRC_DIR) -w $(SRC_DIR) $(IMAGE):$(TAG) /bin/bash -c "tinygo build -target=polkawasm -o=$(SRC_DIR)/$(BUILD_PATH) $(SRC_DIR)/runtime/"; \
		echo "Compiled with conservative GC..."; \
	fi

test: test_unit test_integration

# TODO: ignore the integration tests
test_unit:
	@go test --tags="nonwasmenv" -v ./...

test_integration:
	@go test --tags="nonwasmenv" -v runtime/tests/runtime_test.go
