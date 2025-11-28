# Inkwash

> A world-class CLI tool for managing FiveM servers with beautiful animations and real-time metrics.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-lightgrey)](https://github.com/VexoaXYZ/inkwash)

Inkwash brings the polish of modern web applications (Vercel, Linear) to terminal interfaces. Built in Go with a focus on performance, elegant design, and developer experience.

## ✨ Features

- 🎨 **Beautiful TUI** - Monochrome design with strategic purple accent
- 🎭 **Auto-Adaptive Animations** - Scales based on terminal capabilities (3 tiers)
- 📊 **Real-Time Metrics** - Live CPU/RAM/Network monitoring with sparklines
- 💾 **Smart Caching** - LRU binary caching (saves bandwidth & time)
- 🔐 **Encrypted Vault** - AES-256 license key storage (machine-bound)
- ⚡ **Blazing Fast** - <50ms startup, <50MB memory footprint
- 🌍 **Cross-Platform** - Windows & Linux support

## 🚀 Quick Start

### Installation

**From Source:**
```bash
git clone https://github.com/VexoaXYZ/inkwash.git
cd inkwash
make build
```

**Or download pre-built binaries from [Releases](https://github.com/VexoaXYZ/inkwash/releases)**

### Basic Usage

```bash
# Add a license key (encrypted storage)
inkwash key add

# Create a server
inkwash create my-server --build 17000 --key <key-id>

# Start the server
inkwash start my-server

# View logs
inkwash logs my-server -n 100

# List all servers
inkwash list

# Stop the server
inkwash stop my-server
```

## 📖 Commands

### Server Management

```bash
# Create a new server
inkwash create <name> [flags]
  --build      FXServer build number (default: 17000)
  --key        License key ID from vault
  --port       Server port (default: 30120)
  --path       Installation path

# Start server
inkwash start <name>

# Stop server
inkwash stop <name>

# List all servers with status
inkwash list

# View server logs
inkwash logs <name>
  -f, --follow    Follow log output (tail -f)
  -n, --lines     Number of lines to show (default: 50)
```

### License Key Management

```bash
# Add new license key
inkwash key add
  -l, --label    Label for the key
  -k, --key      License key (cfxk_...)

# List all keys (masked display)
inkwash key list

# Remove a key
inkwash key remove <key-id>
```

## 🏗️ Architecture

### Project Structure

```
inkwash/
├── cmd/                      # CLI commands
│   ├── root.go               # Root command & config
│   ├── create.go             # Server creation
│   ├── start.go / stop.go    # Process management
│   ├── list.go / logs.go     # Server info
│   └── key.go                # License key vault
│
├── internal/
│   ├── ui/                   # User interface
│   │   ├── animation/        # Easing & transitions
│   │   ├── components/       # Progress, spinner, sparkline
│   │   ├── styles.go         # Monochrome theme
│   │   └── detector.go       # Terminal capability detection
│   │
│   ├── server/               # Server management
│   │   ├── installer.go      # Installation orchestration
│   │   ├── process.go        # Start/stop/status
│   │   ├── config.go         # server.cfg generation
│   │   └── metrics.go        # Real-time metrics
│   │
│   ├── download/             # Download system
│   │   ├── artifacts.go      # FiveM API client
│   │   ├── downloader.go     # Parallel downloads
│   │   └── extractor.go      # Archive extraction
│   │
│   ├── cache/                # Caching & storage
│   │   ├── binary.go         # FXServer cache (LRU)
│   │   ├── metadata.go       # Cache metadata
│   │   └── keys.go           # License vault (AES-256)
│   │
│   └── registry/             # Server registry
│       ├── registry.go       # JSON storage
│       └── config.go         # Path helpers
│
└── pkg/types/                # Shared types
    ├── server.go / build.go / metrics.go
```

### Design Philosophy: "Monochrome Elegance"

**Color Palette:**
- Monochrome foundation (whites → grays → blacks)
- Single purple accent (#7C3AED) for emphasis
- Semantic colors only for status (success, error, warning)

**Animation Tiers:**
| Tier | Description | Features |
|------|-------------|----------|
| **Minimal** | Basic terminals | Simple spinners, no effects |
| **Balanced** | Default | Smooth animations, no shimmer |
| **Full** | Modern terminals | All effects + shimmer + particles |

Auto-detected based on:
- Terminal capabilities (ANSI 256-color support)
- System resources (CPU cores, RAM)
- Terminal emulator (Windows Terminal, iTerm2, etc.)

## 🎨 Screenshots

```
$ inkwash list

SERVERS

  ● Running  production-rp
      Port: 30120
      C:\FXServer\production-rp
      RAM: 2.14 GB

  ○ Stopped  dev-server
      Port: 30121
      C:\FXServer\dev-server

Total: 2 server(s)
```

```
$ inkwash key list

LICENSE KEYS

  Production
    ID:  a1b2c3d4-e5f6-7890-abcd-ef1234567890
    Key: cfxk_********************xj2k
    Created: Jan 15, 2025

Total: 1 key(s)
```

## ⚙️ Configuration

**Config Location:**
- Windows: `%APPDATA%\inkwash\config.yaml`
- Linux: `~/.config/inkwash/config.yaml`

**Example Configuration:**
```yaml
version: 1

defaults:
  install_path: "C:\\FXServer"  # Windows
  # install_path: "~/fxserver"  # Linux
  port: 30120

cache:
  enabled: true
  max_builds: 3                 # LRU eviction

ui:
  theme: "purple"               # Accent color
  animations: "auto"            # auto, full, balanced, minimal
  refresh_interval: 2           # Dashboard refresh (seconds)

telemetry:
  enabled: true                 # Opt-out analytics

advanced:
  parallel_downloads: true
  download_chunks: 3
  log_level: "info"
```

## 🔧 Development

### Prerequisites
- Go 1.21+
- Git
- Modern terminal (Windows Terminal, iTerm2, Alacritty recommended)

### Building

```bash
# Development build
make build

# All platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean
```

### Dependencies

| Package | Purpose |
|---------|---------|
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | TUI framework |
| [Lipgloss](https://github.com/charmbracelet/lipgloss) | Styling & layouts |
| [Cobra](https://github.com/spf13/cobra) | CLI framework |
| [Viper](https://github.com/spf13/viper) | Configuration |
| [gopsutil](https://github.com/shirou/gopsutil) | System metrics |

Full dependency list in [go.mod](go.mod)

## 🎯 Performance

**Benchmarks:**
- Cold start: <50ms (p95)
- Memory baseline: <50MB
- Download speed: Matches browser (3-chunk parallel)
- Animation: 60fps on modern terminals

## 🛣️ Roadmap

### v1.0 (Current) ✅
- ✅ Server installation & management
- ✅ Process lifecycle (start/stop)
- ✅ License key vault
- ✅ Binary caching
- ✅ Real-time metrics
- ✅ Cross-platform support

### v1.1 (Planned)
- [ ] Interactive TUI dashboard
- [ ] Interactive creation wizard
- [ ] Log streaming (tail -f)
- [ ] Resource management
- [ ] Server templates

### v2.0 (Future)
- [ ] Resource marketplace
- [ ] Automatic backups
- [ ] Discord integration
- [ ] Web dashboard
- [ ] Multi-server orchestration

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

**TL;DR:** You can freely use, modify, and distribute this software, even commercially. Just keep the copyright notice and license.

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) by Charm
- Inspired by modern CLI tools (Vercel CLI, Linear CLI)
- FiveM server artifacts by [Cfx.re](https://runtime.fivem.net/)

## 📞 Support

- 🐛 [Report a bug](https://github.com/VexoaXYZ/inkwash/issues)
- 💡 [Request a feature](https://github.com/VexoaXYZ/inkwash/issues)
- 💬 [Discussions](https://github.com/VexoaXYZ/inkwash/discussions)

---

**Made with ❤️ by [@VexoaXYZ](https://github.com/VexoaXYZ)**

*Inkwash - Because your FiveM servers deserve better management.*
