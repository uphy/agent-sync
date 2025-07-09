package command

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/user/agent-def/internal/agent"
	"github.com/user/agent-def/internal/model"
	"github.com/user/agent-def/internal/template"
	"github.com/user/agent-def/internal/util"
)

// NewMemoryCommand creates a new memory command
func NewMemoryCommand(registry *agent.Registry) *cobra.Command {
	var (
		agentType string
		basePath  string
	)

	cmd := &cobra.Command{
		Use:   "memory [file]",
		Short: "Process memory context",
		Long: `Process memory context files and format them for the specified agent.
This command takes a markdown file with YAML frontmatter containing context
information and formats it appropriately for the target agent.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the file path from args
			filePath := args[0]
			return RunMemory(filePath, agentType, basePath, registry)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&agentType, "agent", "a", "claude", "Target agent type (claude, roo, etc.)")
	cmd.Flags().StringVarP(&basePath, "base", "b", "", "Base path for resolving relative paths (defaults to input file directory)")

	return cmd
}

// RunMemory processes a memory context file and formats it for the specified agent
func RunMemory(filePath, agentType, basePath string, registry *agent.Registry) error {
	// Use current working directory as base if not specified
	var err error
	if basePath == "" {
		// Use the directory of the input file as the base path
		basePath = filepath.Dir(filePath)
	}

	// Create file system
	fs := &util.RealFileSystem{}

	// Check if agent exists
	selectedAgent, found := registry.Get(agentType)
	if !found {
		return &util.ErrInvalidAgent{Type: agentType}
	}

	// Load context
	ctx, err := model.LoadContext(fs, filePath)
	if err != nil {
		return util.WrapError(err, "failed to load context file")
	}

	// Create template engine
	fileResolver := &MemoryFileResolver{FS: fs, BasePath: basePath}
	engine := template.NewEngine(fileResolver, agentType, basePath, registry)

	// Process context with template engine
	if err := model.ProcessContextWithEngine(engine, ctx); err != nil {
		return util.WrapError(err, "failed to process context")
	}

	// Format memory for the specific agent
	output, err := selectedAgent.FormatMemory(ctx.Processed)
	if err != nil {
		return util.WrapError(err, "failed to format memory for agent")
	}

	// Output to stdout
	fmt.Println(output)
	return nil
}

// MemoryFileResolver adapts FileSystem to template.FileResolver
type MemoryFileResolver struct {
	FS       util.FileSystem
	BasePath string
}

// Read reads a file at the given path
func (r *MemoryFileResolver) Read(path string) ([]byte, error) {
	return r.FS.ReadFile(path)
}

// Exists checks if a file exists
func (r *MemoryFileResolver) Exists(path string) bool {
	return r.FS.FileExists(path)
}

// ResolvePath resolves a relative path to absolute
func (r *MemoryFileResolver) ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(r.BasePath, path)
}
