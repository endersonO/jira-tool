package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/endersonO/jt/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Configure jt interactively (server, credentials, default project)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit()
		},
	}
}

func runInit() error {
	reader := bufio.NewReader(os.Stdin)

	bold := color.New(color.Bold)
	bold.Println("Welcome to jt — Jira Tool")
	fmt.Println("Let's configure your Jira credentials.\n")

	// Pre-fill from existing config if present
	existing, _ := config.Load()

	server := prompt(reader, "Jira server URL", existingOrDefault(existing, func(c *config.Config) string { return c.Server }, "https://your-org.atlassian.net"))
	email := prompt(reader, "Email", existingOrDefault(existing, func(c *config.Config) string { return c.Email }, ""))

	tokenHint := ""
	if existing != nil && existing.Token != "" {
		tokenHint = " [keep existing — press Enter to skip]"
	}
	token := promptSecret("API Token (https://id.atlassian.com/manage-profile/security/api-tokens)" + tokenHint)
	if token == "" && existing != nil && existing.Token != "" {
		token = existing.Token
	}

	project := prompt(reader, "Default project key (optional)", existingOrDefault(existing, func(c *config.Config) string { return c.Project }, ""))

	if server == "" || email == "" || token == "" {
		return fmt.Errorf("server, email, and token are required")
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("finding config directory: %w", err)
	}
	dir := filepath.Join(configDir, "jt")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	path := filepath.Join(dir, "config.yml")
	content := fmt.Sprintf("server: %s\nemail: %s\ntoken: %s\n", server, email, token)
	if project != "" {
		content += fmt.Sprintf("project: %s\n", project)
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	fmt.Println()
	color.Green("✓ Configuration saved to %s", path)
	fmt.Println("\nYou're all set! Try: jt issue list")
	return nil
}

func prompt(r *bufio.Reader, label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}
	val, _ := r.ReadString('\n')
	val = strings.TrimSpace(val)
	if val == "" {
		return defaultVal
	}
	return val
}

func promptSecret(label string) string {
	fmt.Printf("%s: ", label)
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		// Fallback if not a real terminal (e.g. piped input)
		reader := bufio.NewReader(os.Stdin)
		val, _ := reader.ReadString('\n')
		return strings.TrimSpace(val)
	}
	val := strings.TrimSpace(string(b))
	if val == "" {
		return val
	}
	// Keep existing token if user just pressed Enter (can't show default for secrets)
	return val
}

func existingOrDefault(cfg *config.Config, fn func(*config.Config) string, def string) string {
	if cfg == nil {
		return def
	}
	if v := fn(cfg); v != "" {
		return v
	}
	return def
}
