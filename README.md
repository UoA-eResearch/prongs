# prongs

![Build Status](https://img.shields.io/github/actions/workflow/status/thomaslaurenson/prongs/tag.yml?style=flat&logo=github) ![Test Status](https://img.shields.io/github/actions/workflow/status/thomaslaurenson/prongs/tag.yml?style=flat&label=test&logo=github)

![Release Version](https://img.shields.io/github/v/release/thomaslaurenson/prongs?style=flat&logo=github) ![Release downloads](https://img.shields.io/github/downloads/thomaslaurenson/prongs/total?label=downloads&logo=github)

![Go Version](https://img.shields.io/github/go-mod/go-version/thomaslaurenson/prongs?logo=go) ![Code Coverage](https://img.shields.io/badge/Coverage-96.2%25-blue?logo=go)

Fast, custom security scanner.

## Installation

Download a pre-built binary from the [releases page](https://github.com/thomaslaurenson/prongs/releases). For easier install, use the bash installer script:

```sh
curl -fsSL https://github.com/thomaslaurenson/prongs/releases/latest/download/install.sh | bash
```

Or the PowerShell installer script if on Windows:

```ps
irm https://github.com/thomaslaurenson/prongs/releases/latest/download/install.ps1 | iex
```

Install from source:

```sh
go install github.com/thomaslaurenson/prongs@latest
```

## Usage

```
prongs scan --scanner <name>  --target <CIDR|file>
prongs scan --all             --target <CIDR|file>
prongs --version
```

Targets are CIDR ranges or single IPs, supplied via `--target` (repeatable and/or comma-separated) or `--target-file` (one entry per line). The two flags are mutually exclusive. If neither is provided, the `TARGET_CIDRS` environment variable (comma-separated) is used as a fallback.

### Flags

| Flag | Description | Default |
|---|---|---|
| `--target` | CIDR(s) or IP(s) to scan (repeatable and/or comma-separated) | |
| `--target-file` | Path to a file of CIDRs or IPs, one per line | |
| `--scanner` | Scanner to run (repeatable); see [Scanners](#scanners) | |
| `--all` | Run all default-enabled scanners | `false` |
| `--output` | Output format: `text` (TSV) or `pretty` (human-readable) | `text` |
| `--concurrency`, `-c` | Max concurrent probes | `200` |

### Scanners

| Name | Description | Default |
|---|---|---|
| `password-ssh` | Detects SSH servers accepting password authentication | yes |
| `accessible-rdp` | Detects RDP services accepting unauthenticated connections | no |
| `accessible-db` | Detects databases accepting unauthenticated connections | yes |
| `insecure-http` | Detects websites served over plaintext HTTP without an HTTPS redirect | yes |

### Examples

```bash
# Run one scanner against a single network
prongs scan --scanner password-ssh --target 192.168.0.0/24

# Detect websites served over plaintext HTTP
prongs scan --scanner insecure-http --target 192.168.0.0/24

# Run all default scanners against multiple networks
prongs scan --all --target 192.168.0.0/24 --target 10.0.0.0/24

# Multiple networks as a single comma-separated value
prongs scan --all --target 192.168.0.0/24,10.0.0.0/24

# Load targets from a file
prongs scan --all --target-file targets.txt

# Pretty-print output
prongs scan --all --target 192.168.0.0/24 --output pretty

# Limit the number of concurrent probes
prongs scan --all --target 192.168.0.0/24 --concurrency 50

# Use the TARGET_CIDRS environment variable for a single target
TARGET_CIDRS=192.168.0.0/24 prongs scan --all

# ...or multiple comma-separated targets
TARGET_CIDRS=192.168.0.0/24,10.0.0.0/24 prongs scan --all
```
