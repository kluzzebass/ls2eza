package grep2rg

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
			name:     "simple pattern",
			input:    []string{"pattern"},
			expected: []string{"pattern"},
		},
		{
			name:     "pattern with path",
			input:    []string{"pattern", "file.txt"},
			expected: []string{"pattern", "file.txt"},
		},
		{
			name:     "pattern with multiple paths",
			input:    []string{"pattern", "file1.txt", "file2.txt"},
			expected: []string{"pattern", "file1.txt", "file2.txt"},
		},

		// Passthrough flags
		{
			name:     "case insensitive",
			input:    []string{"-i", "pattern"},
			expected: []string{"-i", "pattern"},
		},
		{
			name:     "invert match",
			input:    []string{"-v", "pattern"},
			expected: []string{"-v", "pattern"},
		},
		{
			name:     "word match",
			input:    []string{"-w", "pattern"},
			expected: []string{"-w", "pattern"},
		},
		{
			name:     "line match",
			input:    []string{"-x", "pattern"},
			expected: []string{"-x", "pattern"},
		},
		{
			name:     "count",
			input:    []string{"-c", "pattern"},
			expected: []string{"-c", "pattern"},
		},
		{
			name:     "files with matches",
			input:    []string{"-l", "pattern"},
			expected: []string{"-l", "pattern"},
		},
		{
			name:     "line numbers",
			input:    []string{"-n", "pattern"},
			expected: []string{"-n", "pattern"},
		},
		{
			name:     "combined passthrough",
			input:    []string{"-inl", "pattern"},
			expected: []string{"-i", "-n", "-l", "pattern"},
		},

		// Ignored flags (recursive, etc.)
		{
			name:     "recursive ignored",
			input:    []string{"-r", "pattern", "."},
			expected: []string{"pattern", "."},
		},
		{
			name:     "recursive R ignored",
			input:    []string{"-R", "pattern", "."},
			expected: []string{"pattern", "."},
		},
		{
			name:     "extended regexp ignored",
			input:    []string{"-E", "pattern"},
			expected: []string{"pattern"},
		},
		{
			name:     "combined with ignored",
			input:    []string{"-rni", "pattern", "."},
			expected: []string{"-n", "-i", "pattern", "."},
		},

		// Context flags
		{
			name:     "after context",
			input:    []string{"-A", "3", "pattern"},
			expected: []string{"-A", "3", "pattern"},
		},
		{
			name:     "before context",
			input:    []string{"-B", "3", "pattern"},
			expected: []string{"-B", "3", "pattern"},
		},
		{
			name:     "context",
			input:    []string{"-C", "3", "pattern"},
			expected: []string{"-C", "3", "pattern"},
		},
		{
			name:     "context attached",
			input:    []string{"-A3", "pattern"},
			expected: []string{"-A", "3", "pattern"},
		},

		// Include/exclude
		{
			name:     "include pattern",
			input:    []string{"--include=*.go", "pattern"},
			expected: []string{"-g", "*.go", "pattern"},
		},
		{
			name:     "exclude pattern",
			input:    []string{"--exclude=*.txt", "pattern"},
			expected: []string{"-g", "!*.txt", "pattern"},
		},
		{
			name:     "exclude dir",
			input:    []string{"--exclude-dir=vendor", "pattern"},
			expected: []string{"-g", "!vendor/", "pattern"},
		},
		{
			name:     "exclude dir with slash",
			input:    []string{"--exclude-dir=vendor/", "pattern"},
			expected: []string{"-g", "!vendor/", "pattern"},
		},
		{
			name:     "include separate",
			input:    []string{"--include", "*.go", "pattern"},
			expected: []string{"-g", "*.go", "pattern"},
		},

		// Multiple patterns with -e
		{
			name:     "single -e pattern",
			input:    []string{"-e", "pattern"},
			expected: []string{"pattern"},
		},
		{
			name:     "multiple -e patterns",
			input:    []string{"-e", "foo", "-e", "bar"},
			expected: []string{"-e", "foo", "-e", "bar"},
		},
		{
			name:     "attached -e pattern",
			input:    []string{"-epattern"},
			expected: []string{"pattern"},
		},

		// Null separator
		{
			name:     "null short",
			input:    []string{"-Z", "pattern"},
			expected: []string{"-0", "pattern"},
		},
		{
			name:     "null long",
			input:    []string{"--null", "pattern"},
			expected: []string{"-0", "pattern"},
		},

		// Color
		{
			name:     "color always",
			input:    []string{"--color=always", "pattern"},
			expected: []string{"--color=always", "pattern"},
		},
		{
			name:     "color never",
			input:    []string{"--color=never", "pattern"},
			expected: []string{"--color=never", "pattern"},
		},

		// Fixed strings
		{
			name:     "fixed strings",
			input:    []string{"-F", "pattern"},
			expected: []string{"-F", "pattern"},
		},

		// Pattern from file
		{
			name:     "pattern file",
			input:    []string{"-f", "patterns.txt"},
			expected: []string{"-f", "patterns.txt"},
		},

		// Max count
		{
			name:     "max count",
			input:    []string{"-m", "5", "pattern"},
			expected: []string{"-m", "5", "pattern"},
		},

		// Long options
		{
			name:     "long recursive ignored",
			input:    []string{"--recursive", "pattern"},
			expected: []string{"pattern"},
		},
		{
			name:     "long extended ignored",
			input:    []string{"--extended-regexp", "pattern"},
			expected: []string{"pattern"},
		},

		// Pattern starting with dash
		{
			name:     "pattern with dash",
			input:    []string{"-e", "-pattern"},
			expected: []string{"--", "-pattern"},
		},

		// Complex combinations
		{
			name:     "typical grep usage",
			input:    []string{"-rn", "--include=*.go", "TODO", "."},
			expected: []string{"-n", "-g", "*.go", "TODO", "."},
		},
		{
			name:     "grep with context",
			input:    []string{"-rniA3", "pattern", "src/"},
			expected: []string{"-n", "-i", "-A", "3", "pattern", "src/"},
		},

		// Empty input
		{
			name:     "empty input",
			input:    []string{},
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

func TestTranslatorInterface(t *testing.T) {
	tr := &Translator{}

	if tr.Name() != "grep2rg" {
		t.Errorf("Name() = %q, want %q", tr.Name(), "grep2rg")
	}
	if tr.SourceTool() != "grep" {
		t.Errorf("SourceTool() = %q, want %q", tr.SourceTool(), "grep")
	}
	if tr.TargetTool() != "rg" {
		t.Errorf("TargetTool() = %q, want %q", tr.TargetTool(), "rg")
	}
	if tr.EnvVarName() != "GREP2RG_MODE" {
		t.Errorf("EnvVarName() = %q, want %q", tr.EnvVarName(), "GREP2RG_MODE")
	}
}
