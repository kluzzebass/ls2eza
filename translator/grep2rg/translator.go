package grep2rg

import (
	"strings"

	"github.com/kluzzebass/reflag/translator"
)

func init() {
	translator.Register(&Translator{})
}

// Translator implements the grep to ripgrep flag translation
type Translator struct{}

func (t *Translator) Name() string       { return "grep2rg" }
func (t *Translator) SourceTool() string { return "grep" }
func (t *Translator) TargetTool() string { return "rg" }
func (t *Translator) EnvVarName() string { return "GREP2RG_MODE" }

// Translate converts grep arguments to ripgrep arguments
func (t *Translator) Translate(args []string) []string {
	return translateFlags(args)
}

// Flags that pass through unchanged (same in grep and rg)
var passthroughFlags = map[rune]bool{
	'i': true, // case insensitive
	'v': true, // invert match
	'w': true, // word match
	'x': true, // line match
	'c': true, // count
	'l': true, // files with matches
	'L': true, // files without matches
	'n': true, // line numbers
	'H': true, // print filename
	'h': true, // no filename
	'o': true, // only matching
	'q': true, // quiet
	's': true, // suppress errors
	'F': true, // fixed strings
	'P': true, // perl regex
	'a': true, // text (treat binary as text)
	'U': true, // binary
}

// Flags that take a value and pass through
var passthroughWithValue = map[rune]bool{
	'A': true, // after context
	'B': true, // before context
	'C': true, // context
	'm': true, // max count
	'e': true, // pattern
	'f': true, // file
}

// Flags to ignore (behavior is default in rg or not applicable)
var ignoredFlags = map[rune]bool{
	'r': true, // recursive (rg is recursive by default)
	'R': true, // recursive
	'E': true, // extended regexp (rg default)
	'G': true, // basic regexp (no rg equivalent, close enough)
	'I': true, // skip binary (rg default)
	'b': true, // byte offset (rg has it but different)
	'T': true, // initial tab
	'd': true, // directory handling
	'D': true, // device handling
	'u': true, // unix byte offsets
}

// Long flags that pass through
var longPassthrough = map[string]bool{
	"--color":               true,
	"--colour":              true,
	"--line-number":         true,
	"--with-filename":       true,
	"--no-filename":         true,
	"--count":               true,
	"--files-with-matches":  true,
	"--files-without-match": true,
	"--only-matching":       true,
	"--quiet":               true,
	"--silent":              true,
	"--invert-match":        true,
	"--word-regexp":         true,
	"--line-regexp":         true,
	"--fixed-strings":       true,
	"--perl-regexp":         true,
	"--text":                true,
	"--binary":              true,
	"--max-count":           true,
	"--after-context":       true,
	"--before-context":      true,
	"--context":             true,
}

// Long flags to ignore
var longIgnored = map[string]bool{
	"--recursive":             true,
	"--dereference-recursive": true,
	"--extended-regexp":       true,
	"--basic-regexp":          true,
	"--binary-files":          true,
	"--directories":           true,
	"--devices":               true,
	"--no-messages":           true,
}

func translateFlags(args []string) []string {
	var rgArgs []string
	var patterns []string
	var paths []string
	skipNext := false

	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		if arg == "--" {
			// Everything after -- is paths
			paths = append(paths, args[i+1:]...)
			break
		}

		if strings.HasPrefix(arg, "--") {
			// Handle --option=value format
			if idx := strings.Index(arg, "="); idx != -1 {
				opt := arg[:idx]
				val := arg[idx+1:]

				switch opt {
				case "--include":
					rgArgs = append(rgArgs, "-g", val)
				case "--exclude":
					rgArgs = append(rgArgs, "-g", "!"+val)
				case "--exclude-dir":
					// Ensure directory pattern
					if !strings.HasSuffix(val, "/") {
						val = val + "/"
					}
					rgArgs = append(rgArgs, "-g", "!"+val)
				case "--color", "--colour":
					rgArgs = append(rgArgs, arg)
				case "--regexp":
					patterns = append(patterns, val)
				case "--file":
					rgArgs = append(rgArgs, "-f", val)
				case "--max-count":
					rgArgs = append(rgArgs, "-m", val)
				case "--after-context":
					rgArgs = append(rgArgs, "-A", val)
				case "--before-context":
					rgArgs = append(rgArgs, "-B", val)
				case "--context":
					rgArgs = append(rgArgs, "-C", val)
				case "--label":
					rgArgs = append(rgArgs, arg)
				default:
					if longPassthrough[opt] {
						rgArgs = append(rgArgs, arg)
					}
					// Ignore unknown long options with values
				}
				continue
			}

			// Handle --option format
			switch arg {
			case "--null", "--null-data":
				rgArgs = append(rgArgs, "-0")
			case "--include", "--exclude", "--exclude-dir":
				// These need a value
				if i+1 < len(args) {
					val := args[i+1]
					skipNext = true
					switch arg {
					case "--include":
						rgArgs = append(rgArgs, "-g", val)
					case "--exclude":
						rgArgs = append(rgArgs, "-g", "!"+val)
					case "--exclude-dir":
						if !strings.HasSuffix(val, "/") {
							val = val + "/"
						}
						rgArgs = append(rgArgs, "-g", "!"+val)
					}
				}
			case "--regexp":
				if i+1 < len(args) {
					patterns = append(patterns, args[i+1])
					skipNext = true
				}
			case "--file":
				if i+1 < len(args) {
					rgArgs = append(rgArgs, "-f", args[i+1])
					skipNext = true
				}
			default:
				if longPassthrough[arg] {
					rgArgs = append(rgArgs, arg)
				} else if longIgnored[arg] {
					// Skip
				} else {
					// Pass through unknown
					rgArgs = append(rgArgs, arg)
				}
			}
			continue
		}

		if strings.HasPrefix(arg, "-") && len(arg) > 1 && arg[1] != '-' {
			// Short flags
			flags := arg[1:]
			for j, c := range flags {
				if passthroughFlags[c] {
					rgArgs = append(rgArgs, "-"+string(c))
					continue
				}

				if passthroughWithValue[c] {
					// Check if value is attached or next arg
					remaining := flags[j+1:]
					var val string
					if len(remaining) > 0 {
						val = string(remaining)
					} else if i+1 < len(args) {
						// For -e, always take next arg as pattern (even if starts with -)
						// For others, only take if doesn't start with -
						if c == 'e' || !strings.HasPrefix(args[i+1], "-") {
							val = args[i+1]
							skipNext = true
						}
					}

					if c == 'e' {
						patterns = append(patterns, val)
					} else {
						rgArgs = append(rgArgs, "-"+string(c), val)
					}
					break
				}

				if c == 'Z' {
					rgArgs = append(rgArgs, "-0")
					continue
				}

				if ignoredFlags[c] {
					continue
				}

				// Unknown flag - pass through
				rgArgs = append(rgArgs, "-"+string(c))
			}
			continue
		}

		// Non-flag argument
		if len(patterns) == 0 && !strings.HasPrefix(arg, "-") {
			// First non-flag is the pattern (if no -e was used)
			patterns = append(patterns, arg)
		} else {
			paths = append(paths, arg)
		}
	}

	// Build final command - ensure we return empty slice not nil
	result := make([]string, 0)
	result = append(result, rgArgs...)

	// Add patterns
	if len(patterns) == 1 {
		// Single pattern - add directly
		pat := patterns[0]
		if strings.HasPrefix(pat, "-") {
			result = append(result, "--", pat)
		} else {
			result = append(result, pat)
		}
	} else if len(patterns) > 1 {
		// Multiple patterns - use -e for each
		for _, pat := range patterns {
			result = append(result, "-e", pat)
		}
	}

	// Add paths
	result = append(result, paths...)

	return result
}
