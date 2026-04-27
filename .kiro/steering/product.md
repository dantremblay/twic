# Product: TWIC

TWIC (part of the Harbormaster / HBM suite) is a CLI tool for managing Docker TLS certificates and client profiles.

## What it does

- Manages TLS certificates used to connect to Docker daemons securely
- Communicates with a TSA (Trust Service Adapter) server for certificate issuance and authentication
- Stores certificates and profiles locally in an SQLite database
- Lets users create named profiles that bundle a certificate with a Docker host endpoint
- Outputs shell environment variables (`DOCKER_HOST`, `DOCKER_TLS_VERIFY`, `DOCKER_CERT_PATH`) so users can connect to Docker hosts via TLS

## Core workflow

1. User adds a TSA URL and authenticates
2. TSA issues a token if the user is authorized
3. User creates a certificate via the TSA/CA
4. User creates a profile linking the certificate to a Docker host
5. User runs `twic profile env <name>` to get Docker env vars for their shell

## Key entities

- **Cert** — a TLS certificate (type: `client` or `engine`) tied to a TSA URL
- **Profile** — a named pairing of a client cert + Docker host URL
- **Engine** — represents a Docker engine/host managed through TSA

## License

GPL-3.0
