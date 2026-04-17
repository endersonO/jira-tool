#!/bin/sh
# install.sh — Cross-platform installer for jt (Jira Tool)
# Usage: curl -fsSL https://raw.githubusercontent.com/endersonO/jira-tool/main/install.sh | sh
set -e

REPO="endersonO/jira-tool"
BINARY="jt"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS="$(uname -s)"
case "$OS" in
  Linux*)  OS="linux" ;;
  Darwin*) OS="darwin" ;;
  MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
  *)
    echo "Error: Unsupported operating system: $OS"
    exit 1
    ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Error: Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Windows + arm64 is not supported
if [ "$OS" = "windows" ] && [ "$ARCH" = "arm64" ]; then
  echo "Error: Windows ARM64 is not supported. Use Windows AMD64 instead."
  exit 1
fi

# Determine file extension
if [ "$OS" = "windows" ]; then
  EXT="zip"
else
  EXT="tar.gz"
fi

# Get latest release tag
echo "Detecting latest version..."
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "Error: Could not determine latest version. Check https://github.com/${REPO}/releases"
  exit 1
fi

echo "Latest version: v${LATEST}"
echo "Platform:       ${OS}/${ARCH}"

# Build download URL
FILENAME="${BINARY}_${LATEST}_${OS}_${ARCH}.${EXT}"
URL="https://github.com/${REPO}/releases/download/v${LATEST}/${FILENAME}"

# Download
TMPDIR=$(mktemp -d)
echo "Downloading ${URL}..."
curl -fsSL "$URL" -o "${TMPDIR}/${FILENAME}"

# Extract
cd "$TMPDIR"
if [ "$EXT" = "zip" ]; then
  unzip -q "$FILENAME"
else
  tar -xzf "$FILENAME"
fi

# Install
if [ "$OS" = "windows" ]; then
  INSTALL_DIR="$HOME/bin"
  mkdir -p "$INSTALL_DIR"
  mv "${BINARY}.exe" "$INSTALL_DIR/"
  echo ""
  echo "Installed to ${INSTALL_DIR}/${BINARY}.exe"
  echo "Make sure ${INSTALL_DIR} is in your PATH."
else
  if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY" "$INSTALL_DIR/"
  else
    echo "Installing to ${INSTALL_DIR} (requires sudo)..."
    sudo mv "$BINARY" "$INSTALL_DIR/"
  fi
  echo ""
  echo "Installed to ${INSTALL_DIR}/${BINARY}"
fi

# Cleanup
rm -rf "$TMPDIR"

# Verify
echo ""
if command -v "$BINARY" >/dev/null 2>&1; then
  echo "Success! Installed version:"
  "$BINARY" --version
  echo ""
  echo "Get started with: jt init"
else
  echo "Installed, but '${BINARY}' is not in your PATH."
  echo "Add ${INSTALL_DIR} to your PATH, then run: jt init"
fi
