#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const https = require("https");
const { createWriteStream } = require("fs");
const { pipeline } = require("stream");
const { promisify } = require("util");

const pipelineAsync = promisify(pipeline);

const REPO = "oopsunix/wii";
const { version } = require("./package.json");

function getPlatform() {
  const platform = process.platform;
  switch (platform) {
    case "darwin":
      return "darwin";
    case "linux":
      return "linux";
    case "win32":
      return "windows";
    default:
      console.error(`Unsupported platform: ${platform}`);
      process.exit(1);
  }
}

function getArch() {
  const arch = process.arch;
  switch (arch) {
    case "x64":
      return "amd64";
    case "arm64":
      return "arm64";
    case "ia32":
      return "386";
    default:
      console.error(`Unsupported architecture: ${arch}`);
      process.exit(1);
  }
}

function getBinaryName() {
  return process.platform === "win32" ? "wii.exe" : "wii";
}

function getDownloadUrl() {
  const platform = getPlatform();
  const arch = getArch();
  const ext = platform === "windows" ? ".exe" : "";
  return `https://github.com/${REPO}/releases/download/v${version}/wii_${platform}_${arch}${ext}`;
}

async function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    https
      .get(url, (response) => {
        if (response.statusCode === 302 || response.statusCode === 301) {
          downloadFile(response.headers.location, dest)
            .then(resolve)
            .catch(reject);
          return;
        }
        if (response.statusCode !== 200) {
          reject(new Error(`Download failed: ${response.statusCode}`));
          return;
        }
        const file = createWriteStream(dest);
        pipelineAsync(response, file).then(resolve).catch(reject);
      })
      .on("error", reject);
  });
}

async function main() {
  const platform = getPlatform();
  const arch = getArch();
  const binaryName = getBinaryName();

  console.log(`Installing wii for ${platform}/${arch}...`);

  const binDir = path.join(__dirname, "bin");
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  const url = getDownloadUrl();
  const destPath = path.join(binDir, binaryName);

  console.log(`Downloading from ${url}...`);
  try {
    await downloadFile(url, destPath);
  } catch (error) {
    console.error("Download failed:", error.message);
    console.error(
      "Please download manually from: https://github.com/oopsunix/wii/releases",
    );
    process.exit(1);
  }

  if (platform !== "windows") {
    fs.chmodSync(destPath, "755");
  }

  console.log("");
  console.log("Installation complete!");
  console.log('Run "wii" to get started.');
}

main().catch((error) => {
  console.error("Installation failed:", error.message);
  process.exit(1);
});
