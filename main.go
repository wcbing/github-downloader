package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wcbing/github-downloader/cmd"
	"github.com/wcbing/github-downloader/config"
	"github.com/wcbing/github-downloader/github"
)

// 拼接得到文件名
func replaceFileName(versionTag, templates string) (fileName string) {
	tmpname := strings.ReplaceAll(templates, "{stripped_version}", versionTag[1:])
	fileName = strings.ReplaceAll(tmpname, "{version_tag}", versionTag)
	return fileName
}

func main() {
	cmd.ReadArgs()

	var repoList = config.ReadRepo()
	var versionList = config.ReadVersion()

	wg := sync.WaitGroup{}
	for name, repo := range repoList {
		wg.Add(1)
		go func() {
			defer wg.Done()
			repoUrl := "https://github.com/" + repo.Repo
			versionTag := github.LatestVersionTag(repoUrl)
			releasesDownloadUrl := repoUrl + "/releases/download"
			fmt.Printf("%s: %s\n", name, versionTag)
			// 判断是否需要更新
			if versionList[name] == "" {
				fmt.Printf("└  Add: %s\n", versionTag)
				for _, templates := range repo.FileList {
					fileName := replaceFileName(versionTag, templates)
					fileUrl := fmt.Sprintf("%s/%s/%s", releasesDownloadUrl, versionTag, fileName)
					wg.Add(1)
					go func() {
						defer wg.Done()
						// 下载
						github.Download(fileUrl, repo.Repo)
					}()
				}
			} else if versionList[name] != versionTag {
				fmt.Printf("└  update: %s -> %s\n", versionList[name], versionTag)
				for _, templates := range repo.FileList {
					fileName := replaceFileName(versionTag, templates)
					fileUrl := fmt.Sprintf("%s/%s/%s", releasesDownloadUrl, versionTag, fileName)
					oldFilePath := filepath.Join("releases", repo.Repo, replaceFileName(versionList[name], templates))
					wg.Add(1)
					go func() {
						defer wg.Done()
						// 下载并删除旧版本文件
						github.Download(fileUrl, repo.Repo)
						if !config.Config["recursive"] && !config.Config["dry-run"] {
							os.Remove(oldFilePath)
						}
					}()
				}
				// 删除旧版本目录
				if config.Config["recursive"] && !config.Config["dry-run"] {
					localFileDir := fmt.Sprintf("%s/%s", releasesDownloadUrl, versionList[name])
					os.RemoveAll(localFileDir)
				}
			}
			// 更新版本号
			versionList[name] = versionTag
		}()
	}
	wg.Wait()

	config.SaveVersion(versionList)
}
