#!/usr/bin/env bash
set -euo pipefail

REPO="oopsunix/wii"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="wii"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    darwin)  OS="darwin" ;;
    linux)   OS="linux" ;;
    freebsd) OS="freebsd" ;;
    openbsd) OS="openbsd" ;;
    netbsd)  OS="netbsd" ;;
    *)
        echo -e "${RED}Error: Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    i686|i386)     ARCH="386" ;;
    armv7l)        ARCH="arm" ;;
    *)
        echo -e "${RED}Error: Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}Detected: ${OS}/${ARCH}${NC}"

# Get latest version from GitHub API
echo "Fetching latest version..."
VERSION=$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed 's/.*"v\(.*\)".*/\1/')

if [ -z "$VERSION" ]; then
    echo -e "${RED}Error: Could not fetch latest version${NC}"
    echo "Falling back to manual download..."
    echo "Please visit: https://github.com/${REPO}/releases"
    exit 1
fi

echo -e "${GREEN}Latest version: v${VERSION}${NC}"

# Construct download URL
FILENAME="${BINARY_NAME}_${VERSION}_${OS}_${ARCH}"
if [ "$OS" = "windows" ]; then
    URL="https://github.com/${REPO}/releases/download/v${VERSION}/${FILENAME}.zip"
    TMP_FILE="/tmp/${BINARY_NAME}.zip"
else
    URL="https://github.com/${REPO}/releases/download/v${VERSION}/${FILENAME}.tar.gz"
    TMP_FILE="/tmp/${BINARY_NAME}.tar.gz"
fi

# Download
echo "Downloading ${URL}..."
if command -v curl &> /dev/null; then
    curl -sL "$URL" -o "$TMP_FILE"
elif command -v wget &> /dev/null; then
    wget -q "$URL" -O "$TMP_FILE"
else
    echo -e "${RED}Error: Neither curl nor wget found${NC}"
    exit 1
fi

# Extract and install
echo "Installing to ${INSTALL_DIR}..."
mkdir -p "$INSTALL_DIR"

if [ "$OS" = "windows" ]; then
    unzip -o "$TMP_FILE" -d "$INSTALL_DIR"
    rm "$TMP_FILE"
else
    tar xzf "$TMP_FILE" -C /tmp
    mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    rm "$TMP_FILE"
fi

echo -e "${GREEN}Installed: ${INSTALL_DIR}/${BINARY_NAME}${NC}"

# Check if INSTALL_DIR is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo -e "${YELLOW}NOTE: ${INSTALL_DIR} is not in your PATH.${NC}"
    echo "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
    echo ""
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo ""
    echo "Then reload your shell or run:"
    echo "  source ~/.bashrc  # or ~/.zshrc"
fi

# Verify installation
if [ -x "${INSTALL_DIR}/${BINARY_NAME}" ]; then
    echo ""
    echo -e "${GREEN}Installation successful!${NC}"
    echo "Run '${BINARY_NAME}' to get started."
else
    echo ""
    echo -e "${YELLOW}Warning: Binary installed but may not be executable.${NC}"
    echo "Run: chmod +x ${INSTALL_DIR}/${BINARY_NAME}"
fi
