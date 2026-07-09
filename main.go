package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// findFreePort 获取一个可用的本地端口
func findFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// openBrowserKiosk 以全屏(kiosk)模式打开默认浏览器
func openBrowserKiosk(url string) error {
	switch runtime.GOOS {
	case "windows":
		// 优先尝试 Edge，其次 Chrome，最后用默认浏览器
		if path, err := exec.LookPath("msedge"); err == nil {
			return exec.Command(path, "--kiosk", "--autoplay-policy=no-user-gesture-required", url).Start()
		}
		if path, err := exec.LookPath("chrome"); err == nil {
			return exec.Command(path, "--kiosk", "--autoplay-policy=no-user-gesture-required", url).Start()
		}
		// 退而求其次：用默认浏览器打开（非 kiosk）
		return exec.Command("cmd", "/c", "start", "", url).Start()
	case "darwin":
		// macOS: 用 Safari 全屏
		return exec.Command("open", "-a", "Safari", url).Start()
	default:
		// Linux: 尝试 xdg-open
		return exec.Command("xdg-open", url).Start()
	}
}

func main() {
	// 获取可执行文件所在目录，作为静态资源根目录
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("无法获取可执行文件路径: %v", err)
	}
	rootDir := filepath.Dir(exePath)

	// 确认 index.html 存在
	indexPath := filepath.Join(rootDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		log.Fatalf("未找到 index.html，请确保它与可执行文件在同一目录: %s", indexPath)
	}

	// 获取可用端口
	port, err := findFreePort()
	if err != nil {
		log.Fatalf("无法获取可用端口: %v", err)
	}
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	url := fmt.Sprintf("http://%s", addr)

	// 静态文件服务器
	fs := http.FileServer(http.Dir(rootDir))
	http.Handle("/", fs)

	// 启动 HTTP 服务器（在 goroutine 中）
	go func() {
		log.Printf("服务器已启动: %s", url)
		log.Printf("资源目录: %s", rootDir)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("服务器错误: %v", err)
		}
	}()

	// 打开全屏浏览器
	log.Printf("正在打开全屏浏览器...")
	if err := openBrowserKiosk(url); err != nil {
		log.Printf("无法自动打开浏览器，请手动访问: %s", url)
	}

	// 阻塞主线程，保持服务器运行
	select {}
}
