# Project Structure

```
.
├── main.go              # Entry point — creates root cobra command, runs it
├── packages.go          # Blank-import for storage driver registration (SQLite)
├── Makefile             # Build, test, cross-compile, Docker image targets
├── Dockerfile           # Multi-stage Docker build
│
├── cli/command/         # CLI command tree (one package per command group)
│   ├── commands/        # Root command wiring — AddCommands() registers all subcommands
│   ├── access/          # `twic access` — TSA authentication
│   ├── cert/            # `twic cert` — add, list, remove, renew certificates
│   ├── engine/          # `twic engine` — create, info, remove, renew engine certs
│   ├── profile/         # `twic profile` — add, list, remove, env, status
│   └── system/          # `twic version`
│
├── pkg/                 # Shared utility packages
│   ├── cert/            # Certificate file I/O helpers
│   ├── date/            # Date formatting
│   ├── format/          # Output formatting
│   ├── input/           # Interactive user input
│   ├── sysutil/         # OS/user utilities
│   ├── urlutil/         # URL parsing
│   └── validate/        # Name and port validation
│
├── storage/             # Persistence layer
│   ├── storage.go       # Driver registry (RegisterDriver / NewDriver)
│   └── driver/
│       ├── driver.go    # Storager interface
│       ├── types.go     # CertResult, ProfileResult DTOs
│       └── sqlite/      # SQLite implementation (GORM)
│           ├── sqlite.go   # init() registers driver, DB open/migrate/close
│           ├── models.go   # GORM models (Cert, Profile)
│           ├── cert.go     # Cert CRUD
│           └── profile.go  # Profile CRUD
│
├── version/             # Version info struct, populated via ldflags at build time
│
├── local/               # Forked/local dependencies (replace directives in go.mod)
│   └── github.com/
│       ├── howeyc/gopass/
│       ├── juliengk/{go-cert, go-utils, stack}/
│       └── kassisol/tsa/
│
├── vendor/              # Vendored third-party modules
├── gen/                 # Code generators (man pages, shell completions)
├── scripts/             # Build/packaging/dev helper scripts
└── docs/                # Hugo-style documentation (Harbormaster site)
```

## Conventions

- **CLI commands**: Each command group is a package under `cli/command/`. The package exposes `NewCommand() *cobra.Command` (parent) and unexported `newXxxCommand()` for subcommands. Business logic lives in `runXxx()` functions.
- **Storage driver pattern**: Drivers register themselves via `init()` calling `storage.RegisterDriver()`. The `Storager` interface in `storage/driver/driver.go` defines the contract. New drivers go in `storage/driver/<name>/`.
- **Blank imports for registration**: `packages.go` imports storage drivers for their `init()` side effects.
- **pkg/ utilities**: Small, focused packages. No cross-dependencies between `pkg/` packages.
- **Local forks**: Dependencies under `local/` are project-controlled forks. Changes to these are committed directly; they are not synced from upstream.
