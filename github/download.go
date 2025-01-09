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
	// 使用流式复制减少内存占用
	out, err := os.Create(filePath)
	if err != nil {
		log.Print(err)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Print(err)
	}
}
