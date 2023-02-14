# m3u8-downloader

[![构建状态](https://github.com/cxjava/m3u8-downloader/actions/workflows/build-and-test.yml/badge.svg)](https://github.com/cxjava/m3u8-downloader/actions/workflows/build-and-test.yml)
[![img](https://img.shields.io/github/license/cxjava/m3u8-downloader?label=%E8%AE%B8%E5%8F%AF%E8%AF%81)](https://github.com/cxjava/m3u8-downloader)
[![img](https://img.shields.io/github/release/cxjava/m3u8-downloader?label=%E6%9C%80%E6%96%B0%E7%89%88%E6%9C%AC)](https://github.com/cxjava/m3u8-downloader/releases)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cxjava/m3u8-downloader)

[README](README.md) | 简体中文

Golang版本的M3U8下载器

![m3u8-downloader-ScreenShot](https://user-images.githubusercontent.com/802316/133533481-483aa464-2fbe-4a25-9539-4a6345481dcd.png)

## 功能特性

- 支持CDN下载，突破速度限制
- 支持自定义的key，key的格式可以是hex或者base64，或者原始
- 支持自定义的host前缀
- 从文件读取m3u8或者网络地址读取m3u8
- 支持代理
- 支持自定义请求头
- 是否删除下载的ts文件
- 自定义线程数
- 自定义下载文件名
- 进度条
- 支持在Apple M1，路由器上运行

## TODO

- [ ] 支持边下边解密
- [ ] ffmpeg 处理命令
- [ ] 显示总文件大小
- [ ] 删除出错CDN IP

## 安装

###  Linux / macOS

`-b`参数指定安装目录

```bash
wget -O - https://cdn.jsdelivr.net/gh/cxjava/pingtunnel/install.sh | sh -s -- -b /usr/local/bin
```

## Ping命令使用方法

``` shell
./m3u8-downloader p -h

With this command, you can ping multiple IPs to detect IP response speeds. Easily get the fastest CDN IPs.IP is separated by commas, for example: ./m3u8-downloader p 8.8.8.8,1.1.1.1,9.9.9.9, The have two output format, default output is speed of each IP, for example:

./m3u8-downloader p 8.8.8.8,1.1.1.1,9.9.9.9,

9.9.9.9 time: 23.9635ms
8.8.8.8 time: 61.2325ms
1.1.1.1 time: 225.7115ms

Other output is the parameter of download m3u8, for example:

./m3u8-downloader p 8.8.8.8,1.1.1.1,9.9.9.9, -d www.google.com -o p

-C 'www.google.com:8.8.8.8' -C 'www.google.com:1.1.1.1' -C 'www.google.com:9.9.9.9'

Usage:
  m3u8-downloader ping [flags]

Aliases:
  ping, p

Flags:
  -d, --domain string       Only avaliable when output type is 'p', domain to use output as a parameter, such as: -C 'www.google.com:1.1.1.1' -C 'www.google.com:8.8.8.8'  (default "www.google.com")
  -h, --help                help for ping
  -l, --logLevel string     logging level on a Logger,logging levels: Trace, Debug, Info, Warning, Error, Fatal and Panic. (default "Info")
  -o, --outputType string   output type, can be 't' or 'p', if the type is 'p', it will print "-C 'www.example.com:1.1.1.1' -C 'www.example.com:9.9.9.9'". If the type is 't', it will print 8.8.8.8 time: 62.17ms     (default "t")

Global Flags:
      --config string   config file (default is $HOME/.m3u8-downloader.yaml)
```

### 如何利用CDN进行下载

- 打开 [https://ping.chinaz.com/](https://ping.chinaz.com/) 输入需要查找IP的域名地址，Ping结束后，点击复制按钮

![1](https://user-images.githubusercontent.com/802316/133531905-ac398cc4-77da-44e3-a309-351feebd0628.png)

- 运行 `ping` 命令查询每个IP的ping时间

``` shell
./m3u8-downloader p 8.8.8.8,1.1.1.1,9.9.9.9,
INFO[2021-09-16T08:56:33+08:00] Try to ping 8.8.8.8,1.1.1.1,9.9.9.9,         
3 / 3 [-----------------------------------------------------------------------------------------] 100.00% 3 p/s 1.3s

8.8.8.8 time: 55.4175ms
9.9.9.9 time: 71.2685ms
1.1.1.1 time: 196.496ms

```

- 运行 `ping`命令带上参数`-o P -d www.google.com`得到`download`命令需要的参数

``` shell
./m3u8-downloader ping 8.8.8.8,1.1.1.1,9.9.9.9, -d www.google.com -o p 
INFO[2021-09-16T08:57:21+08:00] Try to ping 8.8.8.8,1.1.1.1,9.9.9.9,         
3 / 3 [-----------------------------------------------------------------------------------------] 100.00% 3 p/s 1.3s

-C 'www.google.com:8.8.8.8' -C 'www.google.com:9.9.9.9' -C 'www.google.com:1.1.1.1' 

```

- 复制`-C 'www.google.com:8.8.8.8' -C 'www.google.com:9.9.9.9' -C 'www.google.com:1.1.1.1'`到`download`命令里面

## Download命令使用方法

``` shell
./m3u8-downloader download -h

All ts segments will be downloaded into a folder then be joined into a single TS file.

Usage:
  m3u8-downloader download [flags]

Aliases:
  download, d

Flags:
  -u, --baseUrl string       base url for m3u8.
  -C, --cdn stringArray      CDN(s) for the download domain, eg. -C 'www.google.com:8.8.8.8' -C 'www.google.com:1.1.1.1' -C 'www.google.com:9.9.9.9' .
  -d, --deleteSyncByte       some TS files do not start with SyncByte 0x47, they can not be played after merging, need to remove the bytes before the SyncByte.
  -D, --deleteTS             delete all the downloaded TS file. (default true)
  -f, --downloadDir string   download directory, base on current folder. (default "./outputFolder")
  -H, --header stringArray   custom http header(s), eg. -H 'user-agent: Mozilla/5.0...' -H 'accept: */*' .
  -h, --help                 help for download
      --key string           custom key to decrypt ts data.
      --keyFormat string     format of key, format can be those values: original, hex, base64. (default "original")
  -l, --logLevel string      logging level on a Logger,logging levels: Trace, Debug, Info, Warning, Error, Fatal and Panic. (default "Info")
  -o, --output string        file name for save.
  -x, --proxy string         use proxy. eg. http://127.0.0.1:8080
  -n, --threadNumber int     the number of download thread. (default 10)

Global Flags:
      --config string   config file (default is $HOME/.m3u8-downloader.yaml)
```

### 参数说明

``` shell
  -u, --baseUrl string       覆盖下载m3u8的host
  -C, --cdn stringArray      CDN参数，例子如： -C 'www.google.com:8.8.8.8' -C 'www.google.com:1.1.1.1' -C 'www.google.com:9.9.9.9' .
  -d, --deleteSyncByte       TS二包头，参见： https://en.wikipedia.org/wiki/MPEG_transport_stream .
  -D, --deleteTS             删除下载的TS文件，默认是true
  -f, --downloadDir string   下载的目录. (默认是 "./outputFolder")
  -H, --header stringArray   自定义请求头，参数兼容curl命令的参数，可以用chrome开发者工具，复制为cUrl命令。例如： -H 'user-agent: Mozilla/5.0...' -H 'accept: */*' .
  -h, --help                 帮助命令
      --key string           自定义解密的key.
      --keyFormat string     自定义秘钥的格式，有原始格式，hex格式，base64格式： original, hex, base64. (默认 "original")
  -l, --logLevel string      日志级别: Trace, Debug, Info, Warning, Error, Fatal and Panic. (默认 "Info")
  -o, --output string        下载之后的文件名称，建议设置为ts结尾.
  -x, --proxy string         设置代理下载. 比如：http://127.0.0.1:8080
  -n, --threadNumber int     下载请求的线程数量. (默认 10)
```

### 快速从chrome中复制下载信息

- 打开chrome的开发者工具，找到m3u8那个请求，右键单击`复制`-->`复制为cURL`,如下图

![2](https://user-images.githubusercontent.com/802316/133640083-8a632552-0af5-464f-9720-e5e866f9fbcf.png)

- 替换`curl` 为 `./m3u8-downloader download`, 如果有 `--compressed` 也去掉

``` shell
./m3u8-downloader download 'https://abc.com/def.m3u8' \
  -H 'sec-ch-ua: "Microsoft Edge";v="93", " Not;A Brand";v="99", "Chromium";v="93"' \
  -H 'accept: text/html' \
  -H 'dnt: 1' \
  -H 'x-requested-with: XMLHttpRequest' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36 Edg/93.0.961.47' \
  -H 'accept-language: en,en-US;q=0.9' \
```

## 感谢以下项目

- [https://github.com/Greyh4t/m3u8-Downloader-Go](https://github.com/Greyh4t/m3u8-Downloader-Go)
- [https://github.com/llychao/m3u8-downloader](https://github.com/llychao/m3u8-downloader)
