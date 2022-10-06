CURRENT_DIR = $(shell pwd)

# Build with the standard system installed TinyGo.
old:
	@tinygo build -target=./target.json -o=build/runtime.wasm runtime/runtime.go

# Build with our forked TinyGo.
.PHONY: build
build:
	@docker build --tag polkawasm/tinygo:0.25.0 -f tinygo/Dockerfile.polkawasm tinygo
	@docker run --rm -v $(CURRENT_DIR):/src/examples/wasm/gosemble -w /src/examples/wasm/gosemble polkawasm/tinygo:0.25.0 /bin/bash \
	-c "tinygo build -target=polkawasm -o=/src/examples/wasm/gosemble/build/runtime.wasm /src/examples/wasm/gosemble/runtime/"

test:
	@go test -v runtime/runtime_test.go
