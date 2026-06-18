# Synapse CLI
## Introdução
A Synapse é uma ferramenta de linha de comando (CLI) para gerenciar commits no Git, construída com a biblioteca padrão do Go, focando em simplicidade, zero dependências externas e economia de tokens para garantir mensagens de commit rápidas e contextualmente precisas.

### Arquitetura e Filosofia
A ferramenta é projetada para atuar como uma ponte altamente otimizada entre o repositório local e as APIs de LLMs. Ela não depende de frameworks pesados ou SDKs de terceiros.

#### Decisões Arquiteturais Chave
* **Zero Dependências:** Compilada em um único binário estático. Todas as rotas HTTP, interações com o sistema operacional e sanitizações de regex são tratadas nativamente pelo Go.
* **Economia de Tokens:** Modifica a saída padrão do `git diff` para ignorar arquivos de bloqueio (`package-lock.json`, `composer.lock`, `go.sum`, etc.) e remove metadados internos do Git, garantindo que a janela de contexto da LLM seja populada apenas por alterações de código reais.
* **Segurança de Aplicativos e Sanitização:** Implementa expressões regulares nativas para interceptar a diferença e redigir dados sensíveis (senhas, chaves de API, tokens) antes de iniciar qualquer solicitação HTTP.
* **Resiliência:** Padrão para Groq (Llama 3.3) para respostas em menos de um segundo, com um fallback automático e silencioso para Gemini se o provedor principal experimentar alta carga (HTTP 503).

## Instalação
A Synapse compila diretamente para código de máquina, tornando a instalação direta em qualquer ambiente Linux (incluindo Manjaro e WSL).

### Compilando o Binário
Clone o repositório e compile o binário:
```bash
git clone https://github.com/Elliton-Luis/synapse.git
cd synapse
go build -o syn .
sudo mv syn /usr/local/bin/
```

## Uso
Após preparar as alterações de código, execute a ferramenta CLI em qualquer lugar do repositório:
```bash
git add .
syn
```
Na primeira execução, a Synapse solicitará as chaves de API e gerará automaticamente a configuração de ambiente (.env) e os arquivos de instrução de prompt estrito (commit_pattern.md), garantindo que sejam adicionados com segurança ao .gitignore local.

### Comandos e Flags
#### Configuração de Idioma
Alterna o idioma padrão para os prompts do sistema e as mensagens de commit geradas. Essa configuração persiste no arquivo .env.
```bash
syn lang pt
syn lang en
```
#### Sobrescrita de Provedor
Ignora o roteamento automático e força o uso de um provedor de API específico para a execução atual.
```bash
syn --provider groq
syn --provider gemini
```
#### Modo de Depuração
Exibe logs detalhados sobre o roteamento interno, status de fallback e arquivos carregados.
```bash
syn --debug
```

## Roadmap de Desenvolvimento
- [x] Camada 1 — Protótipo de interface de terminal e loop de interação.
- [x] Camada 2 — Sistema de configuração isolado e configuração de ambiente.
- [x] Camada 3 — Integração RESTful nativa (Gemini + Groq) sem SDKs externos.
- [x] Camada 4 — Otimização de tokens de diferença do Git avançada e sanitização de segurança de regex.
- [x] Camada 5 — Loop de execução principal (Geração, Aprovação, Recreação e execução de commit).
- [x] Camada 6 — Polimento arquitetônico (Flags de CLI, lógica de fallback 503 e subcomandos de idioma).

## Contribuindo
Contribuições que visam melhorar a lógica de extração de diferença, padrões de sanitização de segurança ou desempenho geral são bem-vindas.

### Passos para Contribuir
1. Faça um fork do repositório.
2. Crie uma branch de recurso (git checkout -b feature/NovaFuncionalidade).
3. Faça commit das alterações (usando syn).
4. Faça push para a branch (git push origin feature/NovaFuncionalidade).
5. Abra uma solicitação de pull request.

---

# Synapse CLI
## Introduction
Synapse is a command-line interface (CLI) tool for managing Git commits, built with the Go standard library, focusing on simplicity, zero external dependencies, and token economy to ensure fast and contextually accurate commit messages.

### Architecture and Philosophy
The tool is designed to act as a highly optimized bridge between the local repository and LLM APIs. It does not rely on heavy frameworks or third-party SDKs.

#### Key Architectural Decisions
* **Zero Dependencies:** Compiled into a single static binary. All HTTP routing, OS interactions, and regex sanitizations are handled natively by Go.
* **Token Economy:** Modifies the standard `git diff` output to ignore lock files (`package-lock.json`, `composer.lock`, `go.sum`, etc.) and removes internal Git metadata, ensuring the LLM context window is populated solely by actual code changes.
* **AppSec & Sanitization:** Implements native regular expressions to intercept the diff and redact sensitive data (passwords, API keys, tokens) before initiating any HTTP request.
* **Resiliency:** Defaults to Groq (Llama 3.3) for sub-second responses, with an automatic, silent fallback routing to Gemini if the primary provider experiences high load (HTTP 503).

## Installation
Synapse compiles directly to machine code, making installation straightforward on any Linux environment (including Manjaro and WSL).

### Compiling the Binary
Clone the repository and compile the binary:
```bash
git clone https://github.com/Elliton-Luis/synapse.git
cd synapse
go build -o syn .
sudo mv syn /usr/local/bin/
```

## Usage
After staging your code changes, execute the CLI tool from anywhere in your repository:
```bash
git add .
syn
```
On the first execution, Synapse will prompt you for your API keys and automatically generate its environment configuration (.env) and strict prompt instruction files (commit_pattern.md), ensuring they are safely added to your local .gitignore.

### Commands & Flags
#### Language Configuration
Toggles the default language for the system prompts and the generated commit messages. This setting persists in the .env file.
```bash
syn lang pt
syn lang en
```
#### Provider Override
Bypasses the automatic routing and forces the use of a specific API provider for the current execution.
```bash
syn --provider groq
syn --provider gemini
```
#### Debug Mode
Outputs detailed logs regarding the internal routing, fallback status, and loaded files.
```bash
syn --debug
```

## Development Roadmap
- [x] Layer 1 — Terminal interface prototype and interaction loop.
- [x] Layer 2 — Isolated configuration system and environment setup.
- [x] Layer 3 — Native RESTful integration (Gemini + Groq) without external SDKs.
- [x] Layer 4 — Advanced Git diff token optimization and Regex security sanitization.
- [x] Layer 5 — Main execution loop (Generation, Approval, Recreation, and Commit execution).
- [x] Layer 6 — Architectural polish (CLI Flags, 503 Fallback logic, and Language subcommands).

## Contributing
Contributions focusing on improving the diff extraction logic, security sanitization patterns, or overall performance are welcome.

### Steps to Contribute
1. Fork the repository.
2. Create a feature branch (git checkout -b feature/NewFeature).
3. Commit your changes (using syn).
4. Push to the branch (git push origin feature/NewFeature).
5. Open a Pull Request.