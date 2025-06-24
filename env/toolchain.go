package env

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
	ui "github.com/ZEL-30/zel/ui/select"
)

type Compiler struct {
	Name    string
	CPath   string
	CXXPath string
}

var compilerTypes = map[string]string{
	"Clang":    "clang++.exe",
	"Clang-cl": "clang-cpp.exe",
	"Mingw":    "g++.exe",
}

func SetToolchain() {
	logger.Log.Info("Finding toolchains...")
	toolchains, err := findToolchains()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	selected, cancelled, err := ui.ListOption("Please select a toolchain: ", toolchains, func(t *config.Toolchain) string { return t.Name })
	if err != nil {
		logger.Log.Errorf("Failed to select a toolchain: %v", err)
		return
	}
	if cancelled {
		logger.Log.Info("Cancelled selecting a toolchain.")
		return
	}

	config.Conf.Toolchain = selected
	logger.Log.Successf("Toolchain set to: %s", selected.Name)
}

func findToolchains() (Toolchains []*config.Toolchain, err error) {

	// 查看环境变量中的路径
	pathEnv := os.Getenv("PATH")

	// 将路径字符串按分号分割成切片
	paths := strings.Split(pathEnv, ";")

	// 在环境变量中搜索 C++ 编译器
	for _, path := range paths {
		for key, value := range compilerTypes {
			compilerPath := filepath.Join(path, value)
			if _, err := os.Stat(compilerPath); err == nil {
				var compiler Compiler
				switch key {
				case "Clang":
					compiler = Compiler{key, filepath.Join(path, "clang.exe"), compilerPath}

				case "Clang-cl":
					compiler = Compiler{key, filepath.Join(path, "clang-cl.exe"), compilerPath}

				case "Mingw":
					compiler = Compiler{key, filepath.Join(path, "gcc.exe"), compilerPath}
				}

				Toolchain, err := getToolchain(compiler)
				if err != nil {
					logger.Log.Error(err.Error())
				}

				Toolchains = append(Toolchains, Toolchain)
			}
		}
	}

	return Toolchains, nil

}

func getToolchain(compiler Compiler) (*config.Toolchain, error) {
	cmd := exec.Command(compiler.CXXPath, "-v")
	cxxInfo, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var (
		Toolchain config.Toolchain
		version   string
		target    string
		regExpStr string
	)

	// 编译正则表达式，用于匹配版本号、目标、线程模型和安装目录
	switch compiler.Name {
	case "Clang":
		{
			regExpStr = `\s*clang version (\d+.\d+.\d+)\s*Target: ([^\s]+)\s*`
			re := regexp.MustCompile(regExpStr)

			// 输出提取的信息
			matches := re.FindStringSubmatch(string(cxxInfo))
			if len(matches) == 3 {
				version = matches[1]
				target = matches[2]
			} else {
				logger.Log.Error("No match found.")
			}

			Toolchain.Name = fmt.Sprintf("Clang %s %s", version, target)
			Toolchain.Compiler.C = compiler.CPath
			Toolchain.Compiler.CXX = compiler.CXXPath
			Toolchain.IsTrusted = true
		}

	case "Clang-cl":
		{
			regExpStr = `\s*clang version (\d+.\d+.\d+)\s*Target: ([^\s]+)\s*`
			re := regexp.MustCompile(regExpStr)

			// 输出提取的信息
			matches := re.FindStringSubmatch(string(cxxInfo))
			if len(matches) == 3 {
				version = matches[1]
				target = matches[2]
			} else {
				logger.Log.Error("No match found.")
			}

			Toolchain.Name = fmt.Sprintf("Clang-cl %s %s", version, target)
			Toolchain.Compiler.C = compiler.CPath
			Toolchain.Compiler.CXX = compiler.CPath
			Toolchain.IsTrusted = true
		}

	case "Mingw":
		{
			regExpStr = `(?s)Target:\s*([^\s]+).*gcc version (\d+.\d+.\d+)`
			re := regexp.MustCompile(regExpStr)

			// 输出提取的信息
			matches := re.FindStringSubmatch(string(cxxInfo))
			if len(matches) == 3 {
				target = matches[1]
				version = matches[2]
			} else {
				logger.Log.Error("No match found.")
			}

			Toolchain.Name = fmt.Sprintf("GCC %s %s", version, target)
			Toolchain.Compiler.C = compiler.CPath
			Toolchain.Compiler.CXX = compiler.CXXPath
			Toolchain.IsTrusted = true
		}
	}

	return &Toolchain, nil
}
