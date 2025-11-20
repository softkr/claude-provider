# Security Policy

This document outlines the security practices and policies for Claude Code API Switcher.

## 🛡️ Security Overview

Claude Code API Switcher is designed with security as a primary concern. The tool operates entirely on your local machine and does not transmit any data to external servers.

### Core Security Principles

1. **Local-Only Processing**: All configuration and data processing happens locally
2. **No Data Transmission**: The tool never sends your configuration to external servers
3. **Minimal Privileges**: Requires only basic file system permissions
4. **Transparent Operations**: All actions are clearly logged and visible to the user

## 🔐 Data Protection

### Configuration Files

**Location**: `~/.claude/settings.json` and `~/.claude/settings.json.backup`

**Content**: API tokens, endpoint URLs, and model preferences

**Protection Measures**:
- Stored locally in user home directory
- File permissions restricted to user only (600)
- No automatic synchronization with cloud services
- No transmission to external services

### API Tokens

**Storage**: Plain text in configuration file

**Risks**: If someone gains access to your user account, they can access your API tokens

**Mitigations**:
- Use environment variables for enhanced security
- Regular token rotation recommended
- File permissions set to 600 (user read/write only)

```bash
# Secure your configuration files
chmod 600 ~/.claude/settings.json*
chmod 700 ~/.claude/
```

### Environment Variables (Alternative Storage)

For enhanced security, you can use environment variables instead of storing tokens in files:

```bash
# Set environment variables
export ANTHROPIC_AUTH_TOKEN="your-token-here"
export ANTHROPIC_BASE_URL="https://api.anthropic.com"

# These will override settings.json values
claude
```

## 🚨 Security Considerations

### Potential Risks

1. **File System Access**
   - **Risk**: Malicious actors with file system access can read API tokens
   - **Mitigation**: Restrict file permissions, use full disk encryption

2. **Shared Systems**
   - **Risk**: Other users on shared systems might access configuration files
   - **Mitigation**: Use environment variables, ensure proper file permissions

3. **Version Control**
   - **Risk**: Accidentally committing configuration files with API tokens
   - **Mitigation**: Add `.claude/` to `.gitignore`, never commit configuration files

4. **Log Files**
   - **Risk**: API tokens might be logged in shell history or log files
   - **Mitigation**: Clear shell history, avoid logging sensitive commands

### Best Practices

1. **File Permissions**
   ```bash
   # Secure configuration directory
   chmod 700 ~/.claude/
   chmod 600 ~/.claude/settings.json*
   ```

2. **Version Control**
   ```bash
   # Never commit configuration files
   echo ".claude/" >> ~/.gitignore
   echo ".env" >> ~/.gitignore
   ```

3. **Environment Variables**
   ```bash
   # For production environments
   export CLAUDE_CONFIG_DIR="/secure/path/.claude"
   export ANTHROPIC_AUTH_TOKEN="production-token"
   ```

4. **Token Rotation**
   - Regularly rotate your API tokens
   - Revoke unused tokens immediately
   - Use different tokens for different environments

## 🔍 Security Features

### Input Validation

The tool validates all inputs before processing:

```go
// Configuration validation
func validateConfig(config *Config) error {
    if config.Env == nil {
        return errors.New("environment configuration is required")
    }

    if token := config.Env["ANTHROPIC_AUTH_TOKEN"]; token == "" {
        return errors.New("ANTHROPIC_AUTH_TOKEN is required")
    }

    return nil
}
```

### Safe File Operations

- Atomic file writes to prevent corruption
- Backup creation before modifications
- Permission checks before file operations

### Error Handling

- No sensitive information in error messages
- Safe logging practices
- Graceful failure modes

## 🏢 Enterprise Security

### Compliance Considerations

- **Data Residency**: All data remains on local machine
- **Audit Trail**: Clear logs of configuration changes
- **Access Control**: File system permissions control access

### Deployment in Enterprise Environments

1. **Centralized Configuration**
   ```bash
   # Use centralized configuration management
   export CLAUDE_CONFIG_DIR="/opt/company/.claude"
   ```

2. **Corporate API Tokens**
   - Use service accounts instead of personal tokens
   - Implement token rotation policies
   - Monitor API usage and costs

3. **Network Security**
   - Tool doesn't make external network requests
   - API calls made by Claude Code itself
   - Consider proxy configurations for enterprise networks

## 🐛 Reporting Security Issues

### Private Disclosure

For security issues, please email us privately at:

**📧 security@claude-provider.com**

### What to Include

1. **Description**: Clear description of the security issue
2. **Impact**: Potential impact of the vulnerability
3. **Reproduction**: Steps to reproduce the issue
4. **Environment**: System and configuration details

### Response Timeline

- **Initial Response**: Within 48 hours
- **Assessment**: Within 1 week
- **Resolution**: As soon as possible, based on severity

### Public Disclosure

We will coordinate public disclosure of security issues:
- Fix will be developed and tested privately
- Security advisory will be published
- Users will be notified of updates

## 🔒 Vulnerability Management

### Severity Classification

| Severity | Description | Response Time |
|----------|-------------|---------------|
| **Critical** | Immediate risk to users or data | 24-48 hours |
| **High** | Significant security impact | 1-2 weeks |
| **Medium** | Limited security impact | 1 month |
| **Low** | Minor security issue | Next release |

### Fix Process

1. **Assessment**: Evaluate impact and scope
2. **Development**: Create and test fix
3. **Review**: Security review of changes
4. **Release**: Coordinate security update
5. **Communication**: Notify users and publish advisory

## 🛠️ Security Development Practices

### Code Review

All code changes undergo security review:
- Input validation
- Error handling
- File operations
- Dependency security

### Dependencies

- Minimal dependencies to reduce attack surface
- Regular security updates
- Vulnerability scanning

```bash
# Check for known vulnerabilities
go list -m -u all
govulncheck ./...
```

### Testing

- Security-focused unit tests
- Integration tests with security scenarios
- Penetration testing for critical components

## 📋 Security Checklist

### For Users

- [ ] Secure configuration file permissions
- [ ] Regular API token rotation
- [ ] Never commit configuration files to version control
- [ ] Use environment variables for sensitive data
- [ ] Monitor API usage and costs
- [ ] Keep the tool updated

### For Developers

- [ ] Validate all inputs
- [ ] Use secure file operations
- [ ] Implement proper error handling
- [ ] Follow secure coding practices
- [ ] Regular security reviews
- [ ] Keep dependencies updated

### For Administrators

- [ ] Implement proper access controls
- [ ] Use enterprise token management
- [ ] Monitor for unusual activity
- [ ] Educate users on security practices
- [ ] Plan incident response procedures

## 🚨 Incident Response

### Security Incident Types

1. **Configuration Exposure**: API tokens or configurations exposed
2. **Unauthorized Access**: Unauthorized access to Claude Code accounts
3. **Malicious Code**: Security vulnerabilities in the tool
4. **Data Breach**: Any compromise of user data

### Response Steps

1. **Identification**: Recognize and confirm the incident
2. **Containment**: Limit the impact and prevent further damage
3. **Investigation**: Determine the cause and scope
4. **Resolution**: Fix the issue and prevent recurrence
5. **Communication**: Notify affected users appropriately

### Contact Information

- **Security Issues**: security@claude-provider.com
- **General Support**: support@claude-provider.com
- **GitHub Issues**: For non-security bugs and features

---

## 🔗 Additional Resources

- [Claude Code Security Documentation](https://docs.anthropic.com/claude/docs/security)
- [OWASP Go Secure Coding Practices](https://owasp.org/www-project-go-secure-coding-practices/)
- [API Security Best Practices](https://owasp.org/www-project-api-security/)

---

*This security policy is regularly updated. Last reviewed: January 15, 2024*

For questions about security practices or to report security issues, please contact us at security@claude-provider.com.