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

func SetToolchain() int {

	logger.Log.Info("Finding Toolchain...")
	Toolchains, err := findToolchains()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	var (
		ToolchainIndex int               // 选择的 Toolchain 索引
		exitIndex      = len(Toolchains) // 退出选项的索引
	)

	logger.Log.Infof("Found %d Toolchains, please select one toolchain to use:", exitIndex)

	// 输出所有 Toolchain
	for i, Toolchain := range Toolchains {
		fmt.Printf("\t[%d] %s\n", i+1, Toolchain.Name)
	}
	fmt.Printf("\t[%d] %s\n", exitIndex+1, "Exit")

	// 选择 Toolchain
	_, err = fmt.Scanln(&ToolchainIndex)
	ToolchainIndex--
	if err != nil {
		logger.Log.Error(err.Error())
	}
	if ToolchainIndex < 0 || ToolchainIndex > exitIndex {
		logger.Log.Error("Invalid Toolchain index")
	}

	if ToolchainIndex == exitIndex {
		logger.Log.Infof("Exit")
		os.Exit(0)
	}

	config.Conf.Toolchain = Toolchains[ToolchainIndex]
	config.SaveConfig()
	logger.Log.Successf("Successfully set Toolchain: %s", config.Conf.Toolchain.Name)

	return 0
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
