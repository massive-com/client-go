const fs = require('fs');
const path = require('path');

const specPath = './openapi.json';
if (!fs.existsSync(specPath)) {
  console.error(`Error: ${specPath} not found. Run pull_spec.js first?`);
  process.exit(1);
}

const spec = JSON.parse(fs.readFileSync(specPath, 'utf8'));

console.log('🔍 Scanning spec for REAL single-letter field clashes...\n');

const clashes = [];

function walkSchema(schema, context = 'root') {
  if (!schema || typeof schema !== 'object') return;

  if (schema.type === 'object' && schema.properties) {
    const fieldsByGoName = new Map(); // key: uppercase Go name → list of {jsonKey, desc}

    Object.entries(schema.properties).forEach(([jsonKey, prop]) => {
      if (jsonKey.length !== 1) return; // only single-letter

      const goName = jsonKey.toUpperCase();
      const description = prop.description || '(no description)';

      if (!fieldsByGoName.has(goName)) {
        fieldsByGoName.set(goName, []);
      }
      fieldsByGoName.get(goName).push({ jsonKey, description });
    });

    // Only report when there is an actual clash (multiple entries for same goName)
    for (const [goName, entries] of fieldsByGoName) {
      if (entries.length > 1) {
        clashes.push({ context, goName, entries });
      }
    }
  }

  // Recurse
  if (schema.allOf) schema.allOf.forEach((s, i) => walkSchema(s, `${context} → allOf[${i}]`));
  if (schema.items) walkSchema(schema.items, `${context} → items`);
  if (schema.properties) Object.entries(schema.properties).forEach(([k, p]) => walkSchema(p, `${context} → prop "${k}"`));
}

// Walk response schemas
Object.entries(spec.paths || {}).forEach(([route, methods]) => {
  Object.entries(methods).forEach(([method, op]) => {
    if (!['get','post','put','delete','patch','head','options','trace'].includes(method)) return;
    const opId = op.operationId || `${method.toUpperCase()} ${route}`;
    Object.entries(op.responses || {}).forEach(([code, resp]) => {
      const jsonSchema = resp.content?.['application/json']?.schema;
      if (jsonSchema) {
        walkSchema(jsonSchema, `${opId} → ${code} response`);
      }
    });
  });
});

// Walk components.schemas (common models)
if (spec.components?.schemas) {
  Object.entries(spec.components.schemas).forEach(([name, schema]) => {
    walkSchema(schema, `components.schemas.${name}`);
  });
}

if (clashes.length === 0) {
  console.log('🎉 No real clashes found! (No duplicate single-letter Go names in any struct)');
} else {
  console.log(`Found ${clashes.length} REAL CLASHES:\n`);

  clashes.forEach((item, idx) => {
    console.log(`#${idx + 1}  🔴 ${item.context}`);
    console.log(`   Go field "${item.goName}" is used by multiple JSON keys:`);
    item.entries.forEach(e => {
      console.log(`     • json:"${e.jsonKey}"  →  "${e.description}"`);
    });
    console.log('');
  });

  console.log('Review these carefully — especially where the same letter means different things.');
  console.log('We can now build a safe rename map based on this output.');
}
