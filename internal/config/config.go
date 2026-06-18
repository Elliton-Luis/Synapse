package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const envFile = ".env"

var defaultModels = map[string]string{
	"groq":       "llama-3.3-70b-versatile",
	"openrouter": "meta-llama/llama-3.3-70b-instruct:free",
}

type Config struct {
	Provider string
	APIKey   string
	Model    string
}

// Load tenta carregar o .env. Se não existir, inicia o setup interativo.
func Load() (*Config, error) {
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return setup()
	}
	return readEnv()
}

// setup pergunta as configs no terminal e salva o .env
func setup() (*Config, error) {
	fmt.Println("\n🔧 Primeira execução! Vamos configurar o Synapse.\n")

	provider := prompt("Provider padrão (groq/openrouter)", "groq")
	if provider != "groq" && provider != "openrouter" {
		return nil, fmt.Errorf("provider '%s' inválido. Use 'groq' ou 'openrouter'", provider)
	}

	apiKey := prompt("API Key", "")
	if apiKey == "" {
		return nil, fmt.Errorf("API key não pode ser vazia")
	}

	model := defaultModels[provider]

	if err := writeEnv(provider, apiKey, model); err != nil {
		return nil, fmt.Errorf("erro ao salvar .env: %w", err)
	}

	if err := ensureGitignore(); err != nil {
		return nil, fmt.Errorf("erro ao atualizar .gitignore: %w", err)
	}

	fmt.Printf("\n✓ .env criado\n")
	fmt.Printf("✓ .gitignore atualizado\n\n")

	return &Config{Provider: provider, APIKey: apiKey, Model: model}, nil
}

// readEnv lê o .env e retorna a config
func readEnv() (*Config, error) {
	file, err := os.Open(envFile)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir .env: %w", err)
	}
	defer file.Close()

	values := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			values[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	provider := values["PROVIDER"]
	apiKey := values["API_KEY"]
	model := values["MODEL"]

	if provider == "" || apiKey == "" {
		return nil, fmt.Errorf("PROVIDER ou API_KEY ausente no .env\nDelete o .env e rode novamente para reconfigurar")
	}

	if model == "" {
		model = defaultModels[provider]
	}

	return &Config{Provider: provider, APIKey: apiKey, Model: model}, nil
}

// writeEnv salva o arquivo .env
func writeEnv(provider, apiKey, model string) error {
	content := fmt.Sprintf("PROVIDER=%s\nAPI_KEY=%s\nMODEL=%s\n", provider, apiKey, model)
	return os.WriteFile(envFile, []byte(content), 0600)
}

// ensureGitignore adiciona .env ao .gitignore se necessário
func ensureGitignore() error {
	gitignorePath := ".gitignore"

	content := ""
	if data, err := os.ReadFile(gitignorePath); err == nil {
		content = string(data)
	}

	if strings.Contains(content, ".env") {
		return nil
	}

	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if content != "" && !strings.HasSuffix(content, "\n") {
		f.WriteString("\n")
	}
	f.WriteString("\n# Synapse\n.env\n")
	return nil
}

func prompt(label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("  %s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("  %s: ", label)
	}

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal
	}
	return input
}