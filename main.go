package main

import (
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/yuaotian/go-win-webview2"
)

//go:embed index.html heisemaoyi.mp3 heisemaoyi.lrc favicon.ico
var assets embed.FS

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
	// 获取可用端口
	port, err := findFreePort()
	if err != nil {
		log.Fatalf("无法获取可用端口: %v", err)
	}
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	url := fmt.Sprintf("http://%s", addr)

	// 从嵌入的资源提供静态文件
	http.Handle("/", http.FileServer(http.FS(assets)))

	// 启动 HTTP 服务器（在 goroutine 中）
	go func() {
		log.Printf("服务器已启动: %s", url)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("服务器错误: %v", err)
		}
	}()

	// 设置 WebView2 浏览器启动参数，允许自动播放
	os.Setenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", "--autoplay-policy=no-user-gesture-required")

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

	// 注册 ESC 全局热键退出（modifier=0, VK_ESCAPE=0x1B）
	// WebView2 全屏模式会拦截 ESC，JS 监听不到，必须用全局热键
	w.RegisterHotKey(0, 0x1B, func() {
		log.Println("退出应用...")
		w.Terminate()
	})

	// 页面加载完成后注入 ESC 监听（非全屏时的后备方案）
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
