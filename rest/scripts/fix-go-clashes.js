const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const filePath = path.join(__dirname, '../gen/client.gen.go');

let content = fs.readFileSync(filePath, 'utf8');

const tagToName = {
  'P': 'AskPrice',
  'p': 'BidPrice',
  'S': 'AskSize',
  's': 'BidSize',
  'X': 'AskExchange',
  'x': 'BidExchange',
  'T': 'Ticker',
  't': 'Timestamp',
};

console.log('🔧 Fixing P/p, S/s, X/x clashes in anonymous structs...');

// General regex (catches most cases)
const fieldRegex = /^([ \t]+)([A-Z])\s+(.+?)\s*`json:"([^",]+)(,omitempty)?"`/gm;
let count = 0;

content = content.replace(fieldRegex, (match, indent, oldName, type, jsonKey, omitempty) => {
  const newName = tagToName[jsonKey];
  omitempty = omitempty || '';
  if (newName && newName !== oldName) {
    console.log(`   ${oldName} → ${newName} (json:"${jsonKey}")`);
    count++;
    return `${indent}${newName} ${type.trim()} \`json:"${jsonKey}${omitempty}"\``;
  }
  return match;
});

console.log(`✅ General rename: ${count} fields`);

// Fallback — guaranteed to catch the exact struct you posted
const fallbacks = [
  [/^([ \t]+)P\s+\*float64\s*`json:"P,omitempty"`/gm, '$1AskPrice *float64 `json:"P,omitempty"`'],
  [/^([ \t]+)P\s+\*float64\s*`json:"p,omitempty"`/gm, '$1BidPrice *float64 `json:"p,omitempty"`'],
  [/^([ \t]+)S\s+\*int\s*`json:"S,omitempty"`/gm, '$1AskSize *int `json:"S,omitempty"`'],
  [/^([ \t]+)S\s+\*int\s*`json:"s,omitempty"`/gm, '$1BidSize *int `json:"s,omitempty"`'],
  [/^([ \t]+)X\s+\*int\s*`json:"X,omitempty"`/gm, '$1AskExchange *int `json:"X,omitempty"`'],
  [/^([ \t]+)X\s+\*int\s*`json:"x,omitempty"`/gm, '$1BidExchange *int `json:"x,omitempty"`'],
];

fallbacks.forEach(([regex, repl]) => {
  const before = (content.match(regex) || []).length;
  content = content.replace(regex, repl);
  if (before) console.log(`   Fallback fixed ${before} exact matches`);
});

// Fix any .P / .S / .X usages elsewhere in the file
const usageMap = { 'P': 'AskPrice', 'S': 'AskSize', 'X': 'AskExchange', 'T': 'Ticker' };
Object.entries(usageMap).forEach(([old, neu]) => {
  content = content.replace(new RegExp(`\\.${old}\\b`, 'g'), `.${neu}`);
});

fs.writeFileSync(filePath, content);

// Clean up formatting
try {
  execSync(`gofmt -w "${filePath}"`);
  console.log('✅ Ran gofmt');
} catch (e) {
  console.log('⚠️  gofmt skipped (optional)');
}

console.log('\n🎉 Fix complete! Now run:');
console.log('   go build ./rest/gen');
