SHELL := /bin/bash

.PHONY: help install install_all update clean build lint format format_check test docker_build docker_run tag_release

help:
	@echo "Available targets:"
	@echo "  install         Install dependencies"
	@echo "  install_all     Install all optional dependencies"
	@echo "  update          Update all packages to latest versions"
	@echo "  clean           Remove .venv and dist directories"
	@echo "  build           Build the package"
	@echo "  lint            Check code with ruff"
	@echo "  format          Format code with ruff"
	@echo "  format_check    Check code formatting"
	@echo "  test            Run pytest"
	@echo "  tag_release     Tag git with version from pyproject and push"
	@echo "  docker_build    Build the docker image"
	@echo "  docker_run      Run the docker image"

install:
	uv sync

install_all:
	uv sync --all-extras

update:
	uv lock --upgrade
	uv sync --all-extras

clean:
	rm -rf .venv dist

build:
	uv build

lint:
	uv run ruff check .

format:
	uv run ruff format .

format_check:
	uv run ruff format --check .

test:
	uv run pytest -rP

# DOCKER
docker_build:
	docker build -f app/Dockerfile -t prongs .

docker_run:
	docker run -it prongs

# TAG
tag_release:
	VERSION=$$(grep -m1 'version = ' pyproject.toml | cut -d '"' -f 2); \
	TAG="v$$VERSION"; \
	echo "[*] Current version: $$TAG"; \
	read -p "[*] Tag and push? (y/N) " yn; \
	case $$yn in \
		[yY]*) \
			echo "Running: git tag $$TAG"; \
			echo "Running: git push origin $$TAG"; \
			;; \
		[nN]*) \
			echo "[*] Exiting..."; \
			;; \
		*) \
			echo "[*] Invalid response... Exiting"; \
			exit 1; \
			;; \
	esac
