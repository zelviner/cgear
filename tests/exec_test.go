package tests

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestExec(t *testing.T) {
	cmd := exec.Command("cmake")
	out, err := cmd.Output()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(out))

	val := os.Getenv("CPATH")
	fmt.Println(val)

	t.Error("test")

}
