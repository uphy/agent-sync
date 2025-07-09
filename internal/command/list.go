package command

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/user/agent-def/internal/agent"
)

// NewListCommand creates a new list supported agents command
func NewListCommand(registry *agent.Registry) *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List supported agents",
		Long:  "List all registered agents and their capabilities",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunList(os.Stdout, verbose, registry)
		},
	}

	// Add flags
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed capabilities for each agent")

	return cmd
}

// RunList displays all registered agents and their capabilities
func RunList(out io.Writer, verbose bool, registry *agent.Registry) error {
	// Get all registered agents
	agents := registry.List()
	if len(agents) == 0 {
		fmt.Fprintln(out, "No agents registered.")
		return nil
	}

	// Create a tabwriter for aligned output
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)

	// Print header
	fmt.Fprintln(w, "ID\tNAME\tCAPABILITIES")

	// Print each agent and its capabilities
	for _, a := range agents {
		// Build capabilities list
		capabilities := []string{
			"Memory formatting",
			"Command formatting",
			"File references",
			"MCP support",
		}

		capabilitiesStr := strings.Join(capabilities, ", ")
		if !verbose {
			// Truncate long capability strings if not in verbose mode
			if len(capabilitiesStr) > 40 {
				capabilitiesStr = capabilitiesStr[:37] + "..."
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\n", a.ID(), a.Name(), capabilitiesStr)
	}

	// Flush the tabwriter
	return w.Flush()
}
