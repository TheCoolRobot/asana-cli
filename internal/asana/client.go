package asana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	
)
// "strings"
type Client struct {
	apiToken string
	baseURL  string
	http     *http.Client
}

func NewClient(apiToken string) *Client {
	if apiToken == "" {
		apiToken = os.Getenv("ASANA_TOKEN")
	}

	return &Client{
		apiToken: apiToken,
		baseURL:  "https://app.asana.com/api/1.0",
		http:     &http.Client{},
	}
}

func (c *Client) do(method, endpoint string, body interface{}) ([]byte, error) {
	if c.apiToken == "" {
		return nil, fmt.Errorf("ASANA_TOKEN not set")
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.apiToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GetMe retrieves current user info
func (c *Client) GetMe() (*User, error) {
	body, err := c.do("GET", "/users/me", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data *User `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetWorkspaces retrieves all workspaces
func (c *Client) GetWorkspaces() ([]Workspace, error) {
	body, err := c.do("GET", "/workspaces", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []Workspace `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetProjects retrieves projects in a workspace
func (c *Client) GetProjects(workspaceGID string) ([]Project, error) {
	endpoint := fmt.Sprintf("/projects?workspace=%s", workspaceGID)
	body, err := c.do("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []Project `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetTasks retrieves tasks from a project with optional filters
// Supports filters: completed_since, assignee, modified_since, etc.
func (c *Client) GetTasks(projectGID string, filters map[string]string) ([]Task, error) {
	endpoint := fmt.Sprintf("/projects/%s/tasks", projectGID)
	
	if len(filters) > 0 {
		q := url.Values{}
		for k, v := range filters {
			q.Add(k, v)
		}
		endpoint += "?" + q.Encode()
	}

	body, err := c.do("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []Task `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetTask retrieves a specific task using task GID
// GET /tasks/{task_gid}
func (c *Client) GetTask(taskGID string) (*Task, error) {
	endpoint := fmt.Sprintf("/tasks/%s", taskGID)
	body, err := c.do("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data *Task `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// CreateTask creates a new task
// POST /tasks
func (c *Client) CreateTask(req *TaskCreateRequest) (*Task, error) {
	payload := map[string]interface{}{
		"data": req,
	}

	body, err := c.do("POST", "/tasks", payload)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data *Task `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// UpdateTask updates a task using task GID
// PUT /tasks/{task_gid}
func (c *Client) UpdateTask(taskGID string, req *TaskUpdateRequest) (*Task, error) {
	endpoint := fmt.Sprintf("/tasks/%s", taskGID)
	payload := map[string]interface{}{
		"data": req,
	}

	body, err := c.do("PUT", endpoint, payload)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data *Task `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// CompleteTask marks a task as complete
// Calls UpdateTask with completed: true
func (c *Client) CompleteTask(taskGID string) (*Task, error) {
	completed := true
	return c.UpdateTask(taskGID, &TaskUpdateRequest{
		Completed: &completed,
	})
}

// DeleteTask deletes a task
// DELETE /tasks/{task_gid}
func (c *Client) DeleteTask(taskGID string) error {
	endpoint := fmt.Sprintf("/tasks/%s", taskGID)
	_, err := c.do("DELETE", endpoint, nil)
	return err
}

// GetSections retrieves sections in a project
func (c *Client) GetSections(projectGID string) ([]Section, error) {
	endpoint := fmt.Sprintf("/projects/%s/sections", projectGID)
	body, err := c.do("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []Section `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// Search searches for tasks by text query
// GET /workspaces/{workspace_gid}/tasks/search
func (c *Client) Search(workspaceGID, query string) ([]Task, error) {
	endpoint := fmt.Sprintf("/workspaces/%s/tasks/search?text=%s", workspaceGID, url.QueryEscape(query))
	body, err := c.do("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []Task `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetUserTaskList retrieves a user's "My Tasks" list
// GET /users/{user_gid}/user_task_list
func (c *Client) GetUserTaskList(userGID string) ([]Task, error) {
	endpoint := fmt.Sprintf("/users/%s/user_task_list", userGID)
	body, err := c.do("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []Task `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetUserTeams retrieves teams a specific user belongs to
// GET /users/{user_gid}/teams
func (c *Client) GetUserTeams(userGID string) ([]Team, error) {
	endpoint := fmt.Sprintf("/users/%s/teams", userGID)
	body, err := c.do("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []Team `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// UpdateWorkspace updates a workspace
// PUT /workspaces/{workspace_gid}
func (c *Client) UpdateWorkspace(workspaceGID string, req *WorkspaceUpdateRequest) (*Workspace, error) {
	endpoint := fmt.Sprintf("/workspaces/%s", workspaceGID)
	payload := map[string]interface{}{
		"data": req,
	}

	body, err := c.do("PUT", endpoint, payload)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data *Workspace `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetTasksByWorkspace retrieves tasks from a workspace with optional filters
func (c *Client) GetTasksByWorkspace(workspaceGID string, filters map[string]string) ([]Task, error) {
	endpoint := fmt.Sprintf("/workspaces/%s/tasks", workspaceGID)
	
	if len(filters) > 0 {
		q := url.Values{}
		for k, v := range filters {
			q.Add(k, v)
		}
		endpoint += "?" + q.Encode()
	}

	body, err := c.do("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []Task `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}