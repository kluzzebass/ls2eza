package df2duf

import (
	"strings"

	"github.com/kluzzebass/reflag/translator"
)

func init() {
	translator.Register(&Translator{})
}

// Translator implements the df to duf flag translation
type Translator struct{}

func (t *Translator) Name() string        { return "df2duf" }
func (t *Translator) SourceTool() string  { return "df" }
func (t *Translator) TargetTool() string  { return "duf" }
func (t *Translator) IncludeInInit() bool { return false }

// Translate converts du arguments to duf arguments
func (t *Translator) Translate(args []string, mode string) []string {
	return translateFlags(args)
}

// Flags to ignore or that have no duf equivalent
var ignoredFlags = map[string]bool{
	"-A":                 true, // apparent size - no equivalent
	"--apparent-size":    true,
	"-B":                 true, // block size (with value) - duf uses SI by default
	"--block-size":       true,
	"-c":                 true, // grand total - duf always shows summary
	"--total":            true,
	"-d":                 true, // max depth - duf shows filesystems, not directories
	"--max-depth":        true,
	"-g":                 true, // gigabytes - duf shows human-readable
	"-k":                 true, // kilobytes - duf shows human-readable
	"-m":                 true, // megabytes - duf shows human-readable
	"-H":                 true, // follow symlinks on command line
	"--dereference-args": true,
	"-L":                 true, // follow all symlinks
	"--dereference":      true,
	"-P":                 true, // don't follow symlinks
	"--no-dereference":   true,
	"-s":                 true, // summarize - duf is always a summary
	"--summarize":        true,
	"-t":                 true, // threshold - no direct equivalent
	"--threshold":        true,
	"-I":                 true, // BSD ignore pattern
	"-n":                 true, // BSD nodump flag
	"-r":                 true, // generate error messages (default)
	"-S":                 true, // separate dirs - not applicable
	"--separate-dirs":    true,
	"--time":             true, // no time display in duf
	"--time-style":       true,
	"-0":                 true, // null terminated
	"--null":             true,
	"-D":                 true, // BSD date format
	"-h":                 true, // human-readable is duf default
	"--human-readable":   true,
	"--si":               true, // duf uses SI by default
}

func translateFlags(args []string) []string {
	dufArgs := []string{}
	var paths []string
	skipNext := false

	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		// Handle --option=value format
		if strings.HasPrefix(arg, "--") {
			if idx := strings.Index(arg, "="); idx != -1 {
				opt := arg[:idx]
				val := arg[idx+1:]

				switch opt {
				case "--exclude":
					// Map to hide mount point pattern
					dufArgs = append(dufArgs, "-hide-mp", val)
				case "--block-size", "--threshold", "--max-depth":
					// Skip these with their values
					continue
				default:
					if ignoredFlags[opt] {
						continue
					}
					// Pass through unknown options
					dufArgs = append(dufArgs, arg)
				}
				continue
			}

			// Handle standalone long options
			switch arg {
			case "--all":
				dufArgs = append(dufArgs, "-all")
			case "--one-file-system", "-x":
				// duf shows all filesystems by default, so no direct equivalent
				// Could potentially use -only local but that's not the same
			case "--inodes":
				dufArgs = append(dufArgs, "-inodes")
			default:
				if ignoredFlags[arg] {
					continue
				}
				// Pass through unknown long options
				dufArgs = append(dufArgs, arg)
			}
			continue
		}

		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			flags := arg[1:]
			for j, c := range flags {
				switch c {
				case 'a': // all files - map to -all to include all filesystems
					dufArgs = append(dufArgs, "-all")
				case 'l': // count hard links multiple times
					// This is the flag the user wants to use: "du -lh"
					// Since duf shows filesystem usage not file sizes, we'll ignore this
					// but not fail - just continue
					continue
				case 'x': // one file system
					// duf shows filesystems, not directory trees
					// Could use -only local but not the same
					continue
				case 'I': // BSD exclude pattern
					remaining := flags[j+1:]
					var val string
					if len(remaining) > 0 {
						val = string(remaining)
					} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						val = args[i+1]
						skipNext = true
					}
					if val != "" {
						dufArgs = append(dufArgs, "-hide-mp", val)
					}
					goto nextArg
				case 'B': // block size - skip value
					remaining := flags[j+1:]
					if len(remaining) == 0 && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						skipNext = true
					}
					goto nextArg
				case 't': // threshold - skip value
					remaining := flags[j+1:]
					if len(remaining) == 0 && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						skipNext = true
					}
					goto nextArg
				case 'd': // max depth - skip value
					remaining := flags[j+1:]
					if len(remaining) == 0 && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						skipNext = true
					}
					goto nextArg
				case 'h', 'c', 'P', 'L', 'H', 's', 'A', 'g', 'k', 'm', 'n', 'r', 'S', '0', 'D':
					// Ignored flags that are either duf defaults or not applicable
					continue
				default:
					// Pass through unknown flags
					dufArgs = append(dufArgs, "-"+string(c))
				}
			}
		nextArg:
			continue
		}

		// Non-flag argument (path)
		// duf doesn't take directory arguments like du does
		// It shows filesystem information, so we'll collect paths but they won't be used
		paths = append(paths, arg)
	}

	// duf doesn't use paths the same way du does
	// It shows mounted filesystems, not directory contents
	// So we just return the flags
	return dufArgs
}
