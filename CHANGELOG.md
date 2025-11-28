# 📋 Changelog

All notable changes to InkWash will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [2.0.0] - 2025-01-XX (Coming Soon!)

### 🎉 Complete Rewrite

Version 2.0 is a complete rewrite of InkWash from the ground up! We've rebuilt everything to make FiveM server management easier, faster, and more beautiful than ever.

### ✨ New Features

#### Interactive Wizards
- **Server Creation Wizard** - Step-by-step guide for creating servers
  - Smart validation at every step
  - Helpful tips and hints
  - Pre-filled defaults for beginners
  - Professional server.cfg templates

- **GTA5 Mod Converter** - Turn GTA5 mods into FiveM resources
  - Support for multiple URLs at once
  - Queue system with rate limiting (respects API limits)
  - Parallel downloads for speed
  - Auto-extraction to proper folders ([vehicles], [weapons], etc.)
  - Clear progress tracking for each mod

#### Beautiful Modern UI
- **Bubble Tea TUI Framework** - Smooth, professional terminal interface
- **Adaptive Animations** - Automatically adjusts to your terminal's capabilities
- **Color-Coded Status** - Easy to see what's happening at a glance
- **Real-Time Progress** - Live progress bars and spinners
- **No More Scrolling** - UI updates in place (500ms throttling)

#### Server Management
- **Per-Server Binaries** - Each server gets its own FXServer installation
- **Metadata Tracking** - Build info, creation date, usage stats
- **Smart Registry** - Auto-removes servers that no longer exist
- **Lifecycle Tracking** - Know when servers were started, stopped, etc.

#### Security & Safety
- **Encrypted License Keys** - AES-256-GCM encryption
- **Machine-Bound Keys** - Keys are tied to your PC
- **Filesystem Safety** - Server names automatically converted to safe folder names
- **Path Validation** - Prevents accidental overwrites

### 🔄 Changed

- **Command Structure** - Simpler, more intuitive commands
- **Config Format** - Improved configuration files
- **Registry System** - Better server tracking and validation
- **Error Messages** - More helpful, actionable error messages
- **UI/UX** - Complete redesign with modern styling

### ⚠️ Breaking Changes

If you're upgrading from v1.x, please note:

- **Commands are different** - Check the new command list in README
- **Config files changed** - Old configs won't work (easy to recreate)
- **Registry format updated** - You'll need to recreate servers
- **License keys** - Re-add your keys with `inkwash key add`

**Migration Tip:** Just recreate your servers with the new wizard - it's super quick!

### 🐛 Fixes

- Server selector no longer shows deleted/moved servers
- Git clone output doesn't bleed through UI during setup
- Text inputs properly handle spaces and special characters
- Cursor position works correctly in all inputs
- Rate limiting prevents API errors when converting multiple mods
- Clear visual feedback for all operations

### 🎨 UI/UX Improvements

- Numbered lists for multi-step wizards
- Clear status indicators (✓ complete, spinner for active, ⏳ queued, ✗ error)
- Consistent mod ordering throughout conversion process
- Better keyboard navigation (arrows, j/k, Enter, Esc)
- Professional completion screens with next steps
- Helpful error screens with solutions

---

## [1.x.x] - Legacy Version

The original version of InkWash is preserved on the [v1-legacy branch](https://github.com/VexoaXYZ/InkWash/tree/v1-legacy).

### Note About v1.x

Version 1.x is no longer maintained. We recommend everyone upgrade to v2.0 for the best experience!

**Why upgrade?**
- ✅ Easier to use (interactive wizards)
- ✅ Faster performance
- ✅ Better looking
- ✅ More features (mod converter!)
- ✅ Better security
- ✅ Active development

---

## 🔮 Future Plans

We're always working to make InkWash better! Here's what's coming:

### v2.1.0 (Next Release)
- Server templates (save your favorite configurations)
- Batch operations (start/stop multiple servers at once)
- Server cloning (duplicate existing servers)
- Resource manager (install popular resources with one click)

### v2.2.0
- Discord integration (bot for server status)
- Web dashboard (manage servers from your browser)
- Automatic backups (schedule regular backups)
- Update notifications (know when new InkWash versions are out)

### Later
- Linux support improvements
- MacOS support
- Docker integration
- Cloud hosting integration

Have ideas? [Open an issue](https://github.com/VexoaXYZ/InkWash/issues) and let us know!

---

<div align="center">

**[Download Latest Version](https://github.com/VexoaXYZ/InkWash/releases/latest)** | **[View All Releases](https://github.com/VexoaXYZ/InkWash/releases)**

</div>
