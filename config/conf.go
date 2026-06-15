package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

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

// LoadConfig 加载 cgear 配置。
// 它在当前路径中查找 cgear.json 或 Cgearfile，如果找不到，则保留默认配置。
func LoadConfig() {
	currentPath := utils.GetCgearWorkPath()
	ensureCgearHome()

	jsonPath := filepath.Join(currentPath, "cgear.json")
	yamlPath := filepath.Join(currentPath, "Cgearfile")
	var err error
	switch {
	case utils.IsExist(jsonPath):
		err = parseJSON(jsonPath, &Conf)
	case utils.IsExist(yamlPath):
		err = parseYAML(yamlPath, &Conf)
	}
	if err != nil {
		logger.Log.Errorf("Failed to parse configuration file: %s", err)
	}
	Conf.ProjectPath = currentPath

	// 检查格式版本
	if Conf.Version != confVer {
		logger.Log.Warn("Your configuration file is outdated. Please consider updating it.")
		logger.Log.Hint("Check the latest version of cgear's configuration file.")
	}

}

func ensureCgearHome() {
	if os.Getenv("CGEAR_HOME") != "" {
		return
	}

	executable, err := os.Executable()
	if err != nil {
		logger.Log.Warnf("Failed to determine CGEAR_HOME: %s", err)
		return
	}
	if err := os.Setenv("CGEAR_HOME", cgearHomeFromExecutable(executable)); err != nil {
		logger.Log.Warnf("Failed to set CGEAR_HOME for this process: %s", err)
	}
}

func cgearHomeFromExecutable(executable string) string {
	dir := filepath.Dir(executable)
	if strings.EqualFold(filepath.Base(dir), "bin") {
		return filepath.Dir(dir)
	}
	return dir
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
		return err
	}

	return os.WriteFile(filepath.Join(projectPath, "cgear.json"), configJson, 0644)
}
