# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2024-11-21

### Added
- **Short flags**: `-a`, `-z`, `-s` for quick access
- **Secure token management**: API tokens are never hardcoded
- **Token prompt system**: Interactive token input with save option
- **Environment variable support**: `Z_AI_AUTH_TOKEN` for Z.AI API key
- **Token file storage**: `~/.claude/.z_ai_token` for saved tokens
- **Token masking**: Displays tokens as `xxxx...xxxx` in status
- **Clear token command**: `--clear-token` to remove saved tokens
- **Fish shell support**: Works with bash, zsh, and fish
- **Quick aliases**: `claude-anthropic`, `claude-z_ai`, `claude-status`
- **Enhanced status display**: Box UI with detailed configuration info

### Changed
- **Authentication model**:
  - Anthropic uses web login token (backed up automatically)
  - Z.AI uses API key (prompted or from environment)
- **Atomic file writes**: Prevents configuration corruption
- **Improved backup/restore**: Automatically manages web login tokens
- **File permissions**: Changed to 0600 for security
- **Directory permissions**: Changed to 0700 for security
- **Repository moved**: Now at `github.com/softkr/claude-provider`

### Fixed
- Partial write issues with atomic temp file + rename
- Token security (removed hardcoded tokens)

### Security
- API tokens no longer hardcoded in source
- Secure file permissions (0600)
- Token masking in status display

## [1.0.0] - 2024-01-15

### Added
- **Initial Go Release**
  - Complete rewrite of original Bash script in Go
  - Single binary distribution for easy installation
  - Cross-platform compatibility (macOS, Linux, Windows)

- **API Provider Support**
  - Anthropic API configuration
  - Z.AI GLM model support
  - Easy switching between providers with single command

- **Configuration Management**
  - Automatic backup of existing configurations
  - JSON-based configuration format
  - Environment variable management

- **User Interface**
  - Colorful terminal output with emoji indicators
  - Clear status reporting
  - Comprehensive help system

- **Installation System**
  - Automatic platform detection
  - One-command installation script
  - Shell alias management

## [0.1.0] - 2024-01-10 (Bash Version)

### Added
- Initial Bash implementation
- Basic API switching functionality
- Simple configuration management

---

## Migration Guide

### From v1.x to v2.0.0

#### New Features
```bash
# Short flags (new)
claude-switch -z   # instead of --z_ai
claude-switch -a   # instead of --anthropic
claude-switch -s   # instead of --status

# Quick aliases (new)
claude-z_ai
claude-anthropic
claude-status
```

#### Breaking Changes
- Token is no longer hardcoded - you'll be prompted on first Z.AI switch
- Backup is now mandatory when switching to Z.AI

#### Authentication Changes
- **Anthropic**: Web login token (automatically backed up)
- **Z.AI**: API key (prompted or use `Z_AI_AUTH_TOKEN` env var)

---

## Support

- **GitHub Issues**: [Report issues](https://github.com/softkr/claude-provider/issues)
- **GitHub Discussions**: [Community](https://github.com/softkr/claude-provider/discussions)
- **Documentation**: [Docs](https://github.com/softkr/claude-provider/docs)
