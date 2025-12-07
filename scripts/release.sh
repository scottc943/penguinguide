#!/usr/bin/env bash
set -euo pipefail

echo "Running tests"
go test ./...

echo "Generating man pages"
go run . man --dir=man

echo "Running GoReleaser"
if [[ "${1-}" == "snapshot" ]]; then
  goreleaser release --snapshot --clean
else
  goreleaser release --clean
fi

