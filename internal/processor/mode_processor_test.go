package processor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uphy/agent-sync/internal/agent"
	"github.com/uphy/agent-sync/internal/util"
	"go.uber.org/zap"
)

// TestModeProcessor_Process verifies Roo modes aggregate to a single file by default:
// - Project scope: ".roomodes"
// - User scope: VS Code globalStorage custom_modes.yaml (macOS path suffix check)
func TestModeProcessor_Process(t *testing.T) {
	// Use the real filesystem adapter interface provided via Pipeline/BaseProcessor pathway.
	// Here we directly construct a ModeProcessor with a BaseProcessor using the real util.RealFileSystem
	// and a temporary directory as AbsInputRoot. However, ModeProcessor.Process reads inputs by path,
	// so we can simulate by writing inputs under a temp dir and passing those paths.

	// Minimal Roo mode bodies to exercise aggregation
	mode1 := `---
roo:
  slug: code-reviewer
  name: Code Reviewer
  description: Reviews code
  roleDefinition: You review code.
  whenToUse: When code needs review
  groups: []
---
Body 1
`
	mode2 := `---
roo:
  slug: test-runner
  name: Test Runner
  description: Runs tests
  roleDefinition: You run tests.
  whenToUse: When tests should be executed
  groups: []
---
Body 2
`

	// Helper to create inputs on disk under a temp dir and return absolute paths and a BaseProcessor
	setup := func(t *testing.T) (absRoot string, in1 string, in2 string, base *BaseProcessor) {
		t.Helper()
		dir := t.TempDir()
		absRoot = dir

		// Create files
		in1 = filepath.Join(dir, "modes", "mode1.md")
		in2 = filepath.Join(dir, "modes", "mode2.md")

		if err := os.MkdirAll(filepath.Dir(in1), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(in1, []byte(mode1), 0o644); err != nil {
			t.Fatalf("write %s: %v", in1, err)
		}
		if err := os.WriteFile(in2, []byte(mode2), 0o644); err != nil {
			t.Fatalf("write %s: %v", in2, err)
		}

		// Build a BaseProcessor matching Pipeline usage
		base = NewBaseProcessor(&util.RealFileSystem{}, zap.NewNop(), absRoot, agent.NewRegistry(), false)
		return
	}

	t.Run("project scope aggregates into .roomodes", func(t *testing.T) {
		absRoot, in1, in2, base := setup(t)
		_ = absRoot
		// Ensure project scope on BaseProcessor
		base.userScope = false
		mp := NewModeProcessor(base)

		cfg := &OutputConfig{
			Agent:       &agent.Roo{},
			RelPath:     "",
			IsDirectory: false,
			AgentName:   "roo",
		}
		// Fill default RelPath like production pipeline would
		cfg.RelPath = mp.GetOutputPath(cfg.Agent, cfg.RelPath)

		result, err := mp.Process([]string{in1, in2}, cfg)
		if err != nil {
			t.Fatalf("process error: %v", err)
		}
		if result == nil {
			t.Fatalf("nil result")
		}

		// Expect a single aggregated output file ".roomodes" relative path
		if len(result.Files) != 1 {
			t.Fatalf("expected 1 aggregated file, got %d", len(result.Files))
		}
		gotRel := result.Files[0].relPath
		if gotRel != ".roomodes" {
			t.Fatalf("expected aggregated path '.roomodes', got %q", gotRel)
		}
		content := result.Files[0].Content
		if !strings.Contains(content, "slug: code-reviewer") || !strings.Contains(content, "slug: test-runner") {
			t.Fatalf("aggregated YAML missing expected slugs, got:\n%s", content)
		}
	})

	t.Run("user scope aggregates to VS Code globalStorage custom_modes.yaml (macOS suffix)", func(t *testing.T) {
		absRoot, in1, _, base := setup(t)
		_ = absRoot
		base.userScope = true
		mp := NewModeProcessor(base)

		cfg := &OutputConfig{
			Agent:       &agent.Roo{},
			RelPath:     "",
			IsDirectory: false,
			AgentName:   "roo",
		}
		cfg.RelPath = mp.GetOutputPath(cfg.Agent, cfg.RelPath)

		result, err := mp.Process([]string{in1}, cfg)
		if err != nil {
			t.Fatalf("process error: %v", err)
		}
		if result == nil || len(result.Files) != 1 {
			t.Fatalf("expected 1 output file, got %d", len(result.Files))
		}

		gotRel := result.Files[0].relPath
		wantSuffix := filepath.Join("Library", "Application Support", "Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings", "custom_modes.yaml")
		if !strings.HasSuffix(gotRel, wantSuffix) {
			t.Fatalf("expected user-scope aggregated path to end with %q, got %q", wantSuffix, gotRel)
		}
		if !strings.Contains(result.Files[0].Content, "slug: code-reviewer") {
			t.Fatalf("aggregated YAML missing expected slug content, got:\n%s", result.Files[0].Content)
		}
	})

	// New tests for template path resolution in mode processing

	t.Run("relative include resolves without duplication", func(t *testing.T) {
		dir := t.TempDir()
		// config/base path
		base := NewBaseProcessor(&util.RealFileSystem{}, zap.NewNop(), dir, agent.NewRegistry(), false)
		mp := NewModeProcessor(base)

		// Create fragment under modes/frags/part.md
		fragsDir := filepath.Join(dir, "modes", "frags")
		if err := os.MkdirAll(fragsDir, 0o755); err != nil {
			t.Fatalf("mkdir frags: %v", err)
		}
		fragPath := filepath.Join(fragsDir, "part.md")
		if err := os.WriteFile(fragPath, []byte("INCLUDED_PART"), 0o644); err != nil {
			t.Fatalf("write frag: %v", err)
		}

		// Mode file includes via relative path "./frags/part.md"
		modeBody := `---
roo:
  slug: rel-include
  name: Rel Include
  description: Test
  roleDefinition: Test
  whenToUse: Test
  groups: []
---
{{ include "./frags/part.md" }}
`
		modePath := filepath.Join(dir, "modes", "rel.md")
		if err := os.WriteFile(modePath, []byte(modeBody), 0o644); err != nil {
			t.Fatalf("write mode: %v", err)
		}

		cfg := &OutputConfig{
			Agent:       &agent.Roo{},
			RelPath:     "",
			IsDirectory: false,
			AgentName:   "roo",
		}
		cfg.RelPath = mp.GetOutputPath(cfg.Agent, cfg.RelPath)

		result, err := mp.Process([]string{modePath}, cfg)
		if err != nil {
			t.Fatalf("process error: %v", err)
		}
		if result == nil || len(result.Files) != 1 {
			t.Fatalf("expected aggregated file")
		}
		content := result.Files[0].Content
		// Ensure included content appears and only once (simple duplication guard)
		if !strings.Contains(content, "INCLUDED_PART") {
			t.Fatalf("expected included content")
		}
		if strings.Count(content, "INCLUDED_PART") != 1 {
			t.Fatalf("expected single inclusion, got %d", strings.Count(content, "INCLUDED_PART"))
		}
	})

	t.Run("absolute include path is rejected", func(t *testing.T) {
		dir := t.TempDir()
		base := NewBaseProcessor(&util.RealFileSystem{}, zap.NewNop(), dir, agent.NewRegistry(), false)
		mp := NewModeProcessor(base)

		// Create an absolute path to include
		fragAbs := filepath.Join(dir, "absFrag.md")
		if err := os.WriteFile(fragAbs, []byte("ABS"), 0o644); err != nil {
			t.Fatalf("write abs frag: %v", err)
		}

		modeBody := `---
roo:
  slug: abs-include
  name: Abs Include
  description: Test
  roleDefinition: Test
  whenToUse: Test
  groups: []
---
{{ include "` + filepath.ToSlash(fragAbs) + `" }}
`
		modePath := filepath.Join(dir, "modes", "abs.md")
		if err := os.MkdirAll(filepath.Dir(modePath), 0o755); err != nil {
			t.Fatalf("mkdir modes: %v", err)
		}
		if err := os.WriteFile(modePath, []byte(modeBody), 0o644); err != nil {
			t.Fatalf("write mode: %v", err)
		}

		cfg := &OutputConfig{
			Agent:       &agent.Roo{},
			RelPath:     "",
			IsDirectory: false,
			AgentName:   "roo",
		}
		cfg.RelPath = mp.GetOutputPath(cfg.Agent, cfg.RelPath)

		_, err := mp.Process([]string{modePath}, cfg)
		if err == nil {
			t.Fatalf("expected error for absolute include path")
		}
		wantSub := "absolute paths are not allowed in template include/reference paths"
		if !strings.Contains(err.Error(), wantSub) {
			t.Fatalf("expected error to contain %q, got %v", wantSub, err)
		}
	})

	t.Run("@-prefix include resolves from config base without duplication", func(t *testing.T) {
		dir := t.TempDir()
		// BasePath = dir (config directory)
		base := NewBaseProcessor(&util.RealFileSystem{}, zap.NewNop(), dir, agent.NewRegistry(), false)
		mp := NewModeProcessor(base)

		// Put shared fragment under config base at shared/part.md
		sharedDir := filepath.Join(dir, "shared")
		if err := os.MkdirAll(sharedDir, 0o755); err != nil {
			t.Fatalf("mkdir shared: %v", err)
		}
		sharedFrag := filepath.Join(sharedDir, "part.md")
		if err := os.WriteFile(sharedFrag, []byte("AT_INCLUDED"), 0o644); err != nil {
			t.Fatalf("write shared frag: %v", err)
		}

		// Mode file that includes "@/shared/part.md" from a different subdir
		modeBody := `---
roo:
  slug: at-include
  name: At Include
  description: Test
  roleDefinition: Test
  whenToUse: Test
  groups: []
---
{{ include "@/shared/part.md" }}
`
		modePath := filepath.Join(dir, "modes", "at.md")
		if err := os.MkdirAll(filepath.Dir(modePath), 0o755); err != nil {
			t.Fatalf("mkdir modes: %v", err)
		}
		if err := os.WriteFile(modePath, []byte(modeBody), 0o644); err != nil {
			t.Fatalf("write mode: %v", err)
		}

		cfg := &OutputConfig{
			Agent:       &agent.Roo{},
			RelPath:     "",
			IsDirectory: false,
			AgentName:   "roo",
		}
		cfg.RelPath = mp.GetOutputPath(cfg.Agent, cfg.RelPath)

		result, err := mp.Process([]string{modePath}, cfg)
		if err != nil {
			t.Fatalf("process error: %v", err)
		}
		if result == nil || len(result.Files) != 1 {
			t.Fatalf("expected aggregated file")
		}
		content := result.Files[0].Content
		if !strings.Contains(content, "AT_INCLUDED") {
			t.Fatalf("expected @-included content")
		}
		if strings.Count(content, "AT_INCLUDED") != 1 {
			t.Fatalf("expected single inclusion for @-path, got %d", strings.Count(content, "AT_INCLUDED"))
		}
	})
}
