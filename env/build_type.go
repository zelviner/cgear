package env

import (
	"fmt"

	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
	ui "github.com/ZEL-30/zel/ui/select"
)

func SetBuildType() {
	options := []string{
		"Debug", "Release", "RelWithDebInfo", "MinSizeRel",
	}
	choice, err := ui.SelectOption("请选择构建类型：", options)
	if err != nil {
		msg := fmt.Sprintf("选择构建类型失败：%s", err.Error())
		logger.Log.Error(msg)
		return
	}
	config.Conf.BuildType = choice
	config.SaveConfig()
	logger.Log.Successf("构建类型设置为：%s", config.Conf.BuildType)
}
