package api

// --- Search / List ---

type SearchResult struct {
	Issues        []Issue `json:"issues"`
	NextPageToken string  `json:"nextPageToken"`
	IsLast        bool    `json:"isLast"`
}

// --- Issue ---

type Issue struct {
	ID     string      `json:"id"`
	Key    string      `json:"key"`
	Self   string      `json:"self"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Summary     string      `json:"summary"`
	Description interface{} `json:"description"` // ADF or nil
	Status      Status      `json:"status"`
	Assignee    *User       `json:"assignee"`
	Reporter    *User       `json:"reporter"`
	Priority    *Priority   `json:"priority"`
	IssueType   IssueType   `json:"issuetype"`
	Labels      []string    `json:"labels"`
	Parent      *IssueRef   `json:"parent"`
	Project     ProjectRef  `json:"project"`
	Created     string      `json:"created"`
	Updated     string      `json:"updated"`
}

type Status struct {
	Name string `json:"name"`
}

type User struct {
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

type Priority struct {
	Name string `json:"name"`
}

type IssueType struct {
	Name string `json:"name"`
}

type IssueRef struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
	} `json:"fields"`
}

type ProjectRef struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// --- Transitions ---

type TransitionsResult struct {
	Transitions []Transition `json:"transitions"`
}

type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// --- Projects ---

type Project struct {
	Key            string `json:"key"`
	Name           string `json:"name"`
	ProjectTypeKey string `json:"projectTypeKey"`
}

// --- Create / Update payloads ---

type CreateIssuePayload struct {
	Fields CreateIssueFields `json:"fields"`
}

type CreateIssueFields struct {
	Project     ProjectRef  `json:"project"`
	Summary     string      `json:"summary"`
	IssueType   IssueType   `json:"issuetype"`
	Description interface{} `json:"description,omitempty"`
	Priority    *Priority   `json:"priority,omitempty"`
	Assignee    *UserRef    `json:"assignee,omitempty"`
	Labels      []string    `json:"labels,omitempty"`
	Parent      *IssueRef   `json:"parent,omitempty"`
}

type UpdateIssuePayload struct {
	Fields UpdateIssueFields `json:"fields"`
}

type UpdateIssueFields struct {
	Summary     string      `json:"summary,omitempty"`
	Description interface{} `json:"description,omitempty"`
	Priority    *Priority   `json:"priority,omitempty"`
	Assignee    *UserRef    `json:"assignee,omitempty"`
	Labels      []string    `json:"labels,omitempty"`
}

type UserRef struct {
	EmailAddress string `json:"emailAddress"`
}

type TransitionPayload struct {
	Transition TransitionRef `json:"transition"`
}

type TransitionRef struct {
	ID string `json:"id"`
}

type CreateIssueResponse struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}
