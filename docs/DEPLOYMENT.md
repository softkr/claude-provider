# Deployment and Release Guide

This guide covers how to build, package, and distribute Claude Code API Switcher across different platforms.

## Build System

### Prerequisites

- **Go 1.21+**: Required for building the application
- **Make**: For using the build system (optional but recommended)
- **Git**: For version control and release tagging

### Local Development Build

```bash
# Clone repository
git clone https://github.com/claude-provider/switch.git
cd switch

# Install dependencies
make deps

# Build for current platform
make build

# Run locally
./build/claude-switch --help
```

### Cross-Platform Build

```bash
# Build for all supported platforms
make build-all

# Build specific platform
GOOS=linux GOARCH=amd64 go build -o claude-switch-linux-amd64 .

# Build with version info
make build LDFLAGS="-ldflags \"-X main.version=v1.0.0\""
```

## Build Targets

### Makefile Commands

| Target | Description | Output |
|--------|-------------|--------|
| `make build` | Build for current platform | `build/claude-switch` |
| `make build-all` | Build all platforms | `build/claude-switch-*` |
| `make release` | Create release packages | `build/release/*` |
| `make install` | Install to `/usr/local/bin` | `/usr/local/bin/claude-switch` |
| `make test` | Run tests | Test results |
| `make clean` | Clean build artifacts | Removes `build/` |

### Supported Platforms

| Platform | Architecture | Binary Name |
|----------|--------------|-------------|
| Linux | AMD64 | `claude-switch-linux-amd64` |
| Linux | ARM64 | `claude-switch-linux-arm64` |
| macOS | AMD64 | `claude-switch-darwin-amd64` |
| macOS | ARM64 | `claude-switch-darwin-arm64` |
| Windows | AMD64 | `claude-switch-windows-amd64.exe` |

## Release Process

### Version Management

```bash
# Update version in Makefile
vim Makefile

# Update version in code if needed
vim main.go

# Commit changes
git add .
git commit -m "Bump version to v1.0.0"
git tag v1.0.0
git push origin main --tags
```

### Automated Release (GitHub Actions)

The project uses GitHub Actions for automated releases:

1. **Trigger**: Create a new tag (e.g., `v1.0.0`)
2. **Build**: Cross-platform compilation
3. **Package**: Create archives and checksums
4. **Release**: Upload to GitHub Releases
5. **Update**: Update installation scripts

### Manual Release Process

```bash
# 1. Build all platforms
make clean
make build-all

# 2. Create release packages
make release

# 3. Generate checksums
cd build/release
sha256sum * > checksums.txt
cd ../..

# 4. Create GitHub release
gh release create v1.0.0 \
  --title "Release v1.0.0" \
  --notes "See CHANGELOG.md for details" \
  build/release/*

# 5. Update installation script (if needed)
vim install-go.sh
git add install-go.sh
git commit -m "Update installation script for v1.0.0"
git push origin main
```

## Package Formats

### Release Packages

Each platform binary is packaged differently:

#### Unix Systems (Linux/macOS)

```bash
# Create tar.gz package
tar -czf claude-switch-v1.0.0-linux-amd64.tar.gz \
  claude-switch-linux-amd64 \
  README.md \
  LICENSE

# Contents:
# claude-switch-v1.0.0-linux-amd64/
# ├── claude-switch-linux-amd64
# ├── README.md
# └── LICENSE
```

#### Windows

```bash
# Create zip package
zip claude-switch-v1.0.0-windows-amd64.zip \
  claude-switch-windows-amd64.exe \
  README.md \
  LICENSE

# Contents:
# claude-switch-v1.0.0-windows-amd64/
# ├── claude-switch-windows-amd64.exe
# ├── README.md
# └── LICENSE
```

### Checksum Files

```bash
# Generate SHA256 checksums
sha256sum claude-switch-* > checksums.txt

# Example checksums.txt:
# 5a1b2c3d...  claude-switch-linux-amd64
# e4f5a6b7...  claude-switch-darwin-arm64
# 8c9d0e1f...  claude-switch-windows-amd64.exe
```

## Distribution Channels

### GitHub Releases (Primary)

```bash
# Download latest release
curl -L https://github.com/claude-provider/switch/releases/latest/download/claude-switch-linux-amd64 \
  -o claude-switch

# Download specific version
curl -L https://github.com/claude-provider/switch/releases/download/v1.0.0/claude-switch-linux-amd64 \
  -o claude-switch
```

### Installation Script

The installation script automatically detects platform and downloads the appropriate binary:

```bash
# One-line installation
curl -sSL https://raw.githubusercontent.com/claude-provider/switch/main/install-go.sh | bash

# Installation with specific method
curl -sSL https://raw.githubusercontent.com/claude-provider/switch/main/install-go.sh | bash -s --prebuilt
```

### Package Managers (Future)

#### Homebrew (macOS)

```bash
# Formula file: Formula/claude-switch.rb
class ClaudeSwitch < Formula
  url "https://github.com/claude-provider/switch/archive/v1.0.0.tar.gz"
  sha256 "..."
  license "MIT"

  def install
    bin.install "build/claude-switch-darwin-#{Hardware::CPU.arch}"
  end
end
```

#### APT (Debian/Ubuntu)

```bash
# Create Debian package
dpkg-deb --build claude-switch_1.0.0_amd64.deb

# Install
sudo dpkg -i claude-switch_1.0.0_amd64.deb
```

#### Chocolatey (Windows)

```bash
# Chocolatey package definition
# claude-switch.nuspec
# tools/chocolateyinstall.ps1
```

## CI/CD Pipeline

### GitHub Actions Workflow

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Build
        run: make build-all

      - name: Package
        run: make release

      - name: Create Release
        uses: actions/create-release@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          tag_name: ${{ github.ref }}
          release_name: Release ${{ github.ref }}
          draft: false
          prerelease: false
```

### Build Matrix

```yaml
strategy:
  matrix:
    include:
      - goos: linux
        goarch: amd64
      - goos: linux
        goarch: arm64
      - goos: darwin
        goarch: amd64
      - goos: darwin
        goarch: arm64
      - goos: windows
        goarch: amd64
```

## Security Considerations

### Binary Signing

#### macOS (Notarization)

```bash
# Sign the binary
codesign --force --verify --verbose \
  --sign "Developer ID Application: Your Name" \
  claude-switch-darwin-amd64

# Notarize the app
xcrun altool --notarize-app \
  --primary-bundle-id "com.claude-provider.switch" \
  --username "your@email.com" \
  --password "@keychain:AC_PASSWORD" \
  --file claude-switch-darwin-amd64.tar.gz
```

#### Windows (Code Signing)

```bash
# Sign the binary
signtool sign /f certificate.pfx /p password \
  /t http://timestamp.digicert.com \
  claude-switch-windows-amd64.exe
```

### Checksum Verification

```bash
# Verify download integrity
curl -L https://github.com/claude-provider/switch/releases/latest/download/checksums.txt \
  -o checksums.txt

sha256sum -c checksums.txt
```

## Versioning Strategy

### Semantic Versioning

- **Major (X.0.0)**: Breaking changes
- **Minor (X.Y.0)**: New features, backward compatible
- **Patch (X.Y.Z)**: Bug fixes, security updates

### Release Types

| Type | Description | Example |
|------|-------------|---------|
| **Stable** | Production-ready releases | `v1.0.0`, `v1.1.0` |
| **Beta** | Feature previews | `v1.2.0-beta.1` |
| **Alpha** | Early testing | `v2.0.0-alpha.1` |
| **RC** | Release candidates | `v1.0.0-rc.1` |

### Changelog

Maintain `CHANGELOG.md` with:

```markdown
# Changelog

## [1.0.0] - 2024-01-15

### Added
- Initial release
- Cross-platform support
- Z.AI provider support
- Configuration backup/restore

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- N/A
```

## Performance Optimization

### Build Optimization

```bash
# Optimize binary size
go build -ldflags="-s -w" -o claude-switch .

# Strip debug information
strip claude-switch

# UPX compression (optional)
upx --best claude-switch
```

### Dependencies

```bash
# Check for unused dependencies
go mod tidy

# Vendor dependencies (if needed)
go mod vendor

# Analyze dependency tree
go mod graph
```

## Testing

### Automated Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Integration Testing

```bash
# Test binary on different systems
docker run --rm -v $(pwd):/app ubuntu:20.04 /app/build/claude-switch-linux-amd64 --help

# Test installation script
docker run --rm -v $(pwd):/app ubuntu:20.04 /app/install-go.sh --prebuilt
```

### Performance Testing

```bash
# Benchmark operations
go test -bench=. ./...

# Memory profiling
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

## Troubleshooting

### Common Build Issues

| Issue | Solution |
|-------|----------|
| Go version too old | Update to Go 1.21+ |
| Missing dependencies | Run `go mod tidy` |
| Cross-compilation errors | Install cross-compilation tools |
| Permission denied | Check file permissions |

### Release Issues

| Issue | Solution |
|-------|----------|
| Upload failed | Check GitHub token permissions |
| Checksum mismatch | Regenerate checksums |
| Installation fails | Verify binary integrity |
| Platform detection fails | Update install script |

This deployment guide ensures reliable, secure, and efficient distribution of Claude Code API Switcher across all supported platforms.