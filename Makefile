SHELL := /bin/bash
BINARY  := prongs
MODULE  := github.com/thomaslaurenson/prongs
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X $(MODULE)/internal/config.Version=$(VERSION)

TAG     ?= $(shell git describe --tags --abbrev=0 2>/dev/null)


.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  %-18s %s\n", $$1, $$2}'

# BUILD
.PHONY: build
build: ## Build binary for current platform
	go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY) .

.PHONY: build_snapshot
build_snapshot: ## Build binaries for all platforms using goreleaser (snapshot)
	goreleaser build --snapshot --clean

# LINT
.PHONY: fmt
fmt: ## Format Go source files
	gofmt -w .

.PHONY: fmt_check
fmt_check: ## Check Go formatting
	@unformatted=$$(gofmt -l .); \
	if [[ -n "$$unformatted" ]]; then \
		echo "Unformatted Go files:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

.PHONY: mod_check
mod_check: ## Check go.mod/go.sum tidiness
	go mod tidy
	git diff --exit-code go.mod go.sum

.PHONY: lint
lint: fmt_check mod_check vet ## Run lint checks

.PHONY: vet
vet: ## Run go vet
	go vet ./...

# TEST
.PHONY: test
test: ## Run all tests with race detector (requires network for integration tests)
	go test -race -count=1 ./...

.PHONY: test_coverage
test_coverage: ## Run tests with coverage report over internal packages
	go test -race -count=1 -coverpkg=./internal/... -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

.PHONY: test_unit
test_unit: ## Run unit tests only (no network)
	go test -v -run 'TestExpand' ./internal/target/...

# GET
.PHONY: get_changelog
get_changelog: ## Print release notes for TAG to stdout (default: latest tag; override with TAG=v1.0.0)
	@tag="$(TAG)"; tag="$${tag#v}"; \
	if [[ -z "$$tag" ]]; then \
	  printf 'get_changelog: TAG is empty; pass TAG=v1.0.0 or create a git tag\n' >&2; \
	  exit 1; \
	fi; \
	notes="$$(awk -v tag="$$tag" ' \
	  /^## / { if (found) exit; if (index($$0,"## "tag" ")==1 || $$0=="## "tag) found=1; next } \
	  found { lines[n++]=$$0 } \
	  END { \
	    s=0; while (s<n && lines[s]~/^[[:space:]]*$$/) s++; \
	    e=n-1; while (e>=s && lines[e]~/^[[:space:]]*$$/) e--; \
	    for (i=s;i<=e;i++) print lines[i] \
	  }' CHANGELOG.md)"; \
	if [[ -z "$$notes" ]]; then \
	  printf 'get_changelog: no CHANGELOG entry for %s\n' "$$tag" >&2; \
	  exit 1; \
	fi; \
	printf '%s\n' "$$notes"

.PHONY: ci
ci: lint test ## Run all CI checks locally

# TASKS
.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin/ dist/ install.sh install.ps1 checksums.txt

.PHONY: docker_build
docker_build: ## Build Docker image
	docker build -t prongs .

.PHONY: docker_run
docker_run: ## Run Docker image against scanme.nmap.org
	docker run --rm -e TARGET_CIDRS=45.33.32.156/32 prongs scan --all
