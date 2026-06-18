package ui

import "fmt"

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	gray   = "\033[90m"
)

func Success(label, detail string) {
	fmt.Printf("%-35s %s✓%s %s\n", label, green, reset, detail)
}

func Step(icon, label string) {
	fmt.Printf("\n%s  %s%s%s\n", icon, bold, label, reset)
}

func Suggestion(msg string) {
	fmt.Printf("\nSugestão: %s%s%s\n", bold+green, msg, reset)
}

func Prompt() string {
	fmt.Printf("\n%s[Y]%s Confirmar  %s[R]%s Recriar  %s[N]%s Cancelar: ",
		bold+green, reset,
		bold+yellow, reset,
		bold+gray, reset,
	)
	var input string
	fmt.Scanln(&input)
	return input
}
