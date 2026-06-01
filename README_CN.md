![wii](https://socialify.git.ci/oopsunix/wii/image?font=JetBrains+Mono&forks=1&issues=1&language=1&logo=https%3A%2F%2Favatars.githubusercontent.com%2Fu%2F133087009&name=1&owner=1&pattern=Plus&pulls=1&stargazers=1&theme=Light)

<h1 align="center">wii(what-is-installed)</h1>

<p align="center">一条命令，零配置，快速了解本机安装了哪些开发环境、包管理器和命令行工具。</p>

<p align="center">
  <a href="https://github.com/oopsunix/hubp"><img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go"></a>
  <a href="https://github.com/oopsunix/hubp"><img src="https://img.shields.io/badge/License-MIT-blue.svg"></a>
  <a href="https://github.com/oopsunix/hubp"><img src="https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey"></a>
</p>

<p align="center">
<a href="README.md">English</a> | <strong>中文</strong>
</p>

---

## 功能特性

- **开发环境检测**：自动识别 Python、Java、Go、Rust、Node.js、.NET 等运行时及版本
- **包管理器集成**：查询 brew、winget、scoop、npm、pip、Cargo 等包管理器安装的 CLI 工具
- **智能分类**：按 User Local、System Local、包管理器自动分组
- **并行探测**：使用 goroutine 工具池并行检测版本，提升全平台扫描速度
- **跨平台支持**：Windows、macOS、Linux 原生支持
- **零依赖分发**：安装脚本 + 单一二进制，无需安装额外运行时
- **实时快照**：无磁盘缓存，每次运行都是全新扫描


## 平台适配

| 平台    | 适配                                                         |
| ------- | ------------------------------------------------------------ |
| Windows | User Local、System Local、winget、Scoop、Chocolatey、NuGet(.NET)、npm Global(Node)、Python Scripts、pip(Python)、go mod(Go)、Cargo(Rust) |
| macOS   | User Local、System Local、Homebrew、apt/apt-get、dnf/yum、npm Global(Node)、Python Scripts、pip(Python)、go mod(Go)、Cargo(Rust) |
| Linux   | User Local、System Local、Homebrew、Snap、npm Global(Node)、Python Scripts、pip(Python 包)、go mod(Go)、Cargo(Rust) |



## 开发环境检测

自动检测以下开发环境：

| 类别         | 工具                                                        |
| ------------ | ----------------------------------------------------------- |
| **运行时**   | Python, Java, Go, Rust, Node.js, .NET, PHP, Ruby, Deno, Bun |
| **包管理器** | pip, npm, pnpm, yarn, Cargo                                 |
| **构建工具** | Maven, Gradle, CMake                                        |
| **编译器**   | GCC, Clang                                                  |
| **容器**     | Docker, Podman, Kubectl                                     |
| **版本控制** | Git                                                         |

## 快速开始

### 安装

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

#### 手动下载

从 [GitHub Releases](https://github.com/oopsunix/wii/releases) 下载预编译二进制文件。

### 用法

```bash
wii [flags]

Flags:
  -c          int      并发数 (默认 = CPU count)
  -format     string   Json/CSV输出
  -no-color   string   禁用颜色输出
  -h                   显示帮助信息
```
### 使用示例

```
# 直接运行（自动检测开发环境和 CLI 工具）
wii

# JSON 输出
wii -format json

# CSV 输出
wii -format csv

# 禁用颜色输出
wii -no-color

# 检测并发数
wii -c 20
```

### 输出示例

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

### 输出分区说明

| 分区 | 说明 |
|------|------|
| **Current development environment** | 已安装的运行时、编译器、包管理器 |
| **User Local** | 用户目录下的 CLI 工具（如 `~/.local/bin`） |
| **System Local** | 系统级目录下的 CLI 工具 |
| **Go Tools** | `go install` 安装的工具 |
| **Python Scripts** | `pip install` 安装的 CLI 工具 |
| **npm Global** | `npm install -g` 安装的工具 |
| **Scoop** | Scoop 包管理器安装的工具 |
| **Cargo** | `cargo install` 安装的工具 |




## 项目结构

```
wii/
├── cmd/
│   └── wii/
│       └── main.go              # 程序入口
├── internal/
│   ├── devenv/
│   │   └── devenv.go            # 开发环境检测
│   ├── model/
│   │   └── model.go             # 数据模型定义
│   ├── platform/
│   │   ├── platform.go          # 平台接口
│   │   ├── windows.go           # Windows 平台实现
│   │   ├── darwin.go            # macOS 平台实现
│   │   ├── linux.go             # Linux 平台实现
│   │   └── bsd.go               # BSD 平台实现
│   ├── probe/
│   │   ├── probe.go             # 版本探测逻辑
│   │   └── version.go           # 版本号提取
│   ├── provider/
│   │   ├── provider.go          # 包管理器接口
│   │   ├── npm.go               # npm 支持
│   │   ├── pip.go               # pip 支持
│   │   ├── cargo.go             # Cargo 支持
│   │   └── ...                  # 其他包管理器
│   ├── render/
│   │   ├── render.go            # 渲染接口
│   │   ├── compact.go           # 紧凑格式渲染
│   │   ├── json.go              # JSON 渲染
│   │   └── csv.go               # CSV 渲染
│   └── scan/
│       ├── scan.go              # PATH 扫描逻辑
│       ├── path_windows.go      # Windows 路径处理
│       └── path_unix.go         # Unix 路径处理
├── go.mod
├── go.sum
└── README.md
```

---

## 贡献

欢迎提交 Issue 和 Pull Request！

---

## 致谢

- 灵感来源于 [what-is-installed](https://github.com/what-is-installed/what-is-installed) Bash 版本
- 使用 Go 标准库构建
- 使用 [GoReleaser](https://goreleaser.com/) 进行跨平台构建

---

## 开源协议

本项目基于 [Apache License 2.0](LICENSE) 协议开源。

---

<div align="center">
  <p>如果你觉得这个项目对你有帮助，不妨给它点个⭐，让更多人发现这款工具！</p>
  <p>Built with ❤️ by <a href="https://github.com/oopsunix">oopsunix</a></p>
</div>
