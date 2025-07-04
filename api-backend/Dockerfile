FROM rust:alpine

ARG TARGET_NAME=gradewise-api-backend
ARG DEV=false

# Set and copy working directory
WORKDIR /app
COPY . .

# Install build dependencies for musl
RUN apk add --no-cache musl-dev build-base

EXPOSE 80

# Run based on DEV
COPY <<EOF /app/start.sh
#!/bin/sh
if [ "$DEV" = "true" ]; then
    cargo run --
else
    cargo run --release --
fi
EOF

RUN chmod +x /app/start.sh
ENTRYPOINT ["/app/start.sh"]