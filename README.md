# prongs

![Build Status](https://img.shields.io/github/actions/workflow/status/thomaslaurenson/prongs/tag.yml?style=flat&logo=github) ![Test Status](https://img.shields.io/github/actions/workflow/status/thomaslaurenson/prongs/tag.yml?style=flat&label=test&logo=github)

![Release Version](https://img.shields.io/github/v/release/thomaslaurenson/prongs?style=flat&logo=github) ![Release downloads](https://img.shields.io/github/downloads/thomaslaurenson/prongs/total?label=downloads&logo=github)

![Go Version](https://img.shields.io/github/go-mod/go-version/thomaslaurenson/prongs?logo=go) ![Code Coverage](https://img.shields.io/badge/Coverage-94%25-blue?logo=go)

Fast, custom security scanner.

## Installation

```bash
curl -fsSL https://github.com/thomaslaurenson/prongs/releases/latest/download/install.sh | bash
```

Or download a binary directly from the [releases page](https://github.com/thomaslaurenson/prongs/releases).

### Container

```bash
docker pull ghcr.io/thomaslaurenson/prongs:latest
docker run --rm -e TARGET_CIDRS=192.168.0.0/24 ghcr.io/thomaslaurenson/prongs:latest scan --all
```

Or build and run locally:

```bash
docker build -t prongs .
docker run --rm -e TARGET_CIDRS=192.168.0.0/24 prongs scan --all
```

## Usage

```
prongs scan --scanner <name>  --target <CIDR|file>
prongs scan --all             --target <CIDR|file>
prongs --version
```

Targets are CIDR ranges or single IPs, supplied via `--target` (repeatable) or a file (one entry per line). If `--target` is omitted, the `TARGET_CIDRS` environment variable is used as a fallback.

### Scanners

| Name | Description | Default |
|---|---|---|
| `password-ssh` | Detects SSH servers accepting password authentication | yes |
| `accessible-rdp` | Detects RDP services accepting unauthenticated connections | yes |
| `accessible-db` | Detects databases accepting unauthenticated connections | yes |

### Examples

```bash
# Run one scanner against a single network
prongs scan --scanner password-ssh --target 192.168.0.0/24

# Run all default scanners against multiple networks
prongs scan --all --target 192.168.0.0/24 --target 10.0.0.0/24

# Load targets from a file
prongs scan --all --target targets.txt

# Pretty-print output
prongs scan --all --target 192.168.0.0/24 --output pretty

# Use environment variable for targets
TARGET_CIDRS=192.168.0.0/24 prongs scan --all
```
