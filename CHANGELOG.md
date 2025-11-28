# Changelog

All notable changes to InkWash will be documented in this file.

## [v2.0.0] - Unreleased

### Complete Rewrite

This is a complete rewrite of InkWash from the ground up with a focus on user experience, reliability, and modern design.

### Added

- **Modern TUI Framework**: Built with Bubble Tea for smooth, interactive terminal UI
- **Interactive Server Creation Wizard**: Step-by-step guided setup with real-time validation
- **GTA5 Mod Converter**: Convert GTA5 mods to FiveM resources using convert.cfx.rs
  - Queue-based conversion system with rate limiting (max 2 concurrent)
  - Parallel downloads with progress tracking
  - Auto-extraction to category subfolders ([vehicles], [weapons], etc.)
  - Throttled UI updates to prevent scrolling issues
- **Enhanced Server Management**: Improved lifecycle tracking and status monitoring
- **Better Animation System**: Adaptive animation tiers based on terminal capabilities
- **Progress Tracking**: Real-time progress bars and spinners for all operations
- **Filesystem Safety**: Auto-convert server names to filesystem-safe slugs
- **License Key Vault**: Secure encrypted storage for FiveM license keys
- **Per-Server Binaries**: Each server has isolated FXServer binaries
- **Metadata Tracking**: Build info, lifecycle events, and usage stats

### Changed

- **New Command Structure**: Simplified and more intuitive commands
- **Config Format**: Updated configuration file structure
- **Registry Format**: Improved server registry with path validation
- **UI/UX**: Complete redesign with modern styling and better feedback
- **Error Handling**: More helpful error messages and recovery options

### Breaking Changes

- Command structure has changed (see README for new commands)
- Config files are not compatible with v1.x
- Server registry format updated (migration required)
- Old servers need to be recreated or migrated

### Fixed

- Server selector no longer shows deleted/moved servers
- Git output no longer bleeds through UI during setup
- Text inputs handle spaces and special characters correctly
- Rate limiting prevents API errors during batch operations
- Cursor management in text inputs works properly

---

## [v1.x.x] - Legacy

See the [v1-legacy branch](https://github.com/VexoaXYZ/InkWash/tree/v1-legacy) for the original version history.

### Note

Version 1.x is preserved for reference but is no longer maintained. All new development happens on v2.x.
