# Contributing to Claude Code API Switcher

Thank you for your interest in contributing to Claude Code API Switcher! This document provides guidelines and information for contributors.

## 🤝 How to Contribute

### Reporting Bugs

1. **Check existing issues**: Search for similar bugs before creating a new one
2. **Use bug report template**: Fill out the template completely
3. **Provide reproduction steps**: Include exact steps to reproduce the issue
4. **Include system information**: OS, Go version, and other relevant details
5. **Add logs/error messages**: Include any relevant console output

### Suggesting Features

1. **Check existing issues**: Search for similar feature requests
2. **Use feature request template**: Describe the feature and use case
3. **Provide implementation ideas**: If you have technical ideas, share them
4. **Consider breaking changes**: Note if the feature would require API changes

### Submitting Pull Requests

1. **Fork the repository**: Create a fork on GitHub
2. **Create a feature branch**: Use descriptive branch names
3. **Make your changes**: Follow coding standards and add tests
4. **Test thoroughly**: Ensure all tests pass and functionality works
5. **Submit PR**: Include clear description and link to relevant issues

## 🛠️ Development Setup

### Prerequisites

- **Go 1.21+**: Required for development
- **Git**: For version control
- **Make**: For build automation (recommended)
- **GitHub CLI**: For release management (optional)

### Local Development

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/switch.git
cd switch

# Add upstream remote
git remote add upstream https://github.com/softkr/claude-provider.git

# Install dependencies
make deps

# Run tests
make test

# Build locally
make build
```

### Development Workflow

```bash
# 1. Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main

# 2. Create feature branch
git checkout -b feature/your-feature-name

# 3. Make changes
# Edit files...

# 4. Test your changes
make test
go run . --help

# 5. Commit changes
git add .
git commit -m "feat: add new feature description"

# 6. Push to your fork
git push origin feature/your-feature-name

# 7. Create Pull Request
# Use GitHub web interface
```

## 📝 Coding Standards

### Code Style

Follow Go conventions and best practices:

```go
// Package main provides CLI tool for switching Claude Code API providers
package main

import (
    "context"
    "fmt"
    "os"
)

// Config represents the application configuration
type Config struct {
    SettingsFile string `json:"settings_file"`
    BackupFile   string `json:"backup_file"`
}

// NewConfig creates a new configuration instance
func NewConfig() *Config {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil
    }

    return &Config{
        SettingsFile: filepath.Join(homeDir, ".claude", "settings.json"),
        BackupFile:   filepath.Join(homeDir, ".claude", "settings.json.backup"),
    }
}
```

### Naming Conventions

- **Package names**: Short, lowercase, one word
- **Constants**: `UPPER_SNAKE_CASE`
- **Variables**: `camelCase` for local, `PascalCase` for exported
- **Functions**: `PascalCase` for exported, `camelCase` for private
- **Files**: `snake_case.go` for implementation files

### Error Handling

```go
// Good: Handle errors immediately
config, err := loadConfig(filename)
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}

// Good: Use structured errors
type ConfigError struct {
    Path string
    Err  error
}

func (e *ConfigError) Error() string {
    return fmt.Sprintf("config error at %s: %v", e.Path, e.Err)
}

func (e *ConfigError) Unwrap() error {
    return e.Err
}
```

### Testing

```go
func TestConfigLoad(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *Config
        wantErr bool
    }{
        {
            name:  "valid config",
            input: `{"env": {"KEY": "value"}}`,
            want: &Config{
                Env: map[string]string{"KEY": "value"},
            },
            wantErr: false,
        },
        {
            name:    "invalid json",
            input:   `{invalid json}`,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create temporary file
            tmpFile, err := os.CreateTemp("", "config-*.json")
            if err != nil {
                t.Fatal(err)
            }
            defer os.Remove(tmpFile.Name())

            // Write test data
            if _, err := tmpFile.WriteString(tt.input); err != nil {
                t.Fatal(err)
            }
            tmpFile.Close()

            // Test function
            got, err := loadConfig(tmpFile.Name())

            if (err != nil) != tt.wantErr {
                t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("loadConfig() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## 🧪 Testing Guidelines

### Test Structure

```
.
├── main.go              # Main application
├── main_test.go         # Main application tests
├── config.go            # Configuration logic
├── config_test.go       # Configuration tests
├── providers/
│   ├── anthropic.go     # Anthropic provider
│   ├── anthropic_test.go
│   ├── zai.go          # Z.AI provider
│   └── zai_test.go
└── testdata/           # Test data files
    ├── valid_config.json
    └── invalid_config.json
```

### Test Coverage

- **Unit tests**: Test individual functions and methods
- **Integration tests**: Test component interactions
- **End-to-end tests**: Test complete user workflows

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
go test -coverprofile=coverage.out ./...

# Run specific test
go test -run TestConfigLoad ./...

# Run benchmarks
go test -bench=. ./...

# Race condition testing
go test -race ./...
```

### Test Data

```go
// testdata/config_test.go
func TestLoadConfigFromFiles(t *testing.T) {
    testCases := []struct {
        file       string
        wantConfig *Config
        wantErr    bool
    }{
        {
            file: "testdata/valid_anthropic.json",
            wantConfig: &Config{
                Env: map[string]string{
                    "ANTHROPIC_BASE_URL": "https://api.anthropic.com",
                },
            },
            wantErr: false,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.file, func(t *testing.T) {
            config, err := loadConfig(tc.file)
            // Test implementation...
        })
    }
}
```

## 📋 Pull Request Process

### PR Checklist

Before submitting a PR, ensure:

- [ ] Code follows project style guidelines
- [ ] Self-review of the code completed
- [ ] Code is properly commented
- [ ] Tests added for new functionality
- [ ] All tests pass locally
- [ ] Documentation updated if needed
- [ ] CHANGELOG.md updated (for significant changes)
- [ ] PR description clearly explains changes

### PR Template

```markdown
## Description
Brief description of changes and motivation.

## Type of Change
- [ ] Bug fix (non-breaking change)
- [ ] New feature (non-breaking change)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Unit tests pass locally
- [ ] Integration tests pass locally
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] CHANGELOG.md updated

## Additional Notes
Any additional context or considerations.
```

### Review Process

1. **Automated checks**: CI/CD pipeline runs tests and linting
2. **Code review**: Maintainers review for code quality and functionality
3. **Approval**: At least one maintainer approval required
4. **Merge**: Maintainer merges after all checks pass

## 🏗️ Project Structure

### Directory Layout

```
switch/
├── cmd/
│   └── cli/                 # CLI application entry point
│       └── main.go
├── internal/
│   ├── config/              # Configuration management
│   │   ├── config.go
│   │   └── config_test.go
│   ├── providers/           # API providers
│   │   ├── anthropic.go
│   │   ├── anthropic_test.go
│   │   ├── zai.go
│   │   └── zai_test.go
│   ├── ui/                  # User interface
│   │   ├── color.go
│   │   └── color_test.go
│   └── utils/               # Utility functions
│       ├── file.go
│       └── file_test.go
├── pkg/                     # Public API (if any)
├── testdata/                # Test data
├── docs/                    # Documentation
├── scripts/                 # Build and deployment scripts
├── .github/
│   ├── workflows/           # GitHub Actions
│   └── ISSUE_TEMPLATE/      # Issue templates
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── CHANGELOG.md
├── CONTRIBUTING.md
└── LICENSE
```

### Package Organization

- **cmd/**: Application entry points
- **internal/**: Private application code
- **pkg/**: Public libraries (if reusable)
- **docs/**: Documentation
- **scripts/**: Automation scripts

## 🔄 Release Process

### Version Management

- Follow [Semantic Versioning](https://semver.org/)
- Update `VERSION` constant in code
- Update `CHANGELOG.md`
- Create Git tag with version number

### Release Checklist

- [ ] All tests passing
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version number updated
- [ ] Git tag created
- [ ] Release created on GitHub
- [ ] Installation scripts updated

## 💡 Guidelines for Different Types of Contributions

### Bug Fixes

1. **Reproduce the issue**: Create test that reproduces the bug
2. **Fix the issue**: Implement minimal fix
3. **Test the fix**: Ensure test passes and no regressions
4. **Document**: Add comments explaining the fix

### New Features

1. **Design first**: Create issue with design proposal
2. **Implement**: Write code following project patterns
3. **Test**: Add comprehensive tests
4. **Document**: Update documentation and README

### Documentation

1. **Consistency**: Follow existing documentation style
2. **Clarity**: Write clear, concise explanations
3. **Examples**: Include code examples where helpful
4. **Review**: Check for typos and formatting

### Performance Improvements

1. **Benchmark**: Add benchmarks for current performance
2. **Optimize**: Make improvements
3. **Compare**: Show performance gains
4. **Test**: Ensure no functionality is broken

## 🎯 Areas for Contribution

### High Priority

- **Additional providers**: Support for more API providers
- **Configuration validation**: Enhanced config validation
- **Error handling**: Better error messages and recovery
- **Testing**: Increase test coverage

### Medium Priority

- **UI improvements**: Enhanced CLI interface
- **Performance**: Optimization and caching
- **Documentation**: Additional guides and examples
- **Integration**: Package manager support

### Low Priority

- **Plugins**: Plugin system for extensibility
- **Web interface**: Optional web-based configuration
- **Monitoring**: Usage analytics and reporting
- **Advanced features**: Profile switching, templates

## 🚫 What Not to Contribute

- **Breaking changes** without discussion
- **New dependencies** without justification
- **Code that doesn't follow style guidelines**
- **Features outside project scope**
- **Security vulnerabilities** (report privately)

## 📞 Getting Help

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Email**: For security issues and confidential matters

## 📄 License

By contributing to this project, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Claude Code API Switcher! Your contributions help make this project better for everyone. 🙏