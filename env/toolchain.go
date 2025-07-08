package env

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
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
	"MSVC":     "cl.exe",
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

	// 查找 MSVC 编译器
	toolchains, err := findMSVCCompiler()
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	for _, toolchain := range toolchains {
		Toolchains = append(Toolchains, toolchain)
	}

	return Toolchains, nil

}

func getToolchain(compiler Compiler) (*config.Toolchain, error) {
	cmd := exec.Command(compiler.CXXPath, "-v")
	cxxInfo, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to run %s: %v", compiler.CXXPath, err)
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

func findMSVCCompiler() (toolchains []*config.Toolchain, err error) {
	// 1. 执行 vswhere 获取 VS 安装路径
	vswhere := filepath.Join(os.Getenv("ProgramFiles(x86)"), "Microsoft Visual Studio", "Installer", "vswhere.exe")
	cmd := exec.Command(vswhere,
		"-latest",
		"-products", "*",
		"-requires", "Microsoft.VisualStudio.Component.VC.Tools.x86.x64",
		"-property", "installationPath")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("vswhere 运行失败: %w", err)
	}
	vsPath := strings.TrimSpace(out.String())

	// 2. 构造 MSVC 根路径
	msvcRoot := filepath.Join(vsPath, "VC", "Tools", "MSVC")

	// 3. 读取所有子目录，选最新版本
	entries, err := os.ReadDir(msvcRoot)
	if err != nil {
		logger.Log.Infof("msvcRoot: %s", msvcRoot)
		return nil, fmt.Errorf("读取 MSVC 目录失败: %w", err)
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			versions = append(versions, entry.Name())
		}
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("未找到任何 MSVC 版本")
	}
	sort.Strings(versions) // 从小到大排序
	latest := versions[len(versions)-1]

	tryAddToolchain(&toolchains, "Visual Studio Community 2022 Release - amd64", filepath.Join(msvcRoot, latest, "bin", "Hostx64", "x64", "cl.exe"))
	tryAddToolchain(&toolchains, "Visual Studio Community 2022 Release - amd64_86", filepath.Join(msvcRoot, latest, "bin", "Hostx64", "x86", "cl.exe"))
	tryAddToolchain(&toolchains, "Visual Studio Community 2022 Release - x86", filepath.Join(msvcRoot, latest, "bin", "Hostx86", "x86", "cl.exe"))
	tryAddToolchain(&toolchains, "Visual Studio Community 2022 Release - x86_64", filepath.Join(msvcRoot, latest, "bin", "Hostx86", "x64", "cl.exe"))

	return toolchains, nil
}

func tryAddToolchain(toolchains *[]*config.Toolchain, name, clPath string) {
	if _, err := os.Stat(clPath); err == nil {
		*toolchains = append(*toolchains, &config.Toolchain{
			Name: name,
			Compiler: config.Compiler{
				C:   clPath,
				CXX: clPath,
			},
		})
	} else {
		logger.Log.Warnf("跳过: %s 不存在", clPath)
	}
}
