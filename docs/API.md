# API Documentation

This document provides detailed information about the configuration and API settings used by Claude Code API Switcher.

## Configuration Structure

### Settings File Location

The tool manages configuration files in your user home directory:

```
~/.claude/
├── settings.json           # Current active configuration
└── settings.json.backup    # Backup of original configuration
```

### Configuration Format

The settings file uses JSON format with the following structure:

```json
{
  "env": {
    "VARIABLE_NAME": "value",
    "ANOTHER_VARIABLE": "value"
  }
}
```

## Supported Environment Variables

### Core API Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `ANTHROPIC_AUTH_TOKEN` | API authentication token | - | Yes |
| `ANTHROPIC_BASE_URL` | API base endpoint URL | `https://api.anthropic.com` | No |
| `API_TIMEOUT_MS` | Request timeout in milliseconds | `300000` | No |

### Model Configuration Variables

| Variable | Description | Example Values |
|----------|-------------|----------------|
| `ANTHROPIC_DEFAULT_OPUS_MODEL` | Default model for opus tier requests | `claude-3-5-sonnet-20241022`, `GLM-4.6` |
| `ANTHROPIC_DEFAULT_SONNET_MODEL` | Default model for sonnet tier requests | `claude-3-5-sonnet-20241022`, `GLM-4.6` |
| `ANTHROPIC_DEFAULT_HAIKU_MODEL` | Default model for haiku tier requests | `claude-3-haiku-20240307`, `GLM-4.5-Air` |

## Provider Configurations

### Anthropic Configuration

```json
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "sk-ant-api03-...",
    "ANTHROPIC_BASE_URL": "https://api.anthropic.com",
    "API_TIMEOUT_MS": "300000",
    "ANTHROPIC_DEFAULT_OPUS_MODEL": "claude-3-opus-20240229",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "claude-3-5-sonnet-20241022",
    "ANTHROPIC_DEFAULT_HAIKU_MODEL": "claude-3-haiku-20240307"
  }
}
```

### Z.AI Configuration

```json
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "d33394705bc34f7c872397cc71c5dc01.DXjF6Nx5t09MSTGn",
    "ANTHROPIC_BASE_URL": "https://api.z.ai/api/anthropic",
    "API_TIMEOUT_MS": "3000000",
    "ANTHROPIC_DEFAULT_OPUS_MODEL": "GLM-4.6",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "GLM-4.6",
    "ANTHROPIC_DEFAULT_HAIKU_MODEL": "GLM-4.5-Air"
  }
}
```

## Authentication Tokens

### Getting Anthropic API Token

1. Visit [Anthropic Console](https://console.anthropic.com/)
2. Sign in or create an account
3. Navigate to API Keys
4. Create a new API key
5. Copy the token (starts with `sk-ant-api03-`)

### Getting Z.AI API Token

1. Visit [Z.AI Platform](https://z.ai/)
2. Register or sign in
3. Navigate to API settings
4. Generate or copy your API token

## API Endpoints

### Anthropic Endpoints

- **Base URL**: `https://api.anthropic.com`
- **Version**: `2023-06-01`
- **Models**: Claude 3.5 Sonnet, Claude 3 Opus, Claude 3 Haiku

### Z.AI Endpoints

- **Base URL**: `https://api.z.ai/api/anthropic`
- **Compatibility**: Anthropic-compatible API
- **Models**: GLM-4.6, GLM-4.5-Air

## Model Mappings

### Anthropic Models

| Model | Context | Use Case |
|-------|---------|----------|
| `claude-3-5-sonnet-20241022` | 200K | General purpose, balanced performance |
| `claude-3-opus-20240229` | 200K | Complex tasks, highest quality |
| `claude-3-haiku-20240307` | 200K | Fast responses, simple tasks |

### Z.AI Models

| Model | Context | Use Case |
|-------|---------|----------|
| `GLM-4.6` | ~128K | Advanced reasoning, multilingual |
| `GLM-4.5-Air` | ~128K | Balanced performance, efficient |

## Error Handling

### Common API Errors

| Error | Description | Solution |
|-------|-------------|----------|
| `401 Unauthorized` | Invalid or missing API token | Verify your authentication token |
| `429 Rate Limited` | Too many requests | Wait and retry, check rate limits |
| `500 Server Error` | API server error | Try again later |
| `Timeout` | Request took too long | Increase `API_TIMEOUT_MS` |

### Configuration Errors

| Error | Description | Solution |
|-------|-------------|----------|
| `Invalid JSON` | Malformed settings file | Run `claude-switch --status` to validate |
| `Permission Denied` | Cannot write to config directory | Check file permissions |
| `Missing Variables` | Required env vars not set | Switch to a valid provider |

## Advanced Configuration

### Custom Provider

You can manually configure any Anthropic-compatible API:

```bash
# Create custom configuration
claude-switch --anthropic  # Start with Anthropic base

# Edit settings manually
nano ~/.claude/settings.json

# Add your custom provider
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "your-custom-token",
    "ANTHROPIC_BASE_URL": "https://your-custom-endpoint.com/api",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "your-custom-model"
  }
}
```

### Environment Variable Override

You can override settings using environment variables:

```bash
export ANTHROPIC_BASE_URL="https://custom-endpoint.com"
export ANTHROPIC_AUTH_TOKEN="custom-token"

# These will take precedence over settings.json
claude
```

## Security Considerations

### Token Security

- 🔒 API tokens are stored locally in plain text
- 🔒 Never commit settings files to version control
- 🔒 Use environment variables for sensitive data in production
- 🔒 Regularly rotate your API tokens

### File Permissions

```bash
# Secure your configuration files
chmod 600 ~/.claude/settings.json*
chmod 700 ~/.claude/
```

## Troubleshooting

### Configuration Validation

```bash
# Check current configuration
claude-switch --status

# Validate JSON syntax
cat ~/.claude/settings.json | jq .

# Test API connection
curl -H "x-api-key: $ANTHROPIC_AUTH_TOKEN" \
     -H "content-type: application/json" \
     -d '{"model": "claude-3-5-sonnet-20241022", "max_tokens": 10, "messages": [{"role": "user", "content": "test"}]}' \
     "$ANTHROPIC_BASE_URL/v1/messages"
```

### Debug Mode

Set debug environment variable for verbose output:

```bash
export CLAUDE_SWITCH_DEBUG=1
claude-switch --status
```

## API Compatibility

### Anthropic API Specification

This tool follows the Anthropic API specification:

- **Messages API**: `/v1/messages`
- **Models API**: `/v1/models`
- **Authentication**: `x-api-key` header
- **Content-Type**: `application/json`

### Compatibility Matrix

| Feature | Anthropic | Z.AI | Custom |
|---------|-----------|------|--------|
| Streaming | ✅ | ✅ | Depends |
| Tools/Function Calling | ✅ | ⚠️ Partial | Depends |
| Image Analysis | ✅ | ⚠️ Limited | Depends |
| System Prompts | ✅ | ✅ | Depends |

## Rate Limits

### Anthropic Rate Limits

- **Free Tier**: 1,000 requests per month
- **Build Tier**: 100,000 requests per month
- **Pro Tier**: 1,000,000 requests per month

### Z.AI Rate Limits

- **Free Tier**: 100 requests per day
- **Pro Tier**: 1,000 requests per day
- **Enterprise**: Custom limits

## API Versioning

The tool is compatible with:

- **Anthropic API**: 2023-06-01 and later
- **Message Format**: Current message API format
- **Authentication**: Token-based authentication

For the latest API specifications, visit:
- [Anthropic API Documentation](https://docs.anthropic.com/claude/reference)
- [Z.AI API Documentation](https://docs.z.ai/)