const fs = require('fs');
const path = require('path');
const spec = require('./openapi.json');

const examplesDir = path.join(__dirname, '../examples/go');
const tokenizedDir = path.join(__dirname, '../examples/go-tokenized');

if (!fs.existsSync(examplesDir)) fs.mkdirSync(examplesDir, { recursive: true });
if (!fs.existsSync(tokenizedDir)) fs.mkdirSync(tokenizedDir, { recursive: true });

// === Helpers ===
const toPascalCase = (str) => {
  return str
    .replace(/([a-z0-9])([A-Z])/g, '$1 $2')   // split camelCase
    .split(/[-_ ]+/)
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join('');
};

const toGoFieldName = (str) => {
  return str
    .split('.')
    .map(segment =>
      segment.split(/[-_]/)
        .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
        .join('')
    )
    .join('');
};

const toPolygonToken = (str) => {
  const words = str.replace(/([a-z])([A-Z])/g, '$1_$2').toUpperCase().split('_');
  return `TOKEN_${words.join('_')}`;
};

const toSnakeCase = (str) =>
  str
    .replace(/([a-z])([A-Z])/g, '$1_$2')
    .replace(/[-\/.]/g, '_')
    .toLowerCase();

const isIntField = (param) => {
  const type = param.schema?.type || 'string';
  return type === 'integer' || type === 'number';
};

const isPaginated = (details) => {
  const successResp = details.responses?.['200'] || details.responses?.default;
  if (!successResp?.content?.['application/json']?.schema) return false;
  const schema = successResp.content['application/json'].schema;
  return !!schema.properties?.next_url || schema.allOf?.some(s => !!s.properties?.next_url);
};

const isEnumParam = (param) => {
  if (!param.schema) return false;
  return Array.isArray(param.schema.enum) || !!param.schema.$ref;
};

// NEW helper: decides whether a TOKEN_ should be wrapped in quotes in the tokenized version
const isStringParam = (param) => {
  if (!param?.schema) return true; // most path params + unknown fields default to string
  const schema = param.schema;
  const type = schema.type;
  if (type === 'string' || !type) return true;
  if (type === 'integer' || type === 'number' || type === 'boolean') return false;
  // enums / $refs / arrays-of-strings are string-like
  return Array.isArray(schema.enum) || !!schema.$ref ||
         (schema.type === 'array' && schema.items?.type === 'string');
};

Object.entries(spec.paths).forEach(([route, methods]) => {
  Object.entries(methods).forEach(([httpMethod, details]) => {
    const operationId = details.operationId;
    if (!operationId) return;

    const goBaseName = toPascalCase(operationId);
    const goMethod = goBaseName + 'WithResponse';
    const paramsType = goBaseName + 'Params';
    const fileName = toSnakeCase(operationId);

    // Collect full path param objects (so we have .schema for type detection)
    const pathParams = (details.parameters || []).filter(p => p.in === 'path');

    const generateSnippet = (dir, useTokens = false) => {
      const lines = [
        `package main`,
        ``,
        `import (`,
        `	"context"`,
        `	"fmt"`,
        `	"log"`,
        `	"massive-go-poc/rest"`,
        `	"massive-go-poc/rest/gen"`,
        `)`,
        ``,
        `func main() {`,
        ``,
        `	c := rest.NewWithOptions("YOUR_API_KEY", rest.WithTrace(false), rest.WithPagination(true))`,
        `	ctx := context.Background()`,
        ``,
        //`	params := &gen.${paramsType}{`,
      ];


            // === Build query param lines (extra indent for inline struct) ===
      const paramLines = [];
      if (details.parameters) {
        details.parameters.forEach(param => {
          if (param.in === 'query') {
            const fieldName = toGoFieldName(param.name);
            const schemaType = param.schema?.type || 'string';

            let value;
            if (useTokens) {
              const token = toPolygonToken(param.name);
              value = isStringParam(param) ? `"${token}"` : token;
            } else {
              if (schemaType === 'boolean') value = 'true';
              else if (schemaType === 'integer' || schemaType === 'number') value = '100';
              else value = param.example !== undefined 
                ? (typeof param.example === 'string' ? `"${param.example}"` : param.example) 
                : `"${param.name}"`;
            }

            let line;
            if (isEnumParam(param)) {
              line = `			${fieldName}: rest.Ptr(gen.${paramsType}${fieldName}(${value})),`;
            } else {
              line = `			${fieldName}: rest.Ptr(${value}),`;
            }
            paramLines.push(line);
          }
        });
      }

      // === Multi-line call with inline params (exactly the style you wanted) ===
      const pathArgs = pathParams.map(param => {
        if (!useTokens) return `"${param.name}"`;
        const token = toPolygonToken(param.name);
        return isStringParam(param) ? `"${token}"` : token;
      });

      lines.push(`	resp, err := c.${goMethod}(ctx,`);

      pathArgs.forEach(arg => {
        lines.push(`		${arg},`);
      });

      lines.push(`		&gen.${paramsType}{`);
      if (paramLines.length > 0) {
        lines.push(paramLines.join('\n'));
      }
      lines.push(`		},`);
      lines.push(`	)`);

      lines.push(`	if err != nil {`);
      lines.push(`		log.Fatal(err)`);
      lines.push(`	}`);
      lines.push(``);
      lines.push(`	if err := rest.CheckResponse(resp); err != nil {`);
      lines.push(`		log.Fatal(err)`);
      lines.push(`	}`);

      if (isPaginated(details)) {
        lines.push(`	`);
        lines.push(`	iter := rest.NewIteratorFromResponse(c, resp)`);
        lines.push(`	for iter.Next() {`);
        lines.push(`		item := iter.Item()`);
        lines.push(`		fmt.Printf("%+v\\n", item)`);
        lines.push(`	}`);
        lines.push(`	if err := iter.Err(); err != nil {`);
        lines.push(`		log.Fatal(err)`);
        lines.push(`	}`);
      } else {
        lines.push(`	fmt.Printf("%+v\\n", resp.JSON200)`);
      }

      lines.push(`}`);
      lines.push(``);

      fs.writeFileSync(path.join(dir, `${fileName}.go`), lines.join('\n'));
    };

    generateSnippet(examplesDir, false);
    generateSnippet(tokenizedDir, true);
  });
});

console.log('🎉 All Go examples generated with string tokens properly quoted in the tokenized version!');
