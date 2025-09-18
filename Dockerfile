# Start with an existing image as the base image
# This particular image contains the Go toolchain
FROM golang:1.25.1-alpine AS builder

# Set the container's working directory
WORKDIR /app

# Copy files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the application
RUN go build -o /app/bin/ledger ./cmd/main.go

FROM debian:stable-slim

# Copies the build results into the image's /bin/ledger folder
COPY --from=builder /app/bin/ledger /bin/ledger

CMD ["/bin/ledger"]