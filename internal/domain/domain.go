// Package domain contains core types used across the application.
package domain

import "time"

// Company represents a single-user workspace.
type Company struct {
	ID        string            `json:"id" gorm:"primaryKey"`
	Name      string            `json:"name"`
	Settings  map[string]string `json:"settings" gorm:"serializer:json"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

// Agent represents a Kubernetes Job template.
type Agent struct {
	ID          string      `json:"id" gorm:"primaryKey"`
	CompanyID   string      `json:"companyId" gorm:"index"`
	Name        string      `json:"name"`
	Role        string      `json:"role"`
	Image       string      `json:"image"`
	Description string      `json:"description"`
	SecretRefs  []SecretRef `json:"secretRefs" gorm:"serializer:json"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

// IssueStatus represents the lifecycle state of an issue.
type IssueStatus string

const (
	IssueStatusBacklog    IssueStatus = "backlog"
	IssueStatusTodo       IssueStatus = "todo"
	IssueStatusInProgress IssueStatus = "in_progress"
	IssueStatusInReview   IssueStatus = "in_review"
	IssueStatusDone       IssueStatus = "done"
	IssueStatusBlocked    IssueStatus = "blocked"
	IssueStatusCancelled  IssueStatus = "cancelled"
)

// IssuePriority represents the urgency of an issue.
type IssuePriority string

const (
	IssuePriorityCritical IssuePriority = "critical"
	IssuePriorityHigh     IssuePriority = "high"
	IssuePriorityMedium   IssuePriority = "medium"
	IssuePriorityLow      IssuePriority = "low"
)

// Issue represents a work item.
type Issue struct {
	ID              string        `json:"id" gorm:"primaryKey"`
	CompanyID       string        `json:"companyId" gorm:"index"`
	ProjectID       string        `json:"projectId,omitempty" gorm:"index"`
	GoalID          string        `json:"goalId,omitempty" gorm:"index"`
	ParentID        string        `json:"parentId,omitempty" gorm:"index"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	Status          IssueStatus   `json:"status"`
	Priority        IssuePriority `json:"priority"`
	AssigneeAgentID string        `json:"assigneeAgentId,omitempty" gorm:"index"`
	TrackerRef      string        `json:"trackerRef,omitempty"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

// IssueUpdate is used to partially update an issue.
type IssueUpdate struct {
	Title           *string        `json:"title,omitempty"`
	Description     *string        `json:"description,omitempty"`
	Status          *IssueStatus   `json:"status,omitempty"`
	Priority        *IssuePriority `json:"priority,omitempty"`
	AssigneeAgentID *string        `json:"assigneeAgentId,omitempty"`
}

// RunStatus represents the state of an agent execution.
type RunStatus string

const (
	RunStatusPending   RunStatus = "pending"
	RunStatusRunning   RunStatus = "running"
	RunStatusSucceeded RunStatus = "succeeded"
	RunStatusFailed    RunStatus = "failed"
	RunStatusCancelled RunStatus = "cancelled"
)

// Run represents a single agent execution.
type Run struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	CompanyID  string     `json:"companyId" gorm:"index"`
	AgentID    string     `json:"agentId" gorm:"index"`
	IssueID    string     `json:"issueId" gorm:"index"`
	Status     RunStatus  `json:"status"`
	Logs       string     `json:"logs,omitempty"`
	StartedAt  *time.Time `json:"startedAt,omitempty"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// Routine represents a scheduled task.
type Routine struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	CompanyID     string    `json:"companyId" gorm:"index"`
	Name          string    `json:"name"`
	CronExpr      string    `json:"cronExpr"`
	AgentID       string    `json:"agentId" gorm:"index"`
	IssueTemplate string    `json:"issueTemplate" gorm:"type:text"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// Project represents a grouping of issues.
type Project struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	CompanyID   string     `json:"companyId" gorm:"index"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	TargetDate  *time.Time `json:"targetDate,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// Goal represents a high-level objective.
type Goal struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	CompanyID   string    `json:"companyId" gorm:"index"`
	ParentID    string    `json:"parentId,omitempty" gorm:"index"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Level       string    `json:"level"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// SecretScope defines the visibility of a secret.
type SecretScope string

const (
	SecretScopeGlobal   SecretScope = "global"
	SecretScopePerAgent SecretScope = "per-agent"
)

// SecretRef is a reference to a registered secret. The secret value is never
// stored in or returned via this struct — only metadata is exposed.
type SecretRef struct {
	ID          string      `json:"id"`
	CompanyID   string      `json:"companyId"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Scope       SecretScope `json:"scope"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	// Secret value is NEVER stored here after write.
}

// SecretEntry is the database storage model for an encrypted secret value.
// All fields carry json:"-" to prevent accidental serialisation.
type SecretEntry struct {
	ID           string    `json:"-" gorm:"primaryKey"`
	CompanyID    string    `json:"-" gorm:"index"`
	Name         string    `json:"-"`
	Description  string    `json:"-"`
	Scope        string    `json:"-"`
	EncryptedVal []byte    `json:"-" gorm:"type:blob"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}
