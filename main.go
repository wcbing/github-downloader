package main

import (
	"sync"

	"github.com/wcbing/github-downloader/internal/config"
	"github.com/wcbing/github-downloader/internal/github"
)

func main() {
	config.ReadArgs()

	var repoList = config.ReadRepo()
	var versionList = config.ReadVersion()

	wg := sync.WaitGroup{}
	for name, repo := range repoList {
		wg.Add(1)
		go func(name string, repo config.GithubRepo) {
			defer wg.Done()
			versionList[name] = github.Check(name, repo, versionList[name])
		}(name, repo)
	}
	wg.Wait()

	config.SaveVersion(versionList)
}
