# DockPilot

DockPilot é um orquestrador leve de microsserviços locais com deploy automatizado, health checks, painel de controle e dashboard interativo no terminal.

## Propósito

DockPilot simplifica o desenvolvimento e a gestão local de microsserviços, oferecendo:

- Deploy automatizado de containers baseado em configuração YAML
- Monitoramento de health check para cada serviço
- CLI, painel de controle e dashboard TUI para fácil gerenciamento
- Gerenciamento completo do ciclo de vida dos containers Docker
- Monitoramento contínuo e autocorreção dos serviços
- **Dashboard visual e interativo no terminal**

## Arquitetura

```
DockPilot CLI (start/stop/restart/status/monitor/dashboard) → Leitura do YAML → Container runner (volumes, redes, portas) + Health check loop + Autocorreção + Dashboard TUI
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
- `monitor` — Inicia o monitoramento contínuo dos serviços, com autocorreção e logs detalhados
- `dashboard` — Abre o dashboard visual no terminal (TUI)

Exemplo:

```sh
go run main.go dashboard
go run main.go monitor
go run main.go start all
go run main.go status users-api
```

## Dashboard TUI (Terminal User Interface)

- Visualize todos os serviços, status, porta, health e uptime (placeholder)
- Atalhos interativos:
  - `s` = start serviço selecionado
  - `r` = restart serviço selecionado
  - `l` = logs (placeholder)
  - `q` = sair do dashboard
- Atualização automática a cada 5 segundos
- Experiência visual produtiva para dev local

## Sistema de Health Check e Autocorreção

- O comando `monitor` verifica periodicamente o endpoint `/health` de cada serviço.
- Classificação automática:
  - **healthy**: serviço responde 200 OK
  - **degraded**: responde, mas não 200 OK
  - **unreachable**: não responde
- Serviços `unreachable` são reiniciados automaticamente.
- Logs detalhados por serviço são exibidos na CLI em tempo real.

## Estrutura do Projeto

- `cmd/orchestrator` — Ponto de entrada principal (CLI e dashboard)
- `pkg/health` — Monitoramento, health check e autocorreção
- `pkg/services` — Gerenciamento dos serviços e integração Docker
- `internal/config` — Parser e gestão de configuração
- `docs/` — Documentação e diagramas

## Roadmap

- Integração com container runner (completa)
- Loop de health check e autocorreção (completo)
- Dashboard TUI (completo)
- Painel de controle (UI web, opcional)
- Suporte avançado a volumes, redes e portas

---

MIT License
