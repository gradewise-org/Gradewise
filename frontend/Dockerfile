FROM oven/bun:1.2.11-alpine

WORKDIR /app

# Accept build argument
ARG DEV=false
ARG BASE_URL
ENV PUBLIC_BASE_URL=$BASE_URL

COPY . .

EXPOSE 80

# Run based on DEV, but always install before
COPY <<EOF /app/start.sh
#!/bin/sh
bun install
if [ "$DEV" = "true" ]; then
    bun run dev --host --port 80
else
    bun run build && bun build/index.js --host 0.0.0.0 --port 80
fi
EOF

RUN chmod +x /app/start.sh
ENTRYPOINT ["/app/start.sh"]