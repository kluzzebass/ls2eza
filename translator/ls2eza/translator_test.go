package ls2eza

import (
	"reflect"
	"testing"
)

func TestTranslateFlagsGNU(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		// Basic flags
		{
			name:     "long format",
			input:    []string{"-l"},
			expected: []string{"-l"},
		},
		{
			name:     "all files",
			input:    []string{"-a"},
			expected: []string{"-a"},
		},
		{
			name:     "combined la",
			input:    []string{"-la"},
			expected: []string{"-l", "-a"},
		},
		{
			name:     "almost all",
			input:    []string{"-A"},
			expected: []string{"-A"},
		},

		// Sort flags with reverse handling
		{
			name:     "time sort adds reverse",
			input:    []string{"-t"},
			expected: []string{"--sort=modified", "--reverse"},
		},
		{
			name:     "time sort with user reverse cancels",
			input:    []string{"-tr"},
			expected: []string{"--sort=modified"},
		},
		{
			name:     "size sort adds reverse",
			input:    []string{"-S"},
			expected: []string{"--sort=size", "--reverse"},
		},
		{
			name:     "size sort with user reverse cancels",
			input:    []string{"-Sr"},
			expected: []string{"--sort=size"},
		},
		{
			name:     "change time sort adds reverse",
			input:    []string{"-c"},
			expected: []string{"--sort=changed", "--reverse"},
		},
		{
			name:     "access time sort adds reverse",
			input:    []string{"-u"},
			expected: []string{"--sort=accessed", "--reverse"},
		},
		{
			name:     "creation time sort adds reverse",
			input:    []string{"-U"},
			expected: []string{"--sort=created", "--reverse"},
		},
		{
			name:     "user reverse without sort flag",
			input:    []string{"-r"},
			expected: []string{"--reverse"},
		},
		{
			name:     "long reverse without sort",
			input:    []string{"--reverse"},
			expected: []string{"--reverse"},
		},

		// Combined flags
		{
			name:     "long time sort",
			input:    []string{"-lt"},
			expected: []string{"-l", "--sort=modified", "--reverse"},
		},
		{
			name:     "long time sort reversed",
			input:    []string{"-ltr"},
			expected: []string{"-l", "--sort=modified"},
		},
		{
			name:     "long size sort with path",
			input:    []string{"-lS", "/tmp"},
			expected: []string{"-l", "--sort=size", "--reverse", "/tmp"},
		},

		// Flags that map to empty (defaults in eza)
		{
			name:     "human readable is default",
			input:    []string{"-h"},
			expected: nil,
		},
		{
			name:     "color is default",
			input:    []string{"-G"},
			expected: nil,
		},
		{
			name:     "lh only outputs l",
			input:    []string{"-lh"},
			expected: []string{"-l"},
		},

		// Long format options
		{
			name:     "inode",
			input:    []string{"-i"},
			expected: []string{"--inode"},
		},
		{
			name:     "numeric",
			input:    []string{"-n"},
			expected: []string{"--numeric"},
		},
		{
			name:     "no group",
			input:    []string{"-o"},
			expected: []string{"-l", "--no-group"},
		},
		{
			name:     "no user",
			input:    []string{"-g"},
			expected: []string{"-l", "--no-user"},
		},
		{
			name:     "file flags",
			input:    []string{"-O"},
			expected: []string{"--flags"},
		},
		{
			name:     "extended attributes",
			input:    []string{"-@"},
			expected: []string{"--extended"},
		},

		// Display format
		{
			name:     "one per line",
			input:    []string{"-1"},
			expected: []string{"-1"},
		},
		{
			name:     "grid",
			input:    []string{"-C"},
			expected: []string{"--grid"},
		},
		{
			name:     "across",
			input:    []string{"-x"},
			expected: []string{"--across"},
		},
		{
			name:     "recurse",
			input:    []string{"-R"},
			expected: []string{"--recurse"},
		},

		// Indicators
		{
			name:     "classify F",
			input:    []string{"-F"},
			expected: []string{"-F"},
		},
		{
			name:     "classify p",
			input:    []string{"-p"},
			expected: []string{"--classify"},
		},

		// Symlinks
		{
			name:     "dereference L",
			input:    []string{"-L"},
			expected: []string{"-X"},
		},
		{
			name:     "dereference H",
			input:    []string{"-H"},
			expected: []string{"-X"},
		},

		// Unsorted
		{
			name:     "unsorted",
			input:    []string{"-f"},
			expected: []string{"--sort=none", "-a"},
		},

		// Long options
		{
			name:     "long option all",
			input:    []string{"--all"},
			expected: []string{"-a"},
		},
		{
			name:     "long option almost-all",
			input:    []string{"--almost-all"},
			expected: []string{"-A"},
		},
		{
			name:     "long option directory",
			input:    []string{"--directory"},
			expected: []string{"-d"},
		},
		{
			name:     "long option recursive",
			input:    []string{"--recursive"},
			expected: []string{"--recurse"},
		},
		{
			name:     "long option human-readable is default",
			input:    []string{"--human-readable"},
			expected: nil,
		},
		{
			name:     "color passthrough",
			input:    []string{"--color=auto"},
			expected: []string{"--color=auto"},
		},
		{
			name:     "color always passthrough",
			input:    []string{"--color=always"},
			expected: []string{"--color=always"},
		},
		{
			name:     "sort passthrough",
			input:    []string{"--sort=name"},
			expected: []string{"--sort=name"},
		},
		{
			name:     "time passthrough",
			input:    []string{"--time=accessed"},
			expected: []string{"--time=accessed"},
		},

		// Paths
		{
			name:     "single path",
			input:    []string{"/tmp"},
			expected: []string{"/tmp"},
		},
		{
			name:     "multiple paths",
			input:    []string{"/tmp", "/var"},
			expected: []string{"/tmp", "/var"},
		},
		{
			name:     "flags and paths",
			input:    []string{"-la", "/tmp", "/var"},
			expected: []string{"-l", "-a", "/tmp", "/var"},
		},

		// Deduplication
		{
			name:     "dedup repeated flags",
			input:    []string{"-l", "-l"},
			expected: []string{"-l"},
		},
		{
			name:     "dedup from og combo",
			input:    []string{"-og"},
			expected: []string{"-l", "--no-group", "--no-user"},
		},

		// GNU ls specific flags
		{
			name:     "GNU ignore pattern separate",
			input:    []string{"-I", "*.bak"},
			expected: []string{"--ignore-glob=*.bak"},
		},
		{
			name:     "GNU ignore pattern attached",
			input:    []string{"-I*.tmp"},
			expected: []string{"--ignore-glob=*.tmp"},
		},
		{
			name:     "GNU sort by extension",
			input:    []string{"-lX"},
			expected: []string{"-l", "--sort=extension"},
		},
		{
			name:     "GNU width",
			input:    []string{"-w", "80"},
			expected: []string{"--width=80"},
		},
		{
			name:     "GNU width attached",
			input:    []string{"-w120"},
			expected: []string{"--width=120"},
		},
		{
			name:     "GNU SELinux context",
			input:    []string{"-lZ"},
			expected: []string{"-l", "-Z"},
		},
		{
			name:     "GNU literal/no-quotes",
			input:    []string{"-lN"},
			expected: []string{"-l", "--no-quotes"},
		},
		{
			name:     "GNU group directories first",
			input:    []string{"--group-directories-first"},
			expected: []string{"--group-directories-first"},
		},
		{
			name:     "GNU full-time",
			input:    []string{"--full-time"},
			expected: []string{"-l", "--time-style=full-iso"},
		},
		{
			name:     "GNU ignore long option",
			input:    []string{"--ignore=*.log"},
			expected: []string{"--ignore-glob=*.log"},
		},
		{
			name:     "GNU hyperlink",
			input:    []string{"--hyperlink"},
			expected: []string{"--hyperlink"},
		},

		// Ignored flags (no eza equivalent)
		{
			name:     "ignored flag W",
			input:    []string{"-lW"},
			expected: []string{"-l"},
		},
		{
			name:     "ignored flag Q",
			input:    []string{"-lQ"},
			expected: []string{"-l"},
		},

		// Edge cases
		{
			name:     "empty input",
			input:    []string{},
			expected: nil,
		},
		{
			name:     "unknown flag passthrough",
			input:    []string{"-z"},
			expected: []string{"-z"},
		},
		{
			name:     "unknown long option passthrough",
			input:    []string{"--unknown"},
			expected: []string{"--unknown"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translateFlags(tt.input, ModeGNU)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("translateFlags(%v, ModeGNU) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTranslateFlagsBSD(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		// BSD-specific flag behaviors
		{
			name:     "BSD full time",
			input:    []string{"-T"},
			expected: []string{"--time-style=full-iso"},
		},
		{
			name:     "BSD date format separate",
			input:    []string{"-D", "%Y-%m-%d"},
			expected: []string{"--time-style=+%Y-%m-%d"},
		},
		{
			name:     "BSD date format attached",
			input:    []string{"-D%Y-%m-%d"},
			expected: []string{"--time-style=+%Y-%m-%d"},
		},
		{
			name:     "BSD -I ignored",
			input:    []string{"-lI"},
			expected: []string{"-l"},
		},
		{
			name:     "BSD -X ignored",
			input:    []string{"-lX"},
			expected: []string{"-l"},
		},
		{
			name:     "BSD -w ignored",
			input:    []string{"-lw"},
			expected: []string{"-l"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translateFlags(tt.input, ModeBSD)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("translateFlags(%v, ModeBSD) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestVersionFlag(t *testing.T) {
	// -V and --version should not be translated, they're handled in main()
	// But if they somehow get to translateFlags, they should pass through
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{"-V"}, []string{"-V"}},
		{[]string{"--version"}, []string{"--version"}},
	}

	for _, tt := range tests {
		result := translateFlags(tt.input, ModeGNU)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("translateFlags(%v) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestGetLSModeParameter(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		expected LSMode
	}{
		{"bsd lowercase", "bsd", ModeBSD},
		{"BSD uppercase", "BSD", ModeBSD},
		{"gnu lowercase", "gnu", ModeGNU},
		{"GNU uppercase", "GNU", ModeGNU},
		{"mixed case Bsd", "Bsd", ModeBSD},
		{"mixed case Gnu", "Gnu", ModeGNU},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLSMode(tt.mode)
			if result != tt.expected {
				t.Errorf("getLSMode(%q) = %v, want %v", tt.mode, result, tt.expected)
			}
		})
	}
}

func TestTranslatorInterface(t *testing.T) {
	tr := &Translator{}

	if tr.Name() != "ls2eza" {
		t.Errorf("Name() = %q, want %q", tr.Name(), "ls2eza")
	}
	if tr.SourceTool() != "ls" {
		t.Errorf("SourceTool() = %q, want %q", tr.SourceTool(), "ls")
	}
	if tr.TargetTool() != "eza" {
		t.Errorf("TargetTool() = %q, want %q", tr.TargetTool(), "eza")
	}

	// Test translation via interface
	result := tr.Translate([]string{"-la"}, "")
	expected := []string{"-l", "-a"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Translate(-la) = %v, want %v", result, expected)
	}
}
