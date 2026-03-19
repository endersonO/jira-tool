package format

import (
	"fmt"
	"os"
	"strings"

	"github.com/endersonO/jt/internal/api"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var (
	bold  = color.New(color.Bold)
	cyan  = color.New(color.FgCyan, color.Bold)
	green = color.New(color.FgGreen)
	gray  = color.New(color.FgHiBlack)
)

// IssueTable prints a list of issues as a table
func IssueTable(issues []api.Issue, verbose bool) {
	if len(issues) == 0 {
		fmt.Println("No issues found.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)

	var rows [][]string
	if verbose {
		table.Header([]string{"KEY", "TYPE", "STATUS", "PRIORITY", "ASSIGNEE", "SUMMARY"})
		for _, issue := range issues {
			assignee := "—"
			if issue.Fields.Assignee != nil {
				assignee = issue.Fields.Assignee.DisplayName
			}
			priority := "—"
			if issue.Fields.Priority != nil {
				priority = issue.Fields.Priority.Name
			}
			rows = append(rows, []string{
				issue.Key,
				issue.Fields.IssueType.Name,
				issue.Fields.Status.Name,
				priority,
				assignee,
				truncate(issue.Fields.Summary, 55),
			})
		}
	} else {
		table.Header([]string{"KEY", "STATUS", "ASSIGNEE", "SUMMARY"})
		for _, issue := range issues {
			assignee := "—"
			if issue.Fields.Assignee != nil {
				assignee = issue.Fields.Assignee.DisplayName
			}
			rows = append(rows, []string{
				issue.Key,
				issue.Fields.Status.Name,
				assignee,
				truncate(issue.Fields.Summary, 60),
			})
		}
	}

	table.Bulk(rows)
	table.Render()
}

// IssueDetail prints full details of a single issue
func IssueDetail(issue *api.Issue) {
	fmt.Println()
	cyan.Printf("  %s", issue.Key)
	fmt.Printf("  %s\n", issue.Fields.Summary)
	fmt.Println(strings.Repeat("─", 70))

	printField("Type", issue.Fields.IssueType.Name)
	printField("Status", issue.Fields.Status.Name)
	if issue.Fields.Priority != nil {
		printField("Priority", issue.Fields.Priority.Name)
	}
	if issue.Fields.Assignee != nil {
		printField("Assignee", issue.Fields.Assignee.DisplayName+" <"+issue.Fields.Assignee.EmailAddress+">")
	}
	if issue.Fields.Reporter != nil {
		printField("Reporter", issue.Fields.Reporter.DisplayName)
	}
	if len(issue.Fields.Labels) > 0 {
		printField("Labels", strings.Join(issue.Fields.Labels, ", "))
	}
	if issue.Fields.Parent != nil {
		printField("Parent", issue.Fields.Parent.Key+" — "+issue.Fields.Parent.Fields.Summary)
	}
	printField("Created", issue.Fields.Created)
	printField("Updated", issue.Fields.Updated)

	fmt.Println()
	bold.Println("  Description")
	fmt.Println(strings.Repeat("─", 70))
	desc := extractDescription(issue.Fields.Description)
	if desc == "" {
		gray.Println("  (no description)")
	} else {
		fmt.Println(indent(desc, "  "))
	}
	fmt.Println()
}

// TransitionTable prints available transitions
func TransitionTable(transitions []api.Transition) {
	fmt.Println()
	bold.Println("Available transitions:")
	for _, t := range transitions {
		fmt.Printf("  %s  %s\n", gray.Sprint(t.ID), t.Name)
	}
	fmt.Println()
}

// ProjectTable prints a list of projects
func ProjectTable(projects []api.Project) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"KEY", "NAME", "TYPE"})

	var rows [][]string
	for _, p := range projects {
		rows = append(rows, []string{p.Key, p.Name, p.ProjectTypeKey})
	}
	table.Bulk(rows)
	table.Render()
}

// Success prints a success message
func Success(msg string) {
	green.Println("✓ " + msg)
}

// Error prints an error message
func Error(msg string) {
	color.New(color.FgRed).Fprintln(os.Stderr, "✗ "+msg)
}

func printField(label, value string) {
	gray.Printf("  %-12s", label)
	fmt.Printf("%s\n", value)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}

// extractDescription converts ADF description to readable plain text
func extractDescription(desc interface{}) string {
	if desc == nil {
		return ""
	}
	switch v := desc.(type) {
	case string:
		return v
	case map[string]interface{}:
		return renderADFNode(v)
	}
	return fmt.Sprintf("%v", desc)
}

func renderADFNode(node map[string]interface{}) string {
	nodeType, _ := node["type"].(string)
	var sb strings.Builder

	switch nodeType {
	case "doc":
		for _, child := range getContent(node) {
			sb.WriteString(renderADFNode(child))
		}
	case "paragraph":
		for _, child := range getContent(node) {
			sb.WriteString(renderADFNode(child))
		}
		sb.WriteString("\n")
	case "heading":
		level := 2
		if attrs, ok := node["attrs"].(map[string]interface{}); ok {
			if l, ok := attrs["level"].(float64); ok {
				level = int(l)
			}
		}
		prefix := strings.Repeat("#", level) + " "
		for _, child := range getContent(node) {
			sb.WriteString(prefix + renderADFNode(child))
		}
		sb.WriteString("\n")
	case "text":
		t, _ := node["text"].(string)
		marks, _ := node["marks"].([]interface{})
		for _, m := range marks {
			if mmap, ok := m.(map[string]interface{}); ok {
				switch mmap["type"] {
				case "strong":
					t = "**" + t + "**"
				case "em":
					t = "*" + t + "*"
				case "code":
					t = "`" + t + "`"
				}
			}
		}
		sb.WriteString(t)
	case "bulletList":
		for _, child := range getContent(node) {
			sb.WriteString("• " + renderListItem(child))
		}
	case "orderedList":
		for i, child := range getContent(node) {
			sb.WriteString(fmt.Sprintf("%d. ", i+1) + renderListItem(child))
		}
	case "codeBlock":
		lang := ""
		if attrs, ok := node["attrs"].(map[string]interface{}); ok {
			lang, _ = attrs["language"].(string)
		}
		sb.WriteString("```" + lang + "\n")
		for _, child := range getContent(node) {
			sb.WriteString(renderADFNode(child))
		}
		sb.WriteString("\n```\n")
	case "rule":
		sb.WriteString("---\n")
	}

	return sb.String()
}

func renderListItem(node map[string]interface{}) string {
	var sb strings.Builder
	for _, child := range getContent(node) {
		sb.WriteString(renderADFNode(child))
	}
	return strings.TrimRight(sb.String(), "\n") + "\n"
}

func getContent(node map[string]interface{}) []map[string]interface{} {
	raw, ok := node["content"].([]interface{})
	if !ok {
		return nil
	}
	var result []map[string]interface{}
	for _, item := range raw {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result
}
