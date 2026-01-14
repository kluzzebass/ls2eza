package df2duf

import (
	"reflect"
	"testing"
)

func TestTranslateFlags(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		// Basic usage
		{
			name:     "no args",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "path only",
			input:    []string{"/tmp"},
			expected: []string{},
		},
		{
			name:     "multiple paths",
			input:    []string{"/tmp", "/var"},
			expected: []string{},
		},

		// Human readable ignored (duf default)
		{
			name:     "human readable ignored",
			input:    []string{"-h", "/tmp"},
			expected: []string{},
		},
		{
			name:     "human readable long ignored",
			input:    []string{"--human-readable", "/tmp"},
			expected: []string{},
		},

		// User's specific case: du -lh
		{
			name:     "du -lh (count links, human-readable)",
			input:    []string{"-lh"},
			expected: []string{},
		},
		{
			name:     "du -lh with path",
			input:    []string{"-lh", "/tmp"},
			expected: []string{},
		},

		// All filesystems
		{
			name:     "all files short",
			input:    []string{"-a", "/tmp"},
			expected: []string{"-all"},
		},
		{
			name:     "all files long",
			input:    []string{"--all", "/tmp"},
			expected: []string{"-all"},
		},

		// Inodes
		{
			name:     "inodes long",
			input:    []string{"--inodes", "/tmp"},
			expected: []string{"-inodes"},
		},

		// Exclude patterns
		{
			name:     "exclude pattern",
			input:    []string{"--exclude=*.tmp", "/tmp"},
			expected: []string{"-hide-mp", "*.tmp"},
		},
		{
			name:     "BSD exclude pattern short",
			input:    []string{"-I", "*.log", "/tmp"},
			expected: []string{"-hide-mp", "*.log"},
		},
		{
			name:     "BSD exclude pattern attached",
			input:    []string{"-I*.log", "/tmp"},
			expected: []string{"-hide-mp", "*.log"},
		},

		// Ignored flags
		{
			name:     "summarize ignored",
			input:    []string{"-s", "/tmp"},
			expected: []string{},
		},
		{
			name:     "max depth ignored",
			input:    []string{"-d", "2", "/tmp"},
			expected: []string{},
		},
		{
			name:     "threshold ignored",
			input:    []string{"-t", "1000", "/tmp"},
			expected: []string{},
		},
		{
			name:     "block size ignored",
			input:    []string{"-B", "1M", "/tmp"},
			expected: []string{},
		},

		// Combined flags
		{
			name:     "combined flags",
			input:    []string{"-ah", "/tmp"},
			expected: []string{"-all"},
		},
		{
			name:     "combined with ignored",
			input:    []string{"-sha", "/tmp"},
			expected: []string{"-all"},
		},

		// One file system
		{
			name:     "one file system short",
			input:    []string{"-x", "/"},
			expected: []string{},
		},
		{
			name:     "one file system long",
			input:    []string{"--one-file-system", "/"},
			expected: []string{},
		},

		// Symlinks (all ignored for duf)
		{
			name:     "follow symlinks",
			input:    []string{"-L", "/tmp"},
			expected: []string{},
		},
		{
			name:     "no follow symlinks",
			input:    []string{"-P", "/tmp"},
			expected: []string{},
		},

		// Size units ignored (duf is human-readable by default)
		{
			name:     "kilobytes ignored",
			input:    []string{"-k", "/tmp"},
			expected: []string{},
		},
		{
			name:     "megabytes ignored",
			input:    []string{"-m", "/tmp"},
			expected: []string{},
		},
		{
			name:     "gigabytes ignored",
			input:    []string{"-g", "/tmp"},
			expected: []string{},
		},

		// Total ignored (duf always shows summary)
		{
			name:     "total ignored",
			input:    []string{"-c", "/tmp"},
			expected: []string{},
		},

		// Real-world examples
		{
			name:     "common du -sh",
			input:    []string{"-sh", "/var/log"},
			expected: []string{},
		},
		{
			name:     "common du -ah",
			input:    []string{"-ah", "."},
			expected: []string{"-all"},
		},
		{
			name:     "du with multiple flags",
			input:    []string{"-shc", "/tmp", "/var"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translateFlags(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("translateFlags(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTranslator(t *testing.T) {
	tr := &Translator{}

	if tr.Name() != "df2duf" {
		t.Errorf("Name() = %v, want df2duf", tr.Name())
	}

	if tr.SourceTool() != "df" {
		t.Errorf("SourceTool() = %v, want df", tr.SourceTool())
	}

	if tr.TargetTool() != "duf" {
		t.Errorf("TargetTool() = %v, want duf", tr.TargetTool())
	}

	if tr.IncludeInInit() {
		t.Errorf("IncludeInInit() = true, want false")
	}

	// Test Translate method
	result := tr.Translate([]string{"-lh", "/tmp"}, "")
	expected := []string{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Translate(['-lh', '/tmp'], '') = %v, want %v", result, expected)
	}
}
