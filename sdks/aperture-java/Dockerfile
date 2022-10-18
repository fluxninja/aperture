# syntax=docker/dockerfile:1

FROM --platform=linux/amd64 gradle:7.5.1-jdk18-alpine AS builder

WORKDIR /src

COPY --link . .

HEALTHCHECK --interval=5s --timeout=60s --retries=3 --start-period=20s \
    CMD wget --no-verbose --tries=1 --spider 127.0.0.1:8080/health || exit 1

CMD ["gradle", "run", "--no-daemon"]
