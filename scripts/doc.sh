#!/usr/bin/env bash

# generate swagger documentation
echo "Generate swagger documentation..."
swag fmt
swag init -g cmd/main.go -o cmd/swagger/docs
