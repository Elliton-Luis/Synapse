# aicommit

Gerador de mensagens de commit via IA (Groq e OpenRouter).

## Estrutura

```
synapse/
├── main.go
├── go.mod
└── internal/
    └── ui/
        └── ui.go       # helpers de terminal
```

## Roadmap

- [x] Camada 1 — Protótipo da interface (mock estático)
- [ ] Camada 2 — Leitura de config (~/.aicommit.toml)
- [ ] Camada 3 — Interface de IA + Groq + OpenRouter
- [ ] Camada 4 — Git diff + sanitização
- [ ] Camada 5 — Loop principal real
- [ ] Camada 6 — Polimento (--debug, flags, fallback)

## Como rodar

```bash
go run .
```
