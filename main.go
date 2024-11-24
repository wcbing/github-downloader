package main

import (
	"sync"

	"github.com/wcbing/github-downloader/config"
	"github.com/wcbing/github-downloader/github"
)

func main() {
	config.ReadArgs()

	var repoList = config.ReadRepo()
	var versionList = config.ReadVersion()

	wg := sync.WaitGroup{}
	for name, repo := range repoList {
		wg.Add(1)
		go func() {
			defer wg.Done()
			versionList[name] = github.Check(name, repo, versionList[name])
		}()
	}
	wg.Wait()

	config.SaveVersion(versionList)
}
