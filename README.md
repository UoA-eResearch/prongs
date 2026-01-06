# prongs

Fast, custom security scanner

## Requirements

- Python 3.9
- [Paramiko](https://www.paramiko.org/)

## Quickstart

### Linux

- Use `uv` to install:

```
uv install
```

- Use `uv` to install, and include development dependencies:

```
uv sync --all-extras
```

- Run using `python app/run.py`

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
