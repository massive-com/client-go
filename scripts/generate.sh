#!/usr/bin/env bash
#
# generate.sh — one-command regeneration of the Massive.com Go REST client.
#
# Reproduces the committed repo layout deterministically:
#   1. Pull the OpenAPI spec         -> rest/scripts/openapi.json
#   2. Pre-process the spec (jq)      (fix invalid numeric schemas)
#   3. Generate the Go REST client with oapi-codegen into the isolated
#      package  rest/gen/client.gen.go, REPLACING only that generated file.
#   4. Post-process (fix single-letter JSON field clashes + gofmt).
#
# oapi-codegen only understands the REST endpoints. Hand-written code
# (rest/client.go, rest/iterator.go and the entire websocket/ package) and
# curated files (README.md, go.mod, LICENSE) live OUTSIDE rest/gen/ and are
# never touched by this script.
#
# The generator version is PINNED below (OAPI_CODEGEN_VERSION) and invoked via
# `go run <module>@<version>` so a run is reproducible regardless of what
# `oapi-codegen` happens to be on PATH — diffs then reflect spec changes, not
# generator upgrades. (There is no openapitools.json here: that file configures
# the Java openapi-generator, which this Go SDK does not use.)
#
# Usage (from anywhere):
#   bash scripts/generate.sh
#
# Requirements: bash, Go (1.21+), Node.js (18+), jq.
#
set -euo pipefail

# Pinned generator version — keep in lockstep with the cache comment in
# .github/workflows/sync-openapi.yml.
OAPI_CODEGEN_VERSION="v2.5.1"
OAPI_CODEGEN_PKG="github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT"

REST_DIR="$ROOT/rest"
SCRIPTS_DIR="$REST_DIR/scripts"
SPEC_FILE="$SCRIPTS_DIR/openapi.json"
GEN_FILE="$REST_DIR/gen/client.gen.go"
GEN_CONFIG="./scripts/oapi-codegen.yaml"   # relative to REST_DIR

echo "==> [1/4] Pulling OpenAPI spec -> ${SPEC_FILE#$ROOT/}"
# pull_spec.js writes ./openapi.json relative to its cwd, so run it from there.
( cd "$SCRIPTS_DIR" && node pull_spec.js )

# Safety gate: never regenerate against a missing/empty spec (pull_spec.js only
# logs fetch failures, it does not exit non-zero — so a failed pull would
# otherwise leave a stale or empty spec in place).
if [ ! -s "$SPEC_FILE" ]; then
  echo "ERROR: $SPEC_FILE is missing or empty after pull_spec.js; aborting." >&2
  exit 1
fi

echo "==> [2/4] Pre-processing spec (fixing invalid 'number'+'int32' schemas)"
FIXED_SPEC="$(mktemp -t openapi-fixed.XXXXXX.json)"
trap 'rm -f "$FIXED_SPEC"' EXIT
jq '
  walk(
    if type == "object" and .type == "number" and .format == "int32" then
      .type = "integer"
    else
      .
    end
  )
' "$SPEC_FILE" > "$FIXED_SPEC"

echo "==> [3/4] Generating Go client with oapi-codegen ${OAPI_CODEGEN_VERSION}"
rm -rf "$REST_DIR/gen"
# oapi-codegen resolves the config's `output:` relative to its cwd, so run it
# from REST_DIR to land the file at rest/gen/client.gen.go.
( cd "$REST_DIR" && go run "${OAPI_CODEGEN_PKG}@${OAPI_CODEGEN_VERSION}" \
    -config "$GEN_CONFIG" \
    "$FIXED_SPEC" )

# Safety gate: never leave a partial/empty generated client committed.
if [ ! -s "$GEN_FILE" ]; then
  echo "ERROR: generation did not produce $GEN_FILE; aborting." >&2
  exit 1
fi

echo "==> [4/4] Post-processing (fix JSON field clashes + gofmt)"
node "$SCRIPTS_DIR/fix-go-clashes.js"

echo "Done. Regenerated ${GEN_FILE#$ROOT/} from ${SPEC_FILE#$ROOT/}"
echo "  (hand-written rest/client.go, rest/iterator.go and websocket/ left untouched)"
