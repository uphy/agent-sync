package model

import (
	"testing"
)

func TestProcessContext_SetsContentAndProcessed(t *testing.T) {
	// Use nil engine; stub should ignore engine for now
	ctx, err := ProcessContext(nil, "raw content")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ctx.Content != "raw content" {
		t.Errorf("expected Content %q, got %q", "raw content", ctx.Content)
	}
	if ctx.Processed != "raw content" {
		t.Errorf("expected Processed %q, got %q", "raw content", ctx.Processed)
	}
	if ctx.Path != "" {
		t.Errorf("expected Path empty, got %q", ctx.Path)
	}
}
