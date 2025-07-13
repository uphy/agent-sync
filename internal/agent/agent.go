package agent

import "github.com/user/agent-def/internal/model"

// Agent represents a target AI agent that has specific format requirements
type Agent interface {
	// Name returns the name of the agent
	Name() string

	// ID returns the unique identifier of the agent
	ID() string

	// FormatFile converts a file reference for this agent's format
	FormatFile(path string) string

	// FormatMCP formats an MCP command for this agent
	FormatMCP(agent, command string, args ...string) string

	// FormatMemory processes a memory context for this agent
	FormatMemory(content string) (string, error)

	// FormatCommand processes a command definition for this agent
	FormatCommand(commands []model.Command) (string, error)

	// DefaultMemoryPath determines the default output path for memory tasks
	DefaultMemoryPath(outputBaseDir string, userScope bool, fileName string) (string, error)

	// DefaultCommandPath determines the default output path for command tasks
	DefaultCommandPath(outputBaseDir string, userScope bool, fileName string) (string, error)

	// ShouldConcatenate determines whether content should be concatenated
	// for this agent based on the task type (memory or command)
	ShouldConcatenate(taskType string) bool
}
