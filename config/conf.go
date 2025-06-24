package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/ZEL-30/zel/logger"
	"github.com/ZEL-30/zel/utils"
)

const confVer = 0

const (
	Version       = "1.0.0"
	GitRemotePath = "github.com/beego/v2"
)

type Config struct {
	Version         int
	Toolchain       *Toolchain `json:"toolchain" yaml:"toolchain"`       // 编译工具链
	Generator       string     `json:"generator" yaml:"generator"`       // 生成器
	Platform        string     `json:"platform" yaml:"platform"`         // 编译架构
	BuildType       string     `json:"build_type" yaml:"build_type"`     // 编译类型
	ProjectType     string     `json:"project_type" yaml:"project_type"` // 项目类型
	ProjectPath     string     `json:"project_path" yaml:"project_path"` // 项目路径
	WatchExts       []string   `json:"watch_exts" yaml:"watch_exts"`
	WatchExtsStatic []string   `json:"watch_exts_static" yaml:"watch_exts_static"`
}

type Toolchain struct {
	Name      string   `json:"name" yaml:"name"`
	Compiler  Compiler `json:"compilers" yaml:"compilers"`
	IsTrusted bool     `json:"isTrusted" yaml:"isTrusted"`
}

// 编译器
type Compiler struct {
	C   string `json:"C" yaml:"C"`
	CXX string `json:"CXX" yaml:"CXX"`
}

var Conf = Config{
	WatchExts:       []string{".h", ".hpp", ".cpp"},
	WatchExtsStatic: []string{".html", ".tpl", ".js", ".css"},
	BuildType:       "Debug",
	Toolchain:       nil,
}

// LoadConfig 加载 Zel tool配置。
// 它在当前路径中查找Zelfile或zel.json，如果找不到，则返回默认配置
func LaodConfig() {
	currentPath := utils.GetZelWorkPath()

	dir, err := os.Open(currentPath)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	for _, file := range files {
		switch file.Name() {
		case "zel.json":
			{
				err = parseJSON(filepath.Join(currentPath, file.Name()), &Conf)
				if err != nil {
					logger.Log.Errorf("Failed to parse JSON file %s", err)
				}
				break
			}
		case "Zelfile":
			{
				err = parseYAML(filepath.Join(currentPath, file.Name()), &Conf)
				if err != nil {
					logger.Log.Errorf("Failed to parse YAML file: %s", err)
				}
				break
			}
		}
	}

	// 检查格式版本
	if Conf.Version != confVer {
		logger.Log.Warn("Your configuartion file is outdated. Please do consider updating is.")
		logger.Log.Hint("Check the latest version of zel's configuration file.")
	}

	// 设置 ZEL_HOME 环境变量
	if zelHome := os.Getenv("ZEL_HOME"); zelHome == "" {
		// 获取程序所在的路径
		programPath := utils.GetZelWorkPath()

		cmd := exec.Command("SETX", "ZEL_HOME", programPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			logger.Log.Fatal(err.Error())
		}
	}

}

func parseJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	return err
}

func parseYAML(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, v)
	return err
}

func SaveConfig(projectPath string) error {
	configJson, err := json.MarshalIndent(Conf, "", "\t")
	if err != nil {
		logger.Log.Error(err.Error())
	}

	err = os.WriteFile(filepath.Join(projectPath, "zel.json"), configJson, 0644)
	return err
}
