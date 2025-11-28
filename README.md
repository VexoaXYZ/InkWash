# InkWash - FiveM Server Manager

> A professional CLI tool for creating and managing FiveM servers with optimized configuration and automated mod conversion.

[![Download](https://img.shields.io/github/v/release/VexoaXYZ/InkWash?label=Download&style=for-the-badge&logo=github)](https://github.com/VexoaXYZ/InkWash/releases/latest)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-blue?style=for-the-badge)](https://github.com/VexoaXYZ/InkWash)
[![Go Version](https://img.shields.io/github/go-mod/go-version/VexoaXYZ/InkWash?style=for-the-badge)](go.mod)

---

## Quick Install

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/VexoaXYZ/InkWash/master/install.ps1 | iex
```

### Linux & macOS
```bash
curl -fsSL https://raw.githubusercontent.com/VexoaXYZ/InkWash/master/install.sh | bash
```

The installer automatically:
- Downloads the latest release
- Extracts and installs the binary
- Adds InkWash to your PATH
- Configures your environment

---

## Features

**Server Management**
- **Interactive Setup Wizard** - Step-by-step server creation with intelligent defaults
- **Multi-Server Support** - Create and manage unlimited FiveM servers
- **Automated Installation** - Downloads and configures FiveM binaries automatically
- **Process Management** - Start, stop, and monitor servers with built-in process control

**Mod Conversion**
- **GTA5 to FiveM Converter** - Automatically converts GTA5 mods to FiveM resources
- **Batch Processing** - Queue multiple mods for conversion
- **Rate Limiting** - Respects API limits with intelligent throttling (max 2 concurrent)
- **Parallel Downloads** - Optimized download performance

**Security & Configuration**
- **Encrypted Key Storage** - AES-256-GCM encryption for license keys
- **Secure Configuration** - Environment-based configuration management
- **Optimized Defaults** - Production-ready server configurations out of the box

---

## Getting Started

### Prerequisites

- **FiveM License Key**: Obtain from [Cfx.re Portal](https://portal.cfx.re/servers/registration-keys)
- **Operating System**: Windows 10+, Linux, or macOS
- **Disk Space**: ~500MB minimum per server

### Installation

#### Automated Install (Recommended)

**Windows:**
```powershell
irm https://raw.githubusercontent.com/VexoaXYZ/InkWash/master/install.ps1 | iex
```

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/VexoaXYZ/InkWash/master/install.sh | bash
```

#### Manual Install

1. Download the latest release for your platform from the [Releases page](https://github.com/VexoaXYZ/InkWash/releases/latest)
2. Extract the archive to your desired location
3. Add the binary location to your PATH (optional)
4. Run `inkwash` from your terminal

#### Build from Source

```bash
git clone https://github.com/VexoaXYZ/InkWash.git
cd InkWash
go build -o inkwash .
```

---

## Usage

### Creating a Server

```bash
inkwash create
```

The interactive wizard will guide you through:
1. Server name and directory
2. FiveM version selection
3. License key configuration
4. Port and network settings

### Managing Servers

```bash
# List all servers
inkwash list

# Start a server
inkwash start <server-name>

# Stop a server
inkwash stop <server-name>

# View server logs
inkwash logs <server-name>
```

### Converting GTA5 Mods

```bash
inkwash convert
```

Supports:
- Direct URLs from [gta5-mods.com](https://www.gta5-mods.com/)
- Batch conversion with URL lists
- Automatic resource installation

### License Key Management

```bash
# Add a license key
inkwash key add

# List configured keys (masked)
inkwash key list

# Remove a key
inkwash key remove <key-id>
```

---

## Command Reference

### Server Commands

| Command | Description |
|---------|-------------|
| `inkwash create` | Launch server creation wizard |
| `inkwash start <name>` | Start a FiveM server |
| `inkwash stop <name>` | Stop a running server |
| `inkwash list` | List all configured servers |
| `inkwash logs <name>` | Stream server logs in real-time |

### Mod Converter

| Command | Description |
|---------|-------------|
| `inkwash convert` | Launch GTA5 mod converter wizard |

### License Keys

| Command | Description |
|---------|-------------|
| `inkwash key add` | Add a new FiveM license key |
| `inkwash key list` | List all stored keys (masked) |
| `inkwash key remove <id>` | Remove a license key |

---

## Configuration

InkWash stores configuration in:
- **Windows**: `%APPDATA%\inkwash\`
- **Linux/macOS**: `~/.config/inkwash/`

### Configuration Files

- `config.json` - Global settings
- `keys.encrypted` - Encrypted license keys
- `servers/` - Per-server configurations

---

## FAQ

### Where do I get a FiveM license key?

Visit the [Cfx.re Portal](https://portal.cfx.re/servers/registration-keys) to:
1. Log in with your Cfx.re account
2. Register a new server
3. Copy your license key (starts with `cfxk_`)

### How do I update InkWash?

Run the install script again, or download the latest release manually. Your configuration and servers will be preserved.

### My server won't start

**Troubleshooting steps:**
1. Verify your license key: `inkwash key list`
2. Check server logs: `inkwash logs <server-name>`
3. Ensure port 30120 is available
4. Verify FiveM binary integrity

### Can I use custom server configurations?

Yes. InkWash creates standard FiveM server directories. You can manually edit `server.cfg` and other configuration files in your server directory.

### Does InkWash work on headless servers?

Yes. InkWash is fully compatible with headless Linux servers and can be used in automated deployment pipelines.

---

## Architecture

InkWash is built with:
- **Language**: Go 1.24+
- **TUI Framework**: Bubble Tea
- **Encryption**: AES-256-GCM
- **Build Tool**: GoReleaser
- **CI/CD**: GitHub Actions

### Project Structure

```
inkwash/
├── cmd/           # CLI commands
├── internal/      # Core business logic
│   ├── config/    # Configuration management
│   ├── server/    # Server management
│   ├── converter/ # Mod converter
│   └── crypto/    # Encryption utilities
├── pkg/types/     # Shared types
└── main.go        # Entry point
```

---

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Setup

```bash
# Clone repository
git clone https://github.com/VexoaXYZ/InkWash.git
cd InkWash

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o inkwash .
```

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history and release notes.

---

## License

InkWash is open source software licensed under the [MIT License](LICENSE).

---

## Support

- **Documentation**: [GitHub Wiki](https://github.com/VexoaXYZ/InkWash/wiki)
- **Issues**: [GitHub Issues](https://github.com/VexoaXYZ/InkWash/issues)
- **Discussions**: [GitHub Discussions](https://github.com/VexoaXYZ/InkWash/discussions)

---

## Acknowledgments

Created by [Vexoa](https://github.com/VexoaXYZ) for the FiveM community.

**Version 2.0** - Complete rewrite with enhanced performance, security, and usability.

---

<div align="center">

### Ready to get started?

**[Download InkWash](https://github.com/VexoaXYZ/InkWash/releases/latest)** | **[View Documentation](https://github.com/VexoaXYZ/InkWash/wiki)** | **[Report Issue](https://github.com/VexoaXYZ/InkWash/issues)**

</div>
