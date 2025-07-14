# DockPilot

DockPilot é um orquestrador leve de microsserviços locais com deploy automatizado, health checks e painel de controle.

## Propósito

DockPilot simplifica o desenvolvimento e a gestão local de microsserviços, oferecendo:

- Deploy automatizado de containers baseado em configuração YAML
- Monitoramento de health check para cada serviço
- CLI e (futuramente) painel de controle para fácil gerenciamento

## Arquitetura

```
DockPilot CLI → Leitura do YAML → Container runner + Health check loop
```

## Primeiros Passos

1. **Clone o repositório**
2. **Edite o `config.yaml`** para definir seus serviços:

```yaml
services:
  - name: users-api
    image: users-api:latest
    port: 8080
    healthcheck: /health
```

3. **Execute o orquestrador:**

```sh
cd cmd/orchestrator
 go run main.go
```

## Estrutura do Projeto

- `cmd/orchestrator` — Ponto de entrada principal (CLI)
- `pkg/health` — Lógica de health check (a implementar)
- `pkg/services` — Gerenciamento dos serviços (a implementar)
- `internal/config` — Parser e gestão de configuração
- `docs/` — Documentação e diagramas

## Roadmap

- Integração com container runner
- Loop de health check
- Painel de controle (UI)

---

MIT License
