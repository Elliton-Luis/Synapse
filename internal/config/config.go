package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const envFile = ".env"

type Config struct {
	GroqKey   string
	GeminiKey string
	Lang      string
}

func Load() (*Config, error) {
	// Prioriza variáveis de ambiente do sistema
	groqEnv := os.Getenv("GROQ_API_KEY")
	geminiEnv := os.Getenv("GEMINI_API_KEY")
	langEnv := os.Getenv("SYNAPSE_LANG")

	if groqEnv != "" || geminiEnv != "" {
		if langEnv == "" {
			langEnv = "en"
		}
		return &Config{GroqKey: groqEnv, GeminiKey: geminiEnv, Lang: langEnv}, nil
	}

	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return setup()
	}
	return readEnv()
}

func setup() (*Config, error) {
	fmt.Println("\n🔧 Primeira execução! Vamos configurar o Synapse.\n")
	fmt.Println("  Forneça as chaves de API necessárias (pressione Enter para pular caso não use o provider).")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("  GROQ_API_KEY: ")
	groqKey, _ := reader.ReadString('\n')
	groqKey = strings.TrimSpace(groqKey)

	fmt.Print("  GEMINI_API_KEY: ")
	geminiKey, _ := reader.ReadString('\n')
	geminiKey = strings.TrimSpace(geminiKey)

	if groqKey == "" && geminiKey == "" {
		return nil, fmt.Errorf("você precisa fornecer pelo menos uma API Key")
	}

	// Salva com o idioma inglês por padrão
	content := fmt.Sprintf("GROQ_API_KEY=%s\nGEMINI_API_KEY=%s\nSYNAPSE_LANG=en\n", groqKey, geminiKey)
	if err := os.WriteFile(envFile, []byte(content), 0600); err != nil {
		return nil, err
	}

	fmt.Printf("\n✓ .env criado com sucesso.\n\n")
	return &Config{GroqKey: groqKey, GeminiKey: geminiKey, Lang: "en"}, nil
}

func readEnv() (*Config, error) {
	file, err := os.Open(envFile)
	if err != nil {
		return nil, err
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

	lang := values["SYNAPSE_LANG"]
	if lang == "" {
		lang = "en"
	}

	return &Config{
		GroqKey:   values["GROQ_API_KEY"],
		GeminiKey: values["GEMINI_API_KEY"],
		Lang:      lang,
	}, nil
}

// SetLang atualiza o idioma preferido diretamente no arquivo .env
func SetLang(lang string) error {
	file, err := os.Open(envFile)
	var lines []string
	found := false

	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(strings.TrimSpace(line), "SYNAPSE_LANG=") {
				lines = append(lines, fmt.Sprintf("SYNAPSE_LANG=%s", lang))
				found = true
			} else {
				lines = append(lines, line)
			}
		}
		file.Close()
	}

	// Se não achou a linha (ou se o arquivo não existia), adiciona a configuração
	if !found {
		lines = append(lines, fmt.Sprintf("SYNAPSE_LANG=%s", lang))
	}

	return os.WriteFile(envFile, []byte(strings.Join(lines, "\n")+"\n"), 0600)
}