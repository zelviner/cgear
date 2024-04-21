package cmake

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
)

const (
	RELEASE = iota
	DEBUG
)

// cmake 配置命令参数
type ConfigArg struct {
	NoWarnUnusedCli       bool       // 不警告在命令行声明但未使用的变量
	BuildType             int        // 构建类型
	ExportCompileCommands bool       // 导出编译命令
	Kit                   config.Kit // 编译器
	AppPath               string     // 源代码路径
	BuildPath             string     // 构建目录
	Generator             string     // 生成器
}

// cmake 构建命令参数
type BuildArg struct {
	BuildPath string // 构建路径
	BuildType int    // 构建类型
}

var (
	appName string // 应用程序名称
	appPath string // 应用程序路径
)

func init() {
	appPath, _ = os.Getwd()
	appName = filepath.Base(appPath)
}

func Run(target string) {

}

func Build(configArg *ConfigArg, buildArg *BuildArg, rebuild bool) {

	logger.Log.Infof("Using '%s' as the kit", configArg.Kit.Name)

	// 检查编译目录是否存在，如果存在则删除，然后创建
	if rebuild {
		if _, err := os.Stat(configArg.BuildPath); err == nil {
			os.RemoveAll(configArg.BuildPath)
		}
	}

	// 配置 C++ 项目
	cmd := exec.Command("cmake", configArg.toStringSlice()...)
	logger.Log.Infof("Running 'cmake %s'", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	// 编译 C++ 项目
	cmd = exec.Command("cmake", buildArg.toStringSlice()...)
	// cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}

func (c *ConfigArg) toStringSlice() []string {
	var result []string

	if c.NoWarnUnusedCli {
		result = append(result, "--no-warn-unused-cli")
	}

	if c.BuildType == RELEASE {
		result = append(result, "-DCMAKE_BUILD_TYPE:STRING=Release")
	} else if c.BuildType == DEBUG {
		result = append(result, "-DCMAKE_BUILD_TYPE:STRING=Debug")
	}

	if c.ExportCompileCommands {
		result = append(result, "-DCMAKE_EXPORT_COMPILE_COMMANDS:BOOL=TRUE")
	}

	if c.Kit.Compiler.C != "" {
		result = append(result, "-DCMAKE_C_COMPILER:FILEPATH="+c.Kit.Compiler.C)
	}

	if c.Kit.Compiler.CXX != "" {
		result = append(result, "-DCMAKE_CXX_COMPILER:FILEPATH="+c.Kit.Compiler.CXX)
	}

	if c.AppPath != "" {
		result = append(result, "-S"+c.AppPath)
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

	if b.BuildType == RELEASE {
		result = append(result, "--config Release")
	} else if b.BuildType == DEBUG {
		result = append(result, "--config Debug")
	}

	result = append(result, "--")

	return result
}
