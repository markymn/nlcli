# Build script for nlcli
$distDir = "dist"
if (Test-Path $distDir) { Remove-Item -Recurse -Force $distDir }
New-Item -ItemType Directory -Path $distDir | Out-Null

$platforms = @(
    @{os = "windows"; arch = "amd64"; ext = ".exe" },
    @{os = "windows"; arch = "arm64"; ext = ".exe" },
    @{os = "linux"; arch = "amd64"; ext = "" },
    @{os = "linux"; arch = "arm64"; ext = "" },
    @{os = "darwin"; arch = "amd64"; ext = "" },
    @{os = "darwin"; arch = "arm64"; ext = "" }
)

foreach ($p in $platforms) {
    $name = "nlcli-$($p.os)-$($p.arch)$($p.ext)"
    Write-Host "Building $name..." -ForegroundColor Cyan
    $env:GOOS = $p.os
    $env:GOARCH = $p.arch
    go build -o "$distDir/$name" ./cmd/nlcli
}

Write-Host "Build complete. Binaries are in '$distDir/'" -ForegroundColor Green
