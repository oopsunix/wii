$ErrorActionPreference = "Stop"

$Repo = "oopsunix/wii"
$BinaryName = "wii"
$InstallDir = "$env:USERPROFILE\.local\bin"

# Colors
function Write-Success { Write-Host $args -ForegroundColor Green }
function Write-Warning { Write-Host $args -ForegroundColor Yellow }
function Write-Error { Write-Host $args -ForegroundColor Red }

# Detect architecture
if ([Environment]::Is64BitOperatingSystem) {
    $Arch = "amd64"
} else {
    $Arch = "386"
}

Write-Success "Detected: windows/$Arch"

# Get latest version from GitHub API
Write-Host "Fetching latest version..."
try {
    $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
    $Version = $Release.tag_name.TrimStart('v')
} catch {
    Write-Error "Error: Could not fetch latest version"
    Write-Host "Please visit: https://github.com/$Repo/releases"
    exit 1
}

Write-Success "Latest version: v$Version"

# Construct download URL
$Filename = "${BinaryName}_${Version}_windows_${Arch}"
$Url = "https://github.com/$Repo/releases/download/v$Version/${Filename}.zip"
$TmpZip = "$env:TEMP\${BinaryName}.zip"

# Download
Write-Host "Downloading $Url..."
try {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    Invoke-WebRequest -Uri $Url -OutFile $TmpZip -UseBasicParsing
} catch {
    Write-Error "Error: Download failed"
    Write-Host $_.Exception.Message
    exit 1
}

# Extract and install
Write-Host "Installing to $InstallDir..."
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

try {
    Expand-Archive -Path $TmpZip -DestinationPath $InstallDir -Force
    Remove-Item $TmpZip -Force
} catch {
    Write-Error "Error: Extraction failed"
    Write-Host $_.Exception.Message
    exit 1
}

$ExePath = "$InstallDir\$BinaryName.exe"

if (Test-Path $ExePath) {
    Write-Success "Installed: $ExePath"
} else {
    Write-Warning "Warning: Binary may not have been extracted correctly"
}

# Add to PATH if not already present
$CurrentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($CurrentPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$CurrentPath;$InstallDir", "User")
    Write-Success "Added $InstallDir to user PATH"
    Write-Warning "Restart your terminal to apply PATH changes"
} else {
    Write-Host "$InstallDir is already in PATH"
}

# Verify installation
Write-Host ""
Write-Success "Installation successful!"
Write-Host "Run '$BinaryName' to get started."
Write-Host ""
Write-Host "Usage examples:"
Write-Host "  $BinaryName                    # Scan and display tools"
Write-Host "  $BinaryName --probe            # Include version probing"
Write-Host "  WII_FORMAT=json $BinaryName    # JSON output"
