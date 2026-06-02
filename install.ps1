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

# Construct download URL (goreleaser produces raw binaries without version in filename)
$Filename = "${BinaryName}_windows_${Arch}.exe"
$Url = "https://github.com/$Repo/releases/latest/download/$Filename"
$MirrorPrefix = "https://hubp.llkk.cc/"
$TmpFile = "$env:TEMP\${BinaryName}.exe"

# Check GitHub accessibility, fall back to mirror if unreachable
Write-Host "Checking GitHub connectivity..."
try {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    $Request = [System.Net.WebRequest]::Create("https://github.com")
    $Request.Timeout = 10000
    $Request.Method = "HEAD"
    $Response = $Request.GetResponse()
    $Response.Close()
} catch {
    Write-Warning "GitHub unreachable, using mirror..."
    $Url = "${MirrorPrefix}${Url}"
}

# Download
Write-Host "Downloading $Url..."
try {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    Invoke-WebRequest -Uri $Url -OutFile $TmpFile -UseBasicParsing
} catch {
    Write-Error "Error: Download failed"
    Write-Host $_.Exception.Message
    exit 1
}

# Install
Write-Host "Installing to $InstallDir..."
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

try {
    Move-Item -Path $TmpFile -Destination "$InstallDir\$BinaryName.exe" -Force
} catch {
    Write-Error "Error: Installation failed"
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
