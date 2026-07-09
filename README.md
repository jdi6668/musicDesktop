# 夏夜草语 · Windows 桌面播放器

用 Go 语言编写的 Windows 桌面音乐播放器。启动本地 HTTP 服务器加载 `index.html` 及相关资源，并以全屏（Kiosk）模式打开浏览器播放。

## 项目结构

```
musicDesk/
├── main.go              # Go 主程序
├── go.mod               # Go 模块文件
├── index.html           # 播放器页面
├── heisemaoyi.mp3       # 音频文件
├── heisemaoyi.lrc       # 歌词文件
└── favicon.ico          # 图标
```

## 编译与运行教程

### 1. 安装 Go

从官网下载并安装 Go：https://go.dev/dl/

安装完成后，打开命令提示符（CMD）或 PowerShell，验证安装：

```bash
go version
```

应输出类似 `go version go1.21.x windows/amd64`。

### 2. 编译

在项目目录下执行：

```bash
cd C:\Users\87411\Desktop\musicDesk
go build -o musicdesk.exe main.go
```

编译完成后会生成 `musicdesk.exe`。

### 3. 运行

双击 `musicdesk.exe`，或在命令行中执行：

```bash
.\musicdesk.exe
```

程序会：
1. 自动查找一个可用的本地端口
2. 启动 HTTP 服务器托管当前目录下的所有文件
3. 以全屏（Kiosk）模式打开浏览器加载播放页面

### 4. 退出

- **Kiosk 模式**：按 `Alt + F4` 或 `Ctrl + W` 关闭浏览器窗口
- **服务器**：在命令行窗口按 `Ctrl + C` 终止服务器进程

## 跨平台编译

如果需要在其他平台编译 Windows 版本：

```bash
# 在 Linux/macOS 上交叉编译 Windows 版本
GOOS=windows GOARCH=amd64 go build -o musicdesk.exe main.go
```

## 隐藏命令行窗口（可选）

如果希望运行时不显示命令行黑窗口，可以编译为窗口模式：

```bash
go build -ldflags "-H windowsgui" -o musicdesk.exe main.go
```

> 注意：使用 `-H windowsgui` 后将看不到日志输出。如果需要调试，建议先用普通方式编译运行。

## 浏览器说明

- 程序优先使用 **Microsoft Edge** 的 Kiosk 模式
- 如果未安装 Edge，则尝试 **Chrome**
- 如果两者都未找到，则用系统默认浏览器打开（非全屏）

## 自定义

- 替换 `heisemaoyi.mp3` 和 `heisemaoyi.lrc` 为你自己的音乐和歌词
- 修改 `index.html` 中的 `src="heisemaoyi.mp3"` 为对应文件名
