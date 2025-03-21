# Makefile for Solana Testnet Faucet

.PHONY: local dev prod down clean help

help:
	@echo "Solana Testnet Faucet Commands:"
	@echo "  make local  - Run local development setup with direct port mapping"
	@echo "  make dev    - Run development setup with Nginx proxy"
	@echo "  make prod   - Run production setup with Cloudflare Tunnel"
	@echo "  make down   - Stop all containers"
	@echo "  make clean  - Remove all containers, networks, and volumes"
	@echo "  make help   - Show this help message"

# Local development setup
local:
	docker compose -f docker-compose.local.yml up -d --build
	@echo "Local development setup running at:"
	@echo "  - Frontend: http://localhost:80"
	@echo "  - Backend API: http://localhost:8080"

# Development setup with Nginx
dev:
	docker compose up -d
	@echo "Development setup with Nginx running at:"
	@echo "  - Frontend and API: http://localhost:80"

# Production setup with Cloudflare Tunnel
prod:
	docker compose -f docker-compose.prod.yml up -d
	@echo "Production setup with Cloudflare Tunnel running"
	@echo "Access via your configured Cloudflare domains"

# Stop all containers
down:
	docker compose -f docker-compose.local.yml down 2>/dev/null || true
	docker compose down 2>/dev/null || true
	docker compose -f docker-compose.prod.yml down 2>/dev/null || true
	@echo "All containers stopped"

# Remove all containers, networks, and volumes
clean:
	docker compose -f docker-compose.local.yml down -v 2>/dev/null || true
	docker compose down -v 2>/dev/null || true
	docker compose -f docker-compose.prod.yml down -v 2>/dev/null || true
	@echo "All containers, networks, and volumes removed" 