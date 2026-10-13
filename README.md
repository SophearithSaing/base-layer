# Base Layer

A Go API server for authentication, users, projects, groups, and AI receipt extraction.

## Requirements

- Go 1.26.5 or later
- MongoDB
- A Together API key

## Repository structure

```text
.
├── api/                 # OpenAPI specifications (JSON and YAML)
├── cmd/server/          # Server entry point
├── internal/            # Application packages and their Go tests
│   ├── ai/
│   ├── api/
│   ├── auth/
│   ├── config/
│   ├── db/
│   ├── group/
│   ├── project/
│   └── user/
├── .env.example         # Example runtime configuration
├── Dockerfile
├── go.mod
└── go.sum
```

Run all commands from the repository root.

## Run locally

Copy the example configuration and set your MongoDB URI, JWT secret (at least
32 characters), and Together API key:

```sh
cp .env.example .env
```

The server reads configuration from environment variables. Load the file into
your shell before starting it:

```sh
set -a
. ./.env
set +a
go run ./cmd/server
```

The default port is `8000`. `CLIENT_ORIGINS` is a comma-separated list of allowed
browser origins. See `.env.example` for all configuration variables.

## Build and check

```sh
mkdir -p bin
go build -o bin/server ./cmd/server
go test ./...
go vet ./...
```

Go tests live alongside the packages they cover in `internal/`. GitHub Actions
runs the tests, vet, and build from the root module on pushes and pull requests.

## Docker

```sh
docker build -t base-layer .
docker run --rm --env-file .env -p 8000:8000 base-layer
```

Set `MONGODB_URI` to a MongoDB address reachable from the container; `localhost`
inside the container refers to the container itself. If you change `PORT`, update
the published container port to match.
