# Synapse

Gerador de mensagens de commit de alta precisão via IA (Groq e Gemini). Focado em simplicidade, zero dependências externas e economia extrema de tokens.

## Estrutura

```text
synapse/
├── main.go
├── go.mod
├── .env                    # Gerado automaticamente (Chaves de API e configs)
├── commit_pattern.md       # Gerado automaticamente (Regras em Inglês)
├── commit_pattern_pt.md    # Gerado automaticamente (Regras em Português)
└── internal/
    ├── ui/
    │   └── ui.go           # Helpers de interface do terminal
    └── config/
        └── config.go       # Gerenciamento de variáveis e setup inicial

```

## Roadmap

- [x] Camada 1 — Protótipo da interface (mock estático)
- [X] Camada 2 — Criaçao dos arquivos de config
- [X] Camada 3 — Interface de Gemini + Groq
- [X] Camada 4 — Git diff + sanitização
- [ ] Camada 5 — Loop principal real
- [ ] Camada 6 — Polimento (--debug, flags, fallback)

## Como rodar

```bash

go run .

```

ou
```bash

go build -o syn .

sudo mv syn /usr/local/bin/

```