FROM python:3.9-alpine

WORKDIR /usr/src/app

COPY pyproject.toml README.md ./
COPY app/ ./app/

RUN pip install --no-cache-dir .

ENTRYPOINT ["prongs"]
CMD ["--enable-all", "--envvars"]
