package cmd

import (
	"fmt"
	"strings"

	"github.com/cxjava/m3u8-downloader/downloader"
	"github.com/golang-module/carbon"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	downloadDir    string
	downloadUrl    string
	baseUrl        string
	deleteSyncByte bool
	deleteTS       bool
	proxy          string
	output         string
	threadNumber   int
	headers        []string
	cdns           []string
	logLevel       string
	key            string
	keyFormat      string
	useFFmpeg      bool
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:     "download",
	Aliases: []string{"d"},
	Short:   "download the m3u8 and save it",
	Long:    `All ts segments will be downloaded into a folder then be joined into a single TS file.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initLog()
		downloadUrl = args[0]
		if len(output) == 0 {
			output = fmt.Sprintf("%d.ts", carbon.Now().TimestampWithMillisecond())
		}
		log.Info("output file name is :", output)
		options := downloader.Options{
			DownloadUrl:    downloadUrl,
			Output:         output,
			DownloadDir:    downloadDir,
			BaseUrl:        baseUrl,
			Proxy:          proxy,
			DeleteSyncByte: deleteSyncByte,
			DeleteTS:       deleteTS,
			ThreadNumber:   threadNumber,
			Headers:        headers,
			CDNs:           cdns,
			Key:            key,
			KeyFormat:      keyFormat,
			UseFFmpeg:      useFFmpeg,
		}
		downloader.SetOptions(options)
		downloader.Download()
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&output, "output", "o", "", "file name for save.")
	downloadCmd.Flags().StringVarP(&downloadDir, "downloadDir", "f", "./outputFolder", "download directory, base on current folder.")
	downloadCmd.Flags().StringVarP(&baseUrl, "baseUrl", "u", "", "base url for m3u8.")
	downloadCmd.Flags().StringVarP(&proxy, "proxy", "x", "", "use proxy. eg. http://127.0.0.1:8080")
	downloadCmd.Flags().BoolVarP(&deleteSyncByte, "deleteSyncByte", "d", false, "some TS files do not start with SyncByte 0x47, they can not be played after merging, need to remove the bytes before the SyncByte.")
	downloadCmd.Flags().BoolVarP(&deleteTS, "deleteTS", "D", true, "delete all the downloaded TS file.")
	downloadCmd.Flags().IntVarP(&threadNumber, "threadNumber", "n", 10, "the number of download thread.")
	downloadCmd.Flags().StringArrayVarP(&headers, "header", "H", nil, "custom http header(s), eg. -H 'user-agent: Mozilla/5.0...' -H 'accept: */*' .")
	downloadCmd.Flags().StringArrayVarP(&cdns, "cdn", "C", nil, "CDN(s) for the download domain, eg. -C 'www.google.com:8.8.8.8' -C 'www.google.com:1.1.1.1' -C 'www.google.com:9.9.9.9' .")
	downloadCmd.Flags().StringVarP(&logLevel, "logLevel", "l", "Info", "logging level on a Logger,logging levels: Trace, Debug, Info, Warning, Error, Fatal and Panic.")
	downloadCmd.Flags().StringVarP(&key, "key", "", "", "custom key to decrypt ts data.")
	downloadCmd.Flags().StringVarP(&keyFormat, "keyFormat", "", "original", "format of key, format can be those values: original, hex, base64.")
	downloadCmd.Flags().BoolVarP(&useFFmpeg, "UseFFmpeg", "", false, "use FFmpeg for merging TS files.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	// log.SetFormatter(&log.JSONFormatter{})
	switch strings.ToLower(logLevel) {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
