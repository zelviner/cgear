package tests

import (
	"fmt"
	"testing"
	"zel/utils"
)

func TestFileTrim(t *testing.T) {
	content := utils.FileTrim("test.txt")
	fmt.Println(content)
}
