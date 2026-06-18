package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Elliton-Luis/synapse/internal/ui"
)

func mockLoading(label, detail string) {
	fmt.Printf("%-35s", label)
	time.Sleep(300 * time.Millisecond)
	ui.Success("", detail)
}

func main() {
	fmt.Println()

	mockLoading("⚙ Carregando configuração...", "groq (llama-3.3-70b-versatile)\n")

	mockLoading("Analisando alterações...", "3 arquivos modificados")

	fmt.Printf("\nGerando sugestão...")
	time.Sleep(500 * time.Millisecond)
	fmt.Println()

	suggestion := "feat(auth): add JWT token validation middleware"

	for {
		ui.Suggestion(suggestion)

		input := ui.Prompt()

		switch strings.ToLower(strings.TrimSpace(input)) {
		case "y":
			fmt.Println("\nCommit realizado com sucesso!")
			os.Exit(0)
		case "r":
			fmt.Printf("\nGerando nova sugestão...")
			time.Sleep(500 * time.Millisecond)
			fmt.Println()
			suggestion = "refactor(auth): extract token validation into separate middleware"
		case "n":
			fmt.Println("\nAbortado.")
			os.Exit(0)
		default:
			fmt.Println("  Digite Y, R ou N.")
		}
	}
}
