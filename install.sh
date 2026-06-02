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

# Construct download URL (goreleaser produces raw binaries without version in filename)
FILENAME="${BINARY_NAME}_${OS}_${ARCH}"
URL="https://github.com/${REPO}/releases/latest/download/${FILENAME}"
MIRROR_PREFIX="https://hubp.llkk.cc/"
TMP_FILE="/tmp/${BINARY_NAME}"

# Check GitHub accessibility, fall back to mirror if unreachable
echo "Checking GitHub connectivity..."
if command -v curl &> /dev/null; then
    if ! curl -sL --connect-timeout 5 --max-time 10 -o /dev/null "https://github.com" 2>/dev/null; then
        echo -e "${YELLOW}GitHub unreachable, using mirror...${NC}"
        URL="${MIRROR_PREFIX}${URL}"
    fi
elif command -v wget &> /dev/null; then
    if ! wget -q --timeout=10 --spider "https://github.com" 2>/dev/null; then
        echo -e "${YELLOW}GitHub unreachable, using mirror...${NC}"
        URL="${MIRROR_PREFIX}${URL}"
    fi
else
    echo -e "${RED}Error: Neither curl nor wget found${NC}"
    exit 1
fi

# Download
echo "Downloading ${URL}..."
if command -v curl &> /dev/null; then
    curl -sL "$URL" -o "$TMP_FILE"
elif command -v wget &> /dev/null; then
    wget -q "$URL" -O "$TMP_FILE"
fi

# Install
echo "Installing to ${INSTALL_DIR}..."
mkdir -p "$INSTALL_DIR"
mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

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
