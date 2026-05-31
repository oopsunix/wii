# wii (what-is-installed)

一条命令，零配置，快速检测本机已安装的开发环境、包管理器和命令行工具。

## 安装

```bash
npm install -g @oopsunix/wii
```

## 使用

```bash
# 直接运行
wii

# JSON 输出
wii -format json

# CSV 输出
wii -format csv

# 禁用颜色输出
wii -no-color

# 设置并发数
wii -c 20

# 显示版本
wii -version
```

## 功能特性

- **开发环境检测**：自动识别 Python、Java、Go、Rust、Node.js、.NET 等运行时及版本
- **包管理器集成**：查询 npm、pip、Cargo、Scoop 等包管理器安装的 CLI 工具
- **智能分类**：按 User Local、System Local、包管理器自动分组
- **跨平台支持**：Windows、macOS、Linux 原生支持

## 更多信息

- [GitHub 仓库](https://github.com/oopsunix/wii)
- [问题反馈](https://github.com/oopsunix/wii/issues)

## 许可证

MIT
