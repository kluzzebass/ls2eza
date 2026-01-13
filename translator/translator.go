package translator

// Translator defines the interface for converting flags between tools
type Translator interface {
	// Name returns the translator identifier (e.g., "ls2eza")
	Name() string

	// SourceTool returns the name of the source tool (e.g., "ls")
	SourceTool() string

	// TargetTool returns the name of the target tool (e.g., "eza")
	TargetTool() string

	// Translate converts source tool arguments to target tool arguments
	// The mode parameter allows dialect selection (e.g., "bsd" or "gnu" for ls2eza)
	Translate(args []string, mode string) []string

	// IncludeInInit returns true if this translator should be included in --init by default
	// Translators returning false can still be explicitly included via --init <translator>
	IncludeInInit() bool
}
