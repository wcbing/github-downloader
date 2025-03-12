package github

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/wcbing/github-downloader/internal/config"
)

// 下载文件
// url: 下载地址
// filePath: 保存路径
func Download(url, filePath string) {
	var method = http.Get
	// 检查是否为 dry-run 模式
	if config.Config["dry-run"] {
		method = http.Head
	}
	// 发出请求
	resp, err := method(url)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("can't download %s because receive %s", url, resp.Status)
		return
	}

	if config.Config["dry-run"] {
		log.Printf("dry-run download: %s", url)
		return
	}
	// 下载文件，使用流式复制减少内存占用
	log.Printf("downloading: %s", url)
	out, err := os.Create(filePath)
	if err != nil {
		log.Print(err)
		return
	}
	defer out.Close()
	if _, err = io.Copy(out, resp.Body); err != nil {
		log.Print(err)
	}
}
