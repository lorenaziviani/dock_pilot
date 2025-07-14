# DockPilot

DockPilot é um orquestrador leve de microsserviços locais com deploy automatizado, health checks e painel de controle.

## Propósito

DockPilot simplifica o desenvolvimento e a gestão local de microsserviços, oferecendo:

- Deploy automatizado de containers baseado em configuração YAML
- Monitoramento de health check para cada serviço
- CLI e (futuramente) painel de controle para fácil gerenciamento
- Gerenciamento completo do ciclo de vida dos containers Docker

## Arquitetura

```
DockPilot CLI (start/stop/restart/status) → Leitura do YAML → Container runner (volumes, redes, portas) + Health check loop
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
    volumes:
      - ./data:/app/data
    networks:
      - dockpilot-net
    ports:
      - 8080:8080
```

3. **Configure as variáveis de ambiente** (exemplo):

```
DOCKPILOT_ENV=development
DOCKER_HOST=unix:///var/run/docker.sock
DOCKPILOT_NETWORK=dockpilot-net
DOCKPILOT_DATA_PATH=./data
```

4. **Execute o orquestrador:**

```sh
cd cmd/orchestrator
 go run main.go <comando> [serviço|all]
```

### Comandos disponíveis

- `start <serviço|all>` — Inicia um ou todos os containers definidos no YAML
- `stop <serviço|all>` — Para um ou todos os containers
- `restart <serviço|all>` — Reinicia um ou todos os containers
- `status <serviço|all>` — Mostra o status de um ou todos os containers

Exemplo:

```sh
go run main.go start all
go run main.go status users-api
```

## Estrutura do Projeto

- `cmd/orchestrator` — Ponto de entrada principal (CLI)
- `pkg/health` — Lógica de health check (a implementar)
- `pkg/services` — Gerenciamento dos serviços e integração Docker
- `internal/config` — Parser e gestão de configuração
- `docs/` — Documentação e diagramas

## Roadmap

- Integração com container runner (completa)
- Loop de health check
- Painel de controle (UI)
- Suporte avançado a volumes, redes e portas

---

MIT License
