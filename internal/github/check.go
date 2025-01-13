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
func replaceFileName(latestVersion, templates string) (fileName string) {
	tmpname := strings.ReplaceAll(templates, "{stripped_version}", latestVersion[1:])
	fileName = strings.ReplaceAll(tmpname, "{version_tag}", latestVersion)
	return fileName
}

func Check(name string, repo config.GithubRepo, localVersion string) (versionTag string) {
	wg := sync.WaitGroup{}
	thread := config.Thread
	sem := make(chan struct{}, thread) // Semaphore

	repoUrl := config.Proxy + "https://github.com/" + repo.Repo
	versionTag = LatestVersionTag(repoUrl)
	log.Printf("%s = %s\n", name, versionTag)
	// 判断是否需要更新
	if localVersion != versionTag {
		releasesDownloadUrl := repoUrl + "/releases/download"
		// 确定本地文件目录并确保目录存在
		var fileDir string
		if config.Config["recursive"] {
			fileDir = filepath.Join(config.OutputDir, releasesDownloadUrl, versionTag)
		} else {
			fileDir = filepath.Join(config.OutputDir, strings.ReplaceAll(repo.Repo, "/", "__"))
		}
		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			if err := os.MkdirAll(fileDir, 0755); err != nil {
				log.Print(err)
			}
		}
		// 下载新版本文件
		for _, templates := range repo.FileList {
			fileName := replaceFileName(versionTag, templates)
			fileUrl := fmt.Sprintf("%s/%s/%s", releasesDownloadUrl, versionTag, fileName)
			filePath := filepath.Join(fileDir, fileName)
			wg.Add(1)
			go func(fileUrl, filePath string) {
				defer wg.Done()
				sem <- struct{}{}        // 获取许可
				defer func() { <-sem }() // 释放许可
				Download(fileUrl, filePath)
			}(fileUrl, filePath)
		}
		// 判断是否是新添加应用
		if localVersion == "" {
			fmt.Printf("AddNew: %s (%s)\n", name, versionTag)
		} else {
			fmt.Printf("Update: %s (%s -> %s)\n", name, localVersion, versionTag)
			// 删除旧版本文件
			if !config.Config["dry-run"] && !config.Config["recursive"] {
				// 对非 "recursive" 的，依次删除旧版本文件
				for _, templates := range repo.FileList {
					// 如果模板中不包含 {stripped_version} 或 {version_tag}，则不删除
					if !strings.Contains(templates, "{stripped_version}") && !strings.Contains(templates, "{version_tag}") {
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
	}
	wg.Wait()
	// 返回版本号
	return versionTag
}
