![wii](https://socialify.git.ci/oopsunix/wii/image?font=JetBrains+Mono&forks=1&issues=1&language=1&logo=https%3A%2F%2Favatars.githubusercontent.com%2Fu%2F133087009&name=1&owner=1&pattern=Plus&pulls=1&stargazers=1&theme=Light)

<h1 align="center">wii (what-is-installed)</h1>

<p align="center">One command, zero configuration. Quickly discover installed development environments, package managers, and CLI tools.</p>

<p align="center">
  <a href="https://github.com/oopsunix/hubp"><img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go"></a>
  <a href="https://github.com/oopsunix/hubp"><img src="https://img.shields.io/badge/License-MIT-blue.svg"></a>
  <a href="https://github.com/oopsunix/hubp"><img src="https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey"></a>
</p>

<p align="center">
  <a href="README.md">English</a> •
  <a href="README_CN.md">中文</a>
</p>

---

## Features

- **Development Environment Detection**: Automatically identify runtimes like Python, Java, Go, Rust, Node.js, .NET and their versions
- **Package Manager Integration**: Query CLI tools installed via brew, winget, scoop, npm, pip, Cargo and more
- **Smart Grouping**: Auto-categorize by User Local, System Local, and package managers
- **Parallel Probing**: Goroutine worker pool for concurrent version detection across all platforms
- **Cross-Platform**: Native support for Windows, macOS, and Linux
- **Zero-Dependency Distribution**: Install scripts + single binary, no runtime dependencies required
- **Live Snapshot**: No disk cache, fresh scan on every run


## Platform Support

| Platform | Support |
| -------- | ------- |
| Windows  | User Local, System Local, winget, Scoop, Chocolatey, NuGet(.NET), npm Global(Node), Python Scripts, pip(Python), go mod(Go), Cargo(Rust) |
| macOS    | User Local, System Local, Homebrew, apt/apt-get, dnf/yum, npm Global(Node), Python Scripts, pip(Python), go mod(Go), Cargo(Rust) |
| Linux    | User Local, System Local, Homebrew, Snap, npm Global(Node), Python Scripts, pip(Python), go mod(Go), Cargo(Rust) |



## Development Environment Detection

Automatically detects the following development environments:

| Category | Tools |
| -------- | ----- |
| **Runtimes** | Python, Java, Go, Rust, Node.js, .NET, PHP, Ruby, Deno, Bun |
| **Package Managers** | pip, npm, pnpm, yarn, Cargo |
| **Build Tools** | Maven, Gradle, CMake |
| **Compilers** | GCC, Clang |
| **Containers** | Docker, Podman, Kubectl |
| **Version Control** | Git |

## Quick Start

### Installation

#### macOS / Linux

```bash
curl -fsSL https://github.com/oopsunix/wii/raw/main/install.sh | bash
```

#### Windows (PowerShell)

```powershell
irm https://github.com/oopsunix/wii/raw/main/install.ps1 | iex
```

#### Go install

```bash
go install github.com/oopsunix/wii@latest
```

#### npm

```bash
npm install -g @oopsunix/wii
```

#### Manual Download

Download pre-built binaries from [GitHub Releases](https://github.com/oopsunix/wii/releases).

### Usage

```bash
wii [flags]

Flags:
  -c          int      Number of concurrent workers (default: CPU count)
  -format     string   Output format: table, json, csv
  -no-color            Disable color output
  -version             Show version information
  -h                   Show help message
```

### Examples

```bash
# Run directly (auto-detect development environments and CLI tools)
wii

# JSON output
wii -format json

# CSV output
wii -format csv

# Disable color output
wii -no-color

# Set concurrency
wii -c 20

# Show version
wii -version
```

### Output Example

```
[*] Current development environment
 ├─ .NET                  10.0.204     C:\Program Files\dotnet\dotnet.exe
 ├─ Cargo                 1.95.0       C:\Users\user\.cargo\bin\cargo.exe
 ├─ Git                   2.51.0       C:\Program Files\Git\cmd\git.exe
 ├─ Go                    1.25.0       D:\Dev\Go\bin\go.exe
 ├─ Java                  1.8.0        D:\Dev\Java\jdk1.8.0_181\bin\java.exe
 ├─ Node                  22.18.0      D:\Dev\nodejs\node.exe
 ├─ Python                3.11.9       D:\Dev\Python\Python311\python.exe
 ├─ Rust                  1.95.0       C:\Users\user\.cargo\bin\rustc.exe
 ├─ npm                   11.6.4       D:\Dev\nodejs\npm.cmd
 ├─ pip                   24.0         D:\Dev\Python\Python311\Scripts\pip.exe
 ├─ pnpm                  11.1.2       C:\Users\user\AppData\Local\pnpm\pnpm.cmd
 └─ yarn                  1.22.22      C:\Users\user\AppData\Roaming\npm\yarn.cmd

[*] User Local
 └─ uv                    0.9.21       C:\Users\user\.local\bin\uv.exe

[*] System Local
 ├─ bz                    7.03         C:\Program Files\Bandizip\bz.exe
 ├─ nmap                  7.80         C:\Program Files (x86)\Nmap\nmap.exe
 └─ upx                   3.96         D:\Dev\upx\upx.exe

[*] Go Tools
 ├─ dlv                   1.25.1       C:\Users\user\go\bin\dlv.exe
 ├─ gopls                 0.20.0       C:\Users\user\go\bin\gopls.exe
 ├─ staticcheck           2025.1.1     C:\Users\user\go\bin\staticcheck.exe
 └─ wails                 2.12.0       C:\Users\user\go\bin\wails.exe

[*] npm Global
 ├─ claude                2.1.153      C:\Users\user\AppData\Roaming\npm\claude.cmd
 ├─ codex                 0.118.0      C:\Users\user\AppData\Roaming\npm\codex.cmd
 └─ gemini                0.42.0       C:\Users\user\AppData\Roaming\npm\gemini.cmd
```

### Output Sections

| Section | Description |
|---------|-------------|
| **Current development environment** | Installed runtimes, compilers, and package managers |
| **User Local** | CLI tools in user directories (e.g., `~/.local/bin`) |
| **System Local** | CLI tools in system-level directories |
| **Go Tools** | Tools installed via `go install` |
| **Python Scripts** | CLI tools installed via `pip install` |
| **npm Global** | Tools installed via `npm install -g` |
| **Scoop** | Tools installed via Scoop package manager |
| **Cargo** | Tools installed via `cargo install` |

---

## Contributing

Issues and Pull Requests are welcome!

---

## Acknowledgments

- Inspired by the [what-is-installed](https://github.com/what-is-installed/what-is-installed) Bash version
- Built with Go standard library
- Cross-platform builds powered by [GoReleaser](https://goreleaser.com/)

---

## License

This project is licensed under the [MIT License](LICENSE).

---

<div align="center">
  <p>If you find this project helpful, please give it a ⭐ to help others discover it!</p>
  <p>Built with ❤️ by <a href="https://github.com/oopsunix">oopsunix</a></p>
</div>
