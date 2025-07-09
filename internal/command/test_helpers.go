package command

import (
	"github.com/user/agent-def/internal/agent"
)

// testRegistry is a package-level variable to store the agent registry for tests
var testRegistry *agent.Registry

// init initializes the test registry
func init() {
	testRegistry = agent.NewRegistry()
	testRegistry.Register(&agent.Claude{})
	testRegistry.Register(&agent.Roo{})
}

// createTestRegistry returns the shared test registry
func createTestRegistry() *agent.Registry {
	return testRegistry
}
