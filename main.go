package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

const (
	Version = "2.0.0"
)

// Config represents the Claude configuration structure
type Config struct {
	Env map[string]string `json:"env"`
}

// Application holds the application state
type Application struct {
	settingsFile string
	backupFile   string
	configDir    string
	green        *color.Color
	yellow       *color.Color
	cyan         *color.Color
	red          *color.Color
}

// Z.AI specific environment keys (excluding ANTHROPIC_AUTH_TOKEN which is shared)
var zaiEnvKeys = []string{
	"ANTHROPIC_BASE_URL",
	"API_TIMEOUT_MS",
	"ANTHROPIC_DEFAULT_OPUS_MODEL",
	"ANTHROPIC_DEFAULT_SONNET_MODEL",
	"ANTHROPIC_DEFAULT_HAIKU_MODEL",
}

// NewApplication creates a new application instance
func NewApplication() *Application {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	configDir := filepath.Join(homeDir, ".claude")

	return &Application{
		settingsFile: filepath.Join(configDir, "settings.json"),
		backupFile:   filepath.Join(configDir, "settings.json.backup"),
		configDir:    configDir,
		green:        color.New(color.FgGreen),
		yellow:       color.New(color.FgYellow),
		cyan:         color.New(color.FgCyan),
		red:          color.New(color.FgRed),
	}
}

// printHeader prints the application header
func (app *Application) printHeader() {
	app.cyan.Printf("🤖 Claude Code API Switcher v%s\n", Version)
	fmt.Println()
}

// loadConfig loads configuration from file
func (app *Application) loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Env: make(map[string]string)}, nil
		}
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if config.Env == nil {
		config.Env = make(map[string]string)
	}

	return &config, nil
}

// saveConfigAtomic saves configuration to file atomically
func (app *Application) saveConfigAtomic(filename string, config *Config) error {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to temp file first
	tempFile := filename + ".tmp"
	err = os.WriteFile(tempFile, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write temp config: %w", err)
	}

	// Atomic rename
	err = os.Rename(tempFile, filename)
	if err != nil {
		os.Remove(tempFile) // Clean up temp file on error
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// promptForToken prompts user for API token
func (app *Application) promptForToken() (string, error) {
	// Check environment variable first
	if token := os.Getenv("ZAI_AUTH_TOKEN"); token != "" {
		app.cyan.Println("📌 Using token from ZAI_AUTH_TOKEN environment variable")
		return token, nil
	}

	// Check if token file exists
	tokenFile := filepath.Join(app.configDir, ".zai_token")
	if data, err := os.ReadFile(tokenFile); err == nil {
		token := strings.TrimSpace(string(data))
		if token != "" {
			app.cyan.Println("📌 Using token from saved token file")
			return token, nil
		}
	}

	// Prompt user for token
	app.yellow.Println("⚠️  No API token found")
	fmt.Println()
	app.cyan.Println("Please enter your Z.AI API token:")
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	token, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read token: %w", err)
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return "", fmt.Errorf("token cannot be empty")
	}

	// Ask if user wants to save the token
	app.cyan.Println("\nSave token for future use? (y/n)")
	fmt.Print("> ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "y" || answer == "yes" {
		err = os.WriteFile(tokenFile, []byte(token), 0600)
		if err != nil {
			app.yellow.Printf("⚠️  Failed to save token: %v\n", err)
		} else {
			app.green.Println("✅ Token saved successfully")
		}
	}

	return token, nil
}

// switchToAnthropic switches to Anthropic configuration
func (app *Application) switchToAnthropic() error {
	app.green.Println("🔄 Switching to Anthropic API...")

	// Check if backup exists (required for web login token)
	if _, err := os.Stat(app.backupFile); os.IsNotExist(err) {
		app.red.Println("❌ No backup found!")
		app.yellow.Println("⚠️  Cannot restore Anthropic web login token without backup.")
		app.yellow.Println("   You may need to re-login to Claude Code.")
		fmt.Println()

		// Create empty config without Z.AI keys
		config := &Config{Env: make(map[string]string)}
		err = app.saveConfigAtomic(app.settingsFile, config)
		if err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		app.yellow.Println("⚠️  Created empty configuration (re-login required)")
		return nil
	}

	// Restore from backup (contains web login token)
	backupConfig, err := app.loadConfig(app.backupFile)
	if err != nil {
		return fmt.Errorf("failed to load backup: %w", err)
	}

	// Remove any Z.AI specific keys that might be in backup
	for _, key := range zaiEnvKeys {
		delete(backupConfig.Env, key)
	}

	err = app.saveConfigAtomic(app.settingsFile, backupConfig)
	if err != nil {
		return fmt.Errorf("failed to restore config: %w", err)
	}

	app.green.Println("✅ Anthropic configuration restored from backup")
	app.cyan.Println("   Web login token has been restored")
	return nil
}

// switchToZAI switches to Z.AI configuration
func (app *Application) switchToZAI() error {
	app.green.Println("🔄 Switching to Z.AI API...")

	// Load current config
	config, err := app.loadConfig(app.settingsFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if already using Z.AI
	if strings.Contains(config.Env["ANTHROPIC_BASE_URL"], "z.ai") {
		app.yellow.Println("⚠️  Already using Z.AI configuration")
		app.cyan.Println("   Use --status to check current settings")
		return nil
	}

	// Backup current Anthropic configuration (CRITICAL: contains web login token)
	if _, err := os.Stat(app.settingsFile); err == nil && len(config.Env) > 0 {
		err = app.createBackupAtomic()
		if err != nil {
			return fmt.Errorf("failed to backup Anthropic config (web login token): %w", err)
		}
		app.green.Println("✅ Anthropic configuration backed up (web login token saved)")
	} else {
		app.yellow.Println("⚠️  No existing Anthropic configuration to backup")
		app.yellow.Println("   You may need to re-login when switching back")
	}

	// Get Z.AI API token
	token, err := app.promptForToken()
	if err != nil {
		return err
	}

	// Create new config for Z.AI (fresh start, don't mix with Anthropic settings)
	newConfig := &Config{
		Env: map[string]string{
			"ANTHROPIC_AUTH_TOKEN":           token,
			"ANTHROPIC_BASE_URL":             "https://api.z.ai/api/anthropic",
			"API_TIMEOUT_MS":                 "3000000",
			"ANTHROPIC_DEFAULT_OPUS_MODEL":   "GLM-4.6",
			"ANTHROPIC_DEFAULT_SONNET_MODEL": "GLM-4.6",
			"ANTHROPIC_DEFAULT_HAIKU_MODEL":  "GLM-4.5-Air",
		},
	}

	err = app.saveConfigAtomic(app.settingsFile, newConfig)
	if err != nil {
		return fmt.Errorf("failed to save Z.AI configuration: %w", err)
	}

	app.green.Println("✅ Z.AI configuration applied successfully")
	fmt.Println()
	app.cyan.Println("💡 To switch back to Anthropic: claude-switch --anthropic")
	return nil
}

// isZAIKey checks if a key is a Z.AI specific key
func isZAIKey(key string) bool {
	for _, zaiKey := range zaiEnvKeys {
		if key == zaiKey {
			return true
		}
	}
	return false
}

// createBackupAtomic creates a backup of the current configuration atomically
func (app *Application) createBackupAtomic() error {
	source, err := os.ReadFile(app.settingsFile)
	if err != nil {
		return err
	}

	tempFile := app.backupFile + ".tmp"
	err = os.WriteFile(tempFile, source, 0600)
	if err != nil {
		return err
	}

	err = os.Rename(tempFile, app.backupFile)
	if err != nil {
		os.Remove(tempFile)
		return err
	}

	return nil
}

// showStatus shows current configuration status
func (app *Application) showStatus() error {
	app.cyan.Println("📊 Current Configuration Status")
	fmt.Println()

	config, err := app.loadConfig(app.settingsFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(config.Env) == 0 {
		app.yellow.Println("⚠️  No configuration found (empty or missing)")
		return nil
	}

	baseURL := config.Env["ANTHROPIC_BASE_URL"]
	if strings.Contains(baseURL, "z.ai") {
		app.green.Println("┌─────────────────────────────────────┐")
		app.green.Println("│  🔗 Provider: Z.AI (GLM Models)     │")
		app.green.Println("└─────────────────────────────────────┘")
		fmt.Println()
		app.cyan.Printf("  Base URL: %s\n", baseURL)

		if model := config.Env["ANTHROPIC_DEFAULT_SONNET_MODEL"]; model != "" {
			app.cyan.Printf("  Sonnet Model: %s\n", model)
		}
		if model := config.Env["ANTHROPIC_DEFAULT_OPUS_MODEL"]; model != "" {
			app.cyan.Printf("  Opus Model: %s\n", model)
		}
		if model := config.Env["ANTHROPIC_DEFAULT_HAIKU_MODEL"]; model != "" {
			app.cyan.Printf("  Haiku Model: %s\n", model)
		}
		if timeout := config.Env["API_TIMEOUT_MS"]; timeout != "" {
			app.cyan.Printf("  Timeout: %s ms\n", timeout)
		}

		// Show masked token
		if token := config.Env["ANTHROPIC_AUTH_TOKEN"]; token != "" {
			maskedToken := maskToken(token)
			app.cyan.Printf("  Auth Token: %s\n", maskedToken)
		}
	} else if baseURL == "" {
		app.green.Println("┌─────────────────────────────────────┐")
		app.green.Println("│  🔗 Provider: Anthropic (Default)   │")
		app.green.Println("└─────────────────────────────────────┘")
		fmt.Println()
		app.cyan.Println("  Base URL: api.anthropic.com (default)")
	} else {
		app.green.Println("┌─────────────────────────────────────┐")
		app.green.Println("│  🔗 Provider: Custom                │")
		app.green.Println("└─────────────────────────────────────┘")
		fmt.Println()
		app.cyan.Printf("  Base URL: %s\n", baseURL)
	}

	fmt.Println()

	// Show other environment variables
	otherEnvCount := 0
	for key := range config.Env {
		if !isZAIKey(key) && key != "ANTHROPIC_BASE_URL" {
			otherEnvCount++
		}
	}
	if otherEnvCount > 0 {
		app.cyan.Printf("  Other env vars: %d\n", otherEnvCount)
	}

	// Check for backup
	if _, err := os.Stat(app.backupFile); err == nil {
		app.cyan.Println("  💾 Backup: Available")
	} else {
		app.yellow.Println("  💾 Backup: Not found")
	}

	// Check for saved token
	tokenFile := filepath.Join(app.configDir, ".zai_token")
	if _, err := os.Stat(tokenFile); err == nil {
		app.cyan.Println("  🔑 Saved Token: Available")
	}

	return nil
}

// maskToken masks an API token for display
func maskToken(token string) string {
	if len(token) <= 8 {
		return "********"
	}
	return token[:4] + "..." + token[len(token)-4:]
}

// clearToken removes the saved token
func (app *Application) clearToken() error {
	tokenFile := filepath.Join(app.configDir, ".zai_token")

	if _, err := os.Stat(tokenFile); os.IsNotExist(err) {
		app.yellow.Println("⚠️  No saved token found")
		return nil
	}

	err := os.Remove(tokenFile)
	if err != nil {
		return fmt.Errorf("failed to remove token: %w", err)
	}

	app.green.Println("✅ Saved token removed successfully")
	return nil
}

// install installs the application
func (app *Application) install() error {
	app.green.Println("🚀 Installing Claude Code API Switcher...")
	fmt.Println()

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve symlinks
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	// Install binary to /usr/local/bin
	installPath := "/usr/local/bin/claude-switch"
	if execPath != installPath {
		app.cyan.Println("📦 Installing binary to /usr/local/bin...")

		// Read source binary
		sourceData, err := os.ReadFile(execPath)
		if err != nil {
			return fmt.Errorf("failed to read source binary: %w", err)
		}

		// Write to destination (requires sudo, so we'll use a temp approach)
		tempFile := "/tmp/claude-switch-install"
		err = os.WriteFile(tempFile, sourceData, 0755)
		if err != nil {
			return fmt.Errorf("failed to write temp file: %w", err)
		}

		// Try direct copy first, then sudo if needed
		err = os.Rename(tempFile, installPath)
		if err != nil {
			// Need sudo - execute cp and chmod
			app.yellow.Println("⚠️  Need sudo permission to install to /usr/local/bin")

			cmd := fmt.Sprintf("sudo cp %s %s && sudo chmod +x %s", tempFile, installPath, installPath)
			fmt.Printf("Running: %s\n", cmd)

			// Use os/exec to run sudo command
			import_cmd := exec.Command("bash", "-c", cmd)
			import_cmd.Stdin = os.Stdin
			import_cmd.Stdout = os.Stdout
			import_cmd.Stderr = os.Stderr

			if err := import_cmd.Run(); err != nil {
				os.Remove(tempFile)
				return fmt.Errorf("failed to install binary (try running with sudo): %w", err)
			}
			os.Remove(tempFile)
		}

		app.green.Println("✅ Binary installed to /usr/local/bin/claude-switch")
		execPath = installPath
	} else {
		app.cyan.Println("📦 Binary already installed at /usr/local/bin/claude-switch")
	}

	// Determine shell configuration files
	shellConfigs := app.detectShellConfigs()
	if len(shellConfigs) == 0 {
		return fmt.Errorf("no supported shell configuration found")
	}

	// Create alias block
	aliasBlock := fmt.Sprintf(`
# Claude Code API Switcher
alias claude-switch='%s'
alias claude-anthropic='%s --anthropic'
alias claude-zai='%s --zai'
alias claude-status='%s --status'
`, execPath, execPath, execPath, execPath)

	// Fish shell uses different syntax
	fishAliasBlock := fmt.Sprintf(`
# Claude Code API Switcher
alias claude-switch '%s'
alias claude-anthropic '%s --anthropic'
alias claude-zai '%s --zai'
alias claude-status '%s --status'
`, execPath, execPath, execPath, execPath)

	installedCount := 0
	for _, shellRC := range shellConfigs {
		isFish := strings.Contains(shellRC, "fish")
		block := aliasBlock
		if isFish {
			block = fishAliasBlock
		}

		// Read existing shell config
		content, err := os.ReadFile(shellRC)
		if err != nil && !os.IsNotExist(err) {
			app.yellow.Printf("⚠️  Failed to read %s: %v\n", shellRC, err)
			continue
		}

		// Check if aliases already exist
		if strings.Contains(string(content), "Claude Code API Switcher") {
			app.yellow.Printf("⚠️  Aliases already exist in %s\n", shellRC)
			continue
		}

		// Append aliases
		f, err := os.OpenFile(shellRC, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			app.yellow.Printf("⚠️  Failed to open %s: %v\n", shellRC, err)
			continue
		}

		_, err = f.WriteString(block)
		f.Close()
		if err != nil {
			app.yellow.Printf("⚠️  Failed to write to %s: %v\n", shellRC, err)
			continue
		}

		app.green.Printf("✅ Aliases added to %s\n", shellRC)
		installedCount++
	}

	if installedCount == 0 {
		app.yellow.Println("⚠️  No new aliases were installed")
	}

	fmt.Println()
	app.green.Println("🎉 Installation complete!")
	fmt.Println()
	app.cyan.Println("Available commands after reload:")
	fmt.Println("  claude-switch --anthropic  # Use Anthropic Claude")
	fmt.Println("  claude-switch --zai        # Use Z.AI GLM")
	fmt.Println("  claude-switch --status     # Check current config")
	fmt.Println("  claude-anthropic           # Quick switch to Anthropic")
	fmt.Println("  claude-zai                 # Quick switch to Z.AI")
	fmt.Println("  claude-status              # Quick status check")
	fmt.Println()
	app.cyan.Println("Reload your shell:")
	for _, shellRC := range shellConfigs {
		fmt.Printf("  source %s\n", shellRC)
	}

	return nil
}

// detectShellConfigs detects available shell configuration files
func (app *Application) detectShellConfigs() []string {
	homeDir := os.Getenv("HOME")
	var configs []string

	// Check current shell
	shell := os.Getenv("SHELL")

	// Common shell config files
	candidates := []struct {
		path     string
		forShell string
	}{
		{filepath.Join(homeDir, ".zshrc"), "zsh"},
		{filepath.Join(homeDir, ".bashrc"), "bash"},
		{filepath.Join(homeDir, ".bash_profile"), "bash"},
		{filepath.Join(homeDir, ".config", "fish", "config.fish"), "fish"},
	}

	// Add config for current shell first
	for _, c := range candidates {
		if strings.Contains(shell, c.forShell) {
			if _, err := os.Stat(c.path); err == nil {
				configs = append(configs, c.path)
				break
			}
		}
	}

	// If no config found for current shell, check all
	if len(configs) == 0 {
		for _, c := range candidates {
			if _, err := os.Stat(c.path); err == nil {
				configs = append(configs, c.path)
				break
			}
		}
	}

	return configs
}

// printUsage prints usage information
func (app *Application) printUsage() {
	app.printHeader()
	app.cyan.Println("Usage:")
	fmt.Println()
	fmt.Println("  claude-switch [command]")
	fmt.Println()
	app.cyan.Println("Commands:")
	fmt.Println("  -a, --anthropic  Switch to Anthropic API (restore web login token)")
	fmt.Println("  -z, --zai        Switch to Z.AI API (use API key)")
	fmt.Println("  -s, --status     Show current configuration")
	fmt.Println("  --clear-token    Remove saved Z.AI API token")
	fmt.Println("  --install        Install aliases to shell")
	fmt.Println("  -v, --version    Show version")
	fmt.Println("  -h, --help       Show this help message")
	fmt.Println()
	app.cyan.Println("Authentication:")
	fmt.Println("  Anthropic  Uses web login token (automatically backed up)")
	fmt.Println("  Z.AI       Uses API key (prompted or from ZAI_AUTH_TOKEN env)")
	fmt.Println()
	app.cyan.Println("Environment Variables:")
	fmt.Println("  ZAI_AUTH_TOKEN   Z.AI API key (optional)")
	fmt.Println()
	app.cyan.Println("Examples:")
	fmt.Println("  claude-switch --zai        # Backup web token, switch to Z.AI")
	fmt.Println("  claude-switch --anthropic  # Restore web token from backup")
	fmt.Println("  claude-switch --status     # Check current provider")
	fmt.Println()
	app.yellow.Println("Note: Switching to Z.AI automatically backs up your Anthropic")
	fmt.Println("      web login token. Use --anthropic to restore it later.")
	fmt.Println()
}

func main() {
	var (
		anthropic  = flag.Bool("anthropic", false, "Switch to Anthropic API")
		a          = flag.Bool("a", false, "Switch to Anthropic API (short)")
		zai        = flag.Bool("zai", false, "Switch to Z.AI API")
		z          = flag.Bool("z", false, "Switch to Z.AI API (short)")
		status     = flag.Bool("status", false, "Show current configuration")
		s          = flag.Bool("s", false, "Show current configuration (short)")
		clearToken = flag.Bool("clear-token", false, "Remove saved Z.AI token")
		install    = flag.Bool("install", false, "Install aliases to shell")
		version    = flag.Bool("version", false, "Show version")
		v          = flag.Bool("v", false, "Show version")
		help       = flag.Bool("help", false, "Show help message")
		h          = flag.Bool("h", false, "Show help message")
	)

	flag.Parse()

	app := NewApplication()

	// Show help if no arguments or help flag
	if len(os.Args) == 1 || *help || *h {
		app.printUsage()
		return
	}

	// Show version
	if *version || *v {
		fmt.Printf("claude-switch v%s (%s/%s)\n", Version, runtime.GOOS, runtime.GOARCH)
		return
	}

	// Execute command
	switch {
	case *anthropic || *a:
		if err := app.switchToAnthropic(); err != nil {
			app.red.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case *zai || *z:
		if err := app.switchToZAI(); err != nil {
			app.red.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case *status || *s:
		if err := app.showStatus(); err != nil {
			app.red.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case *clearToken:
		if err := app.clearToken(); err != nil {
			app.red.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case *install:
		if err := app.install(); err != nil {
			app.red.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		app.printUsage()
		os.Exit(1)
	}
}
