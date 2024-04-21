package tests

import (
	"fmt"
	"testing"

	"golang.org/x/sys/windows/registry"
)

func TestGetMSVC(t *testing.T) {
	// 定义MSVC编译器在注册表中的路径
	msvcRegistryPath := `SOFTWARE\WOW6432Node\Microsoft\VisualStudio\SxS\VC7`

	// 打开注册表项
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, msvcRegistryPath, registry.READ)
	if err != nil {
		fmt.Println("无法打开注册表项:", err)
		return
	}
	defer k.Close()

	// 读取默认的MSVC编译器版本
	defaultCompilerVersion, _, err := k.GetStringValue("Default")
	if err != nil {
		fmt.Println("无法读取默认编译器版本:", err)
		return
	}

	// 构建MSVC编译器路径
	msvcPath := fmt.Sprintf("C:\\Program Files (x86)\\Microsoft Visual Studio\\2019\\Community\\VC\\Tools\\MSVC\\%s", defaultCompilerVersion)

	// 打印MSVC编译器路径
	fmt.Println("MSVC编译器路径:", msvcPath)
}
