package adf

// Doc is the root ADF document node
type Doc struct {
	Type    string `json:"type"`
	Version int    `json:"version"`
	Content []Node `json:"content"`
}

// Node represents any ADF node
type Node struct {
	Type    string                 `json:"type"`
	Attrs   map[string]interface{} `json:"attrs,omitempty"`
	Content []Node                 `json:"content,omitempty"`
	Text    string                 `json:"text,omitempty"`
	Marks   []Mark                 `json:"marks,omitempty"`
}

// Mark represents inline formatting (bold, italic, code, link)
type Mark struct {
	Type  string                 `json:"type"`
	Attrs map[string]interface{} `json:"attrs,omitempty"`
}

func NewDoc(nodes []Node) *Doc {
	return &Doc{
		Type:    "doc",
		Version: 1,
		Content: nodes,
	}
}

func paragraph(children ...Node) Node {
	return Node{Type: "paragraph", Content: children}
}

func text(s string, marks ...Mark) Node {
	n := Node{Type: "text", Text: s}
	if len(marks) > 0 {
		n.Marks = marks
	}
	return n
}

func heading(level int, s string) Node {
	return Node{
		Type:    "heading",
		Attrs:   map[string]interface{}{"level": level},
		Content: []Node{text(s)},
	}
}

func bulletList(items []Node) Node {
	return Node{Type: "bulletList", Content: items}
}

func orderedList(items []Node) Node {
	return Node{Type: "orderedList", Content: items}
}

func listItem(children ...Node) Node {
	return Node{Type: "listItem", Content: children}
}

func codeBlock(language, code string) Node {
	attrs := map[string]interface{}{}
	if language != "" {
		attrs["language"] = language
	}
	return Node{
		Type:    "codeBlock",
		Attrs:   attrs,
		Content: []Node{text(code)},
	}
}

func rule() Node {
	return Node{Type: "rule"}
}
