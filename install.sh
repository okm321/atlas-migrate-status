#!/bin/sh
set -e

# atlas-migrate-status install script
# Usage: curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/atlas-migrate-status/main/install.sh | sh

REPO="okm321/atlas-migrate-status"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

case $OS in
    darwin|linux) ;;
    mingw*|msys*|cygwin*) OS="windows" ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Get latest version
VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "Failed to get latest version"
    exit 1
fi

# Construct download URL
BINARY_NAME="atlas-migrate-status"
if [ "$OS" = "windows" ]; then
    ARCHIVE_EXT="zip"
    BINARY_NAME="${BINARY_NAME}.exe"
else
    ARCHIVE_EXT="tar.gz"
fi

ARCHIVE_NAME="atlas-migrate-status_${VERSION}_${OS}_${ARCH}.${ARCHIVE_EXT}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"

echo "Downloading atlas-migrate-status ${VERSION} for ${OS}/${ARCH}..."

# Create temporary directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

cd "$TMP_DIR"

# Download and extract
if command -v curl >/dev/null 2>&1; then
    curl -sL "$DOWNLOAD_URL" -o "$ARCHIVE_NAME"
elif command -v wget >/dev/null 2>&1; then
    wget -q "$DOWNLOAD_URL" -O "$ARCHIVE_NAME"
else
    echo "Error: curl or wget is required"
    exit 1
fi

if [ "$ARCHIVE_EXT" = "zip" ]; then
    unzip -q "$ARCHIVE_NAME"
else
    tar -xzf "$ARCHIVE_NAME"
fi

# Install binary
echo "Installing to ${INSTALL_DIR}..."

if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo "Requesting sudo permission to install to ${INSTALL_DIR}..."
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
    sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
fi

echo "✓ atlas-migrate-status ${VERSION} installed successfully!"
echo ""
echo "Usage:"
echo "  atlas-migrate-status --url \"postgres://user:pass@localhost:5432/dbname\""
echo ""
echo "For more information:"
echo "  atlas-migrate-status --help"
