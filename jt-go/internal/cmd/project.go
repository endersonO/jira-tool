package cmd

import (
	"github.com/endersonO/jt/internal/api"
	"github.com/endersonO/jt/internal/format"
	"github.com/spf13/cobra"
)

func NewProjectCmd(load ConfigLoader) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"proj", "p"},
		Short:   "Manage Jira projects",
	}

	cmd.AddCommand(newProjectListCmd(load))
	return cmd
}

func newProjectListCmd(load ConfigLoader) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := load()
			if err != nil {
				return err
			}
			client := api.New(cfg)

			projects, err := client.ListProjects()
			if err != nil {
				return err
			}
			format.ProjectTable(projects)
			return nil
		},
	}
}
