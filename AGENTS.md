# Agent development guide

## Project

This repository is a Go 1.27.4 application using Fiber, Viper, GORM, and SQLite. Keep the existing hexagonal package layout and module dependencies unless a task explicitly asks to change them.

## Commands

- Start the server: `go run .` (requires a local `.env`; use `env-sample` as the template).
- Check compilation: `go build ./...`
- Run tests: `go test ./...`
- Format changed Go files: `gofmt -w <files>`

## Working agreements

- Keep local pagination helpers in `pkg/pagination`.
- Keep domain and application logic in `core`; keep transport and persistence adapters in their existing packages.
- Do not commit `.env` files or generated SQLite databases.
- Do not add dependencies or change the Go version without a task requirement.
- Update this guide when project commands or structure change.
