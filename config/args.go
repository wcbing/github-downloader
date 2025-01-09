package config

import (
	"flag"
	"os"
)

var helpMessage = `Usage: 
    -d, --data <data_dir>   Read repo config from <data_dir>
    -o, --output <dir>      Save files to <dir>, default to current dir
    -h, --help              Show this help message
    -p, --proxy <url>       Download files from github proxy <url>
    -r, --recursive         Recursive create directory, file save path like: 
    -t, --thread <number>   The number of concurrent download threads, default is 5
                            'https:/github.com/<user>/<repo>/releases/
                            download/<version-tag>/<filename>'
                            Default path like: '<user>__<repo>/<filename>'
    --dry-run               Dry run with HTTP head method (do not download)

用法: 
    -d, --data <data_dir>   从 <data_dir> 读取仓库配置
    -o, --output <dir>      将文件保存到 <dir>，默认为当前文件夹
    -h, --help              显示该帮助信息
    -p, --proxy <url>       从 Github 代理 <url> 下载文件
    -r, --recursive         递归的创建目录，文件保存路径: 
    -t, --thread <number>   并发下载线程数量，默认为 5
                            'https:/github.com/<user>/<repo>/releases/
                            download/<version-tag>/<filename>'
                            默认路径: '<user>__<repo>/<filename>'
    --dry-run               用 http 的 head 方法试运行（不下载文件）
`

// 读取命令行参数
func ReadArgs() {
	// data_dir
	flag.StringVar(&DataDir, "d", "data", "")
	flag.StringVar(&DataDir, "data", "data", "")
	// output_dir
	flag.StringVar(&OutputDir, "o", ".", "")
	flag.StringVar(&OutputDir, "output", ".", "")
	// proxy
	flag.StringVar(&Proxy, "p", "", "")
	flag.StringVar(&Proxy, "proxy", "", "")
	// recursive
	var recursive bool
	flag.BoolVar(&recursive, "r", false, "")
	flag.BoolVar(&recursive, "recursive", false, "")
	// dryRun
	dryRun := flag.Bool("dry-run", false, "")
	// thread
	flag.IntVar(&Thread, "c", 5, "") // 默认并发数为5
	flag.IntVar(&Thread, "thread", 5, "")
	// help
	var help bool
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")

	flag.Parse()

	if help {
		print(helpMessage)
		os.Exit(0)
	}
	Config["recursive"] = recursive
	Config["dry-run"] = *dryRun
}
