package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

// ValidationIssue represents a single schema validation issue.
type ValidationIssue struct {
	InstancePath string
	Message      string
}

// ValidationError aggregates JSON Schema validation issues for a config file.
type ValidationError struct {
	SourcePath string
	Issues     []ValidationIssue
}

func (e *ValidationError) Error() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "Validation failed for: %s\n", e.SourcePath)
	for _, iss := range e.Issues {
		fmt.Fprintf(&b, "- %s: %s\n", iss.InstancePath, iss.Message)
	}
	return b.String()
}

// ValidateConfigFile validates the provided YAML configuration file against the JSON Schema.
//
// Behavior:
// 1) Resolves configPath to a concrete YAML file (if a directory, probes agent-sync.yml then agent-sync.yaml)
// 2) Resolves schemaPath if empty via default locations (repo CWD, then alongside config file)
// 3) Loads schema via jsonschema/v5 and validates YAML (converted to JSON) against it
//
// Returns nil on success. Returns *ValidationError on schema violations.
// Returns other error types for I/O or setup errors.
func ValidateConfigFile(configPath string) error {
	resolvedConfigPath, err := resolveConfigPath(configPath)
	if err != nil {
		return err
	}

	// Always use embedded schema; no local filesystem schema loading.
	if len(embeddedSchema) == 0 {
		return errors.New("embedded schema is empty")
	}

	compiler := jsonschema.NewCompiler()
	const schemaURL = "mem://schema.json"
	if err := compiler.AddResource(schemaURL, bytes.NewReader(embeddedSchema)); err != nil {
		return fmt.Errorf("failed to add embedded schema resource: %w", err)
	}

	schema, err := compiler.Compile(schemaURL)
	if err != nil {
		return fmt.Errorf("failed to compile schema: %w", err)
	}

	yamlBytes, err := os.ReadFile(resolvedConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Convert YAML to JSON using goccy/go-yaml by decoding YAML into interface{} then encoding to JSON.
	var intermediate interface{}
	if err := yaml.Unmarshal(yamlBytes, &intermediate); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}
	jsonBytes, err := json.Marshal(intermediate)
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	// Decode JSON into any using UseNumber to avoid float coercion
	dec := json.NewDecoder(bytes.NewReader(jsonBytes))
	dec.UseNumber()
	var data any
	if err := dec.Decode(&data); err != nil {
		return fmt.Errorf("failed to decode JSON for validation: %w", err)
	}

	// Validate
	if err := schema.Validate(data); err != nil {
		// Convert jsonschema error tree into aggregated issues
		if verr, ok := err.(*jsonschema.ValidationError); ok {
			issues := flattenValidationErrors(verr)
			return &ValidationError{
				SourcePath: resolvedConfigPath,
				Issues:     issues,
			}
		}
		return fmt.Errorf("schema validation error: %w", err)
	}

	return nil
}

// resolveConfigPath normalizes a config path to a specific YAML file.
// If given a directory, it probes for agent-sync.yml then agent-sync.yaml.
func resolveConfigPath(p string) (string, error) {
	path := filepath.Clean(p)
	if path == "" {
		path = "."
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("config path not found: %s: %w", path, err)
	}
	if info.IsDir() {
		candidates := []string{
			filepath.Join(path, "agent-sync.yml"),
			filepath.Join(path, "agent-sync.yaml"),
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				return c, nil
			}
		}
		return "", fmt.Errorf("no agent-sync.yml or agent-sync.yaml found in %s", path)
	}
	// It's a file; return as-is
	return path, nil
}

// flattenValidationErrors flattens the jsonschema.ValidationError tree into a list of issues.
func flattenValidationErrors(root *jsonschema.ValidationError) []ValidationIssue {
	var issues []ValidationIssue
	var walk func(e *jsonschema.ValidationError)
	walk = func(e *jsonschema.ValidationError) {
		inst := e.InstanceLocation
		if inst == "" {
			inst = "/"
		}
		issues = append(issues, ValidationIssue{
			InstancePath: inst,
			Message:      e.Message,
		})
		for _, c := range e.Causes {
			walk(c)
		}
	}
	walk(root)
	return issues
}
