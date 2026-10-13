# AGENTS.md

Guidelines for AI coding agents working on the Massive Go client (`github.com/massive-com/client-go/v3`).

## Generated code: read this first

The REST client is generated from the OpenAPI spec at https://api.massive.com/openapi. A bot regenerates it every day and opens a PR, so **any hand edit to generated files will be overwritten.**

| Path | Status | Notes |
| --- | --- | --- |
| `rest/gen/client.gen.go` | Generated | REST methods and models, produced by oapi-codegen. Never edit. |
| `rest/scripts/openapi.json` | Generated | The filtered spec the client is built from. Never edit. |
| `rest/client.go`, `rest/iterator.go` | Hand-written | Client constructor, options, pagination iterator. |
| `websocket/` | Hand-written | The whole WebSocket client, including `websocket/models`. |
| `scripts/generate.sh`, `rest/scripts/*.js`, `rest/scripts/oapi-codegen.yaml` | Tooling | The generation pipeline. Edit these to change generated output. |
| `.massive/` | Legacy | Old spec cache. Not used by the current pipeline. |

To change what ends up in `rest/gen/`, edit the right input instead of the output:

- **Endpoint, parameter, schema, or doc changes**: these come from the upstream spec. They can't be fixed in this repo; report them to the Massive team.
- **Public Go method names**: `rest/scripts/operation-mappings.js`. Renaming an entry renames a public function, which is a breaking change. Don't "fix" or reorder entries, and don't copy this file from another SDK.
- **Spec filtering or reshaping before generation**: `rest/scripts/pull_spec.js`, or the `jq` step in `scripts/generate.sh`.
- **Generator options**: `rest/scripts/oapi-codegen.yaml`.
- **Generator version**: `OAPI_CODEGEN_VERSION` in `scripts/generate.sh`. Keep the matching comment in `.github/workflows/sync-openapi.yml` in sync.
- **Fixups after generation** (field-name clashes such as `P`/`p`): `rest/scripts/fix-go-clashes.js`.
- **Client ergonomics** (options, pagination, helpers): `rest/client.go` and `rest/iterator.go`.

After changing any tooling, regenerate and commit both the tooling change and the regenerated output.

See [`scripts/readme.md`](./scripts/readme.md) for the full pipeline.

## Essential commands

```bash
bash scripts/generate.sh      # Regenerate rest/gen/ from the live spec (needs Go 1.21+, Node 18+, jq)
go build ./...                # Build
go test -race -v ./...        # Run tests (what CI runs)
go fmt ./...                  # Format (CI fails on a diff)
go mod tidy                   # CI fails if go.mod/go.sum change
golangci-lint run             # Lint (CI uses v1.64.8; rest/gen/ is excluded)
```

There's no Makefile. CI tests on Go 1.23 and 1.24 across Linux, macOS, and Windows.

## Layout

- `rest/`: REST client. Hand-written wrapper (`client.go`, `iterator.go`) around the generated `gen` package.
- `rest/gen/`: generated client (package `gen`).
- `rest/scripts/`: spec, generator config, and JS tooling.
- `scripts/`: `generate.sh` and pipeline docs.
- `websocket/`: hand-written streaming client (package `massivews`), with `models/` and `example/`.

## Conventions

- API keys come from the `MASSIVE_API_KEY` env var. `rest.New("")` falls back to it. Never hardcode a key in code, tests, or examples.
- The module path carries the major version (`/v3`). A breaking change to the public API means a new major version, so avoid one.
- Releases are git tags (`v3.x.y`). There is no changelog file.
- Tests must not call the live API.

## Pull requests

- Open or link an issue first. PRs not tied to an issue are generally declined.
- Don't hand-edit or commit unrelated changes to `rest/gen/` or `rest/scripts/openapi.json`. Spec-driven updates arrive through the daily `[bot] Sync client-go with OpenAPI spec` PRs.
- Run `go build ./...`, `go test ./...`, and `go fmt ./...` before opening a PR.
