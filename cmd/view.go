package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheCoolRobot/asana-cli/internal/asana"
	"github.com/TheCoolRobot/asana-cli/internal/ui"
)

var viewCmd = &cobra.Command{
	Use:   "view [task-id]",
	Short: "View task details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskGID := args[0]
		client := asana.NewClient(token)

		task, err := client.GetTask(taskGID)
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Println("Error:", err)
			}
			return err
		}

		if jsonOutput {
			ui.PrintJSON(task, nil)
		} else {
			fmt.Printf("📋 %s\n", task.Name)
			fmt.Printf("   GID: %s\n", task.GID)
			fmt.Printf("   Status: %v\n", task.Completed)
			if task.DueDate != nil && !task.DueDate.IsZero() {
				fmt.Printf("   Due: %s\n", task.DueDate.Format("2006-01-02"))
			}
			if task.Description != "" {
				fmt.Printf("   Description: %s\n", task.Description)
			}
			if task.Assignee != nil {
				fmt.Printf("   Assigned to: %s\n", task.Assignee.Name)
			}
			if len(task.Tags) > 0 {
				fmt.Printf("   Tags: ")
				for i, tag := range task.Tags {
					if i > 0 {
						fmt.Printf(", ")
					}
					fmt.Printf("%s", tag.Name)
				}
				fmt.Printf("\n")
			}
		}

		return nil
	},
}