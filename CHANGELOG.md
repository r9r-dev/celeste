# Changelog

Toutes les modifications notables de ce projet seront documentees dans ce fichier.

Le format est base sur [Keep a Changelog](https://keepachangelog.com/fr/1.0.0/),
et ce projet adhère au [Versioning Semantique](https://semver.org/lang/fr/).

## [0.1.0] - 2026-01-06

### Nouveautes

- Interface web d'administration Docker avec theme sci-fi "Mission Control"
- Backend Go/Gin avec API REST complete
- Frontend SvelteKit 5 avec design retro-futuriste
- Gestion des stacks Docker Compose (start, stop, restart, pull)
- Gestion des conteneurs (start, stop, restart, logs)
- Visualisation des volumes, reseaux et images
- Statistiques systeme en temps reel via WebSocket (CPU, memoire, disque)
- Statistiques des conteneurs en temps reel via WebSocket
- Editeur de fichiers docker-compose.yml integre
- Support multi-plateforme (linux/amd64, linux/arm64)

### Technique

- Go 1.25.5 avec Docker SDK v27.5.1
- SvelteKit 5 avec Svelte 5 runes ($state, $derived, $effect)
- WebSocket pour les mises a jour temps reel
- Docker Compose CLI wrapper pour les operations sur les stacks
- gopsutil pour les statistiques systeme hote
