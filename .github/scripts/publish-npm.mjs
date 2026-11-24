#!/usr/bin/env node

import { execFileSync } from "node:child_process";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";

const version = process.env.VERSION;
const distDir = path.resolve(process.env.DIST_DIR || "dist");

if (!version) {
  console.error("VERSION environment variable is required");
  process.exit(1);
}

const PLATFORMS = [
  { os: "darwin", cpu: "arm64", asset: "wii_darwin_arm64", bin: "wii" },
  { os: "darwin", cpu: "x64", asset: "wii_darwin_amd64", bin: "wii" },
  { os: "linux", cpu: "arm64", asset: "wii_linux_arm64", bin: "wii" },
  { os: "linux", cpu: "x64", asset: "wii_linux_amd64", bin: "wii" },
  { os: "win32", cpu: "arm64", asset: "wii_windows_arm64.exe", bin: "wii.exe" },
  { os: "win32", cpu: "x64", asset: "wii_windows_amd64.exe", bin: "wii.exe" },
];

const npmDir = path.resolve("npm");
const rootPkgPath = path.join(npmDir, "package.json");
const rootPkg = JSON.parse(fs.readFileSync(rootPkgPath, "utf8"));
const stageRoot = fs.mkdtempSync(path.join(os.tmpdir(), "wii-npm-"));

try {
  const optionalDeps = {};

  for (const target of PLATFORMS) {
    const suffix = `${target.os}-${target.cpu}`;
    const pkgVersion = `${version}-${suffix}`;
    const srcBinary = path.join(distDir, target.asset);

    if (!fs.existsSync(srcBinary)) {
      throw new Error(`Missing release binary: ${srcBinary}`);
    }

    const pkgDir = path.join(stageRoot, suffix);
    const binDir = path.join(pkgDir, "bin");
    fs.mkdirSync(binDir, { recursive: true });

    const destBinary = path.join(binDir, target.bin);
    fs.copyFileSync(srcBinary, destBinary);
    if (target.os !== "win32") {
      fs.chmodSync(destBinary, 0o755);
    }

    const platformPkg = {
      name: rootPkg.name,
      version: pkgVersion,
      description: `${rootPkg.description} (${suffix} binary)`,
      homepage: rootPkg.homepage,
      license: rootPkg.license,
      repository: rootPkg.repository,
      os: [target.os],
      cpu: [target.cpu],
      files: ["bin/"],
    };

    fs.writeFileSync(
      path.join(pkgDir, "package.json"),
      JSON.stringify(platformPkg, null, 2) + "\n",
    );

    console.log(`Publishing ${rootPkg.name}@${pkgVersion}...`);
    execFileSync(
      "npm",
      ["publish", "--access", "public", "--tag", `platform-${suffix}`],
      { cwd: pkgDir, stdio: "inherit" },
    );

    optionalDeps[`@oopsunix/wii-${suffix}`] = `npm:${rootPkg.name}@${pkgVersion}`;
  }

  rootPkg.version = version;
  rootPkg.optionalDependencies = optionalDeps;
  fs.writeFileSync(rootPkgPath, JSON.stringify(rootPkg, null, 2) + "\n");

  console.log(`Publishing ${rootPkg.name}@${version}...`);
  execFileSync("npm", ["publish", "--access", "public"], {
    cwd: npmDir,
    stdio: "inherit",
  });
} finally {
  fs.rmSync(stageRoot, { recursive: true, force: true });
}
