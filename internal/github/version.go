package github

import (
	"net/http"
	"regexp"
)

// 获取最新版本标签
// repoUrl: "https://github.com/user/repo"
func latestReleasesTag(repoUrl string) (releasesTag string) {
	// 禁止重定向
	clientWithoutRedirect := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url := repoUrl + "/releases/latest"
	resp, err := clientWithoutRedirect.Head(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	// 正则获取最新版本标签
	re := regexp.MustCompile(`.*releases/tag/([^/]+)`)
	if match := re.FindStringSubmatch(resp.Header.Get("Location")); len(match) > 1 {
		releasesTag = match[1]
	}
	return
}

// 提取版本号
func stripVerison(releasesTag string) (version string) {
	re := regexp.MustCompile(`[0-9].*`)
	if match := re.FindString(releasesTag); len(match) > 0 {
		version = match
	}
	return
}
