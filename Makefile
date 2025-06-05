# DOCKER
docker_build:
	docker build -t prongs ./app

docker_run:
	docker run -it prongs

# PYTHON
venv_create:
	cd app; \
	python3 -m venv ./venv; \
	. ./venv/bin/activate && \
	pip3 install -r requirements.txt \

venv_create_dev:
	python3 -m venv ./venv; \
	. ./venv/bin/activate && \
	pip3 install -r requirements-dev.txt

update_packages:
	python3 -m venv ./venv; \
	. ./venv/bin/activate && \
	pip3 install -r requirements-dev.txt && \
	echo "[*] Checking: app/requirements.txt" && \
	pur -r app/requirements.txt && \
	echo "[*] Checking: requirements-dev.txt" && \
	pur -r requirements-dev.txt && \
	deactivate;

# LINTING
lint: \
	venv_create_dev \
	lint_python

lint_python:
	. ./venv/bin/activate && \
	echo "[*] Linting Python..." && \
	ruff check
