# Set the Go version from your go.mod file
GO_VERSION := 1.24
# The name of your final binary
BINARY_NAME := rom-manager
# The path to your main package
MAIN_PACKAGE := ./cmd/app

DEVICE_HOST := 192.168.2.114
DEVICE_USER := root
DEVICE_PORTS_FOLDER := /roms/ports

.PHONY: build-aarch64-binary build-aarch64-dist build-aarch64 clean check-deps upload

upload: build-aarch64-dist
	@echo "Uploading to device"
	scp -r dist/* $(DEVICE_USER)@$(DEVICE_HOST):$(DEVICE_PORTS_FOLDER)

# Build the application for linux/aarch64 using Docker
build-aarch64-binary:
	@echo "Building for linux/arm64..."
	@mkdir -p bin
	@docker build -t rom-manager-builder -f Dockerfile.build .
	@docker run --rm \
		-v $(shell pwd):/app \
		-v rom-manager-gomod:/go/pkg/mod \
		-v rom-manager-gocache:/root/.cache/go-build \
		-e GOOS=linux \
		-e GOARCH=arm64 \
		-e CGO_ENABLED=1 \
		-e CC=aarch64-linux-gnu-gcc \
		-e PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig \
		rom-manager-builder \
		go build -buildvcs=false -o bin/$(BINARY_NAME)-linux-aarch64 $(MAIN_PACKAGE)
	@echo "Build complete: bin/$(BINARY_NAME)-linux-aarch64"

build-aarch64-dist: build-aarch64-binary
	@echo "Building distribution for linux/arm64..."
	@cp assets/RomManager.sh dist/RomManager.sh
	@mkdir -p dist/RomManager
	@cp bin/$(BINARY_NAME)-linux-aarch64 dist/RomManager/


# Check dynamic dependencies of the aarch64 binary
check-deps: build-aarch64
	@echo "Checking dynamic dependencies for bin/$(BINARY_NAME)-linux-aarch64..."
	@docker run --rm \
		-v $(shell pwd)/bin:/app/bin \
		rom-manager-builder \
		aarch64-linux-gnu-readelf -d /app/bin/$(BINARY_NAME)-linux-aarch64 | grep NEEDED || true

# Clean up the build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf bin
	@echo "Clean complete."

