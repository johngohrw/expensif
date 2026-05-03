# Stage 1: Build the UI
FROM node:22-alpine AS ui-builder
WORKDIR /app
COPY ui/package*.json ./ui/
RUN cd ui && npm ci
COPY ui/ ./ui/
RUN cd ui && npm run build

# Stage 2: Build the Go binary
FROM golang:1.25-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/server ./cmd/server

# Stage 3: Runtime
FROM alpine:latest
WORKDIR /app

# Copy built artifacts
COPY --from=ui-builder /app/static ./static
COPY --from=go-builder /app/bin/server ./bin/server
COPY templates/ ./templates/

# Create data directory for SQLite
RUN mkdir -p /data
ENV DATA_DIR=/data

EXPOSE 8080

CMD ["./bin/server"]
