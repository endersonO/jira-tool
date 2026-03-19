package cmd

import (
	"github.com/endersonO/jt/internal/api"
	"github.com/endersonO/jt/internal/format"
	"github.com/spf13/cobra"
)

func NewSearchCmd(load ConfigLoader) *cobra.Command {
	var (
		maxResult int
		verbose   bool
		jsonOut   bool
	)

	cmd := &cobra.Command{
		Use:   "search <JQL>",
		Short: "Search issues with a JQL query",
		Example: `  jt search "project=SCRUM AND assignee=currentUser()"
  jt search "project=SCRUM AND status=\"In Progress\""
  jt search "project=SCRUM AND issuetype=Bug AND priority=High"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := load()
			if err != nil {
				return err
			}
			client := api.New(cfg)

			result, err := client.SearchIssues(args[0], maxResult, nil)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(result.Issues)
			}
			format.IssueTable(result.Issues, verbose)
			return nil
		},
	}

	cmd.Flags().IntVar(&maxResult, "max", 30, "Max results")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show more columns")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output raw JSON")

	return cmd
}
