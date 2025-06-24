package env

import (
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
	ui "github.com/ZEL-30/zel/ui/select"
)

func SetBuildType() {
	buildTypes := []string{"Debug", "Release", "RelWithDebInfo", "MinSizeRel"}

	buildType, cancelled, err := ui.ListOption("Please select build type: ", buildTypes, func(s string) string { return s })
	if err != nil {
		logger.Log.Errorf("Failed to select build type: %v", err)
		return
	}
	if cancelled {
		logger.Log.Info("Build type setting cancelled")
		return
	}

	config.Conf.BuildType = buildType
	logger.Log.Successf("Build type set to: %s", buildType)
}
