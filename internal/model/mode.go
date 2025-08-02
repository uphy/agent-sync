package model

// Mode represents a subagent definition
type Mode struct {
	Claude  *ClaudeCodeMode `yaml:"claude"`
	Roo     *RooMode        `yaml:"roo"`
	Content string
}

// ClaudeCodeMode defines the properties for a Claude subagent
type ClaudeCodeMode struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Tools       []string `yaml:"tools,omitempty"`
}

// RooMode defines the properties for a Roo subagent (mode)
type RooMode struct {
	Slug           string   `yaml:"slug"`
	Name           string   `yaml:"name"`
	Description    string   `yaml:"description"`
	RoleDefinition string   `yaml:"roleDefinition"`
	WhenToUse      string   `yaml:"whenToUse"`
	Groups         []string `yaml:"groups"`
}
