package cmd

import (
	"fmt"
	"os"

	"github.com/wcbing/github-downloader/config"
)

var helpMessage = `Usage: 
    -d, --dir <data_dir>    Read repo config from <data_dir>
    -h, --help              Show this help message
    -p, --proxy <url>       Download files from github proxy <url>
    -r, --recursive         Recursive create directory, file save path like: 
                            'https:/github.com/<user>/<repo>/releases/
                            download/<version-tag>/<filename>'
                            Default path like: 'releases/<user>__<repo>/<filename>'
    --dry-run               Dry run with HTTP head method (do not download)

用法: 
    -d, --dir <data_dir>    从 <data_dir> 读取仓库配置
    -h, --help              显示该帮助信息
    -p, --proxy <url>       从 Github 代理 <url> 下载文件
    -r, --recursive         递归的创建目录，文件保存路径: 
                            'https:/github.com/<user>/<repo>/releases/
                            download/<version-tag>/<filename>'
                            默认路径: 'releases/<user>__<repo>/<filename>'
    --dry-run               用 http 的 head 方法试运行（不下载文件）
`

// 读取命令行参数
func ReadArgs() {
	args := os.Args
	if len(args) > 1 {
		for i, arg := range args {
			switch arg {
			case "-d", "--dir":
				config.Config["dir"] = true
				// 读取下一个参数作为文件路径
				if i+1 < len(args) {
					config.DataDir = args[i+1]
				}
			case "-h", "--help":
				// 显示帮助信息并退出
				fmt.Print(helpMessage)
				os.Exit(0)
			case "-p", "--proxy":
				config.Config["proxy"] = true
				// 读取下一个参数作为代理地址
				if i+1 < len(args) {
					config.Proxy = args[i+1]
				}
			case "-r", "--recursive":
				config.Config["recursive"] = true
			case "--dry-run":
				config.Config["dry-run"] = true
			}
		}
	}
}
