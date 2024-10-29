package github

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wcbing/github-downloader/config"
)

// 下载文件
// url: 下载地址
// repo: "user/repo"
func Download(url, repo string) {
	// 检查是否为 dry-run 模式
	if config.Config["dry-run"] {
		resp, err := http.Head(url)
		if err != nil {
			log.Printf("dry-run head request failed: %s", url)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Printf("download failed: %s", url)
		} else {
			log.Printf("dry-run download: %s", url)
		}
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("download failed: %s", url)
		return
	}
	log.Printf("downloading: %s", url)
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return
	}
	// 创建目录并写入文件
	var fileDir, filePath string
	if config.Config["recursive"] {
		fileDir = filepath.Dir(url)
		filePath = url
	} else {
		fileDir = filepath.Join("releases", repo)
		filePath = filepath.Join(fileDir, filepath.Base(url))
	}
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		if err := os.MkdirAll(fileDir, 0755); err != nil {
			log.Print(err)
		}
	}
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		log.Print(err)
	}
}
