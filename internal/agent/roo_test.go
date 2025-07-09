package agent

import "testing"

func TestRoo_ID_Name(t *testing.T) {
	r := &Roo{}
	if r.ID() != "roo" {
		t.Errorf("expected ID 'roo', got %q", r.ID())
	}
	if r.Name() != "Roo" {
		t.Errorf("expected Name 'Roo', got %q", r.Name())
	}
}
