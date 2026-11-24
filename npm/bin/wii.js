#!/usr/bin/env node

const { execFileSync } = require("child_process");
const path = require("path");
const fs = require("fs");

const pkgName = `@oopsunix/wii-${process.platform}-${process.arch}`;
const binaryName = process.platform === "win32" ? "wii.exe" : "wii";
const binaryRelPath = `bin/${binaryName}`;

let binaryPath;
try {
  binaryPath = require.resolve(`${pkgName}/${binaryRelPath}`);
} catch {
  const fallback = path.join(__dirname, binaryName);
  if (fs.existsSync(fallback)) {
    binaryPath = fallback;
  } else {
    console.error(`Unsupported platform or missing binary package (${pkgName}).`);
    console.error("Please reinstall @oopsunix/wii or download manually from:");
    console.error("  https://github.com/oopsunix/wii/releases");
    process.exit(1);
  }
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
