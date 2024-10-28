package github

import (
	"net/http"
	"regexp"
)

// 获取最新版本标签
// repo: "user/repo"
func LatestVersionTag(repoUrl string) (version_tag string) {
	// 禁止重定向
	// http.DefaultClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
	// 	return http.ErrUseLastResponse
	// }
	clientWithoutRedirect := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url := repoUrl + "/releases/latest"
	resp, _ := clientWithoutRedirect.Head(url)
	// 正则获取最新版本标签
	re := regexp.MustCompile(`.*releases/tag/([^/]+)`)
	if match := re.FindStringSubmatch(resp.Header.Get("Location")); len(match) > 1 {
		version_tag = match[1]
	}
	return version_tag
}
