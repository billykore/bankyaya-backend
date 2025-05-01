#!/usr/bin/env bash

echo "Run linter..."
golangci-lint run --config .golangci.yaml
