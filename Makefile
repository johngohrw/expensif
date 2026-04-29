.PHONY: dev build-ui build-server prod run

# Development: start both servers (run in separate terminals)
dev:
	@echo "Start Go server:   DEV=true go run ./cmd/server"
	@echo "Start Vite server: cd ui && npm run dev"

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
