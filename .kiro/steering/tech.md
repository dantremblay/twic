# Tech Stack

## Language & Runtime

- Go (module: `github.com/kassisol/twic`)
- `CGO_ENABLED=0` — pure Go, no CGo dependencies

## Key Dependencies

| Library | Purpose |
|---|---|
| `spf13/cobra` | CLI framework (commands, flags, help) |
| `gorm.io/gorm` + `glebarez/sqlite` | ORM with pure-Go SQLite driver |
| `golang.org/x/term` | Terminal password input |
| `juliengk/go-cert` | Certificate helpers (CA, PKIX, CSR) |
| `kassisol/tsa` | TSA client for auth, cert issuance, CA interaction |
| `juliengk/go-utils` | File/dir utils, input reading, validation |
| `juliengk/stack` | HTTP client and JSON:API response handling |

## Vendored / Local Dependencies

Several dependencies are forked locally under `local/` and redirected via `replace` directives in `go.mod`. These are **not** fetched from upstream — edits go directly in `local/github.com/...`.

## Build System

GNU Make. The `Makefile` handles versioning via git tags/commits and injects build metadata through `-ldflags`.

## Common Commands

```bash
# Build binary (output: bin/twic)
make build

# Run tests
make test

# Run go vet
make vet

# Tidy and vendor deps
make tidy

# Cross-compile (darwin/linux, amd64/arm64)
make cross

# Build Docker image
make image

# Clean build artifacts
make clean

# Print version info
make version
```

## Docker

Multi-stage build: `golang:1.25-alpine` builder → `alpine:3` runtime. Binary is statically linked.
