package frontmatter

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParse_Basic(t *testing.T) {
	input := []byte(`---
title: Test Title
description: Test Description
---
Content goes here`)

	expectedFrontmatter := map[string]interface{}{
		"title":       "Test Title",
		"description": "Test Description",
	}
	expectedContent := "Content goes here"

	frontmatter, content, err := Parse(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(frontmatter, expectedFrontmatter) {
		t.Errorf("expected frontmatter %v, got %v", expectedFrontmatter, frontmatter)
	}

	if content != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, content)
	}
}

func TestParse_NoFrontmatter(t *testing.T) {
	input := []byte(`Content only, no frontmatter`)

	expectedFrontmatter := map[string]interface{}{}
	expectedContent := "Content only, no frontmatter"

	frontmatter, content, err := Parse(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(frontmatter, expectedFrontmatter) {
		t.Errorf("expected frontmatter %v, got %v", expectedFrontmatter, frontmatter)
	}

	if content != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, content)
	}
}

func TestParse_MalformedFrontmatter(t *testing.T) {
	input := []byte(`---
title: "Unclosed quote
---
Content`)

	_, _, err := Parse(input)
	if err == nil {
		t.Fatal("expected error for malformed frontmatter, got nil")
	}
}

func TestParse_MissingClosingDelimiter(t *testing.T) {
	input := []byte(`---
title: Test
Content without closing delimiter`)

	_, _, err := Parse(input)
	if err == nil {
		t.Fatal("expected error for missing closing delimiter, got nil")
	}
}

func TestParse_NestedMap(t *testing.T) {
	input := []byte(`---
roo:
  slug: test-command
  name: Test Command
  roleDefinition: A test command
  whenToUse: For testing
  groups:
    - group1
    - group2
---
Command content here`)

	expectedContent := "Command content here"

	frontmatter, content, err := Parse(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	roo, ok := frontmatter["roo"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected roo to be map[string]interface{}, got %T", frontmatter["roo"])
	}

	if roo["slug"] != "test-command" {
		t.Errorf("expected roo.slug %q, got %q", "test-command", roo["slug"])
	}

	if roo["name"] != "Test Command" {
		t.Errorf("expected roo.name %q, got %q", "Test Command", roo["name"])
	}

	if content != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, content)
	}
}

func TestExtractFrontmatter_Basic(t *testing.T) {
	input := []byte(`---
title: Test Title
description: Test Description
---
Content goes here`)

	expectedFrontmatter := []byte(`title: Test Title
description: Test Description`)
	expectedContent := []byte(`Content goes here`)

	frontmatter, content, err := ExtractFrontmatter(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !bytes.Equal(frontmatter, expectedFrontmatter) {
		t.Errorf("expected frontmatter %q, got %q", expectedFrontmatter, frontmatter)
	}

	if !bytes.Equal(content, expectedContent) {
		t.Errorf("expected content %q, got %q", expectedContent, content)
	}
}

func TestExtractFrontmatter_NoFrontmatter(t *testing.T) {
	input := []byte(`Content only, no frontmatter`)

	expectedFrontmatter := []byte{}
	expectedContent := []byte(`Content only, no frontmatter`)

	frontmatter, content, err := ExtractFrontmatter(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !bytes.Equal(frontmatter, expectedFrontmatter) {
		t.Errorf("expected frontmatter %q, got %q", expectedFrontmatter, frontmatter)
	}

	if !bytes.Equal(content, expectedContent) {
		t.Errorf("expected content %q, got %q", expectedContent, content)
	}
}

func TestExtractFrontmatter_MissingClosingDelimiter(t *testing.T) {
	input := []byte(`---
title: Test
Content without closing delimiter`)

	_, _, err := ExtractFrontmatter(input)
	if err == nil {
		t.Fatal("expected error for missing closing delimiter, got nil")
	}
}

func TestExtractFrontmatter_NestedMap(t *testing.T) {
	input := []byte(`---
roo:
  slug: test-command
  name: Test Command
  roleDefinition: A test command
  whenToUse: For testing
  groups:
    - group1
    - group2
---
Command content here`)

	expectedContent := []byte(`Command content here`)

	frontmatter, content, err := ExtractFrontmatter(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// We're just testing extraction here, not parsing
	if !bytes.Contains(frontmatter, []byte("roo:")) {
		t.Errorf("frontmatter should contain 'roo:' but got: %q", frontmatter)
	}

	if !bytes.Equal(content, expectedContent) {
		t.Errorf("expected content %q, got %q", expectedContent, content)
	}
}
