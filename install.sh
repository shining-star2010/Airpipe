#!/bin/sh
set -e

# Check if the user is root
if [ "$(id -u)" -eq 0 ]; then
    SUDO=""
else
    SUDO="sudo"
fi

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

URL="https://github.com/Sanyam-G/Airpipe/releases/latest/download/airpipe-${OS}-${ARCH}"

echo "Downloading airpipe for ${OS}-${ARCH}..."
curl -sL "$URL" -o /tmp/airpipe
chmod +x /tmp/airpipe

# Use the dynamic variable here
$SUDO mv /tmp/airpipe /usr/local/bin/airpipe

echo "Installed! Run: airpipe send <file>"
