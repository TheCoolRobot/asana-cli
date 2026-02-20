package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheCoolRobot/asana-cli/internal/asana"
	"github.com/TheCoolRobot/asana-cli/internal/ui"
)

var (
	taskName        string
	taskDescription string
	taskAssignee    string
	taskDueDate     string
	taskPriority    string
	taskSection     string
)

var createCmd = &cobra.Command{
	Use:   "create [project-id]",
	Short: "Create a new task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if taskName == "" {
			if jsonOutput {
				ui.PrintJSON(nil, fmt.Errorf("task name is required"))
			} else {
				fmt.Println("Error: task name is required")
			}
			return fmt.Errorf("task name required")
		}

		projectGID := args[0]
		client := asana.NewClient(token)

		req := &asana.TaskCreateRequest{
			Name:        taskName,
			Description: taskDescription,
			Projects:    []string{projectGID},
			Assignee:    taskAssignee,
			DueOn:       taskDueDate,
			Priority:    taskPriority,
		}

		if taskSection != "" {
			req.Section = taskSection
		}

		task, err := client.CreateTask(req)
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Println("Error:", err)
			}
			return err
		}

		if jsonOutput {
			meta := map[string]interface{}{
				"action":     "created",
				"project_id": projectGID,
			}
			ui.PrintJSONWithMeta(task, meta, nil)
		} else {
			fmt.Printf("✓ Task created: %s\n", task.Name)
			fmt.Printf("  GID: %s\n", task.GID)
		}

		return nil
	},
}

func init() {
	createCmd.Flags().StringVar(&taskName, "name", "", "Task name (required)")
	createCmd.Flags().StringVar(&taskDescription, "description", "", "Task description")
	createCmd.Flags().StringVar(&taskAssignee, "assignee", "", "Assignee user GID")
	createCmd.Flags().StringVar(&taskDueDate, "due", "", "Due date (YYYY-MM-DD)")
	createCmd.Flags().StringVar(&taskPriority, "priority", "", "Priority (1=high, 2=medium, 3=low)")
	createCmd.Flags().StringVar(&taskSection, "section", "", "Section GID")
	if err := createCmd.MarkFlagRequired("name"); err=nil;{return err}
}