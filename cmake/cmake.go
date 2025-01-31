package cmake

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/env"
	"github.com/ZEL-30/zel/logger"
	"github.com/ZEL-30/zel/utils"
)

// cmake 配置命令参数
type ConfigArg struct {
	NoWarnUnusedCli       bool        // 不警告在命令行声明但未使用的变量
	BuildType             string      // 构建类型
	ExportCompileCommands bool        // 导出编译命令
	Kit                   *config.Kit // 编译器
	ProjectPath           string      // 源代码路径
	BuildPath             string      // 构建目录
	Generator             string      // 生成器
	CXXFlags              string      // C++ 编译参数
}

// cmake 构建命令参数
type BuildArg struct {
	BuildPath string // 构建路径
	BuildType string // 构建类型
	Target    string // 构建目标
}

var (
	appName string // 应用程序名称
	appPath string // 应用程序路径
)

func init() {
	appPath = utils.GetZelWorkPath()
	appName = filepath.Base(appPath)
}

func Run(configArg *ConfigArg, buildArg *BuildArg, target string, rebuild bool) error {
	err := Build(configArg, buildArg, rebuild, false)
	if err != nil {
		return err
	}

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
		return err
	}

	return err
}

func Build(configArg *ConfigArg, buildArg *BuildArg, rebuild bool, showInfo bool) error {

	if configArg.Kit == nil {
		env.SetBuildKit()
		configArg.Kit = config.Conf.Kit
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

	if c.NoWarnUnusedCli {
		result = append(result, "--no-warn-unused-cli")
	}

	result = append(result, "-DCMAKE_BUILD_TYPE:STRING="+c.BuildType)

	if c.ExportCompileCommands {
		result = append(result, "-DCMAKE_EXPORT_COMPILE_COMMANDS:BOOL=TRUE")
	}

	if c.Kit.Compiler.C != "" {
		result = append(result, "-DCMAKE_C_COMPILER:FILEPATH="+c.Kit.Compiler.C)
	}

	if c.Kit.Compiler.CXX != "" {
		result = append(result, "-DCMAKE_CXX_COMPILER:FILEPATH="+c.Kit.Compiler.CXX)
	}

	if c.ProjectPath != "" {
		result = append(result, "-S"+c.ProjectPath)
	}

	if c.BuildPath != "" {
		result = append(result, "-B"+c.BuildPath)
	}

	if c.Generator != "" {
		result = append(result, "-G="+c.Generator)
	}

	return result
}

func (b *BuildArg) toStringSlice() []string {
	var result []string

	result = append(result, "--build")
	result = append(result, b.BuildPath)

	result = append(result, "--config")
	result = append(result, b.BuildType)

	result = append(result, "--target")
	if len(b.Target) != 0 {
		result = append(result, b.Target)
	} else {
		result = append(result, "all")
	}

	result = append(result, "--")

	return result
}
