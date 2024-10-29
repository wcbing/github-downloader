package github

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wcbing/github-downloader/config"
)

// 拼接得到文件名
func replaceFileName(latestVersion, templates string) (fileName string) {
	tmpname := strings.ReplaceAll(templates, "{stripped_version}", latestVersion[1:])
	fileName = strings.ReplaceAll(tmpname, "{version_tag}", latestVersion)
	return fileName
}

func Check(name string, repo config.GithubRepo, localVersion string) (versionTag string) {
	wg := sync.WaitGroup{}
	repoUrl := config.Proxy + "https://github.com/" + repo.Repo
	versionTag = LatestVersionTag(repoUrl)
	releasesDownloadUrl := repoUrl + "/releases/download"
	fmt.Printf("%s: %s\n", name, versionTag)
	// 判断是否需要更新
	if localVersion == "" {
		fmt.Printf("└  Add: %s\n", versionTag)
		for _, templates := range repo.FileList {
			fileName := replaceFileName(versionTag, templates)
			fileUrl := fmt.Sprintf("%s/%s/%s", releasesDownloadUrl, versionTag, fileName)
			wg.Add(1)
			go func() {
				defer wg.Done()
				// 下载
				Download(fileUrl, repo.Repo)
			}()
		}
	} else if localVersion != versionTag {
		fmt.Printf("└  update: %s -> %s\n", localVersion, versionTag)
		for _, templates := range repo.FileList {
			fileName := replaceFileName(versionTag, templates)
			fileUrl := fmt.Sprintf("%s/%s/%s", releasesDownloadUrl, versionTag, fileName)
			oldFilePath := filepath.Join("releases", repo.Repo, replaceFileName(localVersion, templates))
			wg.Add(1)
			go func() {
				defer wg.Done()
				// 下载并删除旧版本文件
				Download(fileUrl, repo.Repo)
				if !config.Config["recursive"] && !config.Config["dry-run"] {
					os.Remove(oldFilePath)
				}
			}()
		}
		// 删除旧版本目录
		if config.Config["recursive"] && !config.Config["dry-run"] {
			localFileDir := fmt.Sprintf("%s/%s", releasesDownloadUrl, localVersion)
			os.RemoveAll(localFileDir)
		}
	}
	wg.Wait()
	// 返回版本号
	return versionTag
}
