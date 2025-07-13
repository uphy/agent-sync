package agent

// Registry maintains a registry of available agents
type Registry struct {
	agents map[string]Agent
}

// NewRegistry creates a new agent registry
func NewRegistry() *Registry {
	r := &Registry{
		agents: make(map[string]Agent),
	}
	// Register default agents
	r.RegisterDefaults()
	return r
}

// RegisterDefaults registers all default agents
func (r *Registry) RegisterDefaults() {
	r.Register(&Roo{})
	r.Register(&Claude{})
	r.Register(&Cline{})
	r.Register(&Copilot{}) // Register Copilot agent
}

// Register registers an agent
func (r *Registry) Register(agent Agent) {
	r.agents[agent.ID()] = agent
}

// Get returns an agent by ID
func (r *Registry) Get(id string) (Agent, bool) {
	agent, ok := r.agents[id]
	return agent, ok
}

// List returns all registered agents
func (r *Registry) List() []Agent {
	result := make([]Agent, 0, len(r.agents))
	for _, agent := range r.agents {
		result = append(result, agent)
	}
	return result
}
