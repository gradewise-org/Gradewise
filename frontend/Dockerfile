FROM oven/bun:1.2.11-alpine

WORKDIR /app

# Accept build argument
ARG DEV=false
ARG BASE_URL
ENV PUBLIC_BASE_URL=$BASE_URL

COPY . .

RUN bun install

# Builds to /app/build
RUN bun run build

EXPOSE 80

COPY <<EOF /app/start.sh
#!/bin/sh
if [ "$DEV" = "true" ]; then
    bun run dev --host --port 80
else
    bun build/index.js --host 0.0.0.0 --port 80
fi
EOF

RUN chmod +x /app/start.sh
ENTRYPOINT ["/app/start.sh"]