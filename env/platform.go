package env

import (
	"github.com/zelviner/cgear/config"
	"github.com/zelviner/cgear/logger"
	ui "github.com/zelviner/cgear/ui/select"
)

func SetPlatform() {
	platforms := []string{"x86", "x64"}

	platform, cancelled, err := ui.ListOption("Plase select platform: ", platforms, func(p string) string { return p })
	if err != nil {
		logger.Log.Errorf("Failed to select platform: %v", err)
		return
	}

	if cancelled {
		logger.Log.Info("Cancelled selecting platform")
		return
	}

	config.Conf.Platform = platform
	logger.Log.Infof("Selected platform: %s", platform)
}
