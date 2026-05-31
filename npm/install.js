#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const https = require('https');
const { createWriteStream } = require('fs');
const { pipeline } = require('stream');
const { promisify } = require('util');

const pipelineAsync = promisify(pipeline);

const REPO = 'oopsunix/wii';
const VERSION = '1.0.0';

function getPlatform() {
  const platform = process.platform;
  switch (platform) {
    case 'darwin': return 'darwin';
    case 'linux': return 'linux';
    case 'win32': return 'windows';
    default:
      console.error(`Unsupported platform: ${platform}`);
      process.exit(1);
  }
}

function getArch() {
  const arch = process.arch;
  switch (arch) {
    case 'x64': return 'amd64';
    case 'arm64': return 'arm64';
    case 'ia32': return '386';
    default:
      console.error(`Unsupported architecture: ${arch}`);
      process.exit(1);
  }
}

function getBinaryName() {
  return process.platform === 'win32' ? 'wii.exe' : 'wii';
}

function getDownloadUrl() {
  const platform = getPlatform();
  const arch = getArch();
  const ext = platform === 'windows' ? 'zip' : 'tar.gz';
  return `https://github.com/${REPO}/releases/download/v${VERSION}/wii_${VERSION}_${platform}_${arch}.${ext}`;
}

async function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    https.get(url, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        // Follow redirect
        downloadFile(response.headers.location, dest).then(resolve).catch(reject);
        return;
      }
      if (response.statusCode !== 200) {
        reject(new Error(`Download failed: ${response.statusCode}`));
        return;
      }
      const file = createWriteStream(dest);
      pipelineAsync(response, file).then(resolve).catch(reject);
    }).on('error', reject);
  });
}

async function extractArchive(archivePath, destDir) {
  const platform = process.platform;
  if (platform === 'win32') {
    // Use PowerShell for Windows
    execSync(`powershell -Command "Expand-Archive -Path '${archivePath}' -DestinationPath '${destDir}' -Force"`, { stdio: 'inherit' });
  } else {
    // Use tar for Unix
    execSync(`tar xzf "${archivePath}" -C "${destDir}"`, { stdio: 'inherit' });
  }
}

async function main() {
  const platform = getPlatform();
  const arch = getArch();
  const binaryName = getBinaryName();

  console.log(`Installing wii for ${platform}/${arch}...`);

  // Create bin directory
  const binDir = path.join(__dirname, 'bin');
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  // Download binary
  const url = getDownloadUrl();
  const tmpDir = path.join(__dirname, 'tmp');
  if (!fs.existsSync(tmpDir)) {
    fs.mkdirSync(tmpDir, { recursive: true });
  }

  const ext = platform === 'windows' ? 'zip' : 'tar.gz';
  const archivePath = path.join(tmpDir, `wii.${ext}`);

  console.log(`Downloading from ${url}...`);
  try {
    await downloadFile(url, archivePath);
  } catch (error) {
    console.error('Download failed:', error.message);
    console.error('Please download manually from: https://github.com/oopsunix/wii/releases');
    process.exit(1);
  }

  // Extract
  console.log('Extracting...');
  try {
    await extractArchive(archivePath, tmpDir);
  } catch (error) {
    console.error('Extraction failed:', error.message);
    process.exit(1);
  }

  // Move binary to bin directory
  const binaryPath = path.join(tmpDir, binaryName);
  const destPath = path.join(binDir, binaryName);

  if (fs.existsSync(binaryPath)) {
    fs.copyFileSync(binaryPath, destPath);
    if (platform !== 'windows') {
      fs.chmodSync(destPath, '755');
    }
    console.log(`Installed: ${destPath}`);
  } else {
    console.error('Binary not found after extraction');
    process.exit(1);
  }

  // Cleanup
  try {
    fs.rmSync(tmpDir, { recursive: true, force: true });
  } catch (error) {
    // Ignore cleanup errors
  }

  console.log('');
  console.log('Installation complete!');
  console.log('Run "wii" to get started.');
}

main().catch((error) => {
  console.error('Installation failed:', error.message);
  process.exit(1);
});
