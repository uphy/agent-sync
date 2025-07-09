package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List supported agents and defaults",
		Long:  "Display the list of AI agents supported by agent-def and default output paths per type.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("List command invoked")
			// TODO: implement listing logic
			return nil
		},
	}
	return cmd
}
