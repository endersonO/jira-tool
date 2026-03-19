package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/endersonO/jt/internal/adf"
	"github.com/endersonO/jt/internal/api"
	"github.com/endersonO/jt/internal/config"
	"github.com/endersonO/jt/internal/format"
	"github.com/spf13/cobra"
)

type ConfigLoader func() (*config.Config, error)

func NewIssueCmd(load ConfigLoader) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issue",
		Aliases: []string{"i"},
		Short:   "Manage Jira issues",
	}

	cmd.AddCommand(
		newIssueListCmd(load),
		newIssueViewCmd(load),
		newIssueCreateCmd(load),
		newIssueEditCmd(load),
		newIssueTransitionCmd(load),
	)

	return cmd
}

// --- list ---

func newIssueListCmd(load ConfigLoader) *cobra.Command {
	var (
		status    string
		assignee  string
		issueType string
		project   string
		maxResult int
		verbose   bool
		jsonOut   bool
	)

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := load()
			if err != nil {
				return err
			}
			client := api.New(cfg)

			proj := project
			if proj == "" {
				proj = cfg.Project
			}
			if proj == "" {
				return fmt.Errorf("project is required (set in config or use --project)")
			}

			parts := []string{fmt.Sprintf("project=%s", proj)}
			if status != "" {
				parts = append(parts, fmt.Sprintf("status=%q", status))
			}
			if assignee == "me" {
				parts = append(parts, "assignee=currentUser()")
			} else if assignee == "unassigned" {
				parts = append(parts, "assignee is EMPTY")
			} else if assignee != "" {
				parts = append(parts, fmt.Sprintf("assignee=%q", assignee))
			}
			if issueType != "" {
				parts = append(parts, fmt.Sprintf("issuetype=%q", issueType))
			}

			jql := strings.Join(parts, " AND ") + " ORDER BY updated DESC"
			result, err := client.SearchIssues(jql, maxResult, nil)
			if err != nil {
				return err
			}

			if jsonOut {
				return printJSON(result.Issues)
			}

			format.IssueTable(result.Issues, verbose)
			if !result.IsLast {
				fmt.Printf("\n(showing %d results, more available)\n", len(result.Issues))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "Filter by status (e.g. \"In Progress\")")
	cmd.Flags().StringVar(&assignee, "assignee", "", "Filter by assignee (email, \"me\", or \"unassigned\")")
	cmd.Flags().StringVar(&issueType, "type", "", "Filter by issue type (Task, Story, Bug...)")
	cmd.Flags().StringVar(&project, "project", "", "Project key (overrides config default)")
	cmd.Flags().IntVar(&maxResult, "max", 30, "Max results")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show more columns")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output raw JSON")

	return cmd
}

// --- view ---

func newIssueViewCmd(load ConfigLoader) *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:     "view <KEY>",
		Aliases: []string{"show", "get"},
		Short:   "View issue details",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := load()
			if err != nil {
				return err
			}
			client := api.New(cfg)

			issue, err := client.GetIssue(strings.ToUpper(args[0]))
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(issue)
			}
			format.IssueDetail(issue)
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output raw JSON")
	return cmd
}

// --- create ---

func newIssueCreateCmd(load ConfigLoader) *cobra.Command {
	var (
		summary     string
		issueType   string
		description string
		descFile    string
		assignee    string
		priority    string
		labels      []string
		parent      string
		project     string
		edit        bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := load()
			if err != nil {
				return err
			}
			client := api.New(cfg)

			if summary == "" {
				return fmt.Errorf("--summary is required")
			}

			proj := project
			if proj == "" {
				proj = cfg.Project
			}
			if proj == "" {
				return fmt.Errorf("project is required (set in config or use --project)")
			}

			descMD, err := resolveDescription(description, descFile, edit)
			if err != nil {
				return err
			}

			payload := api.CreateIssuePayload{
				Fields: api.CreateIssueFields{
					Project:   api.ProjectRef{Key: proj},
					Summary:   summary,
					IssueType: api.IssueType{Name: issueType},
					Labels:    labels,
				},
			}

			if descMD != "" {
				payload.Fields.Description = adf.FromMarkdown(descMD)
			}
			if priority != "" {
				payload.Fields.Priority = &api.Priority{Name: priority}
			}
			if assignee != "" {
				email := assignee
				if assignee == "me" {
					email = cfg.Email
				}
				payload.Fields.Assignee = &api.UserRef{EmailAddress: email}
			}
			if parent != "" {
				payload.Fields.Parent = &api.IssueRef{Key: strings.ToUpper(parent)}
			}

			resp, err := client.CreateIssue(payload)
			if err != nil {
				return err
			}

			format.Success(fmt.Sprintf("Created %s", resp.Key))
			return nil
		},
	}

	cmd.Flags().StringVarP(&summary, "summary", "s", "", "Issue summary (required)")
	cmd.Flags().StringVar(&issueType, "type", "Task", "Issue type (Task, Story, Bug, Epic)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Description in Markdown")
	cmd.Flags().StringVar(&descFile, "description-file", "", "Read description from a Markdown file")
	cmd.Flags().StringVarP(&assignee, "assignee", "a", "", "Assignee email or \"me\"")
	cmd.Flags().StringVarP(&priority, "priority", "p", "", "Priority (Highest, High, Medium, Low, Lowest)")
	cmd.Flags().StringSliceVarP(&labels, "labels", "l", nil, "Labels (comma-separated)")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent issue key")
	cmd.Flags().StringVar(&project, "project", "", "Project key (overrides config)")
	cmd.Flags().BoolVarP(&edit, "edit", "e", false, "Open $EDITOR to write description")

	return cmd
}

// --- edit ---

func newIssueEditCmd(load ConfigLoader) *cobra.Command {
	var (
		summary     string
		description string
		descFile    string
		assignee    string
		priority    string
		labels      []string
		status      string
		edit        bool
	)

	cmd := &cobra.Command{
		Use:   "edit <KEY>",
		Short: "Edit an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := load()
			if err != nil {
				return err
			}
			client := api.New(cfg)
			key := strings.ToUpper(args[0])

			descMD, err := resolveDescription(description, descFile, edit)
			if err != nil {
				return err
			}

			payload := api.UpdateIssuePayload{}

			if summary != "" {
				payload.Fields.Summary = summary
			}
			if descMD != "" {
				payload.Fields.Description = adf.FromMarkdown(descMD)
			}
			if priority != "" {
				payload.Fields.Priority = &api.Priority{Name: priority}
			}
			if assignee != "" {
				payload.Fields.Assignee = &api.UserRef{EmailAddress: assignee}
			}
			if len(labels) > 0 {
				payload.Fields.Labels = labels
			}

			if err := client.UpdateIssue(key, payload); err != nil {
				return err
			}

			// Status transition is a separate API call
			if status != "" {
				transitions, err := client.GetTransitions(key)
				if err != nil {
					return fmt.Errorf("could not get transitions: %w", err)
				}
				tid := findTransition(transitions, status)
				if tid == "" {
					available := make([]string, len(transitions))
					for i, t := range transitions {
						available[i] = t.Name
					}
					return fmt.Errorf("unknown status %q — available: %s", status, strings.Join(available, ", "))
				}
				if err := client.TransitionIssue(key, tid); err != nil {
					return fmt.Errorf("could not transition: %w", err)
				}
			}

			format.Success(fmt.Sprintf("Updated %s", key))
			return nil
		},
	}

	cmd.Flags().StringVarP(&summary, "summary", "s", "", "New summary")
	cmd.Flags().StringVarP(&description, "description", "d", "", "New description in Markdown")
	cmd.Flags().StringVar(&descFile, "description-file", "", "Read description from a Markdown file")
	cmd.Flags().StringVar(&status, "status", "", "New status (e.g. \"In Progress\")")
	cmd.Flags().StringVarP(&assignee, "assignee", "a", "", "New assignee email")
	cmd.Flags().StringVarP(&priority, "priority", "p", "", "New priority")
	cmd.Flags().StringSliceVarP(&labels, "labels", "l", nil, "New labels (comma-separated)")
	cmd.Flags().BoolVarP(&edit, "edit", "e", false, "Open $EDITOR to write description")

	return cmd
}

// --- transition ---

func newIssueTransitionCmd(load ConfigLoader) *cobra.Command {
	var list bool

	cmd := &cobra.Command{
		Use:     "transition <KEY> [STATUS]",
		Aliases: []string{"move", "status"},
		Short:   "Change issue status",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := load()
			if err != nil {
				return err
			}
			client := api.New(cfg)
			key := strings.ToUpper(args[0])

			transitions, err := client.GetTransitions(key)
			if err != nil {
				return err
			}

			if list || len(args) == 1 {
				format.TransitionTable(transitions)
				return nil
			}

			target := args[1]
			tid := findTransition(transitions, target)
			if tid == "" {
				available := make([]string, len(transitions))
				for i, t := range transitions {
					available[i] = t.Name
				}
				return fmt.Errorf("unknown status %q — available: %s", target, strings.Join(available, ", "))
			}

			if err := client.TransitionIssue(key, tid); err != nil {
				return err
			}

			format.Success(fmt.Sprintf("Moved %s → %s", key, target))
			return nil
		},
	}

	cmd.Flags().BoolVarP(&list, "list", "l", false, "List available transitions")
	return cmd
}

// --- helpers ---

func findTransition(transitions []api.Transition, name string) string {
	nameLower := strings.ToLower(name)
	for _, t := range transitions {
		if strings.ToLower(t.Name) == nameLower {
			return t.ID
		}
	}
	return ""
}

func resolveDescription(inline, file string, edit bool) (string, error) {
	if edit {
		return openEditor("")
	}
	if file != "" {
		data, err := os.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("reading description file: %w", err)
		}
		return string(data), nil
	}
	return inline, nil
}

func openEditor(initial string) (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	f, err := os.CreateTemp("", "jt-*.md")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())

	if initial != "" {
		f.WriteString(initial)
	}
	f.Close()

	c := exec.Command(editor, f.Name())
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return "", fmt.Errorf("editor exited with error: %w", err)
	}

	data, err := os.ReadFile(f.Name())
	if err != nil {
		return "", err
	}
	return string(data), nil
}
