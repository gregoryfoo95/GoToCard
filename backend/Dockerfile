FROM golang:1.21-alpine AS base

WORKDIR /app

# Install air for hot reload - use compatible version for Go 1.21
RUN go install github.com/cosmtrek/air@v1.40.4

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"] 