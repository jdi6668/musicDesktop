package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yuaotian/go-win-webview2"
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

	// 创建 WebView2 全屏窗口
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     false,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:      "夏夜草语 · 助眠音乐",
			Width:      1920,
			Height:     1080,
			Center:     true,
			Frameless:  false,
			Fullscreen: true,
		},
	})
	defer w.Destroy()

	// 绑定退出函数，供 JS 调用
	w.Bind("exitApp", func() {
		log.Println("退出应用...")
		w.Terminate()
	})

	// 页面加载完成后注入 ESC 监听
	w.OnLoadingStateChanged(func(isLoading bool) {
		if !isLoading {
			w.Eval(`
				document.addEventListener('keydown', function(e) {
					if (e.key === 'Escape') {
						exitApp();
					}
				});
			`)
		}
	})

	// 导航到本地服务器
	w.Navigate(url)

	// 运行 WebView 消息循环（阻塞）
	w.Run()
}
