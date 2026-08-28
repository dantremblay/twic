BINARY     := twic
MODULE     := github.com/kassisol/twic
OUTPUT     ?= bin/$(BINARY)

# Version info
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DIRTY      := $(shell git status --porcelain --untracked-files=no 2>/dev/null)
GIT_TAG    := $(shell git tag -l --contains HEAD 2>/dev/null | head -n 1)

ifdef DIRTY
  GITSTATE := dirty
  VERSION  ?= $(COMMIT)-dirty
else ifdef GIT_TAG
  GITSTATE := clean
  VERSION  ?= $(GIT_TAG)
else
  GITSTATE := clean
  VERSION  ?= $(COMMIT)
endif

LDFLAGS := -s -w \
           -X $(MODULE)/version.Version=$(VERSION) \
           -X $(MODULE)/version.GitCommit=$(COMMIT) \
           -X $(MODULE)/version.GitState=$(GITSTATE) \
           -X $(MODULE)/version.BuildDate=$(shell date +%s)

# Default target
.DEFAULT_GOAL := build

## build: Build the binary
build:
	@echo "Building $(VERSION) ($(COMMIT))"
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(OUTPUT) .

## cross: Cross-compile for multiple platforms
cross:
	@echo "Cross-compiling $(VERSION)"
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(OUTPUT)-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(OUTPUT)-darwin-arm64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(OUTPUT)-linux-amd64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(OUTPUT)-linux-arm64 .

## image: Build the Docker image
image: build
	docker build --build-arg version=$(VERSION) -t kassisol/twic-engine:$(VERSION) -f Dockerfile .

## test: Run tests
test:
	go test ./...

## vet: Run go vet
vet:
	go vet ./...

## tidy: Tidy and vendor dependencies
tidy:
	go mod tidy
	go mod vendor

## generate: Generate man pages and shell completions
generate:
	go run gen/man/genman.go
	go run gen/shellcompletion/genshellcompletion.go

check-nfpm:
	@command -v nfpm >/dev/null 2>&1 || { \
		echo "Error: nfpm not found on PATH."; \
		echo "Install it with:"; \
		echo "    go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest"; \
		echo "and make sure \"\$$(go env GOPATH)/bin\" is on your PATH."; \
		exit 1; \
	}

## package-rpm: Build RPM package (requires nfpm)
package-rpm: RELEASE ?= 1
package-rpm: check-nfpm
	@mkdir -p dist
	VERSION=$(VERSION) RELEASE=$(RELEASE) GOARCH=$(or $(TARGETARCH),amd64) \
		nfpm pkg --packager rpm --target dist/

## package-deb: Build DEB package (requires nfpm)
package-deb: RELEASE ?= 1
package-deb: check-nfpm
	@mkdir -p dist
	VERSION=$(VERSION) RELEASE=$(RELEASE) GOARCH=$(or $(TARGETARCH),amd64) \
		nfpm pkg --packager deb --target dist/

## packages: Build all packages
packages: generate package-rpm package-deb

## clean: Remove build artifacts
clean:
	rm -rf bin/ dist/

## version: Print version info
version:
	@echo "Version:  $(VERSION)"
	@echo "Commit:   $(COMMIT)"
	@echo "GitState: $(GITSTATE)"

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'

.PHONY: build build-static cross image test vet tidy clean version help generate check-nfpm package-rpm package-deb packages
