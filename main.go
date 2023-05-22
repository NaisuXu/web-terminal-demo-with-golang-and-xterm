package main

import (
	"embed"
	"net/http"
	"os/exec"

	"github.com/creack/pty"
	"github.com/olahol/melody"
)

//go:embed index.html node_modules/xterm/css/xterm.css node_modules/xterm/lib/xterm.js
var content embed.FS

func main() {
	c := exec.Command("sh") // 系统默认shell交互程序
	f, err := pty.Start(c)  // pty用于调用系统自带的虚拟终端
	if err != nil {
		panic(err)
	}

	m := melody.New() // melody用于实现WebSocket功能

	go func() { // 处理来自虚拟终端的消息
		for {
			buf := make([]byte, 1024)
			read, err := f.Read(buf)
			if err != nil {
				return
			}
			// fmt.Println("f.Read: ", string(buf[:read]))
			m.Broadcast(buf[:read]) // 将数据发送给网页
		}
	}()

	m.HandleMessage(func(s *melody.Session, msg []byte) { // 处理来自WebSocket的消息
		// fmt.Println("m.HandleMessage: ", string(msg))
		f.Write(msg) // 将消息写到虚拟终端
	})

	http.HandleFunc("/webterminal", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r) // 访问 /webterminal 时将转交给melody处理
	})

	fs := http.FileServer(http.FS(content))
	http.Handle("/", http.StripPrefix("/", fs)) // 设置静态文件服务

	http.ListenAndServe("0.0.0.0:22333", nil) // 启动服务器，访问 http://本机(服务器)IP地址:22333/ 进行测试
}
