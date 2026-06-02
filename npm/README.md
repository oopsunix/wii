# wii (what-is-installed)

One command, zero configuration. Quickly discover installed development environments, package managers, and CLI tools.

## Installation

```bash
npm install -g @oopsunix/wii
```

## Usage

```bash
wii [flags]

Flags:
  -c     int      Number of concurrent workers (default: CPU count)
  -f     string   Output format: table, json, csv
  -nc             Disable color output
  -v              Show version information
  -h              Show help message
```

### Examples

```bash
# Run directly (auto-detect development environments and CLI tools)
wii

# JSON output
wii -f json

# CSV output
wii -f csv

# Disable color output
wii -nc

# Show version
wii -v

# Set concurrency
wii -c 20
```

### Output Example

```
[*] Current development environment
 ├─ Go                    1.25.0       /opt/homebrew/bin/go
 ├─ Java                  1.8.0        /Library/Java/JavaVirtualMachines/jdk1.8.0_181.jdk/Contents/Home/bin/java
 ├─ Node                  23.7.0       /opt/homebrew/bin/node
 ├─ Rust                  1.95.0       ~/.cargo/bin/rustc
 ├─ npm                   10.9.2       /opt/homebrew/bin/npm
 └─ pnpm                  10.17.1      /opt/homebrew/bin/pnpm

[*] User Local
 ├─ uv                    0.8.17       ~/.local/bin/uv
 └─ acorn                 ?            ~/Library/pnpm/acorn

[*] System Local
 ├─ docker                27.3.1       /usr/local/bin/docker
 └─ ghostty               1.3.1        /Applications/Ghostty.app/Contents/MacOS/ghostty

[*] npm Global
 ├─ @anthropic-ai/claude-code   2.1.159       /opt/homebrew/bin/claude
 ├─ @google/gemini-cli          0.41.2        /opt/homebrew/bin/gemini
 └─ @openai/codex               0.125.0       /opt/homebrew/bin/codex

[*] Cargo
 ├─ cargo                 1.95.0       ~/.cargo/bin/cargo
 └─ rustc                 1.95.0       ~/.cargo/bin/rustc

[*] Homebrew
 ├─ airbuddy              2.7.3,641    /opt/homebrew/bin/airbuddy
 ├─ dos2unix              7.5.3        /opt/homebrew/bin/dos2unix
 └─ goreleaser            2.4.8        /opt/homebrew/bin/goreleaser
```

## Output Sections

| Section | Description |
|---------|-------------|
| **Current development environment** | Installed runtimes, compilers, and package managers |
| **User Local** | CLI tools in user directories (e.g., `~/.local/bin`) |
| **System Local** | CLI tools in system-level directories |
| **npm Global** | Tools installed via `npm install -g` |
| **Homebrew** | Tools installed via Homebrew |
| **Go Tools** | Tools installed via `go install` |
| **Python Scripts** | CLI tools installed via `pip install` |
| **Cargo** | Tools installed via `cargo install` |

## Features

- **Development Environment Detection**: Automatically identify runtimes like Python, Java, Go, Rust, Node.js, .NET and their versions
- **Package Manager Integration**: Query CLI tools installed via brew, winget, scoop, npm, pip, Cargo and more
- **Smart Grouping**: Auto-categorize by User Local, System Local, and package managers
- **Parallel Probing**: Concurrent version detection across all platforms
- **Cross-Platform**: Native support for Windows, macOS, and Linux
- **Live Snapshot**: No disk cache, fresh scan on every run

## Other Installation Methods

```bash
# macOS / Linux (curl)
curl -fsSL https://github.com/oopsunix/wii/raw/main/install.sh | bash

# Windows (PowerShell)
irm https://github.com/oopsunix/wii/raw/main/install.ps1 | iex

# Go install
go install github.com/oopsunix/wii@latest
```

Pre-built binaries are also available on [GitHub Releases](https://github.com/oopsunix/wii/releases).

## Links

- [GitHub Repository](https://github.com/oopsunix/wii)
- [Issue Tracker](https://github.com/oopsunix/wii/issues)

## License

[MIT](https://github.com/oopsunix/wii/blob/main/LICENSE)
