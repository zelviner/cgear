package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/zelviner/cgear/logger"
	"github.com/zelviner/cgear/utils"
)

const confVer = 0

const (
	Version = "2.0.0"
)

type Config struct {
	Version             int
	Toolchain           *Toolchain `json:"toolchain" yaml:"toolchain"`                       // 编译工具链
	Generator           string     `json:"generator" yaml:"generator"`                       // 生成器
	Platform            string     `json:"platform" yaml:"platform"`                         // 编译架构
	BuildType           string     `json:"build_type" yaml:"build_type"`                     // 编译类型
	ProjectType         string     `json:"project_type" yaml:"project_type"`                 // 项目类型
	ProjectPath         string     `json:"project_path" yaml:"project_path"`                 // 项目路径
	RuntimeDependencies []string   `json:"runtime_dependencies" yaml:"runtime_dependencies"` // 运行时依赖动态库
}

type Toolchain struct {
	Name     string   `json:"name" yaml:"name"`
	Compiler Compiler `json:"compilers" yaml:"compilers"`
	IsMSVC   bool     `json:"is_msvc" yaml:"is_msvc"`
}

// 编译器
type Compiler struct {
	C   string `json:"C" yaml:"C"`
	CXX string `json:"CXX" yaml:"CXX"`
}

var Conf = Config{
	BuildType:           "Debug",
	Toolchain:           nil,
	RuntimeDependencies: []string{"input dynamic libraries here"},
}

// LoadConfig 加载 cgear tool配置。
// 它在当前路径中查找cgear file 或 cgear.json，如果找不到，则返回默认配置
func LaodConfig() {
	currentPath := utils.GetCgearWorkPath()

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
		case "cgear.json":
			{
				err = parseJSON(filepath.Join(currentPath, file.Name()), &Conf)
				if err != nil {
					logger.Log.Errorf("Failed to parse JSON file %s", err)
				}
				break
			}
		case "Cgearfile":
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
		logger.Log.Hint("Check the latest version of cgear's configuration file.")
	}

	// 设置 CGEAR_HOME 环境变量
	if cgearHome := os.Getenv("CGEAR_HOME"); cgearHome == "" {
		// 获取程序所在的路径
		programPath := utils.GetCgearWorkPath()

		cmd := exec.Command("SETX", "CGEAR_HOME", programPath)
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

	err = os.WriteFile(filepath.Join(projectPath, "cgear.json"), configJson, 0644)
	return err
}
