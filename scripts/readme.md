# Massive.com Go Client — generation pipeline

This directory holds the one-command generator for the REST client. The
WebSocket client (`websocket/`) and the REST helpers (`rest/client.go`,
`rest/iterator.go`) are hand-written and are **never** touched by generation.

## One command

```bash
bash scripts/generate.sh
```

That orchestrator does, in order:

1. **Pull the spec** — `rest/scripts/pull_spec.js` downloads
   `https://api.massive.com/openapi`, drops draft paths, forces a single
   `default` tag, applies the `operationId → method-name` renames from
   `rest/scripts/operation-mappings.js`, and writes `rest/scripts/openapi.json`
   (committed, so spec changes show up in PR diffs).
   *Safety gate:* the run aborts if the spec is missing or empty.
2. **Pre-process** — a `jq` pass rewrites invalid `type: number` + `format: int32`
   schemas to `type: integer`.
3. **Generate** — [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
   (config: `rest/scripts/oapi-codegen.yaml`) emits the isolated file
   `rest/gen/client.gen.go`. Because all generated code lives in that one file
   under `rest/gen/`, regeneration only ever replaces it.
   *Safety gate:* the run aborts if no client file was produced.
4. **Post-process** — `rest/scripts/fix-go-clashes.js` renames single-letter
   JSON field clashes (`P/p`, `S/s`, `X/x`, `T/t` → `AskPrice`, `BidPrice`, …)
   and runs `gofmt`.

## Generator version pin

The generator is pinned in `scripts/generate.sh`:

```
OAPI_CODEGEN_VERSION="v2.5.1"
```

and invoked as `go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.1`,
so a run is reproducible regardless of what `oapi-codegen` happens to be on
`PATH`. Diffs then reflect **spec** changes, not generator upgrades. Keep this in
lockstep with the cache comment in `.github/workflows/sync-openapi.yml`.

> There is no `openapitools.json` here — that file configures the Java
> `openapi-generator`, which this Go SDK does not use.

## Daily automation

`.github/workflows/sync-openapi.yml` runs `scripts/generate.sh` on a daily
schedule (and via **Run workflow**). When the regenerated output differs from
what's committed it:

- mints a short-lived **GitHub App** token (the default `GITHUB_TOKEN` can't open
  PRs — org policy blocks it),
- commits as the GitHub App's bot identity,
- pushes a **unique** branch `bot/openapi-sync-<date>-<run-id>` and opens a
  brand-new `[bot]`-prefixed PR (never reusing an existing one, so author ≠
  reviewer),
- posts a Slack notification.

### Required repo/org configuration

| Kind | Name | Purpose |
| --- | --- | --- |
| Variable | `MASSIVE_CLIENT_LIBRARY_AUTOMATION_APP_ID` | GitHub App id (App must be installed on this repo) |
| Secret | `MASSIVE_CLIENT_LIBRARY_AUTOMATION_APP_PRIVATE_KEY` | GitHub App private key |
| Secret | `SLACK_CLIENT_LIBRARY_WEBHOOK` | Slack notification (shared; optional — skipped if unset) |

## Per-file reference

| File | Role |
| --- | --- |
| `scripts/generate.sh` | The single orchestrator (this pipeline). |
| `rest/scripts/pull_spec.js` | Download + filter + rename the spec. |
| `rest/scripts/operation-mappings.js` | **Owns the public Go method names** (`operationId → name`). Language-specific — do not share with other SDKs or "fix" entries, as that renames functions. |
| `rest/scripts/oapi-codegen.yaml` | oapi-codegen config (package `gen`, output `gen/client.gen.go`). |
| `rest/scripts/fix-go-clashes.js` | Post-process single-letter field clashes + gofmt. |
| `rest/scripts/analyze-field-clashes.js` | Standalone diagnostic — **not** part of the pipeline. |
| `rest/scripts/generate-go-examples.js` | Standalone example-snippet generator — **not** part of the pipeline. |
