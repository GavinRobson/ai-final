#!/bin/bash

# Create output directory if it doesn't exist
mkdir -p builds

APP_NAME="server"  # change this to your binary name

echo "Building binaries for all platforms..."

# Windows
echo " -> Windows (amd64)"
GOOS=windows GOARCH=amd64 go build -o "builds/${APP_NAME}-windows-amd64.exe" main.go

# Linux
echo " -> Linux (amd64)"
GOOS=linux GOARCH=amd64 go build -o "builds/${APP_NAME}-linux-amd64" main.go

# macOS
echo " -> macOS (amd64)"
GOOS=darwin GOARCH=amd64 go build -o "builds/${APP_NAME}-macos-amd64" main.go

echo "Done! Binaries saved to ./builds/"

