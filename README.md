# InkWash

<div align="center">

![InkWash Banner](https://i.vgy.me/KcRrgP.png)

**The modern FiveM server management CLI**

*Effortlessly create, manage, and optimize your FiveM servers with enterprise-grade tooling.*

[![Release](https://img.shields.io/github/v/release/VexoaXYZ/InkWash?style=flat-square&color=007aff)](https://github.com/VexoaXYZ/InkWash/releases)
[![Downloads](https://img.shields.io/github/downloads/VexoaXYZ/InkWash/total?style=flat-square&color=34c759)](https://github.com/VexoaXYZ/InkWash/releases)
[![License](https://img.shields.io/badge/License-Custom-ff9500?style=flat-square)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey?style=flat-square)](https://github.com/VexoaXYZ/InkWash/releases)

[Get Started](#-installation) • [Documentation](#-features) • [Examples](#-usage-examples) • [Support](#-support)

</div>

---

## 🎯 **What is InkWash?**

InkWash transforms FiveM server management from a complex, time-consuming process into a simple, elegant experience. Built with modern development practices and a focus on developer productivity, InkWash provides everything you need to create, manage, and scale FiveM servers across any platform.

### **Why InkWash?**

- **🚀 Lightning Fast** — Intelligent caching and optimized operations
- **🎨 Intuitive Design** — Beautiful CLI with clear, actionable feedback  
- **🔒 Enterprise Ready** — Robust error handling and production-grade reliability
- **🌍 Universal** — Native support for Windows, macOS, and Linux
- **📦 Zero Dependencies** — Single binary with everything included

---

## ✨ **Features**

<table>
<tr>
<td width="50%" valign="top">

### **🖥️ Server Management**
- **Interactive Creation Wizard** — Step-by-step server setup
- **Automated Artifact Management** — Latest FiveM builds, automatically
- **txAdmin Integration** — Web-based administration ready
- **Multi-Server Support** — Manage unlimited server instances
- **Cross-Platform Deployment** — Windows, Linux, macOS native

</td>
<td width="50%" valign="top">

### **📦 Resource Management**
- **Smart Auto-Installation** — Popular resources with one command
- **Dynamic GitHub Integration** — Always latest releases
- **Dependency Resolution** — Automatic prerequisite handling
- **Configuration Management** — Seamless server.cfg updates
- **Resource Scanner** — Lightning-fast discovery with caching

</td>
</tr>
<tr>
<td width="50%" valign="top">

### **⚡ Performance & Optimization**
- **Intelligent Caching** — Sub-second resource scanning
- **Alpine Linux Optimization** — Minimal footprint containers
- **Progress Tracking** — Real-time operation feedback
- **Memory Efficient** — Optimized for high-performance servers
- **Concurrent Operations** — Parallel processing for speed

</td>
<td width="50%" valign="top">

### **🔧 Developer Experience**
- **Beautiful CLI** — Clean, consistent interface design
- **Comprehensive Logging** — Detailed debugging information
- **Modular Architecture** — Extensible and maintainable
- **API Integration** — GitHub releases, FiveM updates
- **Error Recovery** — Graceful failure handling

</td>
</tr>
</table>

---

## 🚀 **Installation**

### **Quick Install**

<details>
<summary><strong>Windows</strong> (PowerShell)</summary>

```powershell
irm https://raw.githubusercontent.com/VexoaXYZ/InkWash/master/install.ps1 | iex
```

*Requires PowerShell 5.1+ and internet connection*

</details>

<details>
<summary><strong>macOS / Linux</strong> (Bash)</summary>

```bash
curl -fsSL https://raw.githubusercontent.com/VexoaXYZ/InkWash/master/install.sh | bash
```

*Requires curl and bash*

</details>

### **Manual Installation**

1. **Download** the appropriate binary from [Releases](https://github.com/VexoaXYZ/InkWash/releases)
2. **Extract** to a directory in your `$PATH`
3. **Verify** installation: `inkwash version`

<details>
<summary><strong>Platform-Specific Binaries</strong></summary>

| Platform | Architecture | Binary |
|----------|--------------|--------|
| Windows | x64 | `inkwash-windows-amd64.exe` |
| macOS | Intel | `inkwash-darwin-amd64` |
| macOS | Apple Silicon | `inkwash-darwin-arm64` |
| Linux | x64 | `inkwash-linux-amd64` |
| Linux | ARM64 | `inkwash-linux-arm64` |

</details>

---

## 📖 **Usage Examples**

### **Create Your First Server**

```bash
# Interactive server creation
inkwash create

# Advanced options
inkwash create --name "MyServer" --template "basic"
```

### **Resource Management**

```bash
# Scan for existing resources
inkwash resources scan

# Install popular resources
inkwash resources install vmenu
inkwash resources install esx

# Configure installed resources
inkwash resources configure vmenu
```

### **Server Operations**

```bash
# Check InkWash version
inkwash version

# Get help for any command
inkwash --help
inkwash create --help
```

---

## 🏗️ **Architecture**

InkWash is built with a modular, extensible architecture designed for reliability and performance:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Layer     │    │  Core Services  │    │   Integrations  │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ • Commands      │───▶│ • Server Mgmt   │───▶│ • GitHub API    │
│ • Flags         │    │ • Resource Mgmt │    │ • FiveM API     │
│ • Validation    │    │ • Config Mgmt   │    │ • txAdmin       │
│ • Output        │    │ • Cache System  │    │ • File System   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### **Key Components**

- **Command Layer** — Cobra-based CLI with rich flag support
- **Service Layer** — Business logic and operations
- **Integration Layer** — External API and system interactions
- **Cache System** — High-performance resource discovery
- **Error Handling** — Comprehensive error recovery and reporting

---

## 📋 **Requirements**

### **System Requirements**

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| **OS** | Windows 10, macOS 10.15, Linux (any modern distro) | Latest versions |
| **RAM** | 512 MB available | 1 GB+ available |
| **Storage** | 100 MB free space | 1 GB+ free space |
| **Network** | Internet connection for downloads | Broadband connection |

### **Supported Platforms**

- ✅ **Windows** 10/11 (x64)
- ✅ **macOS** 10.15+ (Intel & Apple Silicon)
- ✅ **Linux** (x64, ARM64) — Ubuntu, Debian, CentOS, Arch, Alpine

---

## 🔧 **Advanced Configuration**

### **Environment Variables**

```bash
# Custom installation directory
export INKWASH_HOME="/opt/inkwash"

# Enable debug logging
export INKWASH_DEBUG=true

# Custom cache directory
export INKWASH_CACHE_DIR="/tmp/inkwash-cache"
```

### **Configuration Files**

InkWash stores configuration in platform-appropriate locations:

- **Windows**: `%APPDATA%\inkwash\config.yaml`
- **macOS**: `~/Library/Application Support/inkwash/config.yaml`
- **Linux**: `~/.config/inkwash/config.yaml`

---

## 🛠️ **Development**

### **Building from Source**

```bash
# Clone repository
git clone https://github.com/VexoaXYZ/InkWash.git
cd InkWash

# Install dependencies
go mod download

# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test
```

### **Project Structure**

```
InkWash/
├── cmd/inkwash/          # CLI entry point
├── internal/
│   ├── commands/         # Command implementations
│   ├── config/           # Configuration management
│   ├── fivem/           # FiveM-specific logic
│   └── utils/           # Utility functions
├── pkg/                 # Public API packages
├── test-servers/        # Testing infrastructure
└── Makefile            # Build automation
```

---

## 📚 **Documentation**

### **Command Reference**

<details>
<summary><strong>inkwash create</strong></summary>

Create new FiveM servers with interactive setup.

```bash
inkwash create [flags]

Flags:
  --name string      Server name
  --path string      Installation path
  --template string  Server template (basic, advanced)
  --interactive      Interactive mode (default true)
```

</details>

<details>
<summary><strong>inkwash resources</strong></summary>

Manage server resources with automatic installation and configuration.

```bash
inkwash resources [command]

Available Commands:
  scan          Scan for installed resources
  install       Install a resource
  configure     Configure an installed resource
  clear-cache   Clear resource cache

Flags:
  --path string   Server path to operate on
```

</details>

<details>
<summary><strong>inkwash version</strong></summary>

Display version and build information.

```bash
inkwash version

Output includes:
- Version number
- Build date
- Git commit hash
- Go version
- Platform/architecture
```

</details>

---

## 🤝 **Contributing**

We welcome contributions! Please read our contributing guidelines:

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Commit** your changes: `git commit -m 'Add amazing feature'`
4. **Push** to the branch: `git push origin feature/amazing-feature`
5. **Open** a Pull Request

### **Development Guidelines**

- Follow Go best practices and conventions
- Write comprehensive tests for new features
- Update documentation for user-facing changes
- Ensure cross-platform compatibility

---

## 🆘 **Support**

### **Getting Help**

- 📖 **Documentation**: [GitHub Wiki](https://github.com/VexoaXYZ/InkWash/wiki)
- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/VexoaXYZ/InkWash/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/VexoaXYZ/InkWash/discussions)

### **Common Issues**

<details>
<summary><strong>Installation Issues</strong></summary>

**Problem**: `command not found: inkwash`

**Solution**: Ensure the binary is in your `$PATH` and restart your terminal.

**Problem**: Permission denied errors

**Solution**: Run `chmod +x inkwash` on Unix systems or run as administrator on Windows.

</details>

<details>
<summary><strong>Resource Installation Issues</strong></summary>

**Problem**: Failed to download resources

**Solution**: Check internet connection and GitHub API rate limits.

**Problem**: Resources not found after installation

**Solution**: Verify server.cfg configuration and resource paths.

</details>

---

## 📄 **License**

This project is licensed under a Custom Attribution License. See the [LICENSE](LICENSE) file for details.

**Summary**: You may use, modify, and distribute this software freely, but must provide appropriate credit to the original authors. Failure to provide attribution voids the license.

---

## 🙏 **Acknowledgments**

Built with love for the FiveM community. Special thanks to:

- **FiveM Team** — For creating an amazing platform
- **Go Community** — For excellent tooling and libraries
- **Open Source Contributors** — For inspiration and best practices

---

<div align="center">

**Made with ❤️ by the Inkwell Team**

[⬆ Back to Top](#inkwash)

</div>
