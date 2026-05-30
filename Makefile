.PHONY: dev build-ui build-server prod run docker-build docker-up docker-down

dev-install:
	cd ui && npm install

# Development: start both Go (with air hot reload) and Vite dev servers;
# kill both on Ctrl+C
dev:
	@bash -c '\
		IP=$$(ipconfig getifaddr en0 2>/dev/null || ipconfig getifaddr en1 2>/dev/null || ifconfig | grep "inet " | grep -v "127.0.0.1" | awk "{print \$$2}" | head -1); \
		if [ -n "$$IP" ]; then \
			export VITE_DEV_HOST=$$IP:8081; \
			echo "LAN: http://$$IP:8080  |  Vite: http://$$IP:8081"; \
		else \
			echo "LAN IP not found, using localhost only"; \
		fi; \
		trap "kill 0" INT; \
		air & cd ui && npm run dev & wait'

# Build UI for production
build-ui:
	cd ui && npm run build

# Build Go server binary
build-server:
	go build -o bin/server ./cmd/server

# Full production build
prod: build-ui build-server

# Run production server
run: prod
	./bin/server

# Docker
# Build Docker image
docker-build:
	docker compose build

# Start Docker container (production)
docker-up:
	docker compose up -d

# Stop Docker container
docker-down:
	docker compose down
