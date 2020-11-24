package utils

import (
	"fmt"
	"strings"
)

func MakeLabelColor(name string) string {
	return intToRGB(hashCode(name))
}

func intToRGB(i int) string {
	color := i & 0x00FFFFFF
	hex := strings.ToUpper(fmt.Sprintf("%x", color))
	zero := "00000"
	return zero[0:6 - len(hex)] + hex
}

func hashCode(str string) int {
	var hash = 0
	for i := 0; i < len(str); i++ {
		hash = int(str[i]) + ((hash << 5) - hash)
	}
	return hash
}