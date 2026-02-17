package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheCoolRobot/asana-cli/internal/asana"
	"github.com/TheCoolRobot/asana-cli/internal/ui"
)

var (
	updateName        string
	updateDescription string
	updateAssignee    string
	updateDueDate     string
	updatePriority    string
)

var updateCmd = &cobra.Command{
	Use:   "update [task-id]",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskGID := args[0]
		client := asana.NewClient(token)

		req := &asana.TaskUpdateRequest{
			Name:        updateName,
			Description: updateDescription,
			Assignee:    updateAssignee,
			DueOn:       updateDueDate,
			Priority:    updatePriority,
		}

		task, err := client.UpdateTask(taskGID, req)
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
				"action":   "updated",
				"task_gid": taskGID,
			}
			ui.PrintJSONWithMeta(task, meta, nil)
		} else {
			fmt.Printf("✓ Task updated: %s\n", task.Name)
		}

		return nil
	},
}

func init() {
	updateCmd.Flags().StringVar(&updateName, "name", "", "New task name")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "New task description")
	updateCmd.Flags().StringVar(&updateAssignee, "assignee", "", "New assignee user GID")
	updateCmd.Flags().StringVar(&updateDueDate, "due", "", "New due date (YYYY-MM-DD)")
	updateCmd.Flags().StringVar(&updatePriority, "priority", "", "New priority (1=high, 2=medium, 3=low)")
}