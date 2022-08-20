#!/usr/bin/env bash

set -eu

echo "=== Unit tests==="
go test ./...

echo "=== Smoke tests ==="
./tests/converter_test.sh
