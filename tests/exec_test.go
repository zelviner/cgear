package tests

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/ZEL-30/zel/logger/colors"
)

func TestExec(t *testing.T) {
	cmd := exec.Command("cmake")
	out, err := cmd.Output()
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(string(out))

	fmt.Println(colors.Black(string(out)))
	fmt.Println(colors.Bold(string(out)))

	val := os.Getenv("CPATH")
	fmt.Println(val)

}
