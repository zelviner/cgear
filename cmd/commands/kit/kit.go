package kit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/cmd/commands/version"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
)

type Compiler struct {
	Name    string
	CPath   string
	CXXPath string
}

type CompilerInfo struct {
	Version      string
	Target       string
	ThreadModel  string
	InstalledDir string
}

var CmdKit = &commands.Command{
	UsageLine: "kit",
	Short:     "Select a kit for your C++ project",
	Long: `▶ {{"To find C++ compilers available on your system"|bold}}

     $ zel kit find
	`,

	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    SetKit,
}

var (
	compilerTypes = map[string]string{
		"Clang": "clang++.exe",
		"Mingw": "g++.exe",
	}
	kits []*config.Kit
)

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdKit)
}

func SetKit(cmd *commands.Command, args []string) int {

	logger.Log.Info("Finding kits...")

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

	config.SaveKitsConfig(kits)

	logger.Log.Successf("Successfully set kit: %s.", "clang")
	return 0

}

func getKit(compiler Compiler) (*config.Kit, error) {

	cmd := exec.Command(compiler.CXXPath, "-v")
	logger.Log.Info(cmd.String())
	cxxInfo, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	fmt.Println(string(cxxInfo))

	var (
		compilerInfo CompilerInfo
		kit          config.Kit
	)

	switch compiler.Name {
	case "Clang":
		{
			// 编译正则表达式，用于匹配版本号、目标、线程模型和安装目录
			clangRegExpStr := `clang version ([^\s]+)\s+Target: ([^\s]+)\s+Thread model: ([^\s]+)\s+InstalledDir: (.+)`
			re := regexp.MustCompile(clangRegExpStr)

			// 使用正则表达式提取信息
			matches := re.FindStringSubmatch(string(cxxInfo))

			// 输出提取的信息
			if len(matches) == 5 {
				compilerInfo.Version = matches[1]
				compilerInfo.Target = matches[2]
				compilerInfo.ThreadModel = matches[3]
				compilerInfo.InstalledDir = matches[4]
			} else {
				logger.Log.Error("未找到匹配项")
			}
			kit.Name = fmt.Sprintf("Clang %s %s", compilerInfo.Version, compilerInfo.Target)
			kit.Compiler.C = compiler.CPath
			kit.Compiler.CXX = compiler.CXXPath
			kit.IsTrusted = true
		}
	case "Mingw":
		{

		}

	}

	return &kit, nil
}
