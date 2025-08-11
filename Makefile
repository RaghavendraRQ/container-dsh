.PHONY: build
build:
	echo "Building container-dsh"
	go build ./cmd/container-dsh --mode="$1"