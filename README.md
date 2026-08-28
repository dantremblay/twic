# HBM TWIC

TWIC is an open source project for managing Docker certificates to connect to the Docker daemon using TLS.

## How it works

1. The user adds a TSA URL and logs in with their credentials.
2. If the user is authorized to use the Docker host, TSA returns a token to TWIC.
3. The user adds a new profile for connecting to the Docker host.
4. TWIC requests a certificate with an auto-generated private key for the profile, using the token from TSA.
5. If authorized, the CA sends the new certificate back to TWIC.
6. The user applies the profile to set the Docker environment variables needed to connect to the Docker host over TLS.

## Requirements

- A running [TSA](https://github.com/kassisol/tsa) server to issue certificates
- Docker on the host you want to connect to
- Client certificate commands (`cert`, `profile`) run as a regular user; engine certificate commands (`engine`) run as root

## Installation

Install the prebuilt package for your distribution.

```bash
# RHEL / Rocky / AlmaLinux / SLES
sudo dnf install ./twic-<version>.x86_64.rpm

# Debian / Ubuntu
sudo dpkg -i ./twic_<version>_amd64.deb
```

Alternatively, build the binary yourself (see [Building](#building)) and copy `bin/twic` into a directory on your `PATH`.

## Usage

A typical session to connect the Docker CLI to a remote host:

```bash
# Request a client certificate from TSA
twic cert add mycert --tsa-url https://tsa.example.com --username alice

# Create a profile pointing at the Docker host
twic profile add myhost --cert-name mycert --docker-host docker.example.com

# Load the Docker environment variables for the profile
eval $(twic profile env myhost)

# The Docker CLI now talks to the remote host over TLS
docker ps
```

Use `twic <command> --help` to see all available options.

### Engine certificates

On the Docker host itself, the `engine` commands manage the daemon's server certificate (run as root):

```bash
# Create the engine certificate (writes to /etc/docker/tls)
sudo twic engine create --common-name docker.example.com --tsa-url https://tsa.example.com --username alice

# Show the current engine certificate (CN, TSA URL, alt names, expiry)
sudo twic engine info

# Renew the engine certificate before it expires, preserving the CN and
# Subject Alternative Names of the existing certificate
sudo twic engine renew --username alice

# Or renew using an existing TSA token instead of username/password
sudo twic engine renew --token <token>
```

After renewing, restart the Docker daemon for the new certificate to take effect:

```bash
sudo systemctl restart docker
```

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
