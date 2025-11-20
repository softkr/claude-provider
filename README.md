# Claude Code API Switcher

<div align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat" alt="License">
  <img src="https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-blue?style=flat" alt="Platforms">
  <img src="https://img.shields.io/badge/Version-2.0.0-orange?style=flat" alt="Version">
</div>

<br>

🤖 **A powerful cross-platform CLI tool to seamlessly switch between different API providers for Claude Code.**

Perfect for developers who want to use alternative models like Z.AI's GLM series with their existing Claude Code setup without complex configuration management.

## ✨ Key Features

- 🚀 **Cross-platform**: Works on macOS, Linux, and Windows
- 🔧 **Single binary**: No dependencies required after installation
- ⚡ **Easy installation**: One-command install with pre-built binaries
- 💾 **Smart backup**: Automatically backs up and restores configurations
- 🎨 **Beautiful UI**: User-friendly interface with colored messages
- 📊 **Status monitoring**: View current configuration at a glance
- 🔒 **Secure**: No data transmission, local-only configuration
- 🛠️ **Developer-friendly**: Clean Go codebase with comprehensive testing

## Quick Start

### Option 1: Install with automatic script (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/softkr/claude-provider/main/install-go.sh | bash
```

### Option 2: Manual installation

1. **Download the latest binary** for your platform:
   - [Linux AMD64](https://github.com/softkr/claude-provider/releases/latest/download/claude-switch-linux-amd64)
   - [Linux ARM64](https://github.com/softkr/claude-provider/releases/latest/download/claude-switch-linux-arm64)
   - [macOS AMD64](https://github.com/softkr/claude-provider/releases/latest/download/claude-switch-darwin-amd64)
   - [macOS ARM64](https://github.com/softkr/claude-provider/releases/latest/download/claude-switch-darwin-arm64)
   - [Windows AMD64](https://github.com/softkr/claude-provider/releases/latest/download/claude-switch-windows-amd64.exe)

2. **Make it executable** (Unix systems):
   ```bash
   chmod +x claude-switch-*
   sudo mv claude-switch-* /usr/local/bin/claude-switch
   ```

3. **Add aliases** to your shell:
   ```bash
   echo 'alias claude-switch="/usr/local/bin/claude-switch"' >> ~/.bashrc
   echo 'alias claude-status="/usr/local/bin/claude-switch --status"' >> ~/.bashrc
   source ~/.bashrc
   ```

### Option 3: Build from source

```bash
git clone https://github.com/softkr/claude-provider.git
cd claude-provider
make install
```

## Usage

After installation, you can use the following commands:

```bash
# Switch to Z.AI GLM models (backs up Anthropic web login token)
claude-switch -z
# or: claude-switch --zai

# Switch to Anthropic Claude (restores web login token from backup)
claude-switch -a
# or: claude-switch --anthropic

# Check current configuration
claude-switch -s
# or: claude-switch --status

# Show help
claude-switch -h
```

### Quick Aliases

After installation, these convenient aliases are available:

```bash
claude-zai        # Same as claude-switch -z
claude-anthropic  # Same as claude-switch -a
claude-status     # Same as claude-switch -s
```

## Authentication

| Provider | Auth Type | Description |
|----------|-----------|-------------|
| **Anthropic** | Web Login Token | Automatically backed up when switching away |
| **Z.AI** | API Key | Prompted on first use, can be saved for future |

**Important**: When switching to Z.AI, your Anthropic web login token is automatically backed up. Use `claude-switch -a` to restore it later.

### Z.AI API Key Setup

**Option 1: Interactive prompt (Recommended)**
```bash
claude-switch -z
# Enter your API key when prompted
# Choose to save for future use
```

**Option 2: Environment variable**
```bash
export ZAI_AUTH_TOKEN="your-api-key-here"
claude-switch -z
```

**Option 3: Token file**
```bash
echo "your-api-key-here" > ~/.claude/.zai_token
chmod 600 ~/.claude/.zai_token
```

### Token Management

```bash
claude-switch --clear-token  # Remove saved token
claude-switch -s             # Check token status (masked)
```

## Supported Providers

### Anthropic
- **Models**: Claude 3.5 Sonnet, Claude 3 Opus, Claude 3 Haiku
- **Endpoint**: `api.anthropic.com` (default)
- **Use when**: You need official Anthropic models with the highest reliability

### Z.AI
- **Models**: GLM-4.6, GLM-4.5-Air
- **Endpoint**: `https://api.z.ai/api/anthropic`
- **Use when**: You want to try alternative models with potentially better performance

## Configuration Files

The tool manages configuration in your home directory:

- **Settings**: `~/.claude/settings.json`
- **Backup**: `~/.claude/settings.json.backup`

### Example Z.AI Configuration
```json
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "your-token-here",
    "ANTHROPIC_BASE_URL": "https://api.z.ai/api/anthropic",
    "API_TIMEOUT_MS": "3000000",
    "ANTHROPIC_DEFAULT_OPUS_MODEL": "GLM-4.6",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "GLM-4.6",
    "ANTHROPIC_DEFAULT_HAIKU_MODEL": "GLM-4.5-Air"
  }
}
```

## Development

### Building from source

```bash
# Clone the repository
git clone https://github.com/softkr/claude-provider.git
cd claude-provider

# Install dependencies
make deps

# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Install locally
make install
```

### Project Structure

```
claude-provider/
├── main.go              # Main application code
├── go.mod               # Go module definition
├── Makefile             # Build system
├── install-go.sh        # Installation script
├── build/               # Build output directory
├── docs/                # Documentation
└── README.md            # This file
```

## Troubleshooting

### Common Issues

1. **Permission denied**
   ```bash
   sudo chmod +x /usr/local/bin/claude-switch
   ```

2. **Command not found**
   ```bash
   # Add to your PATH or create aliases
   echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
   source ~/.bashrc
   ```

3. **Configuration not working**
   ```bash
   # Check Claude Code configuration directory
   ls -la ~/.claude/

   # Verify current settings
   claude-switch --status
   ```

### Getting Help

```bash
# Show all available commands
claude-switch --help

# Check current status
claude-switch --status

# Reinstall aliases
/usr/local/bin/claude-switch --install
```

## Migration from Bash Version

If you're migrating from the original Bash version:

1. **Backup your current settings** (the Go version will do this automatically)
2. **Install the Go version** using the installation script
3. **Run your first switch**:
   ```bash
   claude-switch --anthropic  # or --zai
   ```
4. **Remove the old Bash version** (optional)

The Go version is fully compatible with existing configuration files.

## Security

- 🔒 No sensitive data is transmitted to external servers
- 🔒 API tokens are stored locally in your configuration file
- 🔒 The tool only modifies local configuration files
- 🔒 Open source - you can audit the code yourself

## License

MIT License - see LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## What's New in v2.0.0

- **Short flags**: `-a`, `-z`, `-s` for quick access
- **Secure token management**: API tokens are never hardcoded
- **Atomic file writes**: Prevents configuration corruption
- **Smart backup/restore**: Automatically manages web login tokens
- **Fish shell support**: Works with bash, zsh, and fish
- **Enhanced status display**: Better visibility of current configuration
- **Token management**: Save, clear, and mask tokens securely

## Support

- 📖 [Documentation](https://github.com/softkr/claude-provider/wiki)
- 🐛 [Issues](https://github.com/softkr/claude-provider/issues)
- 💬 [Discussions](https://github.com/softkr/claude-provider/discussions)

---

**Made with ❤️ for the Claude Code community**