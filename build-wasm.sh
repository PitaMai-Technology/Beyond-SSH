#!/bin/bash
cd wasm
GOOS=js GOARCH=wasm go build -o main.wasm main.go
echo "WASM build complete. main.wasm generated."
