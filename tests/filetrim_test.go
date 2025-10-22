package tests

import (
	"fmt"
	"testing"

	"github.com/zelviner/cgear/utils"
)

func TestFileTrim(t *testing.T) {
	content := utils.FileTrim("test.txt")
	fmt.Println(content)
}
