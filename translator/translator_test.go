package translator

import (
	"testing"
)

// mockTranslator is a test implementation of the Translator interface
type mockTranslator struct {
	name          string
	source        string
	target        string
	includeInInit bool
	translateFn   func([]string, string) []string
}

func (m *mockTranslator) Name() string        { return m.name }
func (m *mockTranslator) SourceTool() string  { return m.source }
func (m *mockTranslator) TargetTool() string  { return m.target }
func (m *mockTranslator) IncludeInInit() bool { return m.includeInInit }
func (m *mockTranslator) Translate(args []string, mode string) []string {
	if m.translateFn != nil {
		return m.translateFn(args, mode)
	}
	return args
}

func TestTranslatorIncludeInInit(t *testing.T) {
	tests := []struct {
		name          string
		includeInInit bool
	}{
		{
			name:          "excluded from init",
			includeInInit: false,
		},
		{
			name:          "included in init",
			includeInInit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &mockTranslator{
				name:          "test2test",
				source:        "test",
				target:        "test",
				includeInInit: tt.includeInInit,
			}

			if got := tr.IncludeInInit(); got != tt.includeInInit {
				t.Errorf("IncludeInInit() = %v, want %v", got, tt.includeInInit)
			}
		})
	}
}

func TestTranslatorInterface(t *testing.T) {
	tr := &mockTranslator{
		name:          "test2test",
		source:        "test",
		target:        "test",
		includeInInit: true,
		translateFn: func(args []string, mode string) []string {
			return append([]string{"--translated"}, args...)
		},
	}

	t.Run("Name", func(t *testing.T) {
		if got := tr.Name(); got != "test2test" {
			t.Errorf("Name() = %v, want %v", got, "test2test")
		}
	})

	t.Run("SourceTool", func(t *testing.T) {
		if got := tr.SourceTool(); got != "test" {
			t.Errorf("SourceTool() = %v, want %v", got, "test")
		}
	})

	t.Run("TargetTool", func(t *testing.T) {
		if got := tr.TargetTool(); got != "test" {
			t.Errorf("TargetTool() = %v, want %v", got, "test")
		}
	})

	t.Run("IncludeInInit", func(t *testing.T) {
		if got := tr.IncludeInInit(); got != true {
			t.Errorf("IncludeInInit() = %v, want %v", got, true)
		}
	})

	t.Run("Translate", func(t *testing.T) {
		args := []string{"arg1", "arg2"}
		want := []string{"--translated", "arg1", "arg2"}
		if got := tr.Translate(args, ""); !equalSlices(got, want) {
			t.Errorf("Translate() = %v, want %v", got, want)
		}
	})
}

func TestRegistryWithIncludeInInit(t *testing.T) {
	// Create test translators
	includedInInit := &mockTranslator{
		name:          "included2test",
		source:        "included",
		target:        "test",
		includeInInit: true,
	}
	excludedFromInit := &mockTranslator{
		name:          "excluded2test",
		source:        "excluded",
		target:        "test",
		includeInInit: false,
	}

	// Register them
	Register(includedInInit)
	Register(excludedFromInit)

	t.Run("GetByName returns correct IncludeInInit status", func(t *testing.T) {
		tr := GetByName("included2test")
		if tr == nil {
			t.Fatal("GetByName returned nil for included2test")
		}
		if !tr.IncludeInInit() {
			t.Error("included2test should be included in init")
		}

		tr = GetByName("excluded2test")
		if tr == nil {
			t.Fatal("GetByName returned nil for excluded2test")
		}
		if tr.IncludeInInit() {
			t.Error("excluded2test should be excluded from init")
		}
	})

	t.Run("Get returns correct IncludeInInit status", func(t *testing.T) {
		tr := Get("included", "test")
		if tr == nil {
			t.Fatal("Get returned nil for included -> test")
		}
		if !tr.IncludeInInit() {
			t.Error("included2test should be included in init")
		}

		tr = Get("excluded", "test")
		if tr == nil {
			t.Fatal("Get returned nil for excluded -> test")
		}
		if tr.IncludeInInit() {
			t.Error("excluded2test should be excluded from init")
		}
	})

	t.Run("List includes both included and excluded translators", func(t *testing.T) {
		names := List()
		hasIncluded := false
		hasExcluded := false
		for _, name := range names {
			if name == "included2test" {
				hasIncluded = true
			}
			if name == "excluded2test" {
				hasExcluded = true
			}
		}
		if !hasIncluded {
			t.Error("List should include included2test")
		}
		if !hasExcluded {
			t.Error("List should include excluded2test")
		}
	})
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
