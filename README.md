# HBM TWIC

TWIC is an open source project for managing Docker certificates to connect to the Docker daemon using TLS.

1. Add a TSA url and login using credentials.
4. TWIC request a certificate with the auto-generated private key for the profile using the token provided when authenticated to TSA.
2. If user authorized to use the Docker host, TSA sends a token to TWIC.
5. If authorized, CA sends the new certificate to TWIC.
3. User add a new profile to connect to Docker host.
6. User can use new profile to set Docker environment variables for connecting to Docker host using TLS.

## Building

### Prerequisites

- Go 1.26+ 
- [nFPM](https://nfpm.goreleaser.com/) (for package builds only)

```bash
go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
```

### Build the binary

```bash
make build
```

The binary is written to `bin/twic`.

### Cross-compile

```bash
make cross
```

Produces binaries for linux/amd64, linux/arm64, darwin/amd64, and darwin/arm64 under `bin/`.

### Build packages (RPM and DEB)

```bash
# Build both RPM and DEB
make packages

# Or individually
make package-rpm
make package-deb
```

Packages are written to `dist/`. The version is derived from git tags automatically.

To build for a different architecture:

```bash
TARGETARCH=arm64 make packages
```

## Getting Started & Documentation

All documentation is available on the [Harbormaster website](http://harbormaster.io/docs/twic/).

## User Feedback

### Issues

If you have any problems with or questions about this application, please contact us through a [GitHub](https://github.com/kassisol/twic/issues) issue.
