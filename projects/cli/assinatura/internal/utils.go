package internal

import (
	"fmt"
)

func PrintError(format string, a ...any) {
	fmt.Printf(ColorRed+format+ColorReset, a...)
}
