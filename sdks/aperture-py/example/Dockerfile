# syntax=docker/dockerfile:1

FROM python:3.8-slim-buster

WORKDIR /app

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  ca-certificates \
  wget \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

HEALTHCHECK --interval=5s --timeout=60s --retries=3 --start-period=5s \
  CMD wget --no-verbose --tries=1 --spider 127.0.0.1:8080/health || exit 1

RUN pip install poetry==1.4.0
COPY poetry.lock pyproject.toml ./
RUN poetry config virtualenvs.create false \
  && poetry config installer.modern-installation false \
  && poetry install --without=dev --with=example --no-root

COPY . .
RUN poetry install --only-root

ENV QUART_APP=example/main.py
EXPOSE 8080
CMD [ "python3", "-m" , "quart", "run", "--host=0.0.0.0", "--port=8080"]
