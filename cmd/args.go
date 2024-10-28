package cmd

import (
	"fmt"
	"os"

	"github.com/wcbing/github-downloader/config"
)

var helpMessage = `Usage: 
    -h, --help          Show this help message
    -r, --recursive     Recursive create directory, like: 
                        'https:/github.com/<user>/<repo>/releases/
                        download/<version-tag>/<filename>'
                        Default like: 'releases/<user>/<repo>/<filename>'
    --dry-run           Dry run with HTTP head method (do not download)

用法: 
    -h, --help          显示该帮助信息
    -r, --recursive     递归的创建目录，文件路径: 
                        'https:/github.com/<user>/<repo>/releases/
                        download/<version-tag>/<filename>'
                        默认情况: 'releases/<user>/<repo>/<filename>'
    --dry-run           用 http 的 head 方法试运行（不下载文件）
`

// 读取命令行参数
func ReadArgs() {
	args := os.Args
	if len(args) > 1 {
		for _, arg := range args {
			switch arg {
			case "-h", "--help":
				fmt.Print(helpMessage)
				os.Exit(0)
			case "-r", "--recursive":
				config.Config["recursive"] = true
			case "--dry-run":
				config.Config["dry-run"] = true
			}
		}
	}
}
