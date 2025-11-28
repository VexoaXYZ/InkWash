#!/usr/bin/env bash
# InkWash Installer for Linux and macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/VexoaXYZ/InkWash/master/install.sh | bash

set -e

# Configuration
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
REPO="VexoaXYZ/InkWash"
BINARY_NAME="inkwash"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Helper functions
print_header() {
    echo -e "${CYAN}InkWash Installer${NC}"
    echo -e "${CYAN}==================${NC}"
    echo ""
}

print_error() {
    echo -e "${RED}ERROR: $1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_info() {
    echo -e "${YELLOW}$1${NC}"
}

# Detect OS and architecture
detect_platform() {
    local os=""
    local arch=""

    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac

    case "$(uname -m)" in
        x86_64)     arch="amd64" ;;
        amd64)      arch="amd64" ;;
        arm64)      arch="arm64" ;;
        aarch64)    arch="arm64" ;;
        armv7l)     arch="arm" ;;
        i686)       arch="386" ;;
        i386)       arch="386" ;;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac

    echo "${os}_${arch}"
}

# Main installation
main() {
    print_header

    # Detect platform
    print_info "Detecting platform..."
    platform=$(detect_platform)
    echo "   Platform: $platform"

    # Get latest release
    print_info ""
    print_info "Fetching latest release..."

    release_json=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest")
    version=$(echo "$release_json" | grep '"tag_name"' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')

    if [ -z "$version" ]; then
        print_error "Failed to fetch latest release version"
        exit 1
    fi

    echo "   Latest version: $version"

    # Find download URL for platform (supports both naming conventions)
    download_url=$(echo "$release_json" | grep "browser_download_url" | grep -E "${platform}\.(tar\.gz|zip)" | sed -E 's/.*"browser_download_url": "([^"]+)".*/\1/' | head -n 1)

    if [ -z "$download_url" ]; then
        print_error "No release found for platform: $platform"
        print_info "Available releases:"
        echo "$release_json" | grep "browser_download_url" | sed -E 's/.*"browser_download_url": "([^"]+)".*/\1/'
        exit 1
    fi

    print_info "Found asset: $(basename "$download_url")"

    # Create install directory
    print_info ""
    print_info "Installing to: $INSTALL_DIR"
    mkdir -p "$INSTALL_DIR"

    # Download and extract
    print_info ""
    print_info "Downloading InkWash $version..."

    tmp_dir=$(mktemp -d)
    archive_name=$(basename "$download_url")
    archive_path="$tmp_dir/$archive_name"

    if ! curl -fsSL -o "$archive_path" "$download_url"; then
        print_error "Download failed"
        rm -rf "$tmp_dir"
        exit 1
    fi

    print_success "Downloaded successfully!"

    # Extract
    print_info ""
    print_info "Extracting files..."

    # Detect archive type and extract accordingly
    if [[ "$archive_name" == *.tar.gz ]]; then
        tar -xzf "$archive_path" -C "$tmp_dir"
    elif [[ "$archive_name" == *.zip ]]; then
        unzip -q "$archive_path" -d "$tmp_dir"
    else
        print_error "Unsupported archive format: $archive_name"
        rm -rf "$tmp_dir"
        exit 1
    fi

    # Find the binary (might be in a subdirectory)
    binary_path=$(find "$tmp_dir" -name "$BINARY_NAME" -type f | head -n 1)

    if [ -z "$binary_path" ]; then
        print_error "Binary not found in archive"
        print_info "Archive contents:"
        ls -R "$tmp_dir"
        rm -rf "$tmp_dir"
        exit 1
    fi

    # Install binary
    mv "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"

    # Cleanup
    rm -rf "$tmp_dir"

    print_success "Extracted successfully!"

    # Check if install dir is in PATH
    print_info ""
    print_info "Checking PATH..."

    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        print_info "   $INSTALL_DIR is not in your PATH"
        print_info ""
        print_info "   Add it to your shell profile:"

        # Detect shell
        if [ -n "$BASH_VERSION" ]; then
            profile="$HOME/.bashrc"
        elif [ -n "$ZSH_VERSION" ]; then
            profile="$HOME/.zshrc"
        else
            profile="$HOME/.profile"
        fi

        echo -e "   ${CYAN}echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> $profile${NC}"
        echo -e "   ${CYAN}source $profile${NC}"
        print_info ""

        # Ask if we should add it automatically
        read -p "   Add to PATH now? (y/N) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$profile"
            print_success "Added to $profile"
            print_info "   Run: source $profile"
        fi
    else
        print_success "Already in PATH!"
    fi

    # Success message
    echo ""
    print_success "InkWash installed successfully!"
    echo ""
    print_info "Quick Start:"
    echo "   1. Open a new terminal (or run: source $profile)"
    echo "   2. Run: inkwash create"
    echo ""
    print_info "Documentation: https://github.com/$REPO/wiki"
    print_info "Get License Key: https://portal.cfx.re/servers/registration-keys"
    print_info "Report Issues: https://github.com/$REPO/issues"
    echo ""
}

main "$@"
