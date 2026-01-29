$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptDir
$targetDir = $projectRoot

Write-Host "Building nlcli..." -ForegroundColor Cyan
Set-Location $projectRoot
go build -o nlcli.exe ./cmd/nlcli

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed. Please ensure Go is installed and configured correctly." -ForegroundColor Red
    exit 1
}

Write-Host "Build successful." -ForegroundColor Green

$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
$pathParts = $currentPath -split ';'

if ($pathParts -notcontains $targetDir) {
    [Environment]::SetEnvironmentVariable("Path", $currentPath + ";" + $targetDir, "User")
    Write-Host "Success: '$targetDir' has been added to your User PATH." -ForegroundColor Green
    Write-Host "Please restart your terminal for changes to take effect." -ForegroundColor Yellow
}
else {
    Write-Host "Info: '$targetDir' is already in your User PATH." -ForegroundColor Cyan
}
