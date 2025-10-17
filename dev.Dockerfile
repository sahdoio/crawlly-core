# Development Dockerfile with hot reload support
FROM golang:1.25.3-alpine

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Install Air for hot reload (new repository location)
RUN go install github.com/air-verse/air@latest

# Add Go bin to PATH (where Air was installed)
ENV PATH="/go/bin:${PATH}"

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 3000

# Run Air for hot reload
CMD ["air", "-c", ".air.toml"]
