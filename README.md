# Simon

Monorepo do sistema **Simon**: API em Go, aplicação web em React com Vite e, no futuro, cliente desktop com Tauri.

## Estrutura

- `apps/api` — serviço backend (Go)
- `apps/web` — interface web (React + Vite)
- `apps/desktop` — cliente desktop (Tauri), ainda não iniciado
- `packages/shared` — código compartilhado entre apps
- `docker/postgres` — artefatos opcionais do PostgreSQL (por exemplo, scripts em `init/`)

## Infraestrutura local

1. Copie as variáveis de ambiente:

   ```bash
   cp .env.example .env
   ```

2. Ajuste `.env` conforme necessário (especialmente `DB_PASSWORD`).

3. Suba o PostgreSQL:

   ```bash
   docker compose up -d
   ```

A API será adicionada ao Compose em uma etapa posterior.

## Requisitos

- Docker e Docker Compose
- Go e Node.js serão usados quando os apps forem implementados
