#!/usr/bin/env bash

SERVICE=$1

# run wire dependency injection
echo "Run wire dependency injection..."
wire ./cmd
