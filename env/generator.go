package env

import (
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
	ui "github.com/ZEL-30/zel/ui/select"
)

func SetGenerator() {

	generators := []string{
		"Ninja",
		"Visual Studio 17 2022",
	}

	logger.Log.Info("")

	generator, cancelled, err := ui.ListOption("Please select a generator:", generators, func(g string) string { return g })
	if err != nil {
		logger.Log.Errorf("Failed to select generator: %v", err)
		return
	}

	if cancelled {
		logger.Log.Info("Cancelled Selecting generator.")
		return
	}

	config.Conf.Generator = generator
	logger.Log.Infof("Generator set to: %s", generator)
}
