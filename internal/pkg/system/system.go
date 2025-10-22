package system

import (
	"os"
	"os/user"
	"path/filepath"
)

// Cgear System Params ...
var (
	Usr, _     = user.Current()
	CgearHome  = filepath.Join(Usr.HomeDir, "/.cgear")
	CurrentDir = GetCurrentDirectory()
	GoPath     = os.Getenv("GOPATH")
)

func GetCurrentDirectory() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	return ""
}
