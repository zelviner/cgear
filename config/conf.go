package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/ZEL-30/zel/logger"
)

const confVer = 0

const (
	Version       = "1.0.0"
	GitRemotePath = "github.com/beego/v2"
)

type Config struct {
	Version            int
	WatchExts          []string `json:"watch_exts" yaml:"watch_exts"`
	WatchExtsStatic    []string `json:"watch_exts_static" yaml:"watch_exts_static"`
	ZelPath            string   `json:"zel_path" yaml:"zel_path"`
	Kit                *Kit     `json:"kit" yaml:"kit"`
	BuildMode          string   `json:"build_mode" yaml:"build_mode"`
	TestMode           string   `json:"test_mode" yaml:"test_mode"`
	Database           database
	EnableReload       bool `json:"enable_reload" yaml:"enable_reload"`
	EnableNotification bool `json:"enable_notification" yaml:"enable_notification"`
}

type Kit struct {
	Name      string   `json:"name" yaml:"name"`
	Compiler  Compiler `json:"compilers" yaml:"compilers"`
	IsTrusted bool     `json:"isTrusted" yaml:"isTrusted"`
}

// 编译器
type Compiler struct {
	C   string `json:"C" yaml:"C"`
	CXX string `json:"CXX" yaml:"CXX"`
}

// database 保存数据库连接信息
type database struct {
	Driver string
	Conn   string
	Dir    string
}

var Conf = Config{
	WatchExts:          []string{".h", ".hpp", ".cpp"},
	WatchExtsStatic:    []string{".html", ".tpl", ".js", ".css"},
	BuildMode:          "Debug",
	Kit:                nil,
	Database:           database{Driver: "mysql"},
	EnableNotification: true,
}

// LoadConfig 加载 Zel tool配置。
// 它在当前路径中查找Zelfile或zel.json，如果找不到，则返回默认配置
func LaodConfig() {
	currentPath, err := os.Getwd()
	if err != nil {
		logger.Log.Error(err.Error())
	}

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

	// 设置ZelPath环境变量
	if zelPath := os.Getenv("ZELPATH"); zelPath == "" {
		userProfile := os.Getenv("USERPROFILE")
		Conf.ZelPath = filepath.Join(userProfile, "zel")

		cmd := exec.Command("SETX", "ZELPATH", Conf.ZelPath)
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

func SaveConfig() error {
	configJson, err := json.MarshalIndent(Conf, "", "\t")
	if err != nil {
		logger.Log.Error(err.Error())
	}

	err = os.WriteFile("zel.json", configJson, 0644)
	return err
}
