package translator

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"text/tabwriter"
)

var (
	registry = make(map[string]Translator)
	mu       sync.RWMutex
)

// Register adds a translator to the global registry
func Register(t Translator) {
	mu.Lock()
	defer mu.Unlock()
	registry[t.Name()] = t
}

// Get returns a translator for the given source and target tools
// Returns nil if no translator is found
func Get(source, target string) Translator {
	return GetByName(source + "2" + target)
}

// GetByName returns a translator by its name (e.g., "ls2eza")
// Returns nil if no translator is found
func GetByName(name string) Translator {
	mu.RLock()
	defer mu.RUnlock()
	return registry[name]
}

// List returns all registered translator names
func List() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// MustGet returns a translator or panics if not found
func MustGet(source, target string) Translator {
	t := Get(source, target)
	if t == nil {
		panic(fmt.Sprintf("no translator registered for %s to %s", source, target))
	}
	return t
}

// PrintTable writes a formatted table of all translators to the given writer
func PrintTable(w io.Writer) {
	mu.RLock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	mu.RUnlock()

	sort.Strings(names)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TRANSLATOR\tSOURCE\tTARGET\tDEFAULT ENABLED")
	for _, name := range names {
		t := GetByName(name)
		included := "no"
		if t.IncludeInInit() {
			included = "yes"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", name, t.SourceTool(), t.TargetTool(), included)
	}
	tw.Flush()
}
