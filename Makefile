CURRENT_DIR = $(shell pwd)
SRC_DIR = /src/examples/wasm/gosemble

# Build with the standard system installed TinyGo.
old:
	@tinygo build -target=./target.json -o=build/runtime.wasm runtime/runtime.go

# Build with our forked TinyGo.
.PHONY: build
build:
	@docker build --tag polkawasm/tinygo:0.25.0 -f tinygo/Dockerfile.polkawasm tinygo

	@docker run --rm -v $(CURRENT_DIR):$(SRC_DIR) -w $(SRC_DIR) polkawasm/tinygo:0.25.0 /bin/bash \
	-c "tinygo build -target=polkawasm -o=$(SRC_DIR)/build/runtime.wasm $(SRC_DIR)/runtime/"

test: test_unit test_integration

test_unit:
	@go test --tags="nonwasmenv" -v ./...

test_integration:
	@go test --tags="nonwasmenv" -v runtime/runtime_test.go
