package internal

import "fmt"

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

func LogInfo(prefix, format string, a ...any) {
	fmt.Printf(ColorCyan+"[%s] "+ColorReset+format+"\n", append([]any{prefix}, a...)...)
}

func LogSuccess(prefix, format string, a ...any) {
	fmt.Printf(ColorCyan+"[%s] "+ColorGreen+format+ColorReset+"\n", append([]any{prefix}, a...)...)
}

func LogWarn(prefix, format string, a ...any) {
	fmt.Printf(ColorCyan+"[%s] "+ColorYellow+format+ColorReset+"\n", append([]any{prefix}, a...)...)
}

func LogError(prefix, format string, a ...any) {
	fmt.Printf(ColorCyan+"[%s] "+ColorRed+format+ColorReset+"\n", append([]any{prefix}, a...)...)
}
