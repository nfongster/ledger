#!/bin/bash

# If a .env file is present, it will be used to load environment variables.
# For CI/CD pipelines, you should instead configure environment variables in the appropriate yml file.
if [ -f .env ]; then
    echo "Importing environment variables from .env file..."
    set -o allexport; source .env; set +o allexport
fi

echo "DB_HOST set to $DB_HOST"
echo "POSTGRES_USER set to $POSTGRES_USER"
echo "POSTGRES_PASSWORD set to $POSTGRES_PASSWORD"
echo "POSTGRES_DB set to $POSTGRES_DB"
echo "POSTGRES_PORT set to $POSTGRES_PORT"

case "$1" in
    "run")
        go run ./cmd/main.go
        ;;
    "test")
        go test ./...
        ;;
    *)
        echo "Invalid argument: $1"
        echo "Usage: ./ledger.sh [run|test]"
        exit 1
        ;;
esac