# Summer Night Grass Whispers · Windows Desktop Music Player

A Windows desktop music player written in Go. Embeds the **WebView2 engine** (not a system browser), starts a local HTTP server to load `index.html` and related assets, and plays in fullscreen.

[中文文档](README.md)

![Preview](show.png)

## Project Structure

```
musicDesk/
├── main.go              # Go main program (WebView2 fullscreen window)
├── go.mod               # Go module file
├── go.sum               # Dependency checksum
├── index.html           # Player page
├── heisemaoyi.mp3       # Audio file
├── heisemaoyi.lrc       # Lyrics file
├── favicon.ico          # Icon
├── favicon.png          # Icon source (PNG)
├── icon.syso            # Compiled Windows resource (icon embedded in exe)
├── build.sh             # Build script
└── show.png             # Screenshot
```

## Requirements

- **Windows 10+**
- **Go 1.16+**
- **Microsoft Edge WebView2 Runtime** (usually pre-installed on Windows 10+. If not, download from [Microsoft](https://developer.microsoft.com/en-us/microsoft-edge/webview2/))

## Build & Run Guide

### 1. Install Go

Download and install from: https://go.dev/dl/

Verify:

```bash
go version
```

### 2. Install Dependencies

```bash
cd C:\Users\87411\Desktop\musicDesk
go mod tidy
```

### 3. Build

```bash
go build -o musicdesk.exe .
```

> **Note**: Use `go build .` (not `go build main.go`) so the `icon.syso` resource file is automatically linked into the exe, giving it a custom icon.

### 4. Run

Double-click `musicdesk.exe`, or run from command line:

```bash
.\musicdesk.exe
```

The program will:
1. Find an available local port automatically
2. Start an HTTP server serving embedded assets
3. Create a **WebView2 fullscreen window** to display the player

### 5. Exit

Press `ESC` to quit the application.

## Embedding Resources

All resources (`index.html`, `heisemaoyi.mp3`, `heisemaoyi.lrc`, `favicon.ico`) are embedded into the single `musicdesk.exe` using Go's `//go:embed` directive. No external files are needed at runtime.

## Hiding the Console Window (Optional)

Build without the black console window:

```bash
go build -ldflags "-H windowsgui" -o musicdesk.exe .
```

> For debugging, use the normal build to see log output.

## Setting the Exe Icon

The icon is set via `icon.syso`, which is generated from `favicon.ico`:

```bash
# Install rsrc tool
go install github.com/akavel/rsrc@latest

# Generate icon.syso from favicon.ico
rsrc -ico favicon.ico -o icon.syso
```

If you only have a PNG, convert it to ICO first using the included tool:

```bash
go run png2ico/main.go
```

## Technical Details

- Uses `github.com/yuaotian/go-win-webview2` — based on **Microsoft Edge WebView2** engine
- No CGO required — pure Go calling Windows API
- Does not depend on system browser — embedded WebView2 rendering engine
- Audio autoplay enabled via `WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS` environment variable
- ESC key registered as a global hotkey to exit fullscreen

## Customization

- Replace `heisemaoyi.mp3` and `heisemaoyi.lrc` with your own music and lyrics
- Update `src="heisemaoyi.mp3"` in `index.html` to match your filename
