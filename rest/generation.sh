#!/bin/bash
set -euo pipefail

echo "=== Pulling latest spec ==="

cd scripts
node pull_spec.js

# back to /rest
cd ..

echo "=== Regenerating Go client with oapi-codegen ==="

# Pre-process: fix all invalid "type: number" + "format: int32" → "type: integer"
echo "Pre-processing openapi.json (fixing numeric schema bugs)..."
jq '
  walk(
    if type == "object" and .type == "number" and .format == "int32" then
      .type = "integer"
    else
      .
    end
  )
' ./scripts/openapi.json > /tmp/openapi-fixed.json

rm -rf ./gen

oapi-codegen \
  -config ./scripts/oapi-codegen.yaml \
  /tmp/openapi-fixed.json 2>&1 | tee generation.log

rm -f /tmp/openapi-fixed.json

# Post-process: fixing P/p, S/s, T/t, X/x clashes in anonymous structs...
node scripts/fix-go-clashes.js

echo "Generation finished!"
echo "   → Full log saved to generation.log"
echo "   → Generated files are in ./gen/"
