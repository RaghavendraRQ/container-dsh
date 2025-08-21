.PHONY: build
build:
	@echo "Building github.com/raghavendrarq/container-dsh"
	@GOOS=linux CGOENABLED=false go build ./cmd/container-dsh --mode="$(MODE)"

.PHONY: run
run:
	@echo "Running Container-dsh"
	@echo "$(EXTRA)"
	@go run ./cmd/container-dsh/ 

.PHONY: monitor
monitor:
	@echo "Monitoring Container-dsh"
	@go run ./cmd/container-dsh/ --mode="prometheus"

