package tests

import (
	"fmt"
	"testing"

	"github.com/zelviner/cgear/utils"
)

func getTestProgramName(testName string) string {
	var result []byte

	for i, letter := range testName {
		if letter >= 65 && letter <= 90 {
			if i == 0 {
				result = append(result, byte(letter+32))
				continue
			}
			result = append(result, '-')
			result = append(result, byte(letter+32))
		} else {
			result = append(result, byte(letter))
		}
	}
	return string(result)
}

func TestWord(t *testing.T) {
	str := "personal-data-in-out"

	words := utils.CapitalizeFirstLetter(str)

	fmt.Println(words)

	bytes := getTestProgramName(words)
	fmt.Println(bytes)
}
