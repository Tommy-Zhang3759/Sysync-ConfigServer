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

	// set root folder of static files
	staticDir := "./webUI/static/"
	path := filepath.Join(staticDir, reqPath)
	// return local file
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
