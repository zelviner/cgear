package system

import (
	"os"
	"os/user"
	"path/filepath"
)

// Zel System Params ...
var (
	Usr, _     = user.Current()
	ZelHome    = filepath.Join(Usr.HomeDir, "/.zel")
	CurrentDir = GetCurrentDirectory()
	GoPath     = os.Getenv("GOPATH")
)

func GetCurrentDirectory() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	return ""
}
