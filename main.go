package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Elliton-Luis/synapse/internal/config"
	"github.com/Elliton-Luis/synapse/internal/ui"
)

const (
	defaultGroqModel   = "llama-3.3-70b-versatile"
	defaultGeminiModel = "gemini-2.5-flash"
	maxDiffLines       = 2000
	maxCommitLength    = 200
)

const defaultPatternEN = `Strictly follow the Conventional Commits specification.
Allowed types:
- feat: Indicates that your code is adding a new feature.
- fix: Indicates that your code is solving a problem (bug fix).
- docs: Documentation changes (does not include code changes).
- test: Creation, alteration, or deletion of unit tests.
- build: Modifications to build files and dependencies.
- perf: Code changes related to performance improvements.
- style: Code formatting changes, semicolons, trailing spaces, lint... (no logic changes).
- refactor: Code refactoring that does not alter functionality.
- chore: Build tasks, admin configs, packages, .gitignore updates.
- ci: Continuous integration changes.
- raw: Changes related to config files, data, parameters.
- cleanup: Removal of commented code, unnecessary snippets, or general cleanup.
- remove: Deletion of obsolete files, directories, or features.

MANDATORY RULES:
1. Always write in the IMPERATIVE mood (e.g., "add", "remove", "fix", "update").
2. Keep it short, personal, and direct.
3. Ensure extreme precision, NEVER be generic.
4. The message MUST communicate exactly WHAT changed.
5. If necessary for context, mention the affected files or components.
6. Return ONLY the final commit message, without any quotes, markdown formatting, or extra text.`

const defaultPatternPT = `Siga estritamente a especificação do Conventional Commits.
Tipos permitidos:
- feat: Indica que seu trecho de código está incluindo um novo recurso.
- fix: Indica que seu trecho de código está solucionando um problema (bug fix).
- docs: Mudanças na documentação (Não inclui alterações em código).
- test: Criação, alteração ou exclusão de testes unitários.
- build: Modificações em arquivos de build e dependências.
- perf: Alterações de código relacionadas a performance.
- style: Formatações de código, semicolons, trailing spaces, lint... (Não inclui alterações lógicas).
- refactor: Mudanças devido a refatorações que não alteram a funcionalidade.
- chore: Atualizações de tarefas de build, configurações, pacotes, .gitignore.
- ci: Mudanças relacionadas a integração contínua.
- raw: Mudanças relacionadas a arquivos de configurações, dados, parâmetros.
- cleanup: Remoção de código comentado, trechos desnecessários ou limpeza do código-fonte.
- remove: Exclusão de arquivos, diretórios ou funcionalidades obsoletas.

REGRAS OBRIGATÓRIAS:
1. Escreva sempre no tempo verbal IMPERATIVO (ex: "adicione", "remova", "atualize", "corrija").
2. Seja curto, pessoal e direto na comunicação.
3. Garanta extrema precisão, NUNCA seja genérico na descrição.
4. A mensagem DEVE comunicar exatamente QUAL foi a mudança.
5. Se for possível e necessário para o contexto, cite os arquivos ou componentes alterados.
6. Retorne APENAS a mensagem final do commit, sem aspas, formatação markdown ou texto extra.`

var debugMode bool

func logDebug(msg string) {
	if debugMode {
		fmt.Printf("\033[90m[DEBUG] %s\033[0m\n", msg)
	}
}

func main() {
	// 1. Parsing dos argumentos/flags
	providerFlag := flag.String("provider", "auto", "Força um provedor (groq, gemini, auto)")
	modelFlag := flag.String("model", "", "Sobrescreve o modelo padrão")
	patternFlag := flag.String("pattern", "commit_pattern.md", "Arquivo de padrão a ser lido")
	flag.BoolVar(&debugMode, "debug", false, "Ativa logs detalhados de depuração")
	flag.Parse()

	// 2. Prepara o ambiente gerando arquivos necessários
	ensurePatternFiles()
	ensureGitignore()

	// 3. Carrega chaves de API
	cfg, err := config.Load()
	if err != nil {
		ui.Fatal(err.Error())
	}

	// 4. Captura e limpa o Git Diff staged
	diff, err := getGitDiff()
	if err != nil {
		ui.Fatal(err.Error())
	}
	if diff == "" {
		ui.Error("Nenhuma alteração no stage. Use 'git add' primeiro.")
		os.Exit(1)
	}

	diffLines := strings.Split(diff, "\n")
	if len(diffLines) > maxDiffLines {
		diff = strings.Join(diffLines[:maxDiffLines], "\n")
		logDebug(fmt.Sprintf("Diff truncado para as primeiras %d linhas.", maxDiffLines))
	}

	// 5. Carrega o padrão de commit desejado
	pattern := loadPattern(*patternFlag)

	basePrompt := fmt.Sprintf("Padrão exigido:\n%s\n\nAnalise o seguinte git diff e gere UMA mensagem de commit:\n%s", pattern, diff)
	retryCount := 0

	// 6. Loop de geração e interface interativa
	for {
		currentPrompt := basePrompt
		if retryCount > 0 {
			currentPrompt += fmt.Sprintf("\n\nO usuário rejeitou as %d sugestões anteriores. Gere uma mensagem de commit diferente das anteriores. Tente focar em outro aspecto das alterações ou use uma abordagem descritiva diferente.", retryCount)
		}

		// Lógica de Roteamento com Fallback Automático
		var providersToTry []string
		if *providerFlag == "auto" {
			if cfg.GroqKey != "" {
				providersToTry = append(providersToTry, "groq")
			}
			if cfg.GeminiKey != "" {
				providersToTry = append(providersToTry, "gemini")
			}
		} else {
			providersToTry = append(providersToTry, *providerFlag)
		}

		if len(providersToTry) == 0 {
			ui.Fatal("Nenhuma chave de API configurada no ambiente ou no .env")
		}

		var suggestion string
		var apiErr error

		fmt.Print("🤖 Gerando sugestão... ")
		fflush()

		for _, provider := range providersToTry {
			logDebug(fmt.Sprintf("Tentando provider: %s", provider))

			if provider == "groq" {
				if cfg.GroqKey == "" {
					apiErr = fmt.Errorf("GROQ_API_KEY não definida")
					continue
				}
				model := defaultGroqModel
				if *modelFlag != "" {
					model = *modelFlag
				}
				suggestion, apiErr = callGroqAPI(cfg.GroqKey, model, currentPrompt)
			} else {
				if cfg.GeminiKey == "" {
					apiErr = fmt.Errorf("GEMINI_API_KEY não definida")
					continue
				}
				model := defaultGeminiModel
				if *modelFlag != "" {
					model = *modelFlag
				}
				suggestion, apiErr = callGeminiAPI(cfg.GeminiKey, model, currentPrompt)
			}

			// Se a requisição foi bem sucedida, sai do loop de tentativas
			if apiErr == nil {
				break
			}
			logDebug(fmt.Sprintf("Falha no provider %s: %v", provider, apiErr))
		}

		fmt.Print("\r" + strings.Repeat(" ", 50) + "\r") // Limpa a linha do "Gerando..."

		if apiErr != nil {
			ui.Error(fmt.Sprintf("Todos os providers falharam ou estão indisponíveis. Último erro: %v", apiErr))
			os.Exit(1)
		}

		// Formatação rigorosa da resposta da IA
		suggestion = strings.NewReplacer(`"`, "", "`", "").Replace(suggestion)
		suggestion = strings.TrimSpace(suggestion)
		if lines := strings.Split(suggestion, "\n"); len(lines) > 0 {
			suggestion = lines[0]
		}
		if len(suggestion) > maxCommitLength {
			suggestion = suggestion[:maxCommitLength]
		}

		ui.Suggestion(suggestion)
		action := ui.Confirm()

		switch action {
		case "y":
			if err := executeCommit(suggestion); err != nil {
				ui.Error(fmt.Sprintf("Falha ao commitar: %v", err))
				os.Exit(1)
			}
			fmt.Println("\n🚀 Commit realizado com sucesso!")
			os.Exit(0)
		case "r":
			retryCount++
			continue
		case "n":
			fmt.Println("\nAbortado pelo usuário.")
			os.Exit(0)
		default:
			fmt.Println("  Opção inválida. Digite Y, R ou N.")
		}
	}
}

// ensurePatternFiles verifica e cria os arquivos de instrução para a IA se não existirem
func ensurePatternFiles() {
	if _, err := os.Stat("commit_pattern.md"); os.IsNotExist(err) {
		if err := os.WriteFile("commit_pattern.md", []byte(defaultPatternEN), 0644); err == nil {
			logDebug("Arquivo commit_pattern.md criado com sucesso.")
		}
	}

	if _, err := os.Stat("commit_pattern_pt.md"); os.IsNotExist(err) {
		if err := os.WriteFile("commit_pattern_pt.md", []byte(defaultPatternPT), 0644); err == nil {
			logDebug("Arquivo commit_pattern_pt.md criado com sucesso.")
		}
	}
}

// ensureGitignore adiciona os arquivos de padrão e o .env ao .gitignore automaticamente
func ensureGitignore() {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return
	}

	patternsToIgnore := []string{"commit_pattern.md", "commit_pattern_pt.md", ".env"}
	gitignorePath := ".gitignore"

	var content string
	if data, err := os.ReadFile(gitignorePath); err == nil {
		content = string(data)
	}

	var missing []string
	for _, p := range patternsToIgnore {
		if !strings.Contains(content, p) {
			missing = append(missing, p)
		}
	}

	if len(missing) > 0 {
		f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		defer f.Close()

		if content != "" && !strings.HasSuffix(content, "\n") {
			f.WriteString("\n")
		}
		f.WriteString("\n# Ignorado automaticamente pelo Synapse CLI\n")
		for _, m := range missing {
			f.WriteString(m + "\n")
		}
		logDebug("Padrões de commit e .env adicionados ao .gitignore.")
	}
}

// sanitizeDiff varre o texto ocultando dados críticos sensíveis via Regex
func sanitizeDiff(diffText string) string {
	re := regexp.MustCompile(`(?i)(password|token|api_key|secret)(\s*[:=>]\s*)(['"].*?['"]|[^\s\n\r,;]+)`)
	return re.ReplaceAllString(diffText, "$1$2[REDACTED]")
}

// getGitDiff executa o comando Git nativo para pegar alterações do stage
func getGitDiff() (string, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return "", fmt.Errorf("comando 'git' não encontrado no sistema operacional")
	}

	cmd := exec.Command("git", "diff", "--cached", "--", ".", ":(exclude).env", ":(exclude)*.env.*")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("este diretório não parece um repositório Git válido")
	}

	return sanitizeDiff(out.String()), nil
}

// loadPattern lê as regras do arquivo .md especificado ou retorna um padrão rígido padrão
func loadPattern(filename string) string {
	pwd, _ := os.Getwd()
	path := filepath.Join(pwd, filename)

	data, err := os.ReadFile(path)
	if err != nil {
		logDebug(fmt.Sprintf("Arquivo %s não localizado. Aplicando regras hardcoded.", filename))
		return defaultPatternEN
	}
	logDebug(fmt.Sprintf("Padrão de commit carregado com sucesso de: %s", path))
	return string(data)
}

// executeCommit realiza o git commit final no terminal
func executeCommit(msg string) error {
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func fflush() { _ = os.Stdout.Sync() }

// ==========================================
// INTEGRAÇÃO DE CLIENTES HTTP NATIVOS (REST)
// ==========================================

func callGroqAPI(apiKey, model, prompt string) (string, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"

	type message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	type payload struct {
		Model       string    `json:"model"`
		Messages    []message `json:"messages"`
		Temperature float64   `json:"temperature"`
	}

	bodyData := payload{
		Model: model,
		Messages: []message{
			{Role: "system", Content: "Você é um gerador estrito de mensagens de commit Git. Retorne APENAS a mensagem final. Sem aspas, sem markdown extra."},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.5,
	}

	jsonData, _ := json.Marshal(bodyData)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status code %d: %s", resp.StatusCode, string(b))
	}

	type responseSchema struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	var res responseSchema
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("resposta vazia retornada pelo Groq")
	}

	return res.Choices[0].Message.Content, nil
}

func callGeminiAPI(apiKey, model, prompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)

	type part struct {
		Text string `json:"text"`
	}
	type content struct {
		Parts []part `json:"parts"`
	}
	type genConfig struct {
		Temperature     float64 `json:"temperature"`
		MaxOutputTokens int     `json:"maxOutputTokens"`
	}
	type systemInst struct {
		Parts []part `json:"parts"`
	}
	type payload struct {
		Contents          []content  `json:"contents"`
		GenerationConfig  genConfig  `json:"generationConfig"`
		SystemInstruction systemInst `json:"systemInstruction"`
	}

	bodyData := payload{
		Contents: []content{
			{Parts: []part{{Text: prompt}}},
		},
		GenerationConfig: genConfig{
			Temperature:     0.3,
			MaxOutputTokens: 500,
		},
		SystemInstruction: systemInst{
			Parts: []part{{Text: "Você é um gerador estrito de mensagens de commit Git. Retorne APENAS a mensagem final, sem aspas ou explicações adicionais."}},
		},
	}

	jsonData, _ := json.Marshal(bodyData)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status code %d: %s", resp.StatusCode, string(b))
	}

	type responseSchema struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	var res responseSchema
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("o Gemini retornou uma resposta sem candidatos válidos")
	}

	return res.Candidates[0].Content.Parts[0].Text, nil
}