#!/bin/bash
# InkWash Installer for Linux/macOS
# Downloads and installs the latest InkWash CLI

set -e

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
REPO="VexoaXYZ/InkWash"

echo "🚀 InkWash Installer for Linux/macOS"
echo "Installing to: $INSTALL_DIR"

# Create installation directory
mkdir -p "$INSTALL_DIR"
echo "✅ Created installation directory"

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "❌ Unsupported architecture: $ARCH" && exit 1 ;;
esac

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case $OS in
    linux) OS="linux" ;;
    darwin) OS="darwin" ;;
    *) echo "❌ Unsupported OS: $OS" && exit 1 ;;
esac

BINARY_NAME="inkwash-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
fi

echo "📥 Downloading InkWash for ${OS}-${ARCH}..."

# Get latest release URL
DOWNLOAD_URL=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep "browser_download_url.*${BINARY_NAME}" \
    | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "❌ Binary not found for ${OS}-${ARCH}"
    exit 1
fi

# Download binary
curl -sL "$DOWNLOAD_URL" -o "$INSTALL_DIR/inkwash"
chmod +x "$INSTALL_DIR/inkwash"

echo "✅ Downloaded successfully!"

# Add to PATH if not already present
case ":$PATH:" in
    *":$INSTALL_DIR:"*) 
        echo "✅ Already in PATH" 
        ;;
    *) 
        echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.bashrc"
        echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.zshrc" 2>/dev/null || true
        echo "✅ Added to PATH (restart terminal or run: export PATH=\"\$PATH:$INSTALL_DIR\")"
        ;;
esac

echo ""
echo "🎉 InkWash installed successfully!"
echo "Run 'inkwash --help' to get started"
echo "Note: You may need to restart your terminal for PATH changes to take effect"