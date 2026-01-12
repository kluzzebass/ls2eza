VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS = -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: all build clean test install

all: build

build:
	go build $(LDFLAGS) -o ls2eza

clean:
	rm -f ls2eza

test:
	go test -v ./...

install:
	go install $(LDFLAGS)

# Cross-compilation targets (matching eza platforms + macOS)
.PHONY: build-all
build-all: build-linux build-darwin build-windows

build-linux:
	# glibc builds
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/ls2eza-x86_64-unknown-linux-gnu
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/ls2eza-aarch64-unknown-linux-gnu
	GOOS=linux GOARCH=arm GOARM=6 go build $(LDFLAGS) -o dist/ls2eza-arm-unknown-linux-gnueabihf
	# musl builds (static)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/ls2eza-x86_64-unknown-linux-musl
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/ls2eza-aarch64-unknown-linux-musl

build-darwin:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/ls2eza-x86_64-apple-darwin
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/ls2eza-aarch64-apple-darwin

build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/ls2eza-x86_64-pc-windows-gnu.exe
