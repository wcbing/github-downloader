package github

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wcbing/github-downloader/internal/config"
)

// 拼接得到文件名
func replaceFileName(releasesTag, file_template string) (fileName string) {
	version := stripVerison(releasesTag)
	tmpname := strings.ReplaceAll(file_template, "{version}", version)
	fileName = strings.ReplaceAll(tmpname, "{releases_tag}", releasesTag)
	return
}

// Check 检查仓库更新并下载新版本文件
// name: 应用名称
// repo: 仓库配置
// localVersion: 本地记录的版本
// 返回最新的 Releases Tag
func Check(name string, repo config.GithubRepo, localVersion string) (releasesTag string) {
	wg := sync.WaitGroup{}
	thread := config.Thread
	sem := make(chan struct{}, thread) // Semaphore

	repoUrl := config.Proxy + "https://github.com/" + repo.Repo
	releasesTag = latestReleasesTag(repoUrl)
	log.Printf("%s = %s\n", name, releasesTag)
	// 判断是否需要更新
	if localVersion == releasesTag || releasesTag == "" {
		return
	}

	releasesDownloadUrl := repoUrl + "/releases/download"

	// 确定本地文件目录并确保目录存在
	var fileDir string
	if config.Config["recursive"] {
		fileDir = filepath.Join(config.OutputDir, releasesDownloadUrl, releasesTag)
	} else {
		fileDir = filepath.Join(config.OutputDir, name)
	}
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		if err := os.MkdirAll(fileDir, 0755); err != nil {
			log.Print(err)
		}
	}

	// 下载新版本文件
	for _, templates := range repo.FileList {
		fileName := replaceFileName(releasesTag, templates)
		fileUrl := fmt.Sprintf("%s/%s/%s", releasesDownloadUrl, releasesTag, fileName)
		filePath := filepath.Join(fileDir, fileName)

		wg.Add(1)
		go func(fileUrl, filePath string) {
			defer wg.Done()
			sem <- struct{}{}        // 获取许可
			defer func() { <-sem }() // 释放许可
			Download(fileUrl, filePath)
		}(fileUrl, filePath)
	}
	wg.Wait()

	// 判断是否是新添加应用
	if localVersion == "" {
		fmt.Printf("AddNew: %s (%s)\n", name, releasesTag)
	} else {
		fmt.Printf("Update: %s (%s -> %s)\n", name, localVersion, releasesTag)
		// 删除旧版本文件
		if !config.Config["dry-run"] && !config.Config["recursive"] {
			// 对非 "recursive" 的，依次删除旧版本文件
			for _, templates := range repo.FileList {
				// 如果模板中不包含 {version} 或 {releases_tag}，则不删除
				if !strings.Contains(templates, "{version}") && !strings.Contains(templates, "{releases_tag}") {
					continue
				}
				oldFilePath := filepath.Join(fileDir, replaceFileName(localVersion, templates))
				os.Remove(oldFilePath)
			}
		} else if !config.Config["dry-run"] && config.Config["recursive"] {
			// 对 "recursive" 的，直接删除旧版本目录
			oldFileDir := filepath.Join(config.OutputDir, releasesDownloadUrl, localVersion)
			os.RemoveAll(oldFileDir)
		}
	}
	return
}
