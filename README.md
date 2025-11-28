# Inkwash

> A world-class CLI tool for managing FiveM servers with beautiful animations and real-time metrics.

Inkwash brings the polish of modern web applications (Vercel, Linear) to terminal interfaces. Built in Go with a focus on performance, elegant design, and developer experience.

## Features

- **Beautiful TUI** - Monochrome design with strategic accent colors
- **Auto-Adaptive Animations** - Scales based on terminal capabilities and system performance
- **Real-Time Metrics** - Live CPU/RAM/Network monitoring with sparkline graphs
- **Smart Caching** - Binary caching with LRU eviction
- **Encrypted Key Vault** - AES-256 encrypted license key storage
- **Fast & Efficient** - <50ms startup, <50MB memory footprint

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/VexoaXYZ/inkwash.git
cd inkwash

# Build
make build

# Or install to $GOPATH/bin
make install
```

### Pre-built Binaries

Download the latest release from [GitHub Releases](https://github.com/VexoaXYZ/inkwash/releases).

## Usage

```bash
# Launch interactive dashboard
inkwash

# Create a new server
inkwash create

# Start a server
inkwash start <name>

# Stop a server
inkwash stop <name>

# List all servers
inkwash list

# View server logs
inkwash logs <name>

# Manage license keys
inkwash key add
inkwash key list
inkwash key remove <id>

# Cache management
inkwash cache list
inkwash cache clear
```

## Architecture

### Project Structure

```
inkwash/
├── cmd/                      # CLI commands
│   ├── root.go               # Root command & config
│   ├── create.go             # Server creation wizard
│   ├── start.go              # Start server
│   ├── stop.go               # Stop server
│   └── ...                   # Other commands
│
├── internal/
│   ├── ui/                   # User interface
│   │   ├── animation/        # Animation system
│   │   │   ├── easing.go     # Easing functions
│   │   │   └── transition.go # Transition system
│   │   ├── components/       # UI components
│   │   │   ├── progress.go   # Progress bars
│   │   │   ├── spinner.go    # Loading spinners
│   │   │   └── sparkline.go  # Inline graphs
│   │   ├── styles.go         # Lipgloss theme
│   │   ├── dashboard.go      # Main dashboard (TBD)
│   │   ├── wizard.go         # Creation wizard (TBD)
│   │   └── detector.go       # Terminal capability detection
│   │
│   ├── server/               # Server management (TBD)
│   ├── cache/                # Caching system (TBD)
│   ├── download/             # Download manager (TBD)
│   ├── registry/             # Server registry (TBD)
│   └── telemetry/            # Usage tracking (TBD)
│
├── pkg/
│   └── types/                # Shared types
│       ├── server.go         # Server struct
│       ├── build.go          # Build struct
│       └── metrics.go        # Metrics struct
│
└── main.go                   # Entry point
```

### Design Philosophy: "Monochrome Elegance"

**Color Palette:**
- Monochrome foundation (whites, grays, blacks)
- Single purple accent color (#7C3AED)
- Semantic colors used sparingly (success, error, warning)

**Animation Tiers:**
- **Minimal**: Basic terminals, slow systems (simple spinners, no shimmer)
- **Balanced**: Default for most users (smooth animations, no effects)
- **Full**: Modern terminals, fast systems (shimmer, particles, transitions)

**Performance Targets:**
- Cold start: <50ms
- Memory baseline: <50MB
- Animation: 60fps on modern terminals

## Development

### Prerequisites

- Go 1.21 or higher
- Modern terminal (Windows Terminal, iTerm2, Alacritty, etc.)

### Building

```bash
# Development build
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run locally
make run
```

### Dependencies

- **Bubble Tea** - TUI framework
- **Lipgloss** - Styling and layouts
- **Bubbles** - Pre-built components
- **Cobra** - CLI framework
- **Viper** - Configuration management
- **gopsutil** - System metrics
- **sevenzip** - 7z extraction (Windows)
- **xz** - tar.xz extraction (Linux)

## Roadmap

### Phase 1: Foundation ✅
- [x] Go project initialization
- [x] Cobra CLI structure
- [x] Viper config system
- [x] Terminal capability detection
- [x] Animation system (easing, transitions)
- [x] Lipgloss monochrome theme
- [x] Base UI components (progress, spinner, sparkline)

### Phase 2: Core Systems (In Progress)
- [ ] FiveM artifact API client
- [ ] Parallel download manager
- [ ] Archive extraction
- [ ] Binary cache system
- [ ] Server registry

### Phase 3: UI (Upcoming)
- [ ] Interactive dashboard
- [ ] Creation wizard
- [ ] License key manager
- [ ] Real-time metrics display

### Phase 4: Server Management (Upcoming)
- [ ] Process management
- [ ] Metrics collection
- [ ] server.cfg generation
- [ ] Launch scripts

### Phase 5: Polish (Upcoming)
- [ ] Cross-platform testing
- [ ] Performance optimization
- [ ] Documentation
- [ ] Release automation

## Configuration

Default config location: `~/.config/inkwash/config.yaml`

```yaml
version: 1

defaults:
  install_path: "C:\\FXServer"  # Windows
  # install_path: "~/fxserver"  # Linux
  port: 30120

cache:
  enabled: true
  max_builds: 3

ui:
  theme: "purple"               # Accent color
  animations: "auto"            # auto, full, balanced, minimal
  refresh_interval: 2           # Dashboard refresh (seconds)

telemetry:
  enabled: true                 # Opt-out

advanced:
  parallel_downloads: true
  download_chunks: 3
  log_level: "info"
```

## License

MIT License - see [LICENSE](LICENSE) for details

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) first.

## Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Inspired by modern CLI tools like Vercel CLI and Linear CLI
- FiveM server artifacts provided by [Cfx.re](https://runtime.fivem.net/)

---

**Built with ❤️ by [@VexoaXYZ](https://github.com/VexoaXYZ)**
