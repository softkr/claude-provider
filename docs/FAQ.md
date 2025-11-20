# Frequently Asked Questions (FAQ)

This document answers common questions about Claude Code API Switcher.

## General Questions

### Q: What is Claude Code API Switcher?

**A:** Claude Code API Switcher is a command-line tool that allows you to easily switch between different API providers for Claude Code. It supports Anthropic's official API and alternative providers like Z.AI, enabling you to use different AI models with your existing Claude Code setup.

### Q: Why would I need this tool?

**A:** You might need this tool if:
- You want to try alternative models like Z.AI's GLM series
- You need to switch between different providers for testing
- You want to compare performance between different models
- You need to use specific providers for different projects

### Q: Is this tool officially supported by Anthropic?

**A:** No, this is an independent open-source tool. It's not officially supported by Anthropic, but it's designed to work seamlessly with Claude Code.

## Installation and Setup

### Q: How do I install the tool?

**A:** The easiest way is using the installation script:
```bash
curl -sSL https://raw.githubusercontent.com/softkr/claude-provider/main/install-go.sh | bash
```

You can also download pre-built binaries or build from source. See the [Installation Guide](../README.md#quick-start) for detailed instructions.

### Q: What are the system requirements?

**A:**
- **Operating System**: macOS, Linux, or Windows
- **Claude Code**: Must be installed separately
- **Permissions**: Write access to your home directory for configuration files

### Q: Do I need Go installed to use this tool?

**A:** No, you only need Go if you want to build from source. The pre-built binaries are self-contained and don't require Go.

### Q: How do I uninstall the tool?

**A:**
```bash
# Remove the binary
sudo rm /usr/local/bin/claude-switch

# Remove aliases from your shell config (~/.bashrc or ~/.zshrc)
# Edit the file and remove the "Claude Code API Switcher" section

# Remove configuration files (optional)
rm -rf ~/.claude/
```

## Usage

### Q: How do I switch between API providers?

**A:** Use the following commands:
```bash
claude-switch --anthropic  # Switch to Anthropic API
claude-switch --zai        # Switch to Z.AI API
claude-switch --status     # Check current configuration
```

### Q: What happens to my current configuration when I switch providers?

**A:** The tool automatically creates a backup of your existing configuration before making changes. You can restore it by switching back to your previous provider.

### Q: Can I use custom API providers?

**A:** Yes, you can manually edit the configuration file:
```bash
# Switch to Anthropic base first
claude-switch --anthropic

# Edit the configuration
nano ~/.claude/settings.json

# Add your custom provider configuration
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "your-token",
    "ANTHROPIC_BASE_URL": "https://your-provider.com/api",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "your-model"
  }
}
```

### Q: How do I check which provider I'm currently using?

**A:** Run the status command:
```bash
claude-switch --status
```

This will show you the current provider, model, and other configuration details.

## Configuration

### Q: Where are the configuration files stored?

**A:** Configuration files are stored in your home directory:
```
~/.claude/
├── settings.json           # Current configuration
└── settings.json.backup    # Backup of original configuration
```

### Q: Can I use my own API token?

**A:** Yes! You can edit the configuration file to use your own API token:
```json
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "your-own-token-here"
  }
}
```

### Q: How do I configure timeout settings?

**A:** Edit the configuration file and set the timeout:
```json
{
  "env": {
    "API_TIMEOUT_MS": "60000"  // 60 seconds
  }
}
```

### Q: Can I have different configurations for different projects?

**A:** Currently, the tool uses a global configuration. For project-specific configurations, you can:
1. Use environment variables to override settings
2. Manually copy different configuration files
3. Use the tool's backup/restore feature

## API Providers

### Q: Which API providers are supported?

**A:** Currently supported:
- **Anthropic**: Official Claude models (Claude 3.5 Sonnet, Claude 3 Opus, Claude 3 Haiku)
- **Z.AI**: GLM models (GLM-4.6, GLM-4.5-Air)

You can also manually configure any Anthropic-compatible API provider.

### Q: What's the difference between Anthropic and Z.AI models?

**A:**
- **Anthropic models**: Official Claude models with proven reliability and features
- **Z.AI models**: Alternative GLM models that may offer different performance characteristics

### Q: Are there costs associated with using different providers?

**A:** Yes, each provider has its own pricing:
- **Anthropic**: Usage-based pricing (check Anthropic's website for current rates)
- **Z.AI**: May have different pricing tiers (check Z.AI's website)

Always check the pricing before switching providers.

### Q: Do all providers support the same features?

**A:** Not necessarily. While all providers support basic text generation, some features like:
- Image analysis
- Tool/function calling
- Streaming responses

may vary between providers. Check the provider's documentation for feature support.

## Troubleshooting

### Q: The command isn't found after installation. What should I do?

**A:**
1. Make sure your shell configuration is reloaded:
   ```bash
   source ~/.bashrc  # or ~/.zshrc
   ```

2. Check if the binary exists:
   ```bash
   ls -la /usr/local/bin/claude-switch
   ```

3. Verify your PATH includes `/usr/local/bin`:
   ```bash
   echo $PATH | grep -q "/usr/local/bin"
   ```

### Q: I'm getting "permission denied" errors. How do I fix this?

**A:**
```bash
# Fix binary permissions
sudo chmod +x /usr/local/bin/claude-switch

# Fix configuration directory permissions
mkdir -p ~/.claude
chmod 700 ~/.claude
chmod 600 ~/.claude/settings.json*
```

### Q: The tool isn't switching providers correctly. What should I check?

**A:**
1. Check your current status:
   ```bash
   claude-switch --status
   ```

2. Verify your API token is valid
3. Check network connectivity to the API endpoint
4. Look for error messages in the output

### Q: Claude Code isn't using the new configuration. Why?

**A:**
1. Make sure Claude Code is restarted after switching
2. Check that the configuration file exists and is valid:
   ```bash
   cat ~/.claude/settings.json | jq .
   ```
3. Verify environment variables are set correctly in Claude Code

### Q: How do I reset everything to the default state?

**A:**
```bash
# Remove all configurations
rm -rf ~/.claude/

# Switch back to Anthropic
claude-switch --anthropic

# This will create a fresh, empty configuration
```

## Security

### Q: Is my API token secure?

**A:**
- **Storage**: Your API token is stored locally in plain text
- **Network**: The tool doesn't transmit your token to any external servers
- **Recommendation**: Use environment variables for additional security in production

### Q: Should I commit the configuration file to version control?

**A:** **No!** Never commit `~/.claude/settings.json` to version control as it contains your API tokens. Add it to your `.gitignore`:
```bash
echo ".claude/" >> ~/.gitignore
```

### Q: How can I improve security?

**A:**
1. Use environment variables instead of storing tokens in files
2. Set restrictive file permissions:
   ```bash
   chmod 600 ~/.claude/settings.json*
   ```
3. Regularly rotate your API tokens
4. Use different tokens for different environments

## Technical

### Q: What's the difference between the Bash and Go versions?

**A:**
- **Bash version**: Original implementation, limited to Unix-like systems
- **Go version**: Complete rewrite with cross-platform support, better error handling, and enhanced features

The Go version is recommended for all users.

### Q: Can I build the tool myself?

**A:** Yes! See the [Development Guide](CONTRIBUTING.md) for build instructions:
```bash
git clone https://github.com/softkr/claude-provider.git
cd switch
make build
```

### Q: How do I report bugs or request features?

**A:**
- **Bugs**: [Create an issue on GitHub](https://github.com/softkr/claude-provider/issues)
- **Features**: [Create a feature request](https://github.com/softkr/claude-provider/issues/new?template=feature_request.md)
- **Security**: Email security@claude-provider.com (private)

### Q: Is there a command-line interface (CLI) help?

**A:** Yes, run:
```bash
claude-switch --help
```

## Integration

### Q: Can I use this with Docker?

**A:** Yes, you can mount the configuration directory in Docker:
```bash
docker run -v ~/.claude:/root/.claude your-claude-image
```

### Q: Does this work with CI/CD pipelines?

**A:** Yes, you can use environment variables in CI/CD:
```bash
export ANTHROPIC_AUTH_TOKEN="$CI_API_TOKEN"
export ANTHROPIC_BASE_URL="$CI_API_ENDPOINT"
```

### Q: Can I automate provider switching?

**A:** Yes, you can script the commands:
```bash
#!/bin/bash
# Switch to Z.AI for testing
claude-switch --zai
# Run your tests
./run-tests.sh
# Switch back to Anthropic
claude-switch --anthropic
```

## Comparison with Alternatives

### Q: How is this different from manually editing configuration files?

**A:**
- **Automation**: Handles backup and restore automatically
- **Safety**: Validates configurations before applying
- **Convenience**: Single command instead of manual editing
- **Error Handling**: Provides clear error messages

### Q: Should I use this or environment variables directly?

**A:**
- **This tool**: Better for permanent configuration changes
- **Environment variables**: Better for temporary overrides

You can use both together - the tool for base configuration, environment variables for overrides.

---

## Still Need Help?

If you have questions not covered here:

- **Documentation**: Check the [full documentation](../docs/)
- **GitHub Issues**: [Browse existing issues](https://github.com/softkr/claude-provider/issues)
- **GitHub Discussions**: [Join the community discussion](https://github.com/softkr/claude-provider/discussions)
- **Email**: support@claude-provider.com

---

*Last updated: January 15, 2024*