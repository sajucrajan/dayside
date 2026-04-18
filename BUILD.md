# Building Dayside from Source

## Prerequisites

- **Go 1.22+** — [go.dev/dl](https://go.dev/dl)
- **Wails v2** — installed automatically by the build scripts

## Windows

Open PowerShell in the project root and run:

```powershell
.\build.ps1
```

The script will check for Go, install Wails CLI if missing, check for WebView2, and produce `build\bin\Dayside.exe`.

## macOS

Open Terminal in the project root and run:

```bash
./build.sh
```

The script will check for Xcode CLT, Go, and Wails CLI, then produce `build/bin/Dayside.app`.

## Manual build

If you prefer to run the steps yourself:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0
go mod tidy
wails build -clean -platform windows/amd64   # Windows
wails build -clean -platform darwin/universal # macOS
```

## Output

| Platform | Path |
|----------|------|
| Windows  | `build/bin/Dayside.exe` |
| macOS    | `build/bin/Dayside.app` |
