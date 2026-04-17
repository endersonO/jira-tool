package main

import (
	"fmt"
	"os"

	"github.com/endersonO/jt/internal/cmd"
	"github.com/endersonO/jt/internal/config"
	"github.com/endersonO/jt/internal/format"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	versionStr := fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)

	root := &cobra.Command{
		Use:   "jt",
		Short: "Jira Tool — manage Jira from the terminal",
		Long: `jt is a CLI for Jira Cloud, inspired by gh.

Run 'jt init' to configure your credentials.
Config location varies by OS:
  macOS:   ~/Library/Application Support/jt/config.yml
  Linux:   ~/.config/jt/config.yml
  Windows: %APPDATA%\jt\config.yml`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       versionStr,
	}

	root.AddCommand(
		cmd.NewInitCmd(),
		cmd.NewIssueCmd(loadConfig),
		cmd.NewProjectCmd(loadConfig),
		cmd.NewSearchCmd(loadConfig),
		newConfigCmd(),
	)

	if err := root.Execute(); err != nil {
		format.Error(err.Error())
		os.Exit(1)
	}
}

// loadConfig is passed as a factory to commands so config is loaded lazily at run time
func loadConfig() (*config.Config, error) {
	return config.Load()
}

func newConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Show config file location and status",
		Run: func(cmd *cobra.Command, args []string) {
			path := config.ConfigPath()
			fmt.Printf("Config file: %s\n", path)
			cfg, err := config.Load()
			if err != nil {
				format.Error("Not configured: " + err.Error())
				fmt.Printf("\nCreate %s with:\n\n", path)
				fmt.Println("  server:  https://your-org.atlassian.net")
				fmt.Println("  email:   you@example.com")
				fmt.Println("  token:   your-api-token")
				fmt.Println("  project: SCRUM")
				return
			}
			fmt.Printf("Server:  %s\n", cfg.Server)
			fmt.Printf("Email:   %s\n", cfg.Email)
			fmt.Printf("Project: %s\n", cfg.Project)
			fmt.Println("Status:  OK")
		},
	}
}
