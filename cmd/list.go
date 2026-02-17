package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/TheCoolRobot/asana-cli/internal/asana"
	"github.com/TheCoolRobot/asana-cli/internal/config"
	"github.com/TheCoolRobot/asana-cli/internal/ui"
)

var (
	filterCompleted bool
	filterAssignee  string
	filterTag       string
)

var listCmd = &cobra.Command{
	Use:   "list [project-id]",
	Short: "List tasks from a project",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectGID := ""

		// Use provided project ID, or fall back to current project
		if len(args) > 0 {
			projectGID = args[0]
		} else {
			cfg, _ := config.Load()
			currentProj := cfg.GetCurrentProject()
			if currentProj == nil {
				return fmt.Errorf("no project ID provided and no current project set. Use: asana-cli list <project-id> or asana-cli config project switch <name>")
			}
			projectGID = currentProj.ProjectID
		}

		client := asana.NewClient(token)

		filters := make(map[string]string)
		if filterCompleted {
			filters["completed_since"] = "now"
		}
		if filterAssignee != "" {
			filters["assignee"] = filterAssignee
		}

		tasks, err := client.GetTasks(projectGID, filters)
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Println("Error:", err)
			}
			return err
		}

		// Convert to pointers
		taskPtrs := make([]*asana.Task, len(tasks))
		for i := range tasks {
			taskPtrs[i] = &tasks[i]
		}

		if jsonOutput {
			meta := map[string]interface{}{
				"count":      len(tasks),
				"project_id": projectGID,
				"fetched_at": time.Now().Format(time.RFC3339),
			}
			if filterCompleted {
				meta["filter_completed"] = true
			}
			if filterAssignee != "" {
				meta["filter_assignee"] = filterAssignee
			}
			ui.PrintJSONWithMeta(tasks, meta, nil)
		} else {
			ui.StartTUI(taskPtrs, client, projectGID)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().BoolVar(&filterCompleted, "completed", false, "Show only completed tasks")
	listCmd.Flags().StringVar(&filterAssignee, "assignee", "", "Filter by assignee ID")
	listCmd.Flags().StringVar(&filterTag, "tag", "", "Filter by tag")
}