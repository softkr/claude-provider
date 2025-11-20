---
name: Bug Report
about: Create a report to help us improve
title: "[BUG] "
labels: bug
assignees: ''

---

## 🐛 Bug Description

A clear and concise description of what the bug is.

## 🔄 Reproduction Steps

To reproduce this behavior:

1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

## ✅ Expected Behavior

A clear and concise description of what you expected to happen.

## ❌ Actual Behavior

A clear and concise description of what actually happened.

## 📋 Screenshots

If applicable, add screenshots to help explain your problem.

## 🖥️ Environment Information

- **Operating System**: [e.g. macOS 14.0, Ubuntu 22.04, Windows 11]
- **Go Version**: [e.g. 1.21.0]
- **Claude Code Version**: [e.g. latest]
- **Tool Version**: [e.g. 1.0.0]

```bash
# Run this command and paste the output
claude-switch --status
```

## 📄 Configuration

Please provide your current configuration (remove sensitive information):

```json
{
  "env": {
    "ANTHROPIC_BASE_URL": "...",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "..."
  }
}
```

## 📝 Additional Context

Add any other context about the problem here.

## 🔍 Debug Information

```bash
# Please run these commands and paste the output

# Tool version
claude-switch --help | head -n 1

# Configuration file location and permissions
ls -la ~/.claude/

# Current configuration status
claude-switch --status

# If there are errors, include full error message and stack trace
```

## ✅ Checklist

- [ ] I have searched for existing issues that match my description
- [ ] I have provided sufficient information to reproduce the issue
- [ ] I have removed sensitive information from my report
- [ ] I have checked that this is not a security vulnerability (security issues should be reported privately)