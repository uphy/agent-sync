package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/user/agent-def/internal/agent"
	"github.com/user/agent-def/internal/model"
	"github.com/user/agent-def/internal/util"
)

// NewCommandCommand creates a new command transformation command
func NewCommandCommand(registry *agent.Registry) *cobra.Command {
	var (
		agentType string
		basePath  string
	)

	cmd := &cobra.Command{
		Use:   "command [file_or_directory]",
		Short: "Process command definitions",
		Long: `Process command definition files and format them for the specified agent.
This command can take either a single file or a directory containing multiple
command definition files. Each file is processed and formatted for the target agent.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the file or directory path from args
			path := args[0]
			return RunCommand(path, agentType, basePath, registry)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&agentType, "agent", "a", "claude", "Target agent type (claude, roo, etc.)")
	cmd.Flags().StringVarP(&basePath, "base", "b", "", "Base path for resolving relative paths (defaults to input file/directory)")

	return cmd
}

// RunCommand processes command definition files and formats them for the specified agent
func RunCommand(path, agentType, basePath string, registry *agent.Registry) error {
	// Use current working directory as base if not specified
	if basePath == "" {
		// Use the directory of the input file as the base path
		basePath = filepath.Dir(path)
	}

	// Create file system
	fs := &util.RealFileSystem{}

	// Check if agent exists
	selectedAgent, found := registry.Get(agentType)
	if !found {
		return &util.ErrInvalidAgent{Type: agentType}
	}

	// Load commands - either from a single file or a directory
	var commands []model.Command
	var err error

	// Check if path is a directory
	fileInfo, err := os.Stat(path)
	if err != nil {
		return util.WrapError(err, "failed to access path")
	}

	if fileInfo.IsDir() {
		// Path is a directory, load all command files
		commands, err = model.LoadCommandsFromDirectory(fs, path)
		if err != nil {
			return util.WrapError(err, "failed to load commands from directory")
		}
	} else {
		// Path is a single file
		content, err := fs.ReadFile(path)
		if err != nil {
			return util.WrapError(err, "failed to read command file")
		}

		cmd, err := model.ParseCommand(path, content)
		if err != nil {
			return util.WrapError(err, "failed to parse command")
		}

		if err := cmd.Validate(); err != nil {
			return util.WrapError(err, fmt.Sprintf("invalid command in %s", path))
		}

		commands = append(commands, cmd)
	}

	// Check if we found any commands
	if len(commands) == 0 {
		return fmt.Errorf("no valid command files found in %s", path)
	}

	// Format commands for the specific agent
	output, err := selectedAgent.FormatCommand(commands)
	if err != nil {
		return util.WrapError(err, "failed to format commands for agent")
	}

	// Output to stdout
	fmt.Println(output)
	return nil
}
