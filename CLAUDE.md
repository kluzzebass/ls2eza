# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

reflag is a Go CLI tool that translates command-line flags between different tools. It currently supports translating `ls` flags to `eza` equivalents, with an extensible architecture for adding more translators.

## Build and Test Commands

```bash
make build               # Build with version info
make test                # Run all tests
make clean               # Remove binary
make build-all           # Cross-compile for all platforms
go run . ls eza -la      # Run directly without building
go test -run TestName    # Run specific test
```

## Releasing

Push a semver tag to trigger a GitHub Actions release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Architecture

### Package Structure

```
reflag/
├── main.go                       # CLI entry point, symlink detection
├── translator/
│   ├── translator.go             # Translator interface
│   ├── registry.go               # Global translator registry
│   └── ls2eza/
│       ├── translator.go         # ls→eza implementation
│       └── translator_test.go    # ls2eza tests
└── main_test.go                  # CLI tests
```

### Core Components

1. **Translator Interface** (`translator/translator.go`):
   - `Name()` - translator identifier (e.g., "ls2eza")
   - `SourceTool()` - source tool name (e.g., "ls")
   - `TargetTool()` - target tool name (e.g., "eza")
   - `Translate(args)` - converts source args to target args
   - `Optional()` - returns true if excluded from `--init` by default

2. **Registry** (`translator/registry.go`):
   - `Register(t)` - register a translator
   - `Get(source, target)` - lookup by source/target
   - `GetByName(name)` - lookup by name (for symlink detection)
   - `List()` - list all registered translators

3. **CLI** (`main.go`):
   - Symlink detection: parses binary name for `<source>2<target>` pattern
   - Explicit mode: `reflag <source> <target> [flags...]`
   - Built-in flags: `--list`, `--version`, `--help`

### ls2eza Translator

Located in `translator/ls2eza/`:

1. **Mode detection** - `getLSMode()` determines BSD vs GNU ls compatibility:
   - Auto-detects based on OS (darwin/freebsd → BSD, linux/others → GNU)
   - Override with `LS2EZA_MODE=bsd` or `LS2EZA_MODE=gnu`

2. **Flag mappings**:
   - `reverseNeeded` - flags that need sort order correction (`t`, `S`, `c`, `u`, `U`)
   - `flagMap` - short flag translations (30+ flags)
   - `longFlagMap` - long option translations
   - `longFlagPrefixes` - long options with =value

3. **Reverse sort handling** - XOR logic to match ls sort order:
   - `ls -lt` needs `--reverse` (ls shows newest first, eza shows oldest first)
   - `ls -ltr` does NOT need `--reverse` (user explicitly wants oldest first)

### BSD vs GNU Conflicts

These flags have different meanings between BSD and GNU ls:
- `-T`: BSD=full time display, GNU=tab size (ignored)
- `-X`: BSD=don't cross filesystems (ignored), GNU=sort by extension
- `-I`: BSD=prevent auto -A (ignored), GNU=ignore pattern
- `-w`: BSD=raw non-printable (ignored), GNU=output width
- `-D`: BSD=date format, GNU=dired mode (ignored)

## Adding a New Translator

1. Create package `translator/<name>/`
2. Implement `translator.Translator` interface
   - Set `Optional()` to `true` if the translator should be excluded from `--init` by default
   - Set `Optional()` to `false` for core/commonly-used translators
3. Call `translator.Register()` in `init()`
4. Import in `main.go` with blank identifier: `_ "github.com/kluzzebass/reflag/translator/<name>"`
5. Update `README.md` with new translator information and how to install the tool

### Optional Translators

Translators marked as optional (returning `true` from `Optional()`) are excluded from `./reflag --init` by default. This is useful for:
- Experimental or less commonly used translators
- Translators for niche tools
- New translators that need more testing

Optional translators can still be explicitly included:
```bash
reflag --init bash dig2doggo  # Include only dig2doggo
reflag --init zsh ls2eza dig2doggo  # Include specific translators
```
