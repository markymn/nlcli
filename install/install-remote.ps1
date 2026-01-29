
# nlcli Remote Installer for Windows

Write-Host "Installing nlcli..." -ForegroundColor Cyan

# 1. Check Prerequisites
if (-not (Get-Command "git" -ErrorAction SilentlyContinue)) {
    Write-Host "Error: 'git' is not installed. Please install Git and try again." -ForegroundColor Red
    exit 1
}
if (-not (Get-Command "go" -ErrorAction SilentlyContinue)) {
    Write-Host "Error: 'go' is not installed. Please install Go (golang.org) and try again." -ForegroundColor Red
    exit 1
}

# 2. Setup Directories
$installDir = [System.IO.Path]::Combine($env:USERPROFILE, ".nlcli")
$srcDir = Join-Path $installDir "src"
$binDir = Join-Path $installDir "bin"

if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Force -Path $binDir | Out-Null
}

# 3. Clone Repository
if (Test-Path $srcDir) {
    Remove-Item -Recurse -Force $srcDir
}
Write-Host "Cloning repository..." -ForegroundColor Cyan
git clone --depth 1 https://github.com/markymn/nlcli.git $srcDir | Out-Null

if (-not (Test-Path $srcDir)) {
    Write-Host "Error: Failed to clone repository." -ForegroundColor Red
    exit 1
}

# 4. Build
Write-Host "Building nlcli..." -ForegroundColor Cyan
Set-Location $srcDir
go build -o nlcli.exe ./cmd/nlcli

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed." -ForegroundColor Red
    exit 1
}

# 5. Install Binary
Move-Item -Force "nlcli.exe" "$binDir\nlcli.exe"
Write-Host "Installed to $binDir\nlcli.exe" -ForegroundColor Green

# 6. Update PATH
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

# Cleanup source
Set-Location $env:USERPROFILE
Remove-Item -Recurse -Force $srcDir
