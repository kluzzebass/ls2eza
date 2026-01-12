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

# Cross-compilation targets
.PHONY: build-all
build-all: build-linux build-darwin build-windows

build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/ls2eza-linux-amd64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/ls2eza-linux-arm64

build-darwin:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/ls2eza-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/ls2eza-darwin-arm64

build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/ls2eza-windows-amd64.exe
