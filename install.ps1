# InkWash Installer for Windows
# Downloads and installs the latest InkWash CLI

param(
    [string]$InstallPath = "$env:USERPROFILE\.local\bin"
)

$ErrorActionPreference = "Stop"

Write-Host "🚀 InkWash Installer for Windows" -ForegroundColor Cyan
Write-Host "Installing to: $InstallPath" -ForegroundColor Green

# Create installation directory
if (!(Test-Path $InstallPath)) {
    New-Item -ItemType Directory -Path $InstallPath -Force | Out-Null
    Write-Host "✅ Created installation directory" -ForegroundColor Green
}

# Download latest release
$latestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/VexoaXYZ/InkWash/releases/latest"
$windowsAsset = $latestRelease.assets | Where-Object { $_.name -eq "inkwash-windows-amd64.exe" }

if (!$windowsAsset) {
    Write-Error "❌ Windows binary not found in latest release"
    exit 1
}

$downloadUrl = $windowsAsset.browser_download_url
$destinationPath = "$InstallPath\inkwash.exe"

Write-Host "📥 Downloading InkWash..." -ForegroundColor Yellow
Invoke-WebRequest -Uri $downloadUrl -OutFile $destinationPath

# Make executable
Write-Host "✅ Downloaded successfully!" -ForegroundColor Green

# Add to PATH if not already present
$currentPath = [Environment]::GetEnvironmentVariable("PATH", [EnvironmentVariableTarget]::User)
if ($currentPath -notlike "*$InstallPath*") {
    [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$InstallPath", [EnvironmentVariableTarget]::User)
    Write-Host "✅ Added to PATH (restart terminal or run refreshenv)" -ForegroundColor Green
} else {
    Write-Host "✅ Already in PATH" -ForegroundColor Green
}

Write-Host "`n🎉 InkWash installed successfully!" -ForegroundColor Green
Write-Host "Run 'inkwash --help' to get started" -ForegroundColor Cyan
Write-Host "Note: You may need to restart your terminal or run 'refreshenv' for PATH changes to take effect" -ForegroundColor Yellow