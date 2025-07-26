# --- Build stage ---
FROM golang:1.24-alpine AS builder

ARG APP
ARG GIN_MODE=release
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GIN_MODE=$GIN_MODE

WORKDIR /src
COPY . .

# Install make
RUN apk add --no-cache make

# Build the Go application using the Makefile
RUN make ${APP}

# --- Minimal runtime stage ---
FROM alpine:3.20

ARG APP=server
ARG GIN_MODE=release
ENV GIN_MODE=$GIN_MODE

# If you need TLS, install ca-certificates
# RUN apk --no-cache add ca-certificates

COPY --from=builder /src/bin/${APP} /app

EXPOSE 8080

ENTRYPOINT ["/app"] 