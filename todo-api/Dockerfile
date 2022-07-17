# syntax=docker/dockerfile:1

# Stage 1: Build
FROM golang:1.18-bullseye AS build-env

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /todo-api

# Stage 2: Runtime
FROM gcr.io/distroless/base-debian11
ENV GIN_MODE=release

COPY --from=build-env /todo-api /todo-api

USER nonroot:nonroot

EXPOSE 8000

ENTRYPOINT ["/todo-api"]
