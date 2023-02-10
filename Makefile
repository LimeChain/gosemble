SHELL := /bin/bash
CURRENT_DIR = $(shell pwd)
SRC_DIR = /src/examples/wasm/gosemble
VERSION = 0.25.0
BRANCH_CONSERVATIVE_GC = new-polkawasm-target-release-$(VERSION)
BRANCH_EXTALLOC_GC = new-polkawasm-target-extalloc-gc-release-$(VERSION)

# Build with the standard system installed TinyGo.
sys_tinygo_build:
	@tinygo build -target=./target.json -o=build/runtime.wasm runtime/runtime.go

# Build with our forked TinyGo.
.PHONY: build
build:
	@if [ $(GC) = extalloc ]; then \
		cd tinygo; \
		git checkout $(BRANCH_EXTALLOC_GC); \
		cd ..; \
		docker build --tag polkawasm/tinygo:$(VERSION)-ext -f tinygo/Dockerfile.polkawasm tinygo; \
		docker run --rm -v $(CURRENT_DIR):$(SRC_DIR) -w $(SRC_DIR) polkawasm/tinygo:$(VERSION)-ext /bin/bash -c "tinygo build -target=polkawasm -o=$(SRC_DIR)/build/runtime.wasm $(SRC_DIR)/runtime/"; \
		echo "Compiled with extalloc GC..."; \
	else \
		cd tinygo; \
		git checkout $(BRANCH_CONSERVATIVE_GC); \
		cd ..; \
		docker build --tag polkawasm/tinygo:$(VERSION) -f tinygo/Dockerfile.polkawasm tinygo; \
		docker run --rm -v $(CURRENT_DIR):$(SRC_DIR) -w $(SRC_DIR) polkawasm/tinygo:$(VERSION) /bin/bash -c "tinygo build -target=polkawasm -o=$(SRC_DIR)/build/runtime.wasm $(SRC_DIR)/runtime/"; \
		echo "Compiled with conservative GC..."; \
	fi

test: test_unit test_integration

# TODO: ignore the integration tests
test_unit:
	@go test --tags="nonwasmenv" -v ./...

test_integration:
	@go test --tags="nonwasmenv" -v runtime/tests/runtime_test.go
