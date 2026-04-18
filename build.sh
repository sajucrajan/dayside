#!/usr/bin/env bash
# build.sh
# One-command macOS build for Dayside.
# Prompts for permission before installing anything.
# Run from project root: ./build.sh

set -e

# ---------- helpers ----------

C_CYAN='\033[36m'
C_YELLOW='\033[33m'
C_GREEN='\033[32m'
C_RED='\033[31m'
C_RESET='\033[0m'

info()  { printf "${C_CYAN}%s${C_RESET}\n" "$1"; }
step()  { printf "${C_YELLOW}%s${C_RESET}\n" "$1"; }
ok()    { printf "${C_GREEN}%s${C_RESET}\n" "$1"; }
fail()  { printf "${C_RED}%s${C_RESET}\n" "$1"; }

ask_yesno() {
    local prompt="$1"
    local response
    while true; do
        printf "%s [y/N]: " "$prompt"
        read -r response
        case "$response" in
            ''|[nN]*) return 1 ;;
            [yY]*) return 0 ;;
        esac
    done
}

abort() {
    echo ""
    fail "=== BUILD ABORTED ==="
    fail "$1"
    echo ""
    exit 1
}

has_cmd() {
    command -v "$1" >/dev/null 2>&1
}

refresh_path() {
    if has_cmd go; then
        local gopath
        gopath=$(go env GOPATH 2>/dev/null || true)
        if [ -n "$gopath" ] && [ -d "$gopath/bin" ]; then
            case ":$PATH:" in
                *":$gopath/bin:"*) : ;;
                *) export PATH="$gopath/bin:$PATH" ;;
            esac
        fi
    fi
    # Homebrew paths (Apple Silicon + Intel)
    for hb in /opt/homebrew/bin /usr/local/bin; do
        if [ -d "$hb" ]; then
            case ":$PATH:" in
                *":$hb:"*) : ;;
                *) export PATH="$hb:$PATH" ;;
            esac
        fi
    done
}

# ---------- Xcode Command-Line Tools ----------

ensure_xcode_clt() {
    if xcode-select -p >/dev/null 2>&1; then
        ok "  Xcode CLT: found at $(xcode-select -p)"
        return
    fi
    echo ""
    info "Xcode Command-Line Tools are not installed."
    echo "  These are required to compile the CGO portion of the app"
    echo "  (used for macOS window enumeration via CoreGraphics)."
    echo ""
    echo "  Running 'xcode-select --install' will pop Apple's system installer"
    echo "  dialog. You'll need to click Install in that dialog and wait for"
    echo "  the download (~500 MB-1 GB, 5-10 minutes)."
    echo ""
    if ! ask_yesno "Trigger the Xcode CLT installer now?"; then
        abort "Xcode CLT install skipped. Run 'xcode-select --install' manually and re-run this script."
    fi
    xcode-select --install 2>/dev/null || true
    step "Waiting for Xcode CLT install to complete (polls every 10s)..."
    echo "  Click Install in the popup dialog. Ctrl-C to abort this wait."
    while ! xcode-select -p >/dev/null 2>&1; do
        sleep 10
        printf "."
    done
    echo ""
    ok "  Xcode CLT installed"
}

# ---------- Homebrew ----------

install_homebrew() {
    echo ""
    info "Homebrew is not installed. It's the easiest way to install Go on macOS."
    echo ""
    echo "  Install command (official installer from brew.sh):"
    echo "    /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
    echo ""
    echo "  This downloads ~500 MB, asks for your admin password once, and"
    echo "  modifies your shell profile (~/.zprofile) to add brew to PATH."
    echo ""
    if ! ask_yesno "Install Homebrew now?"; then
        return 1
    fi
    step "Running Homebrew installer..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    refresh_path
    has_cmd brew
}

# ---------- install chain: Go ----------

install_go_brew() {
    if ! has_cmd brew; then
        echo "  Homebrew not available."
        if ask_yesno "Install Homebrew first so we can use it for Go?"; then
            if ! install_homebrew; then
                echo "  Homebrew install failed or was declined."
                return 1
            fi
        else
            return 1
        fi
    fi
    if ! ask_yesno "Try: brew install go"; then return 1; fi
    step "Running brew install go..."
    if brew install go; then
        refresh_path
        has_cmd go && return 0
    fi
    return 1
}

install_go_pkg() {
    echo ""
    echo "  This will download the official Go .pkg from https://go.dev/dl/"
    echo "  and run the installer, which asks for your admin password."
    echo "  Download size: ~120 MB."
    if ! ask_yesno "Proceed with .pkg download and install?"; then return 1; fi

    local go_version="1.23.4"
    local arch
    if [ "$(uname -m)" = "arm64" ]; then
        arch="arm64"
    else
        arch="amd64"
    fi
    local pkg_url="https://go.dev/dl/go${go_version}.darwin-${arch}.pkg"
    local pkg_path="/tmp/go${go_version}.darwin-${arch}.pkg"

    step "Downloading $pkg_url ..."
    if ! curl -fL -o "$pkg_path" "$pkg_url"; then
        fail "  Download failed."
        return 1
    fi

    step "Launching installer (enter your admin password when prompted)..."
    if ! sudo installer -pkg "$pkg_path" -target /; then
        fail "  Installer failed."
        return 1
    fi

    export PATH="/usr/local/go/bin:$PATH"
    refresh_path
    has_cmd go && return 0
    fail "  Installer finished but Go is not on PATH. Open a new Terminal and re-run."
    return 1
}

ensure_go() {
    if has_cmd go; then
        local go_ver
        go_ver="$(go version)"
        ok "  Go: $go_ver"

        # Parse "go version go1.22.0 darwin/arm64" and verify >= 1.22
        if [[ "$go_ver" =~ go([0-9]+)\.([0-9]+) ]]; then
            local major="${BASH_REMATCH[1]}"
            local minor="${BASH_REMATCH[2]}"
            if [ "$major" -lt 1 ] || { [ "$major" -eq 1 ] && [ "$minor" -lt 22 ]; }; then
                echo ""
                fail "  Your Go version is too old."
                echo "  Dayside requires Go 1.22 or newer."
                echo "  The Wails CLI install will fail on older Go due to an"
                echo "  incompatibility in golang.org/x/tools v0.17.0."
                echo ""
                if ask_yesno "Attempt to upgrade Go automatically?"; then
                    if install_go_brew; then ok "  Go upgraded via Homebrew"; return; fi
                    if install_go_pkg;  then ok "  Go upgraded via .pkg"; return; fi
                    abort "Go upgrade failed. Install Go 1.22+ manually from https://go.dev/dl/"
                else
                    abort "Go 1.22+ is required. Install it and re-run ./build.sh"
                fi
            fi
        fi
        return
    fi
    echo ""
    info "Go is not installed on this machine. Go 1.22+ is required to build Dayside."
    echo ""
    echo "The script can try two install methods:"
    echo "  1. brew install go  (fastest if Homebrew is present)"
    echo "  2. .pkg from go.dev (direct download)"
    echo ""
    if ! ask_yesno "Begin automatic Go install?"; then
        abort "Go install skipped. Install Go 1.22+ from https://go.dev/dl/ and re-run ./build.sh"
    fi

    if install_go_brew; then ok "  Go installed via Homebrew"; return; fi
    echo ""
    if install_go_pkg;  then ok "  Go installed via .pkg"; return; fi

    abort "All install methods failed or were declined. Install Go 1.22+ manually from https://go.dev/dl/ and re-run this script."
}

# ---------- install: Wails CLI ----------

ensure_wails() {
    if has_cmd wails; then
        ok "  Wails CLI: present"
        return
    fi
    echo ""
    info "Wails CLI is not installed."
    echo "  Wails bundles the Go backend and HTML UI into a single .app."
    echo "  Install is fast (~30 seconds) and local to this user."
    echo ""
    if ! ask_yesno "Install Wails via: go install github.com/wailsapp/wails/v2/cmd/wails@v2.9.1"; then
        abort "Wails CLI is required to build. Install manually with the command above and re-run."
    fi
    step "Running go install..."
    go install github.com/wailsapp/wails/v2/cmd/wails@v2.9.1
    refresh_path
    if ! has_cmd wails; then
        abort "Wails installed but not found on PATH. Open a new Terminal and re-run this script."
    fi
    ok "  Wails CLI installed"
}

# ---------- main ----------

echo ""
info "=== Dayside - macOS Build ==="
echo ""
echo "This script checks for required tools and asks permission before"
echo "installing anything. You may be prompted up to 4 times:"
echo "  1. To install Xcode CLT   (if missing; system dialog popup)"
echo "  2. To install Homebrew    (if missing; optional, for Go install)"
echo "  3. To install Go          (if missing)"
echo "  4. To install Wails CLI   (always required; uses Go)"
echo ""
printf "Press Enter to begin, or Ctrl-C to abort: "
read -r _

echo ""
step "[1/5] Xcode Command-Line Tools"
ensure_xcode_clt

echo ""
step "[2/5] Go"
ensure_go

echo ""
step "[3/5] Wails CLI"
ensure_wails

echo ""
step "[4/5] Fetching Go dependencies..."
go mod tidy
ok "  Done"

echo ""
step "[5/5] Building Dayside.app..."
wails build -clean -platform darwin/universal

APP_PATH="build/bin/Dayside.app"
if [ -d "$APP_PATH" ]; then
    SIZE=$(du -sh "$APP_PATH" | awk '{print $1}')
    FULL=$(cd "$(dirname "$APP_PATH")" && pwd)/"$(basename "$APP_PATH")"
    echo ""
    ok "=== BUILD SUCCEEDED ==="
    ok "  Output: $FULL"
    ok "  Size:   $SIZE"
    echo ""
    info "First launch will prompt for automation permissions for each browser."
    info "Grant those, then re-scan. To run: open \"$APP_PATH\""
    echo ""
else
    abort "Build reported success but .app not found at $APP_PATH"
fi
