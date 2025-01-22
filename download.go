package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// downloadFiles 下载文件到 resource 目录
// func downloadFiles() {
// 	baseURL := "https://a.com?user_id="
// 	resourceDir := "resource"

// 	// 创建 resource 目录
// 	err := os.MkdirAll(resourceDir, os.ModePerm)
// 	if err != nil {
// 		log.Fatalf("Failed to create resource directory: %v", err)
// 	}

// 	for i := 1; i <= totalUser; i++ {
// 		url := baseURL + strconv.Itoa(i)
// 		resp, err := http.Get(url)
// 		if err != nil {
// 			log.Printf("Failed to download file for user_id=%d: %v\n", i, err)
// 			continue
// 		}

// 		filePath := filepath.Join(resourceDir, fmt.Sprintf("%d.html", i))
// 		file, err := os.Create(filePath)
// 		if err != nil {
// 			log.Printf("Failed to create file: %v\n", err)
// 			resp.Body.Close()
// 			continue
// 		}

// 		_, err = io.Copy(file, resp.Body)
// 		resp.Body.Close()
// 		file.Close()

// 		if err != nil {
// 			log.Printf("Failed to save file for user_id=%d: %v\n", i, err)
// 		} else {
// 			log.Printf("Downloaded file for user_id=%d\n", i)
// 		}
// 	}
// }

// downloadFiles 并行下载文件到 resource 目录
func downloadFiles() {
	envTotalUser := os.Getenv("TOTAL_USER")
	if envTotalUser == "" {
		log.Fatal("TOTAL_USER is required")
	}
	totalUser, err := strconv.Atoi(envTotalUser)
	if err != nil {
		log.Fatalf("Failed to parse TOTAL_USER: %v", err)
	}
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Fatal("BASE_URL is required")
	}
	resourceDir := "resource"

	// 创建 resource 目录
	err = os.MkdirAll(resourceDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create resource directory: %v", err)
	}

	const maxWorkers = 1000           // 最大并发数
	jobs := make(chan int, totalUser) // 任务 channel
	results := make(chan string)      // 结果 channel

	// 工作函数
	worker := func(jobs <-chan int, results chan<- string) {
		for userID := range jobs {
			url := baseURL + strconv.Itoa(userID)
			resp, err := http.Get(url)
			if err != nil {
				results <- fmt.Sprintf("Failed to download file for user_id=%d: %v\n", userID, err)
				continue
			}

			filePath := filepath.Join(resourceDir, fmt.Sprintf("%d.html", userID))
			file, err := os.Create(filePath)
			if err != nil {
				results <- fmt.Sprintf("Failed to create file for user_id=%d: %v\n", userID, err)
				resp.Body.Close()
				continue
			}

			_, err = io.Copy(file, resp.Body)
			resp.Body.Close()
			file.Close()

			if err != nil {
				results <- fmt.Sprintf("Failed to save file for user_id=%d: %v\n", userID, err)
			} else {
				results <- fmt.Sprintf("Downloaded file for user_id=%d\n", userID)
			}
		}
	}

	// 启动工作协程
	for i := 0; i < maxWorkers; i++ {
		go worker(jobs, results)
	}

	// 分发任务
	go func() {
		for i := 1; i <= totalUser; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// 收集结果
	for i := 1; i <= totalUser; i++ {
		log.Println(<-results)
	}
}
