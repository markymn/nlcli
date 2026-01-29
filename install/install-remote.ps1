# nlcli Remote Installer for Windows

Write-Host "Installing nlcli..." -ForegroundColor Cyan

# 1. Detect Architecture
$arch = "amd64"
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { $arch = "arm64" }

$binaryName = "nlcli-windows-$arch.exe"
$downloadUrl = "https://github.com/markymn/nlcli/releases/latest/download/$binaryName"

# 2. Setup Directories
$installDir = [System.IO.Path]::Combine($env:USERPROFILE, ".nlcli")
$binDir = Join-Path $installDir "bin"

if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Force -Path $binDir | Out-Null
}

# 3. Download Binary
Write-Host "Downloading $binaryName..." -ForegroundColor Cyan
$targetFile = Join-Path $binDir "nlcli.exe"

try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $targetFile -ErrorAction Stop
}
catch {
    Write-Host "Error: Failed to download binary from GitHub. Please check your connection or ensures a release exists." -ForegroundColor Red
    exit 1
}

Write-Host "Installed to $targetFile" -ForegroundColor Green

# 4. Update PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
$pathParts = $currentPath -split ';'

if ($pathParts -notcontains $binDir) {
    [Environment]::SetEnvironmentVariable("Path", $currentPath + ";" + $binDir, "User")
    Write-Host "Success: Added '$binDir' to your User PATH." -ForegroundColor Green
    Write-Host "Please restart your terminal to use 'nlcli'." -ForegroundColor Yellow
}
else {
    Write-Host "Success: 'nlcli' is ready to use!" -ForegroundColor Green
}
