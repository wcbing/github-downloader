package config

import (
	"encoding/json"
	"log"
	"os"
)

type GithubRepo struct {
	Repo     string   `json:"repo"`
	FileList []string `json:"file_list"`
}

// a github repo info format:
//	"name": {
//		"repo": github repo, user/repo
//		"file_list": file name list, contains the valuables
//	}

// 读取 github.json
func ReadRepo() (repoList map[string]GithubRepo) {
	if content, err := os.ReadFile("data/github.json"); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(content, &repoList); err != nil {
		log.Fatal("data/github.json :", err)
	}
	return
}

// a github version info format:
// "name": "version"

// 读取 github-local.json
func ReadVersion() (versionList map[string]string) {
	if content, err := os.ReadFile("data/github-local.json"); err != nil {
		log.Print(err)
		versionList = make(map[string]string)
	} else if len(content) == 0 || string(content) == "null" {
		versionList = make(map[string]string)
	} else if err := json.Unmarshal(content, &versionList); err != nil {
		log.Fatal("data/github-local.json :", err)
	}
	return
}

// 保存到 github-local.json
func SaveVersion(versionList map[string]string) {
	if content, err := json.MarshalIndent(versionList, "", "    "); err != nil {
		log.Fatal(err)
	} else if err := os.WriteFile("data/github-local.json", content, 0644); err != nil {
		log.Fatal("data/github-local.json", err)
	}
}
