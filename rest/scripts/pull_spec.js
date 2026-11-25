import path from 'path';
import fs from 'fs';
import mappings from './operation-mappings.js'; // Import the operationId mappings

const fetchSpec = async (url) => {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`Failed to fetch spec: ${response.statusText}`);
  }
  return response.json();
};

const saveSpec = (spec, filename) => {
  const dir = path.dirname(filename);
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true });
  }
  fs.writeFileSync(filename, JSON.stringify(spec, null, 2));
};

const pascal = (s) =>
  String(s)
    .replace(/(^|[^a-zA-Z0-9])([a-zA-Z0-9])/g, (_m, _sep, c) => c.toUpperCase())
    .replace(/[^a-zA-Z0-9]/g, '');

// Inline oneOf/anyOf parameter schemas generate anonymous
// `struct{ union json.RawMessage }` fields with no constructors, which callers
// outside the gen package cannot build. Hoisting each distinct union into
// components.schemas makes oapi-codegen emit a named type with From*/As*/Merge*
// helpers instead (issue #579).
const hoistUnionParamSchemas = (spec) => {
  const components = (spec.components = spec.components || {});
  const schemas = (components.schemas = components.schemas || {});
  const hoisted = new Map(); // canonical union JSON -> component name

  const componentFor = (kind, members) => {
    const key = `${kind}:${JSON.stringify(members)}`;
    let name = hoisted.get(key);
    if (name) return name;

    const base = members.map((m) => pascal(m.format || m.type || 'value')).join('Or');
    name = base;
    const unionSchema = { [kind]: members };
    for (let i = 2; schemas[name] && JSON.stringify(schemas[name]) !== JSON.stringify(unionSchema); i++) {
      name = `${base}${i}`;
    }
    schemas[name] = unionSchema;
    hoisted.set(key, name);
    return name;
  };

  for (const pathObj of Object.values(spec.paths || {})) {
    for (const operation of Object.values(pathObj || {})) {
      if (!operation || typeof operation !== 'object' || !Array.isArray(operation.parameters)) continue;
      for (const param of operation.parameters) {
        const schema = param && param.schema;
        if (!schema || typeof schema !== 'object') continue;
        const kind = schema.oneOf ? 'oneOf' : schema.anyOf ? 'anyOf' : null;
        if (!kind || !Array.isArray(schema[kind])) continue;

        const name = componentFor(kind, schema[kind]);
        const { [kind]: _hoisted, ...rest } = schema;
        param.schema = { $ref: `#/components/schemas/${name}`, ...rest };
      }
    }
  }
  return hoisted.size;
};

const main = async () => {
  const specUrl = 'https://api.massive.com/openapi';
  const specFilename = path.join('./', 'openapi.json');

  try {
    const spec = await fetchSpec(specUrl);

    const paths = {};
    for (const [path, pathObj] of Object.entries(spec.paths)) {

      // Skip paths marked as drafts
      if (pathObj["x-polygon-draft"]) continue;


      // Since all endpoints use GET, process the 'get' method
      const operation = pathObj.get;
      const key = `${path}+${operation.operationId}`; // e.g., "/fed/vX/treasury-yields+TreasuryYields"
      const newOperationId = mappings[key] || operation.operationId; // Rename if mapping exists

      paths[path] = {
        ...pathObj,
        get: {
          ...operation,
          tags: ["default"], // Add default tag
          operationId: newOperationId, // Update operationId
        },
      };
    }

    // Create the cleaned spec with updated paths
    const cleanedSpec = { ...spec, paths };

    const hoisted = hoistUnionParamSchemas(cleanedSpec);
    if (hoisted) {
      console.log(`Hoisted ${hoisted} inline union parameter schema(s) into components.schemas`);
    }

    saveSpec(cleanedSpec, specFilename);
    console.log(`Spec saved to ${specFilename}`);
  } catch (error) {
    console.error(`Error: ${error.message}`);
  }
};

main();
