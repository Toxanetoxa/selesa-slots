# syntax=docker/dockerfile:1.6

###########################
# 1) Build stage
###########################
FROM golang:1.23-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN buf generate || true

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/slots ./cmd/server

###########################
# 2) Runtime stage
###########################
FROM alpine:3.20
RUN adduser -S -u 10001 app
USER app
COPY --from=build /bin/slots /slots

EXPOSE 8080 9090
ENTRYPOINT ["/slots"]
