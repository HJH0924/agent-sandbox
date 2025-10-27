FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY bin/api-server ./bin/api-server
COPY configs/ ./configs/

RUN chmod +x /app/bin/api-server && \
    mkdir -p /tmp/agent-sandbox
