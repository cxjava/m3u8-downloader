package downloader

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/cheggaaa/pb/v3"
	"github.com/cxjava/m3u8-downloader/decrypter"
	"github.com/cxjava/m3u8-downloader/utils"
	"github.com/grafov/m3u8"
	log "github.com/sirupsen/logrus"
)

var (
	tmpl       = `{{ "Downloading:" | rndcolor }} {{string . "prefix" | rndcolor }}{{counters . | rndcolor }} {{bar . | rndcolor }} {{percent . | rndcolor }} {{speed . | rndcolor }} {{rtime . "ETA %s"| rndcolor }}{{string . "suffix"| rndcolor }}`
	syncByte   = uint8(71) //0x47
	outputPath string
)

var (
	downloadUrl    string
	output         string
	downloadDir    string
	proxy          string
	deleteSyncByte bool
	deleteTS       bool
	threadNumber   int
	headers        []string
	cdns           []string
	baseUrl        string
	keyStr         string
	keyFormat      string
	useFFmpeg      bool
)

// Options defines common m3u8 options.
type Options struct {
	DownloadUrl    string
	Output         string
	DownloadDir    string
	Proxy          string
	DeleteSyncByte bool
	DeleteTS       bool
	ThreadNumber   int
	Headers        []string
	CDNs           []string
	BaseUrl        string
	Key            string
	KeyFormat      string
	UseFFmpeg      bool
}

// SetOptions sets the common request option.
func SetOptions(opt Options) {
	downloadUrl = opt.DownloadUrl
	output = opt.Output
	downloadDir = ResolveDir(opt.DownloadDir)
	proxy = opt.Proxy
	deleteSyncByte = opt.DeleteSyncByte
	deleteTS = opt.DeleteTS
	threadNumber = opt.ThreadNumber
	headers = opt.Headers
	cdns = opt.CDNs
	baseUrl = opt.BaseUrl
	keyStr = opt.Key
	keyFormat = opt.KeyFormat
	useFFmpeg = opt.UseFFmpeg
}

func ResolveDir(dirStr string) string {
	abs := filepath.IsAbs(dirStr)
	if abs {
		log.Trace("Resolve download directory:" + dirStr)
		return dirStr
	}
	pwd, _ := os.Getwd()
	dir, err := filepath.Abs(filepath.Join(pwd, dirStr))

	if err != nil {
		log.Error("Resolve download directory failed")
		return dirStr
	}

	log.Trace("Resolve download directory:" + dir)
	return dir
}

func Download() {
	initCDN(cdns)
	initHttpClient(proxy, headers)
	checkOutputFolder()
	var data []byte
	var err error
	if strings.HasPrefix(downloadUrl, "http") {
		data, err = download(downloadUrl)
		if err != nil {
			log.Error("Download m3u8 failed:" + downloadUrl + ",Error:" + err.Error())
			return
		}
	} else {
		// read from file
		data, err = os.ReadFile(downloadUrl)
		if err != nil {
			log.Error("Read m3u8 file failed:" + downloadUrl + ",Error:" + err.Error())
			return
		}
		if len(baseUrl) == 0 {
			log.Warn("make sure ts file have full path in the m3u8 file")
		}
	}

	mpl, err := parseM3u8(data, downloadUrl)
	if err != nil {
		log.Error("Parse m3u8 file failed:" + err.Error())
		return
	} else {
		log.Info("Parse m3u8 file successfully")
	}

	downloadM3u8(mpl)
	temp_name := mergeFile()
	renameFile(temp_name)
}

func renameFile(temp_file string) {
	path1 := filepath.Join(outputPath, temp_file)
	path2 := filepath.Join(downloadDir, output)
	err := os.Rename(path1, path2)
	if err != nil {
		log.Println("[error] Rename failed: " + err.Error())
	}
	if deleteTS {
		err = os.RemoveAll(outputPath)
		if err != nil {
			log.Println(err)
		}
	}
}

func mergeFile() string {
	if !useFFmpeg {
		switch runtime.GOOS {
		case "windows":
			return utils.WinMergeFile(outputPath, deleteTS)
		default:
			return utils.UnixMergeFile(outputPath, deleteTS)
		}
	}
	return utils.FFmpegMergeFile(outputPath, deleteTS)
}

func checkOutputFolder() {
	log.Trace("Check output folder")
	if len(downloadDir) == 0 {
		return
	}

	outputPath = filepath.Join(downloadDir, output+"_downloading")
	log.Trace("Output path is : " + outputPath)
	utils.MkAllDir(outputPath)
}

func parseM3u8(data []byte, downloadUrl string) (*m3u8.MediaPlaylist, error) {
	log.Debug("Parse m3u8")
	playlist, listType, err := m3u8.Decode(*bytes.NewBuffer(data), false)
	if err != nil {
		log.Error("Decode m3u8 failed: " + err.Error())
		return nil, err
	}

	if listType == m3u8.MEDIA {
		var baseHost *url.URL
		if len(baseUrl) > 0 {
			baseHost, err = url.Parse(baseUrl)
			if err != nil {
				log.Error("url.Parse(" + baseUrl + ") failed: " + err.Error())
				return nil, errors.New("parse base url failed: " + err.Error())
			}
		} else if len(downloadUrl) > 0 {
			baseHost, err = url.Parse(downloadUrl)
			if err != nil {
				log.Error("url.Parse(" + downloadUrl + ") failed: " + err.Error())
				return nil, errors.New("parse m3u8 url failed: " + err.Error())
			}
		}
		log.Trace("Base host is " + baseHost.String())

		mpl := playlist.(*m3u8.MediaPlaylist)

		if mpl.Key != nil && mpl.Key.URI != "" {
			uri, err := formatURI(baseHost, mpl.Key.URI)
			if err != nil {
				log.Error("formatURI(" + mpl.Key.URI + ") failed: " + err.Error())
				return nil, err
			}
			log.Trace("MPL key URI is " + uri)
			mpl.Key.URI = uri
		}

		total := int(mpl.Count())
		for i := 0; i < total; i++ {
			segment := mpl.Segments[i]

			uri, err := formatURI(baseHost, segment.URI)
			if err != nil {
				log.Error("formatURI(" + segment.URI + ") failed: " + err.Error())
				return nil, err
			}
			log.Trace("Segment URI is " + uri)
			segment.URI = uri

			if segment.Key != nil && segment.Key.URI != "" {
				uri, err := formatURI(baseHost, segment.Key.URI)
				if err != nil {
					log.Error("formatURI(" + segment.Key.URI + ") failed: " + err.Error())
					return nil, err
				}
				log.Trace("Segment key URI is " + uri)
				segment.Key.URI = uri
			}

			mpl.Segments[i] = segment
		}
		return mpl, nil
	}

	return nil, errors.New("unsupport m3u8 type")
}

func downloadM3u8(mpl *m3u8.MediaPlaylist) {

	var wg sync.WaitGroup
	threadLimiter := make(chan struct{}, threadNumber)

	var total = int(mpl.Count())
	bar := pb.ProgressBarTemplate(tmpl).Start64(int64(total))

	for i := 0; i < total; i++ {
		wg.Add(1)
		threadLimiter <- struct{}{}
		go func(index int, segment *m3u8.MediaSegment, globalKey *m3u8.Key) {
			defer func() {
				bar.Increment()
				wg.Done()
				<-threadLimiter
				log.Trace("args ...interface{}")
			}()
			curr_path := fmt.Sprintf("%s/%05d.ts", outputPath, index)
			if utils.IsExist(curr_path) {
				log.Warn("File: " + curr_path + " already exist")
				return
			}
			var keyURL, ivStr string
			if segment.Key != nil && segment.Key.URI != "" {
				keyURL = segment.Key.URI
				ivStr = segment.Key.IV
			} else if globalKey != nil && globalKey.URI != "" {
				keyURL = globalKey.URI
				ivStr = globalKey.IV
			}
			log.Trace("keyURL is " + keyURL + ", ivStr is " + ivStr)

			data, err := download(segment.URI)
			if err != nil {
				log.Error("Download : " + segment.URI + " failed: " + err.Error())
			}

			var originalData []byte

			if len(keyStr) > 0 {
				log.Info("Try to decrypt data by custom key " + keyStr)
				var key, iv []byte
				if ivStr != "" {
					iv, err = hex.DecodeString(strings.TrimPrefix(ivStr, "0x"))
					if err != nil {
						log.Error("Decode iv failed:" + err.Error())
					}
				} else {
					iv = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(index)}
				}
				switch strings.ToLower(keyFormat) {
				case "original":
					key = []byte(keyStr)
				case "hex":
					key = utils.HexDecode(keyStr)
				case "base64":
					var err error
					key, err = base64.StdEncoding.DecodeString(keyStr)
					if err != nil {
						log.Errorf("base64 Decode %s Failed: %s.", keyStr, err.Error())
					}
				default:
					key = []byte(keyStr)
				}

				originalData, err = decrypter.Decrypt(data, key, iv)
				if err != nil {
					log.Errorf("Decrypt failed by own key %s : %s", keyStr, err.Error())
				}
			} else if keyURL == "" {
				originalData = data
			} else {
				log.Info("Try to decrypt data")
				var key, iv []byte
				key, err = download(keyURL)
				if err != nil {
					log.Error("Download : " + keyURL + " failed: " + err.Error())
				}

				if ivStr != "" {
					iv, err = hex.DecodeString(strings.TrimPrefix(ivStr, "0x"))
					if err != nil {
						log.Error("Decode iv failed:" + err.Error())
					}
				} else {
					iv = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(index)}
				}
				originalData, err = decrypter.Decrypt(data, key, iv)
				if err != nil {
					log.Error("Decrypt failed:" + err.Error())
				}
			}

			if deleteSyncByte {
				log.Info("Delete sync byte.")
				// https://en.wikipedia.org/wiki/MPEG_transport_stream
				// Some TS files do not start with SyncByte 0x47, they can not be played after merging,
				// Need to remove the bytes before the SyncByte 0x47(71).
				dataLength := len(originalData)
				for j := 0; j < dataLength; j++ {
					if originalData[j] == syncByte {
						log.Warn("Find sync byte, and delete it.")
						originalData = originalData[j:]
						break
					}
				}
			}

			err = os.WriteFile(curr_path, originalData, 0666)
			if err != nil {
				log.Error("WriteFile failed:" + err.Error())
			}
			log.Trace("Save file '" + curr_path + "' successfully!")
		}(i, mpl.Segments[i], mpl.Key)
	}
	wg.Wait()
	bar.Finish()
}

func formatURI(base *url.URL, u string) (string, error) {
	if strings.HasPrefix(u, "http") {
		return u, nil
	}

	if base == nil {
		return "", errors.New("you must set m3u8 url for file to download")
	}

	obj, err := base.Parse(u)
	if err != nil {
		return "", err
	}

	return obj.String(), nil
}
