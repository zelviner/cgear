package kit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
)

var CmdKit = &commands.Command{
	UsageLine: "kit",
	Short:     "Select a kit for your C++ project",
	Long: `
Select a kit for your C++ project
	`,

	PreRun: nil,
	Run:    SelectKit,
}

var (
	kits []config.Kit
)

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdKit)
}

func SelectKit(cmd *commands.Command, args []string) int {
	userProfile := os.Getenv("USERPROFILE")
	kitPath := filepath.Join(userProfile, "AppData", "Local", "CMakeTools", "cmake-tools-kits.json")
	file, err := os.Open(kitPath)
	if err != nil {
		return -1
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return -1
	}
	kitsData := string(bytes)

	err = json.Unmarshal([]byte(kitsData), &kits)
	if err != nil {
		return -1
	}

	var (
		kitIndex  int         // 选择的 kit 索引
		exitIndex = len(kits) // 退出选项的索引
	)

	// 输出所有 kit
	logger.Log.Infof("Found %d kits:", len(kits))
	for i, kit := range kits {
		fmt.Printf("\t[%d] %s\n", i+1, kit.Name)
		fmt.Println("\t\t" + kit.Compiler.CXX)
	}
	fmt.Printf("\t[%d] %s\n", exitIndex, "Exit")

	// 选择 kit
	logger.Log.Info("Please select one kit to use:")
	_, err = fmt.Scanln(&kitIndex)
	kitIndex--
	if err != nil {
		return -1
	}
	if kitIndex < 0 || kitIndex > exitIndex {
		return -1
	}

	if kitIndex == exitIndex {
		logger.Log.Infof("Exit")
		os.Exit(0)
	}

	kit := kits[kitIndex]

	config.Conf.Kit = kit

	logger.Log.Infof("Selected kit: %s", kit.Name)

	// 保存配置
	config.SaveConfig()

	return 0
}
