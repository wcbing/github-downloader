# Github Releases 更新下载器

根据 `data/github.json` 定义内容检查 Github Releases 更新并下载特定文件。

主要用于服务访问 Github 受限（无法访问、访问速度慢）的人群。

## 下载、安装

若需要预编译的可执行文件，请点击 [releases](https://github.com/wcbing/github-downloader/releases)

若您系统安装有 Go，可以直接执行
```sh
go install github.com/wcbing/github-downloader@latest
```

## 用法
```
Usage: 
    -d, --data <data_dir>   Read repo config from <data_dir>
    -o, --output <dir>      Save files to <dir>, default to current dir
    -h, --help              Show this help message
    -p, --proxy <url>       Download files from github proxy <url>
    -r, --recursive         Recursive create directory, file save path like: 
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
                            'https:/github.com/<user>/<repo>/releases/
                            download/<version-tag>/<filename>'
                            默认路径: '<user>__<repo>/<filename>'
    --dry-run               用 http 的 head 方法试运行（不下载文件）
```

具体用途

- `-r`
    - 不使用 `-r`  
    适用于做文件服务器、镜像站
    - 使用 `-r`  
    适用于 Github 下载加速
        - 作为文件服务器实现下载加速（对标“反代”类加速的形式）
        - 配合 301 重定向，用于转发到其他下载加速服务
- `-p`
    - 使用 `-p`  
    适用于服务器访问 Github 受限


## 配置文件格式

`data/github.json` 样例

```json
{
    "rustdesk": {
        "repo": "rustdesk/rustdesk",
        "file_list": [
            "rustdesk-{version_tag}-x86_64.deb"
        ]
    },
    "draw.io": {
        "repo": "jgraph/drawio-desktop",
        "file_list": [
            "drawio-amd64-{stripped_version}.deb"
        ]
    }
}
```

说明：

- `"repo": xxx`  
github 仓库名，形如：user/repo。
- `"file_list": []`  
文件名列表，可包含以下变量。
- `{version_tag}`  
用于代替文件名中和 Releases Tag 相同的部分。
- `{stripped_version}`  
用于一些 Tag 是 `v1.1.0`，但是文件名中是 `1.1.0` 的情况。

## Todo

- [x] `-p, --proxy <url>` 使用 Github 下载代理
- [x] `-d, --dir <data_dir>` 指定仓库配置所在目录
- [ ] 限制并发数量