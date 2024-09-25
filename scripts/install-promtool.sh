#!/bin/bash

OS=$(uname -s)
ARCH=$(uname -m)

PROM_VERSION="2.54.1"

# Set the base URL for Prometheus
BASE_URL="https://github.com/prometheus/prometheus/releases/download/v${PROM_VERSION}"

# Map uname outputs to Prometheus compatible versions
case "$OS" in
    Linux)
        OS="linux"
        ;;
    Darwin)
        OS="darwin"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

FILENAME="prometheus-${PROM_VERSION}.${OS}-${ARCH}.tar.gz"
URL="${BASE_URL}/${FILENAME}"
curl -LO "$URL"

tar -xzf prometheus-*.tar.gz --strip-components=1 '*/promtool'
