package env

import (
	"github.com/zelviner/cgear/config"
	"github.com/zelviner/cgear/logger"
	ui "github.com/zelviner/cgear/ui/select"
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
