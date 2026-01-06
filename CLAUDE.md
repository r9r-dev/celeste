# Celeste - Aperture Science Network

Interface d'administration Docker avec esthetique sci-fi "Mission Control".

## Stack technique

- **Frontend**: SvelteKit 5 + Svelte 5 runes + Tailwind CSS v4 + TypeScript
- **Backend**: Go 1.24 + Gin + gorilla/websocket + Docker SDK v27.5.1
- **Stats systeme**: gopsutil (avec HOST_PROC/HOST_SYS pour deploiement containerise)

## Structure du projet

```
celeste/
├── frontend/                 # SvelteKit application
│   ├── src/
│   │   ├── lib/
│   │   │   ├── api/            # Client REST et WebSocket
│   │   │   ├── components/     # Composants UI
│   │   │   └── stores/         # Stores Svelte 5 reactifs
│   │   └── routes/             # Pages SvelteKit
│   └── static/
├── backend/                  # Go API server
│   ├── cmd/server/           # Point d'entree
│   └── internal/
│       ├── api/              # Routes et handlers Gin
│       ├── compose/          # Wrapper CLI docker-compose
│       ├── docker/           # Client Docker SDK
│       ├── stats/            # Stats systeme (gopsutil)
│       └── ws/               # WebSocket hub
├── docker-compose/           # Configuration deploiement
│   ├── docker-compose.yml
│   └── .env.example
├── Dockerfile                # Build multi-stage
└── CHANGELOG.md
```

## Commandes

### Backend

```bash
cd backend
go run ./cmd/server              # Dev server
go build -o celeste ./cmd/server # Build
go test ./...                    # Tests
go get -u ./... && go mod tidy   # Update deps
```

### Frontend

```bash
cd frontend
npm install      # Install deps
npm run dev      # Dev server (:5173)
npm run build    # Build production
npm run check    # Type check
```

### Docker

```bash
docker build -t celeste .        # Build image
cd docker-compose && docker compose up -d  # Deploy
```

## Ports

- Frontend dev: 5173
- Backend: 8080

## Variables d'environnement

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | Port HTTP |
| STACKS_PATH | /home/share/docker/dockge/stacks | Repertoire des stacks |
| HOST_PROC | /proc | Mount /proc hote |
| HOST_SYS | /sys | Mount /sys hote |
| GIN_MODE | debug | Mode Gin |

## API Endpoints

### System
- `GET /health` - Health check
- `GET /api/stats` - Stats systeme

### Stacks
- `GET /api/stacks` - Liste des stacks
- `GET /api/stacks/:name` - Details d'une stack
- `POST /api/stacks/:name/start` - Demarrer
- `POST /api/stacks/:name/stop` - Arreter
- `POST /api/stacks/:name/restart` - Redemarrer
- `POST /api/stacks/:name/pull` - Pull images
- `GET /api/stacks/:name/compose` - Contenu docker-compose.yml
- `PUT /api/stacks/:name/compose` - Modifier docker-compose.yml

### Containers
- `GET /api/containers` - Liste
- `GET /api/containers/:id` - Details
- `POST /api/containers/:id/start` - Demarrer
- `POST /api/containers/:id/stop` - Arreter
- `POST /api/containers/:id/restart` - Redemarrer
- `GET /api/containers/:id/logs` - Logs
- `GET /api/containers/:id/stats` - Stats

### Volumes
- `GET /api/volumes` - Liste
- `POST /api/volumes` - Creer
- `DELETE /api/volumes/:name` - Supprimer

### Networks
- `GET /api/networks` - Liste
- `POST /api/networks` - Creer
- `DELETE /api/networks/:id` - Supprimer

### Images
- `GET /api/images` - Liste

### WebSocket
- `WS /ws` - Stats temps reel
  - `stats` - Stats systeme (CPU, memoire, disque)
  - `container_stats` - Stats par conteneur

## Deploiement

L'application necessite:
- Docker socket (`/var/run/docker.sock`) pour la gestion des conteneurs
- `/proc` et `/sys` de l'hote pour les statistiques systeme
- Repertoire des stacks pour les operations docker-compose

Voir `docker-compose/docker-compose.yml` pour la configuration recommandee.

## Design

Esthetique "Mission Control" sci-fi:
- Fond noir (#050508)
- Accent cyan (#00d4aa)
- Police monospace (JetBrains Mono)
- Bordures fines, effets glow
- Graphiques style oscilloscope
