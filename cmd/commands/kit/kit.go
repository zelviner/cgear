package kit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/cmd/commands/version"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
)

type Compiler struct {
	Name string
	Path string
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
	compilers     = []Compiler{}
	compilerTypes = map[string]string{
		"Clang for C":   "clang.exe",
		"Clang for C++": "clang++.exe",
		"Mingw for C":   "gcc.exe",
		"Mingw for C++": "g++.exe",
	}
	kits []config.Kit
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
				err := appendKit(Compiler{key, compilerPath})
				if err != nil {
					logger.Log.Fatal(err.Error())
				}
			}
		}
	}

	logger.Log.Successf("Successfully set kit: %s.", "clang")
	return 0

}

func appendKit(compiler Compiler) error {
	cmd := exec.Command(compiler.Path, "-v")
	logger.Log.Info(cmd.String())
	info, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Println(string(info))
	return err
}
