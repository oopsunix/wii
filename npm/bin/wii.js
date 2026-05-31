#!/usr/bin/env node

const { execFileSync } = require('child_process');
const path = require('path');
const fs = require('fs');

// Get the binary path
const binDir = path.join(__dirname);
const binaryName = process.platform === 'win32' ? 'wii.exe' : 'wii';
const binaryPath = path.join(binDir, binaryName);

// Check if binary exists
if (!fs.existsSync(binaryPath)) {
  console.error('Binary not found. Please reinstall the package:');
  console.error('  npm install -g @oopsunix/wii');
  process.exit(1);
}

// Get command line arguments (skip 'node' and script name)
const args = process.argv.slice(2);

try {
  // Run the binary with arguments
  execFileSync(binaryPath, args, {
    stdio: 'inherit',
    windowsHide: true
  });
} catch (error) {
  // If the binary exits with a non-zero code, exit with the same code
  if (error.status !== undefined) {
    process.exit(error.status);
  }
  // If there's another error, report it
  console.error('Failed to run wii:', error.message);
  process.exit(1);
}
