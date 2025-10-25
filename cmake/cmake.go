package cmake

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/zelviner/cgear/config"
	"github.com/zelviner/cgear/env"
	"github.com/zelviner/cgear/logger"
	"github.com/zelviner/cgear/utils"
)

// cmake 配置命令参数
type ConfigArg struct {
	Toolchain             *config.Toolchain // 工具链
	Platform              string            // 架构
	Generator             string            // 生成器
	BuildType             string            // 构建类型
	ProjectPath           string            // 源代码路径
	BuildPath             string            // 构建目录
	CXXFlags              string            // C++ 编译参数
	NoWarnUnusedCli       bool              // 不警告在命令行声明但未使用的变量
	ExportCompileCommands bool              // 导出编译命令
}

// cmake 构建命令参数
type BuildArg struct {
	BuildPath string // 构建路径
	Target    string // 构建目标
	BuildType string // 构建类型
	IsMSVC    bool   // 是否为 MSVC 工具链
}

var (
	appName string // 应用程序名称
	appPath string // 应用程序路径
)

func init() {
	appPath = utils.GetCgearWorkPath()
	appName = filepath.Base(appPath)
}

func Run(configArg *ConfigArg, buildArg *BuildArg, target string, rebuild bool) error {
	err := Build(configArg, buildArg, rebuild, false)
	if err != nil {
		return err
	}

	// 设置临时环境变量
	var dllPath string
	cgearHome := utils.GetCgearHomePath()
	switch config.Conf.Platform {
	case "x86":
		dllPath = filepath.Join(cgearHome, "installed", "x86-windows")
	case "x64":
		dllPath = filepath.Join(cgearHome, "installed", "x64-windows")
	}

	switch config.Conf.BuildType {
	case "Debug":
		dllPath = filepath.Join(dllPath, "debug", "bin")
	case "Release":
		dllPath = filepath.Join(dllPath, "bin")
	}

	restore, err := utils.SetEnvTemp("PATH", dllPath)
	if err != nil {
		logger.Log.Errorf("Failed to set PATH environment variable: %v", err)
		return err
	}
	defer restore() // 确保在函数结束时恢复原始 PATH

	// 运行应用程序
	if len(target) == 0 {
		target = appName + ".exe"
	} else {
		target = target + ".exe"
	}

	runPath := filepath.Join(appPath, "bin", target)
	cmd := exec.Command(runPath)
	cmd.Dir = filepath.Join(appPath, "bin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logger.Log.Errorf("Failed to run application: %v", err)
		return err
	}

	return err
}

func Build(configArg *ConfigArg, buildArg *BuildArg, rebuild bool, showInfo bool) error {

	if configArg.Toolchain == nil {
		env.SetToolchain()
		configArg.Toolchain = config.Conf.Toolchain
	}

	// 检查是否需要重新构建
	if rebuild {
		if _, err := os.Stat(configArg.BuildPath); err == nil {
			os.RemoveAll(configArg.BuildPath)
		}
	}

	// 配置 C++ 项目
	cmd := exec.Command("cmake", configArg.toStringSlice()...)
	if showInfo {
		logger.Log.Infof("Running '%s'", cmd.String())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	// 编译 C++ 项目
	cmd = exec.Command("cmake", buildArg.toStringSlice()...)
	if showInfo {
		logger.Log.Infof("Running '%s'", cmd.String())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func (c *ConfigArg) toStringSlice() []string {
	var result []string

	if c.Generator != "" {
		result = append(result, "-G", c.Generator)
	}

	if c.Toolchain.Compiler.C != "" {
		if !c.Toolchain.IsMSVC {
			result = append(result, "-DCMAKE_C_COMPILER:FILEPATH="+c.Toolchain.Compiler.C)
		}
	}

	if c.Toolchain.Compiler.CXX != "" {
		if c.Toolchain.IsMSVC {
			result = append(result, "-T", c.Toolchain.Compiler.C)
		} else {
			result = append(result, "-DCMAKE_CXX_COMPILER:FILEPATH="+c.Toolchain.Compiler.CXX)
		}
	}

	switch c.Platform {
	case "x86":
		if c.Toolchain.IsMSVC {
			result = append(result, "-A", "Win32")
			break
		}

		toolchainFile := filepath.Join(c.ProjectPath, "cmake/clang-32bit-toolchain.cmake")
		result = append(result, "-DCMAKE_TOOLCHAIN_FILE="+toolchainFile)

	case "x64":
		if c.Toolchain.IsMSVC {
			result = append(result, "-A", "x64")
			break
		}

		toolchainFile := filepath.Join(c.ProjectPath, "cmake/clang-64bit-toolchain.cmake")
		result = append(result, "-DCMAKE_TOOLCHAIN_FILE="+toolchainFile)
	}

	switch c.BuildType {
	case "Debug":
		result = append(result, "-DCMAKE_BUILD_TYPE=Debug")
	case "Release":
		result = append(result, "-DCMAKE_BUILD_TYPE=Release")
	}

	if c.NoWarnUnusedCli {
		result = append(result, "--no-warn-unused-cli")
	}

	if c.ExportCompileCommands {
		result = append(result, "-DCMAKE_EXPORT_COMPILE_COMMANDS:BOOL=TRUE")
	}

	if c.ProjectPath != "" {
		result = append(result, "-S"+c.ProjectPath)
	}

	if c.BuildPath != "" {
		result = append(result, "-B"+c.BuildPath)
	}

	return result
}

func (b *BuildArg) toStringSlice() []string {
	var result []string

	result = append(result, "--build")
	result = append(result, b.BuildPath)

	if b.IsMSVC {
		result = append(result, "--config")
		result = append(result, b.BuildType)
	}

	if !b.IsMSVC && b.Target != "" {
		result = append(result, "--target")
		if len(b.Target) != 0 {
			result = append(result, b.Target)
		} else {
			result = append(result, "all")
		}
	}

	result = append(result, "--")

	return result
}
