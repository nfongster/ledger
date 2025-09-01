#!/bin/bash

case "$1" in
    "server")
        go run ./cmd/server/main.go
        ;;
    "client")
        go run ./cmd/client/main.go
        ;;
    *)
        echo "Invalid option.  Please use 'server' or 'client'."
        exit 1
        ;;
esac