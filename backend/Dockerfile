FROM rust:1.88.0-alpine3.21

# Accept build arguments (default to PROD if not specified)
ARG TARGET="NOT_SET"
ARG DEV=false

# Exit if TARGET is not set
RUN if [ "$TARGET" = "NOT_SET" ]; then \
      echo "Build argument TARGET is not set"; \
      exit 1; \
    fi

# Set and copy working directory
WORKDIR /app
COPY . .

# Install build dependencies for musl and OpenSSL
RUN apk add --no-cache musl-dev build-base openssl-dev openssl-libs-static pkgconf

EXPOSE 8080

# Build the target
# NOTE: Even in DEV we want to be able to execute immediately (especially if we're not changing this service).
RUN if [ "$DEV" = "true" ]; then \
        cargo build --bin "$TARGET"; \
    else \
        cargo build --bin "$TARGET" --release; \
    fi

# Create a start script
COPY <<EOF /app/start.sh
#!/bin/sh
# Run based on DEV
if [ "$DEV" = "true" ]; then
    # Will build on startup
    cargo run --bin $TARGET
else
    # Will just run the binary
    cargo run --bin $TARGET --release
fi
EOF
RUN chmod +x /app/start.sh

# Use the start script as the entrypoint
ENTRYPOINT ["/app/start.sh"]