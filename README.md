# prongs

Fast, custom security scanner

## Requirements

- Python 3.9

## Quickstart

### Run from source using `uv`

- Requirements:
  - git
  - [uv](https://github.com/astral-sh/uv)

```
# Clone repo
git https://github.com/UoA-eResearch/prongs.git
# Install using uv
uv pip install -e .
# Run
uv run prongs --help
```

### Install from repo using `uv`

- Requirements:
  - [uv](https://github.com/astral-sh/uv)

```
# Install using uv
uv tool install git+https://github.com/UoA-eResearch/prongs.git@v0.2.0
```

### Examples

- Execute password SSH check against two target networks:

```
python app/run.py -s password-ssh -t 192.168.0.0/32,192.168.88.0/32
```

- Execute all scanners against target networks specified in a file:

```
echo -e "192.168.0.0/32\n192.168.88.0/32" > targets.txt
python app/run.py -s password-ssh -f targets.txt
```

- Execute all on password SSH scanner using environment variables:

```
TARGET_CIDRS=192.168.0.0/32,192.168.88.0/24 python app/run.py -s password-ssh -e
```

### Docker

Build the image:

```
docker build -f app/Dockerfile -t prongs .
```

Run the container:

```
docker run -e TARGET_CIDRS=192.168.0.0/32,192.168.88.0/24 -it prongs
```
