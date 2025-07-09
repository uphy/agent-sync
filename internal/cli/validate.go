package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the agent-def.yml configuration",
		Long:  "Check syntax and semantics of agent-def.yml for correctness and completeness.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Validate command invoked")
			// TODO: implement validation logic
			return nil
		},
	}
	return cmd
}
