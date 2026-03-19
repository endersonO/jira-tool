package adf

import (
	"strings"
)

// FromMarkdown converts a Markdown string to an ADF document.
// Supports: headings (#), bold (**), italic (*), inline code (`),
// code blocks (```), bullet lists (-/*), ordered lists (1.), horizontal rules (---)
func FromMarkdown(md string) *Doc {
	lines := strings.Split(strings.ReplaceAll(md, "\r\n", "\n"), "\n")
	nodes := parseLines(lines)
	return NewDoc(nodes)
}

func parseLines(lines []string) []Node {
	var nodes []Node
	i := 0
	for i < len(lines) {
		line := lines[i]

		// Code block
		if strings.HasPrefix(line, "```") {
			lang := strings.TrimPrefix(line, "```")
			lang = strings.TrimSpace(lang)
			var codeLines []string
			i++
			for i < len(lines) && !strings.HasPrefix(lines[i], "```") {
				codeLines = append(codeLines, lines[i])
				i++
			}
			nodes = append(nodes, codeBlock(lang, strings.Join(codeLines, "\n")))
			i++ // skip closing ```
			continue
		}

		// Horizontal rule
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" || trimmed == "***" || trimmed == "___" {
			nodes = append(nodes, rule())
			i++
			continue
		}

		// Headings
		if strings.HasPrefix(line, "# ") {
			nodes = append(nodes, heading(1, strings.TrimPrefix(line, "# ")))
			i++
			continue
		}
		if strings.HasPrefix(line, "## ") {
			nodes = append(nodes, heading(2, strings.TrimPrefix(line, "## ")))
			i++
			continue
		}
		if strings.HasPrefix(line, "### ") {
			nodes = append(nodes, heading(3, strings.TrimPrefix(line, "### ")))
			i++
			continue
		}
		if strings.HasPrefix(line, "#### ") {
			nodes = append(nodes, heading(4, strings.TrimPrefix(line, "#### ")))
			i++
			continue
		}

		// Bullet list
		if isBullet(line) {
			var items []Node
			for i < len(lines) && isBullet(lines[i]) {
				content := bulletContent(lines[i])
				items = append(items, listItem(paragraph(parseInline(content)...)))
				i++
			}
			nodes = append(nodes, bulletList(items))
			continue
		}

		// Ordered list
		if isOrdered(line) {
			var items []Node
			for i < len(lines) && isOrdered(lines[i]) {
				content := orderedContent(lines[i])
				items = append(items, listItem(paragraph(parseInline(content)...)))
				i++
			}
			nodes = append(nodes, orderedList(items))
			continue
		}

		// Empty line → skip (don't emit empty paragraphs)
		if trimmed == "" {
			i++
			continue
		}

		// Regular paragraph
		nodes = append(nodes, paragraph(parseInline(line)...))
		i++
	}
	return nodes
}

func isBullet(line string) bool {
	t := strings.TrimSpace(line)
	return strings.HasPrefix(t, "- ") || strings.HasPrefix(t, "* ")
}

func isOrdered(line string) bool {
	t := strings.TrimSpace(line)
	for j := 0; j < len(t) && j < 4; j++ {
		if t[j] >= '0' && t[j] <= '9' {
			continue
		}
		if t[j] == '.' && j > 0 {
			return true
		}
		break
	}
	return false
}

func bulletContent(line string) string {
	t := strings.TrimSpace(line)
	if strings.HasPrefix(t, "- ") {
		return strings.TrimPrefix(t, "- ")
	}
	return strings.TrimPrefix(t, "* ")
}

func orderedContent(line string) string {
	t := strings.TrimSpace(line)
	idx := strings.Index(t, ". ")
	if idx > 0 {
		return t[idx+2:]
	}
	return t
}

// parseInline handles bold, italic, inline code inside a line of text.
func parseInline(s string) []Node {
	var nodes []Node
	for len(s) > 0 {
		// Inline code: `...`
		if idx := strings.Index(s, "`"); idx >= 0 {
			if idx > 0 {
				nodes = append(nodes, parseStrongEm(s[:idx])...)
			}
			rest := s[idx+1:]
			end := strings.Index(rest, "`")
			if end >= 0 {
				nodes = append(nodes, text(rest[:end], Mark{Type: "code"}))
				s = rest[end+1:]
			} else {
				nodes = append(nodes, text("`"))
				s = rest
			}
			continue
		}

		// Bold: **...**
		if idx := strings.Index(s, "**"); idx >= 0 {
			if idx > 0 {
				nodes = append(nodes, parseEm(s[:idx])...)
			}
			rest := s[idx+2:]
			end := strings.Index(rest, "**")
			if end >= 0 {
				nodes = append(nodes, text(rest[:end], Mark{Type: "strong"}))
				s = rest[end+2:]
			} else {
				nodes = append(nodes, text("**"))
				s = rest
			}
			continue
		}

		nodes = append(nodes, parseEm(s)...)
		break
	}
	return nodes
}

func parseStrongEm(s string) []Node {
	return parseEm(s)
}

// parseEm handles *italic*
func parseEm(s string) []Node {
	var nodes []Node
	for len(s) > 0 {
		if idx := strings.Index(s, "*"); idx >= 0 {
			if idx > 0 {
				nodes = append(nodes, text(s[:idx]))
			}
			rest := s[idx+1:]
			end := strings.Index(rest, "*")
			if end >= 0 {
				nodes = append(nodes, text(rest[:end], Mark{Type: "em"}))
				s = rest[end+1:]
			} else {
				nodes = append(nodes, text("*"))
				s = rest
			}
			continue
		}
		nodes = append(nodes, text(s))
		break
	}
	return nodes
}
