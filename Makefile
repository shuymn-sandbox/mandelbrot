.PHONY: all
all: server

.PHONY: build
build: build-wasm

.PHONY: build-wasm
build-wasm:
	tinygo build -o ./static/wasm/main.wasm ./lib/main.go

.PHONY: server
server: build
	go run main.go
