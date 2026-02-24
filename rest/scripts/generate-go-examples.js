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
    .replace(/([a-z0-9])([A-Z])/g, '$1 $2')
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

const shouldUseTypedEnum = (param) => {
  return !!param.schema?.$ref || param.schema?.default !== undefined;
};

const isStringParam = (param) => {
  if (!param?.schema) return true;
  const schema = param.schema;
  const type = schema.type;
  if (type === 'string' || !type) return true;
  if (type === 'integer' || type === 'number' || type === 'boolean') return false;
  return Array.isArray(schema.enum) || !!schema.$ref ||
         (schema.type === 'array' && schema.items?.type === 'string');
};

const getParamValue = (param, useTokens = false) => {
  if (useTokens) {
    const token = toPolygonToken(param.name);
    return isStringParam(param) ? `"${token}"` : token;
  }

  let value = param.example ?? param.schema?.example ?? param.schema?.default;
  if (value !== undefined) {
    return typeof value === 'string' ? `"${value}"` : value;
  }

  const schemaType = param.schema?.type || 'string';
  if (schemaType === 'boolean') return 'true';
  if (schemaType === 'integer' || schemaType === 'number') return '100';
  return `"${param.name}"`;
};

Object.entries(spec.paths).forEach(([route, methods]) => {
  Object.entries(methods).forEach(([httpMethod, details]) => {
    const operationId = details.operationId;
    if (!operationId) return;

    const goBaseName = toPascalCase(operationId);
    const goMethod = goBaseName + 'WithResponse';
    const paramsType = goBaseName + 'Params';
    const fileName = toSnakeCase(operationId);

    const allParams = details.parameters || [];
    const pathParams = allParams.filter(p => p.in === 'path');
    const queryParams = allParams.filter(p => p.in === 'query');

    const hasQueryParams = queryParams.length > 0;

    const generateSnippet = (dir, useTokens = false) => {
      const lines = [
        `package main`,
        ``,
        `import (`,
        `	"context"`,
        `	"fmt"`,
        `	"log"`,
        `	"github.com/massive-com/client-go/v3/rest"`,
      ];

      if (hasQueryParams) {
        lines.push(`	"github.com/massive-com/client-go/v3/rest/gen"`);
      }

      lines.push(`)`);
      lines.push(``);
      lines.push(`func main() {`);
      lines.push(``);

      lines.push(`	c := rest.NewWithOptions("GLOBAL_TOKEN_API_KEY",`);
      lines.push(`		rest.WithTrace(false),`);
      lines.push(`		rest.WithPagination(true),`);
      lines.push(`	)`);
      lines.push(`	ctx := context.Background()`);
      lines.push(``);

      if (hasQueryParams) {
        lines.push(`	params := &gen.${paramsType}{`);

        const paramLines = [];
        queryParams.forEach(param => {
          const fieldName = toGoFieldName(param.name);
          const value = getParamValue(param, useTokens);

          let line;
          if (isEnumParam(param)) {
            if (shouldUseTypedEnum(param)) {
              line = `		${fieldName}: rest.Ptr(gen.${paramsType}${fieldName}(${value})),`;
            } else {
              line = `		${fieldName}: ${value},`;
            }
          } else {
            line = `		${fieldName}: rest.Ptr(${value}),`;
          }
          paramLines.push(line);
        });

        lines.push(paramLines.join('\n'));
        lines.push(`	}`);
        lines.push(``);
      }

      // === UNIFORM MULTILINE CALL FOR EVERY ENDPOINT (fixes the syntax error) ===
      // Every argument (including ctx) is now on its own line → perfect for token replacement
      lines.push(`	resp, err := c.${goMethod}(`);
      lines.push(`		ctx,`);

      // Path parameters (if any)
      pathParams.forEach(param => {
        const value = getParamValue(param, useTokens);
        lines.push(`		${value},`);
      });

      // Query params struct (if any)
      if (hasQueryParams) {
        lines.push(`		params,`);
      }

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

console.log('🎉 All Go examples generated.');
