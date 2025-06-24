package env

import (
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
	ui "github.com/ZEL-30/zel/ui/select"
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
