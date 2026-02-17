package asana

import (
	"time"
)

// CustomTime handles both date-only and full datetime formats from Asana
type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		return nil
	}

	// Remove quotes
	s = s[1 : len(s)-1]

	// Try RFC3339 first (with time)
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		ct.Time = t
		return nil
	}

	// Try date-only format (YYYY-MM-DD)
	t, err = time.Parse("2006-01-02", s)
	if err == nil {
		ct.Time = t
		return nil
	}

	return err
}

// Task represents an Asana task
type Task struct {
	GID             string      `json:"gid"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	Completed       bool        `json:"completed"`
	DueDate         *CustomTime `json:"due_on,omitempty"`
	DueAt           *CustomTime `json:"due_at,omitempty"`
	Priority        string      `json:"priority_value,omitempty"`
	Status          string      `json:"status,omitempty"`
	AssigneeStatus  string      `json:"assignee_status,omitempty"`
	Assignee        *User       `json:"assignee,omitempty"`
	Projects        []Project   `json:"projects,omitempty"`
	Tags            []Tag       `json:"tags,omitempty"`
	Dependencies    []Task      `json:"dependencies,omitempty"`
	Attachments     []Attachment `json:"attachments,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	ModifiedAt      time.Time   `json:"modified_at"`
}

// Project represents an Asana project
type Project struct {
	GID             string    `json:"gid"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Owner           *User     `json:"owner,omitempty"`
	Status          string    `json:"status"`
	Color           string    `json:"color"`
	CreatedAt       time.Time `json:"created_at"`
	ModifiedAt      time.Time `json:"modified_at"`
	Archived        bool      `json:"archived"`
	TaskCount       int       `json:"task_count,omitempty"`
}

// Section represents a project section
type Section struct {
	GID       string `json:"gid"`
	Name      string `json:"name"`
	ProjectID string `json:"project,omitempty"`
}

// User represents an Asana user
type User struct {
	GID       string `json:"gid"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Avatar    string `json:"photo,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
}

// Team represents an Asana team
type Team struct {
	GID   string `json:"gid"`
	Name  string `json:"name"`
	Users []User `json:"members,omitempty"`
}

// Workspace represents an Asana workspace
type Workspace struct {
	GID   string `json:"gid"`
	Name  string `json:"name"`
	Teams []Team `json:"teams,omitempty"`
}

// Tag represents an Asana tag
type Tag struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

// Attachment represents a task attachment
type Attachment struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// TaskCreateRequest for creating tasks
type TaskCreateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Projects    []string `json:"projects,omitempty"`
	Section     string   `json:"section,omitempty"`
	Assignee    string   `json:"assignee,omitempty"`
	DueOn       string   `json:"due_on,omitempty"`
	DueAt       string   `json:"due_at,omitempty"`
	Priority    string   `json:"priority_value,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// TaskUpdateRequest for updating tasks
type TaskUpdateRequest struct {
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	Completed      *bool  `json:"completed,omitempty"`
	Assignee       string `json:"assignee,omitempty"`
	DueOn          string `json:"due_on,omitempty"`
	DueAt          string `json:"due_at,omitempty"`
	Priority       string `json:"priority_value,omitempty"`
	AssigneeStatus string `json:"assignee_status,omitempty"`
	Status         string `json:"status,omitempty"`
}

// WorkspaceUpdateRequest for updating workspaces
type WorkspaceUpdateRequest struct {
	Name string `json:"name,omitempty"`
}