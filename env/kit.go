package env

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
)

type Compiler struct {
	Name    string
	CPath   string
	CXXPath string
}

var (
	compilerTypes = map[string]string{
		"Clang":    "clang++.exe",
		"Clang-cl": "clang-cpp.exe",
		"Mingw":    "g++.exe",
	}
	kits      []*config.Kit
	findAgain bool
)

func SetKit(cmd *commands.Command, args []string) int {

	cmd.Flag.Parse(args)

	if findAgain || len(config.Conf.Kits) == 0 {
		logger.Log.Info("Finding kits...")
		err := findKits()
		if err != nil {
			logger.Log.Fatal(err.Error())
		}
	}

	var (
		kitIndex  int                     // 选择的 kit 索引
		exitIndex = len(config.Conf.Kits) // 退出选项的索引
	)

	logger.Log.Infof("Found %d Kits, please select one kit to use:", exitIndex)

	// 输出所有 kit
	for i, kit := range config.Conf.Kits {
		fmt.Printf("\t[%d] %s\n", i+1, kit.Name)
	}
	fmt.Printf("\t[%d] %s\n", exitIndex+1, "Exit")

	// 选择 kit
	_, err := fmt.Scanln(&kitIndex)
	kitIndex--
	if err != nil {
		logger.Log.Error(err.Error())
	}
	if kitIndex < 0 || kitIndex > exitIndex {
		logger.Log.Error("Invalid kit index")
	}

	if kitIndex == exitIndex {
		logger.Log.Infof("Exit")
		os.Exit(0)
	}

	config.Conf.Kit = config.Conf.Kits[kitIndex]
	config.SaveConfig()
	logger.Log.Successf("Successfully set kit: %s", config.Conf.Kit.Name)

	return 0
}

func findKits() (err error) {

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

				kit, err := getKit(compiler)
				if err != nil {
					logger.Log.Error(err.Error())
				}

				kits = append(kits, kit)
			}
		}
	}

	config.Conf.Kits = kits

	return err

}

func getKit(compiler Compiler) (*config.Kit, error) {
	cmd := exec.Command(compiler.CXXPath, "-v")
	cxxInfo, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var (
		kit       config.Kit
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

			kit.Name = fmt.Sprintf("Clang %s %s", version, target)
			kit.Compiler.C = compiler.CPath
			kit.Compiler.CXX = compiler.CXXPath
			kit.IsTrusted = true
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

			kit.Name = fmt.Sprintf("Clang-cl %s %s", version, target)
			kit.Compiler.C = compiler.CPath
			kit.Compiler.CXX = compiler.CXXPath
			kit.IsTrusted = true
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

			kit.Name = fmt.Sprintf("GCC %s %s", version, target)
			kit.Compiler.C = compiler.CPath
			kit.Compiler.CXX = compiler.CXXPath
			kit.IsTrusted = true
		}
	}

	return &kit, nil
}
