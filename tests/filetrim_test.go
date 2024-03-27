package tests

import (
	"fmt"
	"testing"

	"github.com/ZEL-30/zel/utils"
)

func TestFileTrim(t *testing.T) {
	content := utils.FileTrim("test.txt")
	fmt.Println(content)
}
