package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	envTotalUser := os.Getenv("TOTAL_USER")
	if envTotalUser == "" {
		log.Fatal("TOTAL_USER is required")
	}
	totalUser, err := strconv.Atoi(envTotalUser)
	if err != nil {
		log.Fatalf("Failed to parse TOTAL_USER: %v", err)
	}
	// mode := os.Getenv("MODE")
	// // 如果是测试模式，触发下载文件
	// if mode == "test" {
	// 	fmt.Println("Running in test mode: downloading files...")
	// 	downloadFiles()
	// }

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-Consumer-Username")
		if userID == "" {
			log.Println("User ID not found in header")
			for key, values := range r.Header {
				log.Printf("Header: %s = %v", key, values)
			}
			// 未授权，重定向到登录页面
			authURL := os.Getenv("AUTH_URL")               // 从环境变量获取 AUTH_URL
			http.Redirect(w, r, authURL, http.StatusFound) // 302 重定向
			return
		}
		userIDint, err := strconv.Atoi(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			log.Printf("Invalid user ID: %s\n", userID)
			return
		}
		if userIDint < 1 || userIDint > totalUser {
			http.Error(w, "User ID out of range", http.StatusBadRequest)
			log.Printf("User ID out of range: %s\n", userID)
			return
		}

		// 用户资源目录
		userDir := filepath.Join("resource", userID)
		if _, err := os.Stat(userDir); os.IsNotExist(err) {
			http.Error(w, "User resources not found", http.StatusNotFound)
			log.Printf("Resource directory not found for user: %s\n", userID)
			return
		}

		// 请求路径（去掉前导斜杠）
		requestedPath := r.URL.Path[1:]
		fullPath := filepath.Join(userDir, requestedPath)

		// 如果是目录请求，返回用户的默认 HTML 文件
		if fileInfo, err := os.Stat(fullPath); err == nil && fileInfo.IsDir() {
			serveDefaultHTML(w, r, userDir, userID)
			return
		}

		// 如果文件存在，返回该文件
		if _, err := os.Stat(fullPath); err == nil {
			http.ServeFile(w, r, fullPath)
			log.Printf("Served file: %s for user: %s\n", fullPath, userID)
			return
		}

		// 如果文件不存在或请求路径为空，返回默认 HTML 文件
		if requestedPath == "" || requestedPath == "/" {
			serveDefaultHTML(w, r, userDir, userID)
			return
		}

		// 如果请求的路径无法找到，返回 404
		http.Error(w, "File not found", http.StatusNotFound)
		log.Printf("File not found: %s for user: %s\n", requestedPath, userID)
	})

	port := "8080"
	log.Printf("Server is running on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// serveDefaultHTML 返回默认的 HTML 文件
func serveDefaultHTML(w http.ResponseWriter, r *http.Request, userDir, userID string) {
	htmlFilePath := filepath.Join(userDir, fmt.Sprintf("%s.html", userID))
	if _, err := os.Stat(htmlFilePath); os.IsNotExist(err) {
		http.Error(w, "HTML file not found", http.StatusNotFound)
		log.Printf("Default HTML file not found for user: %s\n", userID)
		return
	}
	http.ServeFile(w, r, htmlFilePath)
	log.Printf("Served default HTML file: %s for user: %s\n", htmlFilePath, userID)
}
