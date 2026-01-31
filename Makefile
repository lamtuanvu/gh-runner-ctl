VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/lamtuanvu/gh-runner-ctl/internal/cli.Version=$(VERSION)"

.PHONY: build install clean test vet

build:
	go build $(LDFLAGS) -o bin/ghr ./cmd/ghr

install:
	go install $(LDFLAGS) ./cmd/ghr

clean:
	rm -rf bin/

test:
	go test ./...

vet:
	go vet ./...
