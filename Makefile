
# Variables
GO ?= go
CARGO ?= CARGO
PROTOC ?= protoc
MODE ?= server

# Paths
PROTO_GO_OUT := api/gen/go/
PROTO_DIR := api/proto/
BIN := bin

# Exports
export PORT=:8080
export CLIENT_URL=http://localhost:3000 
export PROM_URL=:5050
export REDIS_ADDR=localhost:6379
export REDIS_PASSWD=


# Utils
.PHONY: help
help:
	@echo "Usage:<options> make <target>"
	@echo ""
	@echo "Common targets: "
	@echo "  build-go       Build Go code"
	@echo "  run-go       	Run Go binaries"
	@echo "  test-go        Run Go tests"
	@echo "  build-rs     	Build Rust code"
	@echo "  run-rs       	Run Rust binaries"
	@echo "  test-rs      	Run Rust tests"
	@echo "  proto          Generate protobuf stubs"
	@echo "  docker-build   Build Docker image (WIP)"
	@echo "  docker-run     Run Docker container (WIP)"
	@echo "  clean          Clean build artifacts"
	@echo ""
	@echo "Tip: Run 'make all' to build everything."

.PHONY: clean
clean:
	@echo "Deleting all the artifacts"
	@rm -rf $(PROTO_GO_OUT) $(BIN) ./target

.PHONY: all
all:
	proto
	build-go
	build-rs


# Go Options

.PHONY: build-go
build-go:
	@echo "Building github.com/raghavendrarq/container-dsh"
	@GOOS=linux CGOENABLED=false $(GO) build ./cmd/container-dsh --mode="$(MODE)"

.PHONY: run-go
run-go:
	@echo "Running Container-dsh"
	@$(GO) run ./cmd/container-dsh/ --mode=$(MODE)


# Proto Options

.PHONY: proto
proto:
	@echo "Making Proto files"
	@$(PROTOC) \
	--proto_path=$(PROTO_DIR) \
	--go_opt=module=github.com/raghavendrarq/container-dsh \
	--go-grpc_opt=module=github.com/raghavendrarq/container-dsh \
	--go_out=. \
	--go-grpc_out=. \
	$(PROTO_DIR)/v1/*.proto
	@echo "Done ...."


# Rust Options

.PHONY: build-rs
build-rs:
	@echo "Building rust code"
	@rm -rf ./target
	@$(CARGO) build

.PHONY: run-rs
run-rs:
	@echo "Running rust code"
	@$(CARGO) run

