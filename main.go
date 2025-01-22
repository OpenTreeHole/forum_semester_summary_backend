package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	mode := os.Getenv("MODE")

	// 如果是测试模式，触发下载文件
	if mode == "test" {
		fmt.Println("Running in test mode: downloading files...")
		downloadFiles()
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-Consumer-Username")
		if userID == "" {
			// 未授权，重定向到登录页面
			authURL := os.Getenv("AUTH_URL")               // 从环境变量获取 AUTH_URL
			http.Redirect(w, r, authURL, http.StatusFound) // 302 重定向
			return
		}

		filePath := filepath.Join("resource", fmt.Sprintf("%s.html", userID))
		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			log.Printf("File not found for user: %s\n", userID)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "text/html")
		io.Copy(w, file)
		log.Printf("Served file: %s for user: %s\n", filePath, userID)
	})

	port := "8080"
	log.Printf("Server is running on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
