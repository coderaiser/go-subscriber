package adapters

import "testing"

func TestNewMemory(t *testing.T) {
	m := NewMemory[string]()
	if m == nil {
		t.Fatal("expected non-nil adapter")
	}
}

func TestMemorySetAndGet(t *testing.T) {
	m := NewMemory[string]()
	m.Set("k", "v")
	ptr, err := m.Get("k")
	if err != nil {
		t.Fatal(err)
	}
	if ptr == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *ptr != "v" {
		t.Errorf("want v, got %s", *ptr)
	}
}

func TestMemoryGetNotFound(t *testing.T) {
	m := NewMemory[string]()
	ptr, err := m.Get("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ptr != nil {
		t.Errorf("expected nil, got %v", ptr)
	}
}
