package webUI

import (
	"ConfigServer/webUI/WebAPI"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	//_, err := fmt.Fprintf(w, "Hello from server on port 8080! You've requested: %s\n", r.URL.Path)
	//if err != nil {
	//	return
	//}

	reqPath := r.URL.Path
	query := r.URL.Query()

	if strings.HasPrefix(reqPath, "/api/") {
		WebAPI.Handler(w, r, query)
		return
	}

	// 设置静态文件的根目录
	staticDir := "./webUI/static/"

	// 构造文件的完整路径
	path := filepath.Join(staticDir, reqPath)

	// 返回本地文件
	http.ServeFile(w, r, path)
}

func StartServer(port string, handlerFunc http.HandlerFunc) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerFunc)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Starting server at port %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server on port %s failed to start: %v\n", port, err)
	}
}
