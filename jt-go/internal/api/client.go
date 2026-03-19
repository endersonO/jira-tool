package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/endersonO/jt/internal/config"
)

type Client struct {
	cfg        *config.Config
	httpClient *http.Client
}

func New(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{},
	}
}

func (c *Client) do(method, path string, body interface{}) ([]byte, int, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.cfg.Server+path, bodyReader)
	if err != nil {
		return nil, 0, err
	}

	req.SetBasicAuth(c.cfg.Email, c.cfg.Token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, fmt.Errorf("API error %d: %s", resp.StatusCode, string(data))
	}

	return data, resp.StatusCode, nil
}

var defaultFields = []string{
	"summary", "status", "assignee", "reporter",
	"priority", "issuetype", "labels", "parent",
	"project", "created", "updated", "description",
}

// SearchIssues runs a JQL query
func (c *Client) SearchIssues(jql string, maxResults int, fields []string) (*SearchResult, error) {
	if len(fields) == 0 {
		fields = defaultFields
	}

	params := url.Values{}
	params.Set("jql", jql)
	params.Set("maxResults", fmt.Sprintf("%d", maxResults))
	for _, f := range fields {
		params.Add("fields", f)
	}

	data, _, err := c.do("GET", "/rest/api/3/search/jql?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIssue fetches a single issue
func (c *Client) GetIssue(key string) (*Issue, error) {
	data, _, err := c.do("GET", "/rest/api/3/issue/"+key, nil)
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(payload CreateIssuePayload) (*CreateIssueResponse, error) {
	data, _, err := c.do("POST", "/rest/api/3/issue", payload)
	if err != nil {
		return nil, err
	}

	var resp CreateIssueResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateIssue updates an existing issue (204 No Content on success)
func (c *Client) UpdateIssue(key string, payload UpdateIssuePayload) error {
	_, _, err := c.do("PUT", "/rest/api/3/issue/"+key, payload)
	return err
}

// GetTransitions returns available transitions for an issue
func (c *Client) GetTransitions(key string) ([]Transition, error) {
	data, _, err := c.do("GET", "/rest/api/3/issue/"+key+"/transitions", nil)
	if err != nil {
		return nil, err
	}

	var result TransitionsResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result.Transitions, nil
}

// TransitionIssue moves an issue to a new status by transition ID
func (c *Client) TransitionIssue(key, transitionID string) error {
	payload := TransitionPayload{Transition: TransitionRef{ID: transitionID}}
	_, _, err := c.do("POST", "/rest/api/3/issue/"+key+"/transitions", payload)
	return err
}

// ListProjects returns all accessible projects
func (c *Client) ListProjects() ([]Project, error) {
	data, _, err := c.do("GET", "/rest/api/3/project", nil)
	if err != nil {
		return nil, err
	}

	var projects []Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}
