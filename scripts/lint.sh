#!/usr/bin/env bash

echo "Run linter..."
golangci-lint run korecli/... pkg/... services/...
