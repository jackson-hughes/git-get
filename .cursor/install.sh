#!/usr/bin/env bash
# Idempotent Cloud Agent setup for the git-get Go CLI.
# Prepares module dependencies and the developer tooling used by the
# Taskfile and CI (task, golangci-lint, govulncheck, goreleaser), then
# builds the project to warm the cache and prove the toolchain works.
set -euo pipefail

# The repository targets Go 1.25 (see go.mod). The base image ships an older
# 1.25 patch, so install the pinned patch release into /usr/local/go and put it
# ahead of any system Go on PATH. This keeps `govulncheck` clean, matching CI's
# use of the latest stable Go.
GO_VERSION="go1.25.12"
if ! /usr/local/go/bin/go version 2>/dev/null | grep -q "$GO_VERSION"; then
  echo "==> Installing ${GO_VERSION} into /usr/local/go"
  tarball="${GO_VERSION}.linux-amd64.tar.gz"
  curl -sSfL -o "/tmp/${tarball}" "https://go.dev/dl/${tarball}"
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf "/tmp/${tarball}"
  rm -f "/tmp/${tarball}"
fi
sudo ln -sf /usr/local/go/bin/go /usr/local/bin/go
sudo ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt
hash -r

GOBIN_DIR="$(go env GOPATH)/bin"
mkdir -p "$GOBIN_DIR"

echo "==> go version: $(go version)"

echo "==> Downloading module dependencies"
go mod download

echo "==> Installing developer tooling"
go install github.com/go-task/task/v3/cmd/task@latest
go install golang.org/x/vuln/cmd/govulncheck@v1.5.0
go install github.com/goreleaser/goreleaser/v2@latest
# golangci-lint is pinned to the version used in CI (.github/workflows/go.yaml).
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2

# Expose the Go-installed tools on the system PATH so every shell (login or
# not) can find them. /usr/local/bin is already on PATH in this image.
echo "==> Linking tools onto the system PATH"
sudo install -m 0755 \
  "$GOBIN_DIR/task" \
  "$GOBIN_DIR/govulncheck" \
  "$GOBIN_DIR/goreleaser" \
  "$GOBIN_DIR/golangci-lint" \
  /usr/local/bin/

echo "==> Building the project"
go build -v ./...

echo "==> Setup complete"
