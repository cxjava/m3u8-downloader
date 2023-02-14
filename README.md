# m3u8-downloader

[![Build Status](https://github.com/cxjava/m3u8-downloader/actions/workflows/build-and-test.yml/badge.svg)](https://github.com/cxjava/m3u8-downloader/actions/workflows/build-and-test.yml)
[![License](https://img.shields.io/github/license/cxjava/m3u8-downloader.svg)](https://github.com/cxjava/m3u8-downloader)
[![Release](https://img.shields.io/github/release/cxjava/m3u8-downloader.svg)](https://github.com/cxjava/m3u8-downloader/releases)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cxjava/m3u8-downloader)

README | [简体中文](README.zh-cn.md)

m3u8 downloader by Golang

![m3u8-downloader-ScreenShot](https://user-images.githubusercontent.com/802316/133533481-483aa464-2fbe-4a25-9539-4a6345481dcd.png)

## feature

- Support CDN download, break the speed limit
- Support custom key, the format of the key can be hex or base64, or original
- Support custom host prefix
- Read m3u8 from file or m3u8 from network address
- Support proxy
- Support custom request header
- Whether to delete the downloaded ts file
- Custom number of threads
- Custom download file name
- Progress bar
- Support running on Apple M1, router

## TODO

- [ ] Support decryption while downloading
- [ ] ffmpeg processing command
- [ ] Show total file size
- [ ] Delete error CDN IP


## Installation

###  Linux / macOS

This [godownloader](https://github.com/kamilsk/godownloader) script will query GitHub for the latest release and download the correct binary for your platform into the directory set with the `-b` flag.

#### System-wide Install

```bash
wget -O - https://raw.githubusercontent.com/cxjava/m3u8-downloader/main/install.sh | sh -s -- -b /usr/local/bin
```


## Ping command

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

### How to use CDN for download

- Open [https://ping.chinaz.com/](https://ping.chinaz.com/) Enter the address of the domain name for which you need to find the IP, and click the Copy button after the ping is finished

![1](https://user-images.githubusercontent.com/802316/133531905-ac398cc4-77da-44e3-a309-351feebd0628.png)

- Run the `ping` command to query the ping time of each IP

``` shell
./m3u8-downloader p 8.8.8.8,1.1.1.1,9.9.9.9,
INFO[2021-09-16T08:56:33+08:00] Try to ping 8.8.8.8,1.1.1.1,9.9.9.9,         
3 / 3 [-----------------------------------------------------------------------------------------] 100.00% 3 p/s 1.3s

8.8.8.8 time: 55.4175ms
9.9.9.9 time: 71.2685ms
1.1.1.1 time: 196.496ms

```

- Run the `ping` command with the parameter `-o P -d www.google.com` to get the parameters needed for the `download` command

``` shell
./m3u8-downloader ping 8.8.8.8,1.1.1.1,9.9.9.9, -d www.google.com -o p 
INFO[2021-09-16T08:57:21+08:00] Try to ping 8.8.8.8,1.1.1.1,9.9.9.9,         
3 / 3 [-----------------------------------------------------------------------------------------] 100.00% 3 p/s 1.3s

-C 'www.google.com:8.8.8.8' -C 'www.google.com:9.9.9.9' -C 'www.google.com:1.1.1.1' 

```

- Copy `-C 'www.google.com:8.8.8.8' -C 'www.google.com:9.9.9.9' -C 'www.google.com:1.1.1.1'` to the `download` command

## Download command

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

### Quickly copy download information from chrome

- Open Chrome's developer tools, find the m3u8 request, right click `Copy`->`Copy as cURL`, as follows

![9](https://user-images.githubusercontent.com/802316/133644904-aa049868-e1fc-40d6-b7b8-85f615f86c07.png)

- Replace `curl` with `. /m3u8-downloader download`, and remove `--compressed` if it has it

``` shell
./m3u8-downloader download 'https://abc.com/def.m3u8' \
  -H 'sec-ch-ua: "Microsoft Edge";v="93", " Not;A Brand";v="99", "Chromium";v="93"' \
  -H 'accept: text/html' \
  -H 'dnt: 1' \
  -H 'x-requested-with: XMLHttpRequest' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36 Edg/93.0.961.47' \
  -H 'accept-language: en,en-US;q=0.9' \
```

## Thanks

- [https://github.com/Greyh4t/m3u8-Downloader-Go](https://github.com/Greyh4t/m3u8-Downloader-Go)
- [https://github.com/llychao/m3u8-downloader](https://github.com/llychao/m3u8-downloader)
