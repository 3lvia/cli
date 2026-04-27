#!/bin/bash

set -euo pipefail

LATEST_VERSION=$(curl -fsSl https://raw.githubusercontent.com/3lvia/cli/refs/heads/trunk/VERSION)

# Check if mac or linux
if [[ "$OSTYPE" == "darwin"* ]]; then
    OS="mac"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
else
    echo "Unsupported OS: $OSTYPE"
    exit 1
fi

# Check architecture
ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

cd "$(mktemp -d)"
TARBALL="3lv-$LATEST_VERSION-$OS-$ARCH.tar.gz"

curl -fsSL -o "$TARBALL" "https://github.com/3lvia/cli/releases/download/v$LATEST_VERSION/$TARBALL"
curl -fsSL -o "$TARBALL.md5" "https://github.com/3lvia/cli/releases/download/v$LATEST_VERSION/$TARBALL.md5"

# Check if md5sum is installed
if ! md5sum --version &> /dev/null; then
    echo "Command 'md5sum' could not be found, will skip verifying integrity of the download."
    echo "In the future, please ensure md5sum is installed."
else
    md5sum -c "$TARBALL.md5"
fi

tar -xzf "$TARBALL"

if [[ "$OS" == "mac" ]]; then
    install -Dm755 3lv "$HOME/.local/bin"
else
    install -Dm755 -t "$HOME/.local/bin" 3lv
fi

echo "3lv version $LATEST_VERSION installed successfully!"
