package github

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLatestReleasesTag(t *testing.T) {
	// 创建测试服务器模拟GitHub API响应
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/user/repo/releases/latest" {
			// 设置重定向头，重定向到具体版本
			w.Header().Set("Location", "https://github.com/user/repo/releases/tag/v1.2.3")
			w.WriteHeader(http.StatusFound) // 302重定向
			return
		}
		if r.URL.Path == "/user/no-releases-repo/releases/latest" {
			// 设置重定向头，重定向到Releases
			w.Header().Set("Location", "https://github.com/user/no-releases-repo/releases")
			w.WriteHeader(http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	// 测试用例
	tests := []struct {
		name     string
		repoURL  string
		expected string
	}{
		{
			name:     "有 Releases 仓库",
			repoURL:  ts.URL + "/user/repo",
			expected: "v1.2.3",
		},
		{
			name:     "无 Releases 仓库",
			repoURL:  ts.URL + "/user/no-releases-repo",
			expected: "",
		},
		{
			name:     "无效仓库 URL",
			repoURL:  ts.URL + "/invalid-url",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := latestReleasesTag(tt.repoURL)
			if got != tt.expected {
				t.Errorf("latestReleasesTag() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestStripVerison(t *testing.T) {
	tests := []struct {
		name        string
		releasesTag string
		expected    string
	}{
		{
			name:        "标准v前缀版本号",
			releasesTag: "v1.2.3",
			expected:    "1.2.3",
		},
		{
			name:        "无前缀版本号",
			releasesTag: "1.2.3",
			expected:    "1.2.3",
		},
		{
			name:        "带有额外文本的版本号",
			releasesTag: "release-1.2.3",
			expected:    "1.2.3",
		},
		{
			name:        "带有构建信息的版本号",
			releasesTag: "v1.2.3-beta.1",
			expected:    "1.2.3-beta.1",
		},
		{
			name:        "空字符串",
			releasesTag: "",
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripVerison(tt.releasesTag)
			if got != tt.expected {
				t.Errorf("stripVerison() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
