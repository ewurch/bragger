#!/usr/bin/env node

import puppeteer from 'puppeteer';
import { resolve, dirname, basename } from 'path';
import { existsSync } from 'fs';
import { spawn, exec } from 'child_process';

const htmlPath = process.argv[2];

if (!htmlPath) {
  console.error('Usage: npm run pdf <path-to-html-file>');
  console.error('Example: npm run pdf outputs/company_role/resume.html');
  process.exit(1);
}

const absolutePath = resolve(htmlPath);

if (!existsSync(absolutePath)) {
  console.error(`File not found: ${absolutePath}`);
  process.exit(1);
}

if (!absolutePath.endsWith('.html')) {
  console.error('File must be an HTML file');
  process.exit(1);
}

const dir = dirname(absolutePath);
const filename = basename(absolutePath, '.html');
const pdfPath = resolve(dir, `${filename}.pdf`);

console.log(`Converting: ${absolutePath}`);
console.log(`Output: ${pdfPath}`);

const browser = await puppeteer.launch({ headless: true });
const page = await browser.newPage();

await page.goto(`file://${absolutePath}`, { waitUntil: 'networkidle0' });

await page.pdf({
  path: pdfPath,
  format: 'A4',
  margin: { top: 0, right: 0, bottom: 0, left: 0 },
  printBackground: true,
  preferCSSPageSize: true,
});

await browser.close();

console.log('Done!');

// Copy path to clipboard for easy pasting in file dialogs (Cmd+Shift+G)
exec(`echo "${pdfPath}" | pbcopy`);
console.log('Path copied to clipboard');

// Reveal in Finder for drag-and-drop
spawn('open', ['-R', pdfPath], { detached: true, stdio: 'ignore' });
