.PHONY: dev build-ui build-server prod run

dev-install:
	cd ui && npm install

# Development: start both Go and Vite dev servers; kill both on Ctrl+C
dev:
	@bash -c 'trap "kill 0" INT; DEV=true go run ./cmd/server & cd ui && npm run dev & wait'

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
