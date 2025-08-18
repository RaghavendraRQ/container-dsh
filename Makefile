.PHONY: build
build:
	echo "Building container-dsh"
	GOOS=linux CGOENABLED=false go build ./cmd/container-dsh --mode="$(MODE)"

run:
	echo "Running Container-dsh"
	echo "$(EXTRA)"
	go run ./cmd/container-dsh/ 

monitor:
	echo "Monitoring Container-dsh"
	go run ./cmd/container-dsh/ --mode="prometheus"
