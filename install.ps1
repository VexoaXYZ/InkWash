# InkWash Installer for Windows
# Usage: irm https://raw.githubusercontent.com/VexoaXYZ/InkWash/master/install.ps1 | iex

param(
    [string]$InstallDir = "$env:LOCALAPPDATA\InkWash",
    [switch]$NoPath,
    [switch]$Desktop
)

$ErrorActionPreference = "Stop"

Write-Host "🎨 InkWash Installer" -ForegroundColor Cyan
Write-Host "===================" -ForegroundColor Cyan
Write-Host ""

# Get latest release info from GitHub
Write-Host "📡 Fetching latest release..." -ForegroundColor Yellow
try {
    $apiUrl = "https://api.github.com/repos/VexoaXYZ/InkWash/releases/latest"
    $release = Invoke-RestMethod -Uri $apiUrl -Headers @{
        "User-Agent" = "InkWash-Installer"
    }
    $version = $release.tag_name
    Write-Host "   Latest version: $version" -ForegroundColor Green
} catch {
    Write-Host "❌ Failed to fetch release info. Please check your internet connection." -ForegroundColor Red
    exit 1
}

# Find Windows AMD64 asset
$asset = $release.assets | Where-Object { $_.name -like "*windows-amd64.zip" } | Select-Object -First 1
if (-not $asset) {
    Write-Host "❌ Could not find Windows release asset." -ForegroundColor Red
    exit 1
}

# Create install directory
Write-Host ""
Write-Host "📁 Installing to: $InstallDir" -ForegroundColor Yellow
if (Test-Path $InstallDir) {
    Write-Host "   Removing old installation..." -ForegroundColor Gray
    Remove-Item -Path $InstallDir -Recurse -Force
}
New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null

# Download release
$downloadUrl = $asset.browser_download_url
$zipPath = Join-Path $env:TEMP "inkwash-$version.zip"
Write-Host ""
Write-Host "⬇️  Downloading InkWash $version..." -ForegroundColor Yellow
try {
    $ProgressPreference = 'SilentlyContinue'
    Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -UseBasicParsing
    Write-Host "   Downloaded successfully!" -ForegroundColor Green
} catch {
    Write-Host "❌ Download failed: $_" -ForegroundColor Red
    exit 1
}

# Extract archive
Write-Host ""
Write-Host "📦 Extracting files..." -ForegroundColor Yellow
try {
    Expand-Archive -Path $zipPath -DestinationPath $InstallDir -Force
    Remove-Item $zipPath
    Write-Host "   Extracted successfully!" -ForegroundColor Green
} catch {
    Write-Host "❌ Extraction failed: $_" -ForegroundColor Red
    exit 1
}

# Add to PATH if requested
if (-not $NoPath) {
    Write-Host ""
    Write-Host "🔧 Adding to PATH..." -ForegroundColor Yellow

    # Get current user PATH
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")

    if ($userPath -notlike "*$InstallDir*") {
        $newPath = "$userPath;$InstallDir"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        $env:Path = "$env:Path;$InstallDir"
        Write-Host "   Added to PATH! (You may need to restart your terminal)" -ForegroundColor Green
    } else {
        Write-Host "   Already in PATH!" -ForegroundColor Green
    }
}

# Create desktop shortcut if requested
if ($Desktop) {
    Write-Host ""
    Write-Host "🖥️  Creating desktop shortcut..." -ForegroundColor Yellow
    $WshShell = New-Object -ComObject WScript.Shell
    $Shortcut = $WshShell.CreateShortcut("$env:USERPROFILE\Desktop\InkWash.lnk")
    $Shortcut.TargetPath = "powershell.exe"
    $Shortcut.Arguments = "-NoExit -Command `"cd '$InstallDir'; .\inkwash.exe`""
    $Shortcut.WorkingDirectory = $InstallDir
    $Shortcut.Description = "InkWash - FiveM Server Manager"
    $Shortcut.Save()
    Write-Host "   Desktop shortcut created!" -ForegroundColor Green
}

# Success message
Write-Host ""
Write-Host "✅ InkWash installed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "🚀 Quick Start:" -ForegroundColor Cyan
if (-not $NoPath) {
    Write-Host "   1. Open a NEW terminal (to load updated PATH)" -ForegroundColor White
    Write-Host "   2. Run: inkwash create" -ForegroundColor White
} else {
    Write-Host "   1. Run: cd '$InstallDir'" -ForegroundColor White
    Write-Host "   2. Run: .\inkwash.exe create" -ForegroundColor White
}
Write-Host ""
Write-Host "📚 Documentation: https://github.com/VexoaXYZ/InkWash/wiki" -ForegroundColor Cyan
Write-Host "🐛 Issues: https://github.com/VexoaXYZ/InkWash/issues" -ForegroundColor Cyan
Write-Host ""
