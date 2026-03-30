// Package tracker provides issue tracker adapters.
package tracker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dobbo-ca/lepton/internal/domain"
)

// ErrNotImplemented is returned by tracker stubs.
var ErrNotImplemented = errors.New("not implemented")

// IssueTracker is the interface for interacting with an external issue tracker.
type IssueTracker interface {
	ListIssues(ctx context.Context, project string) ([]domain.Issue, error)
	GetIssue(ctx context.Context, id string) (*domain.Issue, error)
	UpdateIssue(ctx context.Context, id string, update domain.IssueUpdate) error
	CreateComment(ctx context.Context, issueID string, body string) error
}

// Config holds tracker configuration loaded from app config.
type Config struct {
	Type    string // "jira" | "linear"
	BaseURL string // e.g. https://your-org.atlassian.net
	Email   string // Jira user email
	Token   string // API token
	Project string // default project key
}

// New returns the appropriate IssueTracker implementation based on config.
func New(cfg Config) (IssueTracker, error) {
	switch strings.ToLower(cfg.Type) {
	case "jira":
		return NewJiraTracker(cfg), nil
	case "linear":
		return NewLinearTracker(), nil
	default:
		return nil, fmt.Errorf("unknown tracker type: %s", cfg.Type)
	}
}

// ---- Jira ----

// JiraTracker implements IssueTracker using Jira REST API v3.
type JiraTracker struct {
	cfg    Config
	client *http.Client
}

// NewJiraTracker creates a new JiraTracker with the given config.
func NewJiraTracker(cfg Config) *JiraTracker {
	return &JiraTracker{cfg: cfg, client: &http.Client{}}
}

func (j *JiraTracker) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	url := strings.TrimRight(j.cfg.BaseURL, "/") + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(j.cfg.Email, j.cfg.Token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (j *JiraTracker) do(req *http.Request) ([]byte, error) {
	resp, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("jira: HTTP %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

func (j *JiraTracker) ListIssues(ctx context.Context, project string) ([]domain.Issue, error) {
	path := fmt.Sprintf("/rest/api/3/search?jql=project=%s&maxResults=50", project)
	req, err := j.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	data, err := j.do(req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Issues []struct {
			ID  string `json:"id"`
			Key string `json:"key"`
			Fields struct {
				Summary     string `json:"summary"`
				Description interface{} `json:"description"`
				Status      struct {
					Name string `json:"name"`
				} `json:"status"`
				Priority struct {
					Name string `json:"name"`
				} `json:"priority"`
			} `json:"fields"`
		} `json:"issues"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("jira: parse response: %w", err)
	}

	issues := make([]domain.Issue, 0, len(result.Issues))
	for _, ji := range result.Issues {
		issues = append(issues, domain.Issue{
			ID:         ji.ID,
			Title:      ji.Fields.Summary,
			Status:     mapJiraStatus(ji.Fields.Status.Name),
			Priority:   mapJiraPriority(ji.Fields.Priority.Name),
			TrackerRef: ji.Key,
		})
	}
	return issues, nil
}

func (j *JiraTracker) GetIssue(ctx context.Context, id string) (*domain.Issue, error) {
	req, err := j.newRequest(ctx, http.MethodGet, "/rest/api/3/issue/"+id, nil)
	if err != nil {
		return nil, err
	}
	data, err := j.do(req)
	if err != nil {
		return nil, err
	}

	var ji struct {
		ID  string `json:"id"`
		Key string `json:"key"`
		Fields struct {
			Summary  string `json:"summary"`
			Status   struct{ Name string `json:"name"` } `json:"status"`
			Priority struct{ Name string `json:"name"` } `json:"priority"`
		} `json:"fields"`
	}
	if err := json.Unmarshal(data, &ji); err != nil {
		return nil, fmt.Errorf("jira: parse issue: %w", err)
	}
	return &domain.Issue{
		ID:         ji.ID,
		Title:      ji.Fields.Summary,
		Status:     mapJiraStatus(ji.Fields.Status.Name),
		Priority:   mapJiraPriority(ji.Fields.Priority.Name),
		TrackerRef: ji.Key,
	}, nil
}

func (j *JiraTracker) UpdateIssue(ctx context.Context, id string, update domain.IssueUpdate) error {
	fields := map[string]interface{}{}
	if update.Title != nil {
		fields["summary"] = *update.Title
	}
	if update.Status != nil {
		// Jira status transitions require a separate call; this is a simplified update.
		fields["status"] = map[string]string{"name": string(*update.Status)}
	}
	payload, err := json.Marshal(map[string]interface{}{"fields": fields})
	if err != nil {
		return err
	}
	req, err := j.newRequest(ctx, http.MethodPut, "/rest/api/3/issue/"+id, strings.NewReader(string(payload)))
	if err != nil {
		return err
	}
	_, err = j.do(req)
	return err
}

func (j *JiraTracker) CreateComment(ctx context.Context, issueID string, body string) error {
	payload, err := json.Marshal(map[string]interface{}{
		"body": map[string]interface{}{
			"type":    "doc",
			"version": 1,
			"content": []map[string]interface{}{
				{
					"type": "paragraph",
					"content": []map[string]interface{}{
						{"type": "text", "text": body},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	req, err := j.newRequest(ctx, http.MethodPost, "/rest/api/3/issue/"+issueID+"/comment", strings.NewReader(string(payload)))
	if err != nil {
		return err
	}
	_, err = j.do(req)
	return err
}

func mapJiraStatus(s string) domain.IssueStatus {
	switch strings.ToLower(s) {
	case "to do", "open", "backlog":
		return domain.IssueStatusTodo
	case "in progress":
		return domain.IssueStatusInProgress
	case "in review", "code review":
		return domain.IssueStatusInReview
	case "done", "closed", "resolved":
		return domain.IssueStatusDone
	default:
		return domain.IssueStatusTodo
	}
}

func mapJiraPriority(p string) domain.IssuePriority {
	switch strings.ToLower(p) {
	case "critical", "blocker":
		return domain.IssuePriorityCritical
	case "high", "major":
		return domain.IssuePriorityHigh
	case "low", "minor", "trivial":
		return domain.IssuePriorityLow
	default:
		return domain.IssuePriorityMedium
	}
}

// ---- Linear (stub) ----

// LinearTracker is a stub that returns ErrNotImplemented for all methods.
type LinearTracker struct{}

func NewLinearTracker() *LinearTracker { return &LinearTracker{} }

func (l *LinearTracker) ListIssues(_ context.Context, _ string) ([]domain.Issue, error) {
	return nil, ErrNotImplemented
}
func (l *LinearTracker) GetIssue(_ context.Context, _ string) (*domain.Issue, error) {
	return nil, ErrNotImplemented
}
func (l *LinearTracker) UpdateIssue(_ context.Context, _ string, _ domain.IssueUpdate) error {
	return ErrNotImplemented
}
func (l *LinearTracker) CreateComment(_ context.Context, _ string, _ string) error {
	return ErrNotImplemented
}
