# prongs

Fast, custom security scanner

## Requirements

- Python 3.9
- [Paramiko](https://www.paramiko.org/)

## Quickstart

### Linux

- Create virtual environment, activate and install packages:

```
cd app
python3 -m venv venv
source venv/bin/activate
pip3 install -r requirements.txt
```

- Test using `python3 test.py`
- Run using `python3 run.py`

### Examples

- Execute password SSH check against two target networks:

```
python3 run.py -s password-ssh -t 192.168.0.0/32,192.168.88.0/32
```

- Execute all scanners against target networks specified in a file:

```
echo -e "192.168.0.0/32\n192.168.88.0/32" > targets.txt
python3 run.py -s password-ssh -f targets.txt
```

- Execute all on password SSH scanner using environment variables:

```
TARGET_CIDRS=192.168.0.0/32,192.168.88.0/24 python3 run.py -s password-ssh -e
```

### Docker

Build the image:

```
docker build -t prongs ./app
```

Run the container:

```
docker run -e TARGET_CIDRS=192.168.0.0/32,192.168.88.0/24 -it prongs
```

## Development

- Create virtual environment, activate and install packages:

```
python3 -m venv venv
source venv/bin/activate
pip3 install -r requirements-dev.txt
```

- Update packages

```
make update_packages
```

- Python linting:

```
make lint
```
