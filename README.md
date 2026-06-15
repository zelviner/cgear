# Cgear

Cgear 是一个面向 Windows C++/CMake 项目的命令行工具。它用于创建项目、选择编译工具链、调用 CMake 构建和运行程序、执行 GoogleTest、安装第三方 C++ 库以及生成发布压缩包。

当前版本：`2.0.0`

## 功能

- 交互式创建 Application、Qt Application、静态库、动态库和测试用例
- 自动发现 MSVC、Clang、Clang-cl 和 MinGW 工具链
- 管理平台、CMake 生成器和构建类型
- 构建、运行和重新构建 CMake 项目
- 列出或运行 GoogleTest 测试
- 从 GitHub 下载并安装 CMake 第三方库
- 收集可执行文件和 DLL，生成发布 ZIP
- 生成 `include` 目录和 Word 源码文档

## 环境要求

Cgear 当前仅支持 Windows x64。使用前至少需要：

- Windows 10 或 Windows 11
- CMake 3.14 或更高版本
- 一个受支持的 C/C++ 工具链：Visual Studio 2019/2022、Clang、Clang-cl 或 MinGW
- Git：使用 `cgear install author:repository` 或完整版在线安装器时需要
- Qt 的 `windeployqt`：仅在打包 Qt Application 时需要

选择 `Ninja` 生成器时还需要 `ninja.exe`。完整版安装器会自动安装 Ninja 和 vcpkg。

## 安装

### 完整版安装器

`cgear_install.iss` 生成按用户安装的在线安装器。安装器会：

- 安装 `cgear.exe`
- 下载 Ninja `1.12.1`
- 下载并引导 vcpkg `2026.06.01`
- 设置用户级 `CGEAR_HOME`
- 将 `%CGEAR_HOME%\bin` 加入用户 PATH

安装器需要网络连接和 Git。卸载时会移除 Cgear、Ninja、vcpkg 工具文件和环境变量，但保留 `installed` 与 `pkg` 中的用户库。

### 便携版

`cgear_portable.iss` 只解压 `bin\portable`，不会修改注册表、PATH 或系统环境变量。运行：

```powershell
.\bin\cgear.exe version
```

当可执行文件位于 `bin` 目录时，Cgear 会自动将其父目录作为 `CGEAR_HOME`。

### 从源码构建

```powershell
go build -o bin\cgear.exe .
.\bin\cgear.exe version
```

开发环境要求 Go `1.20` 或更高版本。

## 快速开始

```powershell
# 创建项目，并按提示选择项目类型、工具链、平台和生成器
cgear new hello

cd hello

# 查看当前配置
cgear env

# 构建并运行
cgear build
cgear run
```

生成的基础结构如下：

```text
hello/
├── CMakeLists.txt
├── cgear.json
├── cmake/
├── doc/
├── src/
├── test/
└── .vscode/
```

## 配置文件

Cgear 优先读取项目根目录的 `cgear.json`，不存在时读取 `Cgearfile` YAML。典型 JSON 配置：

```json
{
  "Version": 0,
  "toolchain": {
    "name": "Visual Studio Community 2022 Release",
    "compilers": {
      "C": "v143",
      "CXX": "v143"
    },
    "is_msvc": true
  },
  "generator": "Visual Studio 17 2022",
  "platform": "x64",
  "build_type": "Debug",
  "project_type": "Application",
  "project_path": "C:\\Projects\\hello",
  "runtime_dependencies": [
    "example"
  ]
}
```

字段说明：

| 字段 | 说明 |
| --- | --- |
| `Version` | 配置格式版本，当前为 `0` |
| `toolchain` | 编译器名称、C/C++ 编译器路径或 MSVC toolset |
| `generator` | `Ninja` 或 `Visual Studio 17 2022` |
| `platform` | `x86` 或 `x64` |
| `build_type` | `Debug`、`Release`、`RelWithDebInfo` 或 `MinSizeRel` |
| `project_type` | 创建项目时选择的类型 |
| `project_path` | 保存时记录的路径；加载后以当前工作目录为准 |
| `runtime_dependencies` | `pack` 时从安装目录复制的 DLL 名称，不含 `.dll` 后缀 |

## 命令

### `new`

```powershell
cgear new <project_name>
```

创建新项目。项目类型和构建环境通过交互列表选择。选择 `Test cases` 时，会在当前项目中创建测试用例。

### `env`

```powershell
cgear env
cgear env Toolchain
cgear env Generator
cgear env Platform
cgear env BuildType
```

查看或修改当前项目配置。子命令区分大小写。

### `build`

```powershell
cgear build
cgear build -r
cgear build -t <target>
```

配置并构建项目。`-r` 会先删除 `build` 目录，`-t` 指定 CMake target。

### `run`

```powershell
cgear run
cgear run -r
```

构建并运行由顶层 `project(...)` 声明的程序，默认从项目的 `bin` 目录启动。

### `test`

```powershell
cgear test
cgear test FileSystem
cgear test FileSystem.read
cgear test FileSystem -r
```

无参数时列出 `bin\test` 下 GoogleTest 程序中的测试。指定名称时使用 `--gtest_filter` 运行测试，`-r` 表示重新构建。

### `install`

```powershell
# 构建并安装当前 CMake 第三方库
cgear install

# 通过 GitHub SSH 下载、构建并安装
cgear install google:googletest

# 导入已经包含 include 和 lib 的绝对路径
cgear install C:\Libraries\example
```

GitHub 安装使用 `git@github.com:author/repository`，因此需要可用的 GitHub SSH 密钥。源码包会存放在 `%CGEAR_HOME%\pkg`，安装结果位于 `%CGEAR_HOME%\installed\<platform>`。

### `pack`

```powershell
cgear pack
```

输入版本号后构建 Release，将 `bin` 中的应用程序、DLL 和配置的运行时依赖复制到 `bin\release`，最后生成 ZIP。Qt 项目会额外执行 `windeployqt`。

### `generate`

```powershell
cgear generate include
cgear generate docs
```

创建 `include` 目录，或将源码导出到项目根目录的 `source.docx`。

### 其他命令

```powershell
cgear count
cgear version
cgear version -o json
cgear version -o yaml
cgear help <command>
```

## CGEAR_HOME 目录

```text
%CGEAR_HOME%/
├── bin/          # cgear、Ninja 和 vcpkg 命令
├── installed/    # 按平台和构建类型安装的头文件、库和 DLL
├── pkg/          # 从 GitHub 克隆的第三方源码
├── ports/        # vcpkg ports
├── scripts/      # 包含 vcpkg CMake toolchain
├── triplets/
└── versions/
```

## 生成安装包

先准备程序：

```powershell
go build -o bin\cgear.exe .
```

然后使用 Inno Setup 6 编译：

```powershell
ISCC.exe cgear_install.iss
```

便携版脚本不会自动组装目录。编译前应准备完整的 `bin\portable`，至少包含：

```text
bin/portable/
├── bin/cgear.exe
├── installed/
├── pkg/
├── ports/
├── scripts/
├── triplets/
├── versions/
├── .vcpkg-root
└── vcpkg.exe
```

然后执行：

```powershell
ISCC.exe cgear_portable.iss
```

发布新版本时，需要同步更新 `config.Version`、`cgear_install.iss` 和 `cgear_portable.iss` 中的版本号。

## 开发验证

```powershell
go test ./...
go vet ./...
staticcheck ./...
```

## 当前限制

- 仅支持 Windows 项目和 `.exe`/`.dll` 发布流程
- GitHub 库安装目前使用 SSH 地址，不支持 HTTPS 参数
- `pack` 面向当前生成的固定目录结构，不是通用 CPack 替代品
- 部分命令仍直接退出进程，命令级自动化测试覆盖有限
