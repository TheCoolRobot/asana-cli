package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheCoolRobot/asana-cli/internal/asana"
	"github.com/TheCoolRobot/asana-cli/internal/ui"
)

var searchCmd = &cobra.Command{
	Use:   "search [workspace-id] [query]",
	Short: "Search for tasks",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceGID := args[0]
		query := args[1]
		client := asana.NewClient(token)

		tasks, err := client.Search(workspaceGID, query)
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Println("Error:", err)
			}
			return err
		}

		taskPtrs := make([]*asana.Task, len(tasks))
		for i := range tasks {
			taskPtrs[i] = &tasks[i]
		}

		if jsonOutput {
			meta := map[string]interface{}{
				"count": len(tasks),
				"query": query,
			}
			ui.PrintJSONWithMeta(tasks, meta, nil)
		} else {
			// Use empty projectGID for search results since they're from multiple projects
			ui.StartTUI(taskPtrs, client, "")
		}

		return nil
	},
}