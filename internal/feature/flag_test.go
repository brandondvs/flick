package feature

import "testing"

func TestCreate(t *testing.T) {
	f := Create("dark-mode")

	if f.Name() != "dark-mode" {
		t.Errorf("expected name %q, got %q", "dark-mode", f.Name())
	}
	if f.IsEnabled() {
		t.Error("expected new feature to be disabled by default")
	}
}

func TestSet(t *testing.T) {
	f := Create("notifications")

	f.Set(true)
	if !f.IsEnabled() {
		t.Error("expected feature to be enabled after Set(true)")
	}

	f.Set(false)
	if f.IsEnabled() {
		t.Error("expected feature to be disabled after Set(false)")
	}
}

func TestToggle(t *testing.T) {
	f := Create("beta-ui")

	f.Toggle()
	if !f.IsEnabled() {
		t.Error("expected feature to be enabled after first toggle")
	}

	f.Toggle()
	if f.IsEnabled() {
		t.Error("expected feature to be disabled after second toggle")
	}
}

func TestToggleMultiple(t *testing.T) {
	f := Create("experiment")

	for i := 0; i < 10; i++ {
		before := f.IsEnabled()
		f.Toggle()
		if f.IsEnabled() == before {
			t.Fatalf("toggle %d did not change state", i+1)
		}
	}
}

func TestSetIdempotent(t *testing.T) {
	f := Create("cache")

	f.Set(true)
	f.Set(true)
	if !f.IsEnabled() {
		t.Error("expected feature to remain enabled after setting true twice")
	}

	f.Set(false)
	f.Set(false)
	if f.IsEnabled() {
		t.Error("expected feature to remain disabled after setting false twice")
	}
}

func TestCreateMultipleIndependent(t *testing.T) {
	a := Create("feature-a")
	b := Create("feature-b")

	a.Set(true)

	if !a.IsEnabled() {
		t.Error("feature-a should be enabled")
	}
	if b.IsEnabled() {
		t.Error("feature-b should still be disabled")
	}
}
