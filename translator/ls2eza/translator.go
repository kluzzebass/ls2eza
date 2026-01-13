package ls2eza

import (
	"runtime"
	"strings"

	"github.com/kluzzebass/reflag/translator"
)

func init() {
	translator.Register(&Translator{})
}

// Translator implements the ls to eza flag translation
type Translator struct{}

func (t *Translator) Name() string        { return "ls2eza" }
func (t *Translator) SourceTool() string  { return "ls" }
func (t *Translator) TargetTool() string  { return "eza" }
func (t *Translator) IncludeInInit() bool { return true }

// Translate converts ls arguments to eza arguments
func (t *Translator) Translate(args []string, mode string) []string {
	return translateFlags(args, getLSMode(mode))
}

// LSMode determines which ls flavor to emulate
type LSMode int

const (
	ModeBSD LSMode = iota
	ModeGNU
)

// getLSMode returns the ls compatibility mode based on mode string or OS detection
func getLSMode(mode string) LSMode {
	switch strings.ToLower(mode) {
	case "bsd":
		return ModeBSD
	case "gnu":
		return ModeGNU
	}

	// Auto-detect based on OS
	switch runtime.GOOS {
	case "darwin", "freebsd", "openbsd", "netbsd", "dragonfly":
		return ModeBSD
	default:
		return ModeGNU
	}
}

// Flags that need --reverse in eza to match ls default behavior
var reverseNeeded = map[rune]bool{
	't': true, // time sort
	'S': true, // size sort
	'c': true, // change time sort
	'u': true, // access time sort
	'U': true, // creation time sort (BSD)
}

// Simple 1:1 flag mappings
var flagMap = map[rune][]string{
	// Display format
	'l': {"-l"},        // long format
	'1': {"-1"},        // one entry per line
	'C': {"--grid"},    // multi-column output
	'x': {"--across"},  // sort grid across
	'm': {"--oneline"}, // stream output

	// Show/hide entries
	'a': {"-a"},        // show all including . and ..
	'A': {"-A"},        // show hidden but not . and ..
	'd': {"-d"},        // list directories themselves
	'R': {"--recurse"}, // recurse into directories

	// Sorting
	't': {"--sort=modified"},   // sort by modification time
	'S': {"--sort=size"},       // sort by size
	'c': {"--sort=changed"},    // sort by change time
	'u': {"--sort=accessed"},   // sort by access time
	'U': {"--sort=created"},    // sort by creation time (BSD)
	'f': {"--sort=none", "-a"}, // unsorted, show all
	'v': {"--sort=name"},       // natural version sort

	// File size display
	'h': {},              // human-readable (default in eza)
	'k': {},              // 1024-byte blocks
	's': {"--blocksize"}, // show allocated blocks

	// Indicators and classification
	'F': {"-F"},         // append file type indicators
	'p': {"--classify"}, // append / to directories

	// Long format options
	'i': {"--inode"},          // show inode numbers
	'n': {"--numeric"},        // numeric user/group IDs
	'o': {"-l", "--no-group"}, // long format without group (BSD)
	'g': {"-l", "--no-user"},  // long format without owner
	'O': {"--flags"},          // show file flags (BSD/macOS)
	'e': {},                   // show ACL (no eza equivalent)
	'@': {"--extended"},       // show extended attributes

	// Symlink handling
	'L': {"-X"}, // dereference symlinks
	'H': {"-X"}, // follow symlinks on command line
	'P': {},     // don't follow symlinks (default)

	// Color
	'G': {}, // color output (default in eza)

	// Misc BSD
	'q': {}, // replace non-printable with ?
	'b': {}, // C-style escapes
	'B': {}, // octal escapes (BSD) / ignore-backups (GNU)
	'W': {}, // display whiteouts (BSD)

	// GNU ls specific
	'Z': {"-Z"},          // SELinux security context
	'N': {"--no-quotes"}, // print names without quoting
	'Q': {},              // quote names
}

// Long option mappings
var longFlagMap = map[string][]string{
	"--all":             {"-a"},
	"--almost-all":      {"-A"},
	"--directory":       {"-d"},
	"--recursive":       {"--recurse"},
	"--human-readable":  {},
	"--inode":           {"--inode"},
	"--numeric-uid-gid": {"--numeric"},
	"--classify":        {"-F"},
	"--file-type":       {"--classify"},
	"--dereference":     {"-X"},
	"--no-group":        {"--no-group"},

	"--group-directories-first": {"--group-directories-first"},
	"--reverse":                 {"--reverse"},
	"--size":                    {"--blocksize"},
	"--context":                 {"-Z"},
	"--literal":                 {"--no-quotes"},
	"--quote-name":              {},
	"--hide-control-chars":      {},
	"--show-control-chars":      {},
	"--hyperlink":               {"--hyperlink"},
	"--full-time":               {"-l", "--time-style=full-iso"},
	"--author":                  {},
	"--escape":                  {},
	"--ignore-backups":          {},
	"--kibibytes":               {},
	"--si":                      {},
	"--dired":                   {},
	"--zero":                    {},
}

// Long options with =value that need prefix matching
var longFlagPrefixes = []struct {
	prefix string
	pass   bool
}{
	{"--color", true},
	{"--colour", true},
	{"--sort=", true},
	{"--time=", true},
	{"--time-style=", true},
	{"--hyperlink=", true},
	{"--width=", true},
	{"--ignore=", true},
	{"--hide=", false},
	{"--block-size=", false},
	{"--indicator-style=", false},
	{"--quoting-style=", false},
	{"--tabsize=", false},
}

func translateFlags(args []string, mode LSMode) []string {
	var ezaArgs []string
	var paths []string
	userReverse := false
	needsReverse := false
	skipNext := false

	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		if strings.HasPrefix(arg, "--") {
			if arg == "--reverse" {
				userReverse = true
				continue
			}

			handled := false
			for _, pf := range longFlagPrefixes {
				if strings.HasPrefix(arg, pf.prefix) {
					if pf.pass {
						if pattern, ok := strings.CutPrefix(arg, "--ignore="); ok {
							ezaArgs = append(ezaArgs, "--ignore-glob="+pattern)
						} else {
							ezaArgs = append(ezaArgs, arg)
						}
					}
					handled = true
					break
				}
			}
			if handled {
				continue
			}

			if mapped, ok := longFlagMap[arg]; ok {
				ezaArgs = append(ezaArgs, mapped...)
			} else {
				ezaArgs = append(ezaArgs, arg)
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			flags := arg[1:]
			for j, c := range flags {
				if c == 'r' {
					userReverse = true
					continue
				}
				if c == 'D' {
					if mode == ModeBSD {
						remaining := flags[j+1:]
						var format string
						if len(remaining) > 0 {
							format = string(remaining)
						} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
							format = args[i+1]
							skipNext = true
						}
						if format != "" {
							ezaArgs = append(ezaArgs, "--time-style=+"+format)
						}
						break
					}
					continue
				}
				if c == 'I' {
					remaining := flags[j+1:]
					if mode == ModeGNU {
						var pattern string
						if len(remaining) > 0 {
							pattern = string(remaining)
						} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
							pattern = args[i+1]
							skipNext = true
						}
						if pattern != "" {
							ezaArgs = append(ezaArgs, "--ignore-glob="+pattern)
						}
						break
					}
					continue
				}
				if c == 'w' {
					if mode == ModeGNU {
						remaining := flags[j+1:]
						var width string
						if len(remaining) > 0 {
							width = string(remaining)
						} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
							width = args[i+1]
							skipNext = true
						}
						if width != "" {
							ezaArgs = append(ezaArgs, "--width="+width)
						}
						break
					}
					continue
				}
				if c == 'T' {
					if mode == ModeBSD {
						ezaArgs = append(ezaArgs, "--time-style=full-iso")
						continue
					}
					remaining := flags[j+1:]
					if len(remaining) > 0 {
						break
					} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						skipNext = true
					}
					continue
				}
				if c == 'X' {
					if mode == ModeGNU {
						ezaArgs = append(ezaArgs, "--sort=extension")
					}
					continue
				}
				if reverseNeeded[c] {
					needsReverse = true
				}
				if mapped, ok := flagMap[c]; ok {
					ezaArgs = append(ezaArgs, mapped...)
				} else {
					ezaArgs = append(ezaArgs, "-"+string(c))
				}
			}
		} else {
			paths = append(paths, arg)
		}
	}

	if needsReverse != userReverse {
		ezaArgs = append(ezaArgs, "--reverse")
	}

	seen := make(map[string]bool)
	var deduped []string
	for _, f := range ezaArgs {
		if !seen[f] {
			seen[f] = true
			deduped = append(deduped, f)
		}
	}

	return append(deduped, paths...)
}
