#!/bin/bash
VERSION="dev-$(git describe --tags --always --dirty)"
COMMIT=$(git rev-parse HEAD)
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

echo "Building with version: $VERSION, commit: $COMMIT, date: $DATE"
go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"
go install -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"