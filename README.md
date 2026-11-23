# go-hexagonal-example

A small Go REST API example in `examples/` built with Fiber, GORM, and SQLite. It exposes a todo listing endpoint and organizes handlers, services, and repositories in separate packages.

## Run locally

1. Install Go 1.27.4 or a compatible newer Go release.
2. From `examples/`, copy `env-sample` to `.env` and choose an available `APP_PORT` (the sample uses port 80).
3. From `examples/`, start the server with `go run .`.
4. Open `http://localhost:<APP_PORT>/` to call the todo listing endpoint. Query parameters can filter and paginate results; see `examples/core/handlers.go` and `examples/utils/pagination.go` for the supported fields.

Useful development commands:

```sh
go build ./...
go test ./...
gofmt -w <changed-go-files>
```

## Agent-assisted development: learn by doing

Using coding agents can be a worthwhile hands-on way to learn software development when you stay involved in the work. Ask an agent to explain a change, trace a request from HTTP handler through service and repository, or suggest a small exercise. Then inspect the code, make a change yourself, run it, and explain what happened in your own words. Treat generated code as a draft to understand and review, not as a substitute for understanding.

Go and REST API fundamentals remain essential for building reliable systems. Practice Go types, interfaces, errors, context, concurrency, testing, and package design. For REST APIs, learn HTTP methods and status codes, request validation, pagination, error responses, authentication, and compatibility. This repository is a starting point; build on it with focused exercises such as adding a todo endpoint, validating input, and writing tests for success and failure cases. Strong systems also require deliberate attention to security, observability, data integrity, and operational behavior.

See [AGENTS.md](AGENTS.md) for project-specific guidance for coding agents.
