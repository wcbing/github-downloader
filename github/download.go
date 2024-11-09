package github

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/wcbing/github-downloader/config"
)

// 下载文件
// url: 下载地址
// filePath: 保存路径
func Download(url, filePath string) {
	// 检查是否为 dry-run 模式
	if config.Config["dry-run"] {
		resp, err := http.Head(url)
		if err != nil {
			log.Print(err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Printf("can't download %s because receive %s", url, resp.Status)
		} else {
			log.Printf("dry-run download: %s", url)
		}
		return
	}
	// 下载文件
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("can't download %s because receive %s", url, resp.Status)
		return
	}
	log.Printf("downloading: %s", url)
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return
	}
	// 保存文件
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		log.Print(err)
	}
}
