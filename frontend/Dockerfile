FROM oven/bun:1.2.11-alpine

WORKDIR /app

# Accept build argument
ARG DEV=false

COPY . .

RUN bun install

# Builds to /app/build
RUN bun run build

EXPOSE 3000

COPY <<EOF /app/start.sh
#!/bin/sh
if [ "$DEV" = "true" ]; then
    bun run dev --port 3000
else
    bun build/index.js --host 0.0.0.0 --port 3000
fi
EOF

RUN chmod +x /app/start.sh
ENTRYPOINT ["/app/start.sh"]