package agent

import "github.com/uphy/agent-sync/internal/model"

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

	// MemoryPath returns the default path for memory files based on user scope
	MemoryPath(userScope bool) string

	// CommandPath returns the default path for command files based on user scope
	CommandPath(userScope bool) string

	// FormatMode processes mode definitions for this agent
	FormatMode(modes []model.Mode) (string, error)

	// ModePath returns the default path for mode files based on user scope
	ModePath(userScope bool) string
}
