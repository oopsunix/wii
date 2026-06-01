#!/usr/bin/env node

const { execFileSync } = require("child_process");
const path = require("path");
const fs = require("fs");

const binDir = __dirname;
const binaryName = process.platform === "win32" ? "wii.exe" : "wii";
const binaryPath = path.join(binDir, binaryName);

if (!fs.existsSync(binaryPath)) {
  console.error("Binary not found. Please reinstall the package:");
  console.error("  npm install -g @oopsunix/wii");
  process.exit(1);
}

const args = process.argv.slice(2);

try {
  execFileSync(binaryPath, args, {
    stdio: "inherit",
    windowsHide: true,
  });
} catch (error) {
  if (error.status !== undefined) {
    process.exit(error.status);
  }
  console.error("Failed to run wii:", error.message);
  process.exit(1);
}
