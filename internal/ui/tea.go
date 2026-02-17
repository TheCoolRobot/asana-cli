package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/TheCoolRobot/asana-cli/internal/asana"
	"github.com/TheCoolRobot/asana-cli/internal/config"
)

type TaskListItem struct {
	Task     *asana.Task
	Selected bool
}

type Model struct {
	items           []TaskListItem
	cursor          int
	filterCompleted bool
	searchQuery     string
	sortBy          string // name, due_date, priority
	width           int
	height          int
	loading         bool
	message         string
	mode            string // "tasks", "projects", or "confirm"
	projects        []config.ProjectConfig
	projectCursor   int
	currentProject  string
	client          *asana.Client
	projectGID      string
	confirmAction   string // "delete", "complete"
	confirmTaskGID  string
	confirmTaskName string
}

func NewModel(tasks []*asana.Task, client *asana.Client, projectGID string) Model {
	cfg, _ := config.Load()
	items := make([]TaskListItem, len(tasks))
	for i, t := range tasks {
		items[i] = TaskListItem{Task: t, Selected: false}
	}

	return Model{
		items:          items,
		cursor:         0,
		sortBy:         "name",
		width:          80,
		height:         24,
		mode:           "tasks",
		projects:       cfg.ListProjects(),
		projectCursor:  0,
		currentProject: cfg.CurrentProject,
		client:         client,
		projectGID:     projectGID,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if m.mode == "confirm" {
			return m.updateConfirmMode(msg)
		} else if m.mode == "projects" {
			return m.updateProjectMode(msg)
		}
		return m.updateTaskMode(msg)
	}
	return m, nil
}

func (m Model) updateTaskMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.items)-1 {
			m.cursor++
		}

	case " ":
		if m.cursor < len(m.items) {
			m.items[m.cursor].Selected = !m.items[m.cursor].Selected
		}

	case "c":
		if m.cursor < len(m.items) {
			m.confirmAction = "complete"
			m.confirmTaskGID = m.items[m.cursor].Task.GID
			m.confirmTaskName = m.items[m.cursor].Task.Name
			m.mode = "confirm"
		}

	case "d":
		if m.cursor < len(m.items) {
			m.confirmAction = "delete"
			m.confirmTaskGID = m.items[m.cursor].Task.GID
			m.confirmTaskName = m.items[m.cursor].Task.Name
			m.mode = "confirm"
		}

	case "f":
		m.filterCompleted = !m.filterCompleted
		if m.filterCompleted {
			m.message = "Filtering: Hidden completed tasks"
		} else {
			m.message = "Showing: All tasks"
		}

	case "s":
		m.toggleSort()

	case "p":
		m.mode = "projects"
		m.message = "Project selector (↑↓ navigate, enter to switch, q to close)"

	case "e":
		if m.cursor < len(m.items) {
			m.message = fmt.Sprintf("Edit task: %s (not yet implemented in TUI)", m.items[m.cursor].Task.Name)
		}

	case "/":
		m.message = "Enter search query (not implemented in TUI demo)"

	case "enter":
		if m.cursor < len(m.items) {
			m.message = fmt.Sprintf("View details: %s (press 'v' to view full)", m.items[m.cursor].Task.Name)
		}
	}
	return m, nil
}

func (m Model) updateConfirmMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.loading = true
		if m.confirmAction == "complete" {
			_, err := m.client.CompleteTask(m.confirmTaskGID)
			if err != nil {
				m.message = fmt.Sprintf("❌ Error completing task: %v", err)
			} else {
				// Mark as completed in the list
				for i, item := range m.items {
					if item.Task.GID == m.confirmTaskGID {
						m.items[i].Task.Completed = true
						break
					}
				}
				m.message = fmt.Sprintf("✓ Completed: %s", m.confirmTaskName)
			}
		} else if m.confirmAction == "delete" {
			err := m.client.DeleteTask(m.confirmTaskGID)
			if err != nil {
				m.message = fmt.Sprintf("❌ Error deleting task: %v", err)
			} else {
				// Remove from the list
				for i, item := range m.items {
					if item.Task.GID == m.confirmTaskGID {
						m.items = append(m.items[:i], m.items[i+1:]...)
						if m.cursor >= len(m.items) && m.cursor > 0 {
							m.cursor--
						}
						break
					}
				}
				m.message = fmt.Sprintf("✓ Deleted: %s", m.confirmTaskName)
			}
		}
		m.loading = false
		m.mode = "tasks"
		m.confirmAction = ""
		m.confirmTaskGID = ""
		m.confirmTaskName = ""

	case "n", "esc":
		m.mode = "tasks"
		m.confirmAction = ""
		m.confirmTaskGID = ""
		m.confirmTaskName = ""
		m.message = "Cancelled"
	}
	return m, nil
}

func (m Model) updateProjectMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.mode = "tasks"
		m.message = ""

	case "up", "k":
		if m.projectCursor > 0 {
			m.projectCursor--
		}

	case "down", "j":
		if m.projectCursor < len(m.projects)-1 {
			m.projectCursor++
		}

	case "enter":
		if m.projectCursor < len(m.projects) {
			selectedProject := m.projects[m.projectCursor]
			cfg, _ := config.Load()
			cfg.SetCurrentProject(selectedProject.Name)
			m.currentProject = selectedProject.Name
			m.message = fmt.Sprintf("✓ Switched to: %s", selectedProject.Name)
			m.mode = "tasks"
		}
	}
	return m, nil
}

func (m Model) toggleSort() {
	switch m.sortBy {
	case "name":
		m.sortBy = "due_date"
		m.message = "Sorting by: Due Date"
	case "due_date":
		m.sortBy = "priority"
		m.message = "Sorting by: Priority"
	default:
		m.sortBy = "name"
		m.message = "Sorting by: Name"
	}
}

func (m Model) View() string {
	if m.mode == "confirm" {
		return m.viewConfirmation()
	} else if m.mode == "projects" {
		return m.viewProjects()
	}
	return m.viewTasks()
}

func (m Model) viewConfirmation() string {
	var sb strings.Builder

	sb.WriteString(StyleError.Render("⚠️  CONFIRM ACTION") + "\n\n")

	if m.confirmAction == "complete" {
		sb.WriteString(fmt.Sprintf("Mark task as complete?\n\n  %s\n\n", m.confirmTaskName))
	} else if m.confirmAction == "delete" {
		sb.WriteString(fmt.Sprintf("Delete this task permanently?\n\n  %s\n\n", m.confirmTaskName))
	}

	if m.loading {
		sb.WriteString(StyleWarning.Render("Processing...") + "\n")
	} else {
		sb.WriteString(StyleDim.Render("[y] Confirm  [n] Cancel") + "\n")
	}

	return sb.String()
}

func (m Model) viewTasks() string {
	if len(m.items) == 0 {
		return StyleError.Render("No tasks found")
	}

	var sb strings.Builder

	// Title with current project
	title := fmt.Sprintf("📋 Asana Tasks - %s", m.currentProject)
	sb.WriteString(StyleTitle.Render(title) + "\n\n")

	// Status line
	statusLine := fmt.Sprintf("[%d/%d tasks]", len(m.items), len(m.items))
	if m.filterCompleted {
		statusLine += " [filtered]"
	}
	statusLine += fmt.Sprintf(" [sort: %s]", m.sortBy)
	sb.WriteString(StyleDim.Render(statusLine) + "\n\n")

	// Task list
	for i, item := range m.items {
		if m.filterCompleted && item.Task.Completed {
			continue
		}

		cursor := "  "
		if m.cursor == i {
			cursor = "→ "
		}

		checkmark := "☐"
		if item.Task.Completed {
			checkmark = "☑"
		}

		selected := " "
		if item.Selected {
			selected = "✓"
		}

		// Task name with style
		taskName := item.Task.Name
		if item.Task.Completed {
			taskName = StyleCompleted.Render(taskName)
		} else if m.cursor == i {
			taskName = StyleSelected.Render(taskName)
		}

		// Priority indicator
		priority := ""
		if item.Task.Priority != "" {
			if strings.Contains(strings.ToLower(item.Task.Priority), "high") {
				priority = " " + StyleHighPriority.Render("!!!")
			} else if strings.Contains(strings.ToLower(item.Task.Priority), "medium") {
				priority = " " + StyleMediumPriority.Render("!!")
			}
		}

		// Due date
		dueDate := ""
		if item.Task.DueDate != nil && !item.Task.DueDate.IsZero() {
			daysUntil := int(time.Until(item.Task.DueDate.Time).Hours() / 24)
			if daysUntil < 0 {
				dueDate = fmt.Sprintf(" %s", StyleError.Render(fmt.Sprintf("[%d days overdue]", -daysUntil)))
			} else if daysUntil == 0 {
				dueDate = fmt.Sprintf(" %s", StyleWarning.Render("[Due today]"))
			} else if daysUntil <= 3 {
				dueDate = fmt.Sprintf(" %s", StyleWarning.Render(fmt.Sprintf("[%d days left]", daysUntil)))
			}
		}

		line := fmt.Sprintf("%s[%s] %s %s%s%s\n", cursor, selected, checkmark, taskName, priority, dueDate)
		sb.WriteString(line)
	}

	// Footer
	sb.WriteString("\n" + StyleDim.Render(strings.Repeat("─", m.width)) + "\n")
	helpText := "[↑↓] navigate  [space] select  [c] complete  [d] delete  [f] filter  [s] sort  [p] projects  [?] help  [q] quit"
	if m.width < len(helpText) {
		helpText = "[↑↓] nav  [c] done  [d] del  [p] proj  [q] quit"
	}
	sb.WriteString(StyleDim.Render(helpText) + "\n")

	if m.message != "" {
		sb.WriteString("\n" + StyleSuccess.Render(m.message) + "\n")
	}

	return sb.String()
}

func (m Model) viewProjects() string {
	if len(m.projects) == 0 {
		return StyleError.Render("No projects configured")
	}

	var sb strings.Builder

	sb.WriteString(StyleTitle.Render("🗂️  Projects") + "\n\n")
	sb.WriteString(StyleDim.Render("Select a project to switch to:") + "\n\n")

	for i, proj := range m.projects {
		cursor := "  "
		if m.projectCursor == i {
			cursor = "→ "
		}

		style := lipgloss.NewStyle()
		if m.projectCursor == i {
			style = StyleSelected
		}

		current := ""
		if proj.Name == m.currentProject {
			current = " ✓"
		}

		line := fmt.Sprintf("%s%s%s\n", cursor, style.Render(proj.Name), current)
		sb.WriteString(line)

		if proj.Description != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", StyleDim.Render(proj.Description)))
		}
	}

	sb.WriteString("\n" + StyleDim.Render(strings.Repeat("─", m.width)) + "\n")
	sb.WriteString(StyleDim.Render("[↑↓] navigate  [enter] switch  [q] close") + "\n")

	if m.message != "" {
		sb.WriteString("\n" + StyleSuccess.Render(m.message) + "\n")
	}

	return sb.String()
}

func StartTUI(tasks []*asana.Task, client *asana.Client, projectGID string) {
	m := NewModel(tasks, client, projectGID)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}