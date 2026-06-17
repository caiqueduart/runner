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

// exibe uma mensagem no terminal com um prefixo colorido em Ciano.
func LogFeedback(prefix, format string, a ...any) {
	fmt.Printf(ColorCyan+"[%s] "+ColorReset+format+"\n", append([]any{prefix}, a...)...)
}
