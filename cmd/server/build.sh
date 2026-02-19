#!/bin/bash

VERSION=$(git describe --tags 2>/dev/null || echo "develop")
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

go build -ldflags "-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$BUILD_TIME" -o ./build/go-cdn ./cmd/server/main.go
