# build.ps1
# One-command Windows build for Dayside.
# Prompts for permission before installing anything.
# Run from PowerShell in the project root: .\build.ps1

$ErrorActionPreference = "Stop"

# ---------- helpers ----------

function Write-Info   { Write-Host $args[0] -ForegroundColor Cyan }
function Write-Step   { Write-Host $args[0] -ForegroundColor Yellow }
function Write-OK     { Write-Host $args[0] -ForegroundColor Green }
function Write-Fail   { Write-Host $args[0] -ForegroundColor Red }

function Ask-YesNo($prompt) {
    while ($true) {
        $response = Read-Host "$prompt [y/N]"
        if ($response -eq '' -or $response -match '^[nN]') { return $false }
        if ($response -match '^[yY]') { return $true }
    }
}

function Abort($msg) {
    Write-Host ""
    Write-Fail "=== BUILD ABORTED ==="
    Write-Fail $msg
    Write-Host ""
    exit 1
}

function Has-Command($name) {
    return $null -ne (Get-Command $name -ErrorAction SilentlyContinue)
}

function Refresh-Path {
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" +
                [System.Environment]::GetEnvironmentVariable("Path", "User")
    $gopath = (& go env GOPATH 2>$null)
    if ($gopath -and (Test-Path "$gopath\bin")) {
        $env:Path = "$gopath\bin;" + $env:Path
    }
}

# ---------- install chain: Go ----------

function Install-Go-Winget {
    if (-not (Has-Command winget)) {
        Write-Host "  winget is not available on this machine."
        return $false
    }
    if (-not (Ask-YesNo "Try: winget install GoLang.Go --silent")) { return $false }
    Write-Step "Running winget..."
    try {
        & winget install --id GoLang.Go --silent --accept-package-agreements --accept-source-agreements
        if ($LASTEXITCODE -eq 0) {
            Refresh-Path
            if (Has-Command go) { return $true }
        }
        Write-Fail "  winget finished but Go is still not available."
    } catch {
        Write-Fail "  winget failed: $_"
    }
    return $false
}

function Install-Go-Choco {
    if (-not (Has-Command choco)) {
        Write-Host "  Chocolatey (choco) is not installed on this machine."
        return $false
    }
    if (-not (Ask-YesNo "Try: choco install golang -y (requires admin)")) { return $false }
    Write-Step "Running choco..."
    try {
        & choco install golang -y
        if ($LASTEXITCODE -eq 0) {
            Refresh-Path
            if (Has-Command go) { return $true }
        }
        Write-Fail "  choco finished but Go is still not available."
    } catch {
        Write-Fail "  choco failed: $_"
    }
    return $false
}

function Install-Go-MSI {
    Write-Host ""
    Write-Host "  This will download the official Go MSI from https://go.dev/dl/"
    Write-Host "  and run the installer, which asks for administrator permission."
    Write-Host "  Download size: ~150 MB."
    if (-not (Ask-YesNo "Proceed with MSI download and install?")) { return $false }

    $goVersion = "1.23.4"
    $msiUrl    = "https://go.dev/dl/go$goVersion.windows-amd64.msi"
    $msiPath   = Join-Path $env:TEMP "go$goVersion.windows-amd64.msi"

    Write-Step "Downloading $msiUrl ..."
    try {
        Invoke-WebRequest -Uri $msiUrl -OutFile $msiPath -UseBasicParsing
    } catch {
        Write-Fail "  Download failed: $_"
        return $false
    }

    Write-Step "Launching installer (watch for the UAC prompt)..."
    try {
        Start-Process msiexec.exe -ArgumentList "/i `"$msiPath`" /passive" -Wait
    } catch {
        Write-Fail "  Installer failed: $_"
        return $false
    }

    Refresh-Path
    if (Has-Command go) { return $true }
    Write-Fail "  Installer finished but Go is still not on PATH. Open a new PowerShell and re-run this script."
    return $false
}

function Ensure-Go {
    if (Has-Command go) {
        $goVersionOutput = (& go version)
        Write-OK "  Go: $goVersionOutput"

        # Parse "go version go1.22.0 windows/amd64" and verify >= 1.22
        if ($goVersionOutput -match 'go(\d+)\.(\d+)') {
            $major = [int]$Matches[1]
            $minor = [int]$Matches[2]
            if ($major -lt 1 -or ($major -eq 1 -and $minor -lt 22)) {
                Write-Host ""
                Write-Fail "  Your Go version is too old."
                Write-Host "  Dayside requires Go 1.22 or newer."
                Write-Host "  The Wails CLI install will fail on older Go due to an"
                Write-Host "  incompatibility in golang.org/x/tools v0.17.0."
                Write-Host ""
                if (Ask-YesNo "Attempt to upgrade Go automatically?") {
                    if (Install-Go-Winget) { Write-OK "  Go upgraded via winget"; return }
                    if (Install-Go-Choco)  { Write-OK "  Go upgraded via choco"; return }
                    if (Install-Go-MSI)    { Write-OK "  Go upgraded via MSI"; return }
                    Abort "Go upgrade failed. Install Go 1.22+ manually from https://go.dev/dl/"
                } else {
                    Abort "Go 1.22+ is required. Install it and re-run .\build.ps1"
                }
            }
        }
        return
    }

    Write-Host ""
    Write-Info "Go is not installed on this machine. Go 1.22+ is required to build Dayside."
    Write-Host ""
    Write-Host "The script can try up to three install methods:"
    Write-Host "  1. winget  (built into Windows 10 1809+ and Windows 11)"
    Write-Host "  2. choco   (Chocolatey; common on developer machines)"
    Write-Host "  3. MSI     (direct download from go.dev/dl)"
    Write-Host ""
    if (-not (Ask-YesNo "Begin automatic Go install?")) {
        Abort "Go install skipped. Install Go 1.22+ from https://go.dev/dl/ and re-run .\build.ps1"
    }

    if (Install-Go-Winget) { Write-OK "  Go installed via winget"; return }
    Write-Host ""
    if (Install-Go-Choco)  { Write-OK "  Go installed via choco"; return }
    Write-Host ""
    if (Install-Go-MSI)    { Write-OK "  Go installed via MSI"; return }

    Abort "All three install methods failed or were declined. Install Go 1.22+ manually from https://go.dev/dl/ and re-run this script."
}

# ---------- install chain: WebView2 ----------

function Check-WebView2 {
    $paths = @(
        "HKLM:\SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}",
        "HKLM:\SOFTWARE\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}",
        "HKCU:\SOFTWARE\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}"
    )
    foreach ($p in $paths) {
        if (Test-Path $p) {
            $v = (Get-ItemProperty $p -Name pv -ErrorAction SilentlyContinue).pv
            if ($v) { return $v }
        }
    }
    return $null
}

function Install-WebView2-Winget {
    if (-not (Has-Command winget)) { return $false }
    if (-not (Ask-YesNo "Try: winget install Microsoft.EdgeWebView2Runtime --silent")) { return $false }
    Write-Step "Running winget..."
    try {
        & winget install --id Microsoft.EdgeWebView2Runtime --silent --accept-package-agreements --accept-source-agreements
        if ($LASTEXITCODE -eq 0 -and (Check-WebView2)) { return $true }
    } catch {
        Write-Fail "  winget failed: $_"
    }
    return $false
}

function Install-WebView2-Bootstrapper {
    Write-Host ""
    Write-Host "  This will download the Microsoft Evergreen Bootstrapper and run it."
    Write-Host "  Download size: ~2 MB (it fetches the actual runtime online)."
    if (-not (Ask-YesNo "Proceed with bootstrapper download and install?")) { return $false }

    $url  = "https://go.microsoft.com/fwlink/p/?LinkId=2124703"
    $exe  = Join-Path $env:TEMP "MicrosoftEdgeWebview2Setup.exe"

    Write-Step "Downloading bootstrapper..."
    try {
        Invoke-WebRequest -Uri $url -OutFile $exe -UseBasicParsing
    } catch {
        Write-Fail "  Download failed: $_"
        return $false
    }

    Write-Step "Running bootstrapper (watch for the UAC prompt)..."
    try {
        Start-Process -FilePath $exe -ArgumentList "/silent /install" -Wait
    } catch {
        Write-Fail "  Bootstrapper failed: $_"
        return $false
    }

    if (Check-WebView2) { return $true }
    Write-Fail "  Bootstrapper finished but WebView2 still not detected."
    return $false
}

function Ensure-WebView2 {
    $v = Check-WebView2
    if ($v) {
        Write-OK "  WebView2: present (version $v)"
        return
    }

    Write-Host ""
    Write-Info "WebView2 runtime is not installed."
    Write-Host "  WebView2 is what renders Dayside's UI. Without it, the built"
    Write-Host "  exe will open a blank window. It's pre-installed on Windows 11 and"
    Write-Host "  most updated Windows 10 machines, but not this one."
    Write-Host ""
    if (-not (Ask-YesNo "Begin WebView2 install?")) {
        Write-Fail "  Skipped. The built exe will open a blank window until WebView2 is installed."
        Write-Host "  Manual download: https://developer.microsoft.com/microsoft-edge/webview2/"
        return
    }

    if (Install-WebView2-Winget)         { Write-OK "  WebView2 installed via winget"; return }
    Write-Host ""
    if (Install-WebView2-Bootstrapper)   { Write-OK "  WebView2 installed via bootstrapper"; return }

    Write-Fail "  WebView2 install failed. Dayside will still build but won't render."
    Write-Host "  Manual download: https://developer.microsoft.com/microsoft-edge/webview2/"
}

# ---------- install: Wails CLI ----------

function Ensure-Wails {
    if (Has-Command wails) {
        Write-OK "  Wails CLI: present"
        return
    }
    Write-Host ""
    Write-Info "Wails CLI is not installed."
    Write-Host "  Wails is the framework that bundles the Go backend and HTML UI"
    Write-Host "  into a single .exe. Install is fast (~30 seconds) and local to"
    Write-Host "  this user; no admin required."
    Write-Host ""
    if (-not (Ask-YesNo "Install Wails CLI via: go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0")) {
        Abort "Wails CLI is required to build. Install manually with the command above and re-run."
    }
    Write-Step "Running go install..."
    & go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0
    if ($LASTEXITCODE -ne 0) {
        Abort "Wails CLI install failed. Check that `$env:GOPATH\bin is on your PATH."
    }
    Refresh-Path
    if (-not (Has-Command wails)) {
        Abort "Wails installed but not found on PATH. Open a new PowerShell and re-run this script."
    }
    Write-OK "  Wails CLI installed"
}

# ---------- main ----------

Write-Host ""
Write-Info "=== Dayside - Windows Build ==="
Write-Host ""
Write-Host "This script checks for required tools and asks permission before"
Write-Host "installing anything. You may be prompted up to 4 times:"
Write-Host "  1. To install Go          (if missing)"
Write-Host "  2. To install Wails CLI   (always required; uses Go)"
Write-Host "  3. To install WebView2    (if missing; needed to render the UI)"
Write-Host "  4. UAC prompts            (if any install needs admin)"
Write-Host ""
Read-Host "Press Enter to begin, or Ctrl-C to abort"

Write-Host ""
Write-Step "[1/5] Go"
Ensure-Go

Write-Host ""
Write-Step "[2/5] Wails CLI"
Ensure-Wails

Write-Host ""
Write-Step "[3/5] WebView2 runtime"
Ensure-WebView2

Write-Host ""
Write-Step "[4/5] Fetching Go dependencies..."
& go mod tidy
if ($LASTEXITCODE -ne 0) { Abort "go mod tidy failed." }
Write-OK "  Done"

Write-Host ""
Write-Step "[5/5] Building Dayside.exe..."
& wails build -clean -platform windows/amd64
if ($LASTEXITCODE -ne 0) { Abort "Wails build failed." }

$exePath = "build\bin\Dayside.exe"
if (Test-Path $exePath) {
    $full = (Resolve-Path $exePath).Path
    $size = [Math]::Round((Get-Item $full).Length / 1MB, 1)
    Write-Host ""
    Write-OK "=== BUILD SUCCEEDED ==="
    Write-OK "  Output: $full"
    Write-OK "  Size:   $size MB"
    Write-Host ""
    Write-Info "Double-click Dayside.exe to run."
    Write-Host ""
} else {
    Abort "Build reported success but exe not found at $exePath"
}
