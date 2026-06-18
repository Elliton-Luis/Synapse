package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	green  = "\033[32m"
	yellow = "\033[33m"
	gray   = "\033[90m"
	red    = "\033[31m"
)

func Success(label, detail string) {
	fmt.Printf("%-35s %s✓%s %s\n", label, green, reset, detail)
}

func Error(msg string) {
	fmt.Fprintf(os.Stderr, "%s❌ Erro: %s%s\n", red, msg, reset)
}

func Fatal(msg string) {
	fmt.Fprintf(os.Stderr, "%s💥 FATAL: %s%s\n", red, msg, reset)
	os.Exit(1)
}

func Info(msg string) {
	fmt.Printf("  %s%s%s\n", gray, msg, reset)
}

func Suggestion(msg string) {
	fmt.Printf("\n✨ Sugestão da IA: %s%s%s\n", bold+green, msg, reset)
}

func Confirm() string {
	fmt.Printf("\n%s[Y]%s Confirmar  %s[R]%s Recriar  %s[N]%s Cancelar: ",
		bold+green, reset,
		bold+yellow, reset,
		bold+gray, reset,
	)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(input))
}