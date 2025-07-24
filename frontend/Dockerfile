FROM oven/bun:1.2.11-alpine


# Accept build arguments (default to PROD if not specified)
ARG DEV=false
ARG BASE_URL
ENV PUBLIC_BASE_URL=$BASE_URL

# Set and copy working directory
WORKDIR /app
COPY . .

# Install dependencies and build it
# NOTE: Even in DEV we want to be able to execute immediately (especially if we're not changing this service).
RUN bun install && \
    if [ "$DEV" != "true" ]; then \
        bun run build; \
    fi

# Expose the port for the frontend
EXPOSE 80

# Create a start script
COPY <<EOF /app/start.sh
#!/bin/sh
# Run based on DEV, but always install before
if [ "$DEV" = "true" ]; then
    bun install
    bun run dev --host --port 8080
else
    bun build/index.js --port 8080
fi
EOF
RUN chmod +x /app/start.sh

# Use the start script as the entrypoint
ENTRYPOINT ["/app/start.sh"]