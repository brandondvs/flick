package store

import (
	"sort"
	"sync"
	"testing"

	"github.com/brandondvs/flick/internal/feature"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("expected non-nil Memory")
	}
	if m.data == nil {
		t.Fatal("expected data map to be initialized")
	}
}

func TestStoreAndGet(t *testing.T) {
	m := New()
	flag := feature.Create("dark-mode")

	m.Store("dark-mode", flag)

	got := m.Get("dark-mode")
	if got == nil {
		t.Fatal("expected to retrieve stored flag")
	}
	if got.Name() != "dark-mode" {
		t.Errorf("expected name %q, got %q", "dark-mode", got.Name())
	}
}

func TestGetMissing(t *testing.T) {
	m := New()

	got := m.Get("nonexistent")
	if got != nil {
		t.Errorf("expected nil for missing key, got %v", got)
	}
}

func TestDelete(t *testing.T) {
	m := New()
	flag := feature.Create("beta-ui")

	m.Store("beta-ui", flag)
	m.Delete("beta-ui")

	got := m.Get("beta-ui")
	if got != nil {
		t.Errorf("expected nil after delete, got %v", got)
	}
}

func TestDeleteMissing(t *testing.T) {
	m := New()

	// should not panic
	m.Delete("nonexistent")
}

func TestStoreOverwrite(t *testing.T) {
	m := New()
	first := feature.Create("cache")
	second := feature.Create("cache")
	second.Set(true)

	m.Store("cache", first)
	m.Store("cache", second)

	got := m.Get("cache")
	if got == nil {
		t.Fatal("expected to retrieve flag")
	}
	if !got.IsEnabled() {
		t.Error("expected overwritten flag to be enabled")
	}
}

func TestMultipleKeys(t *testing.T) {
	m := New()
	a := feature.Create("feature-a")
	b := feature.Create("feature-b")
	c := feature.Create("feature-c")

	m.Store("feature-a", a)
	m.Store("feature-b", b)
	m.Store("feature-c", c)

	if m.Get("feature-a") == nil {
		t.Error("expected feature-a to exist")
	}
	if m.Get("feature-b") == nil {
		t.Error("expected feature-b to exist")
	}
	if m.Get("feature-c") == nil {
		t.Error("expected feature-c to exist")
	}

	m.Delete("feature-b")

	if m.Get("feature-b") != nil {
		t.Error("expected feature-b to be deleted")
	}
	if m.Get("feature-a") == nil {
		t.Error("expected feature-a to still exist")
	}
	if m.Get("feature-c") == nil {
		t.Error("expected feature-c to still exist")
	}
}

func TestAllKeysEmpty(t *testing.T) {
	m := New()

	keys := m.AllKeys()
	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}

func TestAllKeys(t *testing.T) {
	m := New()
	m.Store("beta", feature.Create("beta"))
	m.Store("cache", feature.Create("cache"))
	m.Store("auth", feature.Create("auth"))

	keys := m.AllKeys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}

	sort.Strings(keys)
	expected := []string{"auth", "beta", "cache"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("expected key %q at index %d, got %q", expected[i], i, k)
		}
	}
}

func TestAllKeysAfterDelete(t *testing.T) {
	m := New()
	m.Store("one", feature.Create("one"))
	m.Store("two", feature.Create("two"))
	m.Store("three", feature.Create("three"))

	m.Delete("two")

	keys := m.AllKeys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys after delete, got %d", len(keys))
	}

	sort.Strings(keys)
	expected := []string{"one", "three"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("expected key %q at index %d, got %q", expected[i], i, k)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	m := New()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			flag := feature.Create("flag")
			m.Store("flag", flag)
		}(i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Get("flag")
		}()
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.AllKeys()
		}()
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Delete("flag")
		}()
	}

	wg.Wait()
}
