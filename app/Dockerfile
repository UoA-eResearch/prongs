FROM python:3.9-alpine

WORKDIR /usr/src/app

COPY pyproject.toml ./
RUN pip install --no-cache-dir .

COPY app/ ./

CMD [ "python", "./run.py", "--enable-all", "--envvars"]
