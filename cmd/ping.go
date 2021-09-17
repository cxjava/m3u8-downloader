package cmd

import (
	"strings"
	"time"

	"github.com/cxjava/m3u8-downloader/ping"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	outputType string
	host       string
	ip         string
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:     "ping",
	Aliases: []string{"p"},
	Short:   "With this command, you can ping multiple IPs to detect IP response speeds. Easily get the fastest CDN IPs.",
	Long: `With this command, you can ping multiple IPs to detect IP response speeds. Easily get the fastest CDN IPs.
IP is separated by commas, for example: ./m3u8-downloader p 8.8.8.8,1.1.1.1,9.9.9.9,
The have two output format, default output is speed of each IP, for example:

./m3u8-downloader p 8.8.8.8,1.1.1.1,9.9.9.9,

9.9.9.9 time: 23.9635ms
8.8.8.8 time: 61.2325ms
1.1.1.1 time: 225.7115ms

Other output is the parameter of download m3u8, for example:

./m3u8-downloader p 8.8.8.8,1.1.1.1,9.9.9.9, -d www.google.com -o p

-C 'www.google.com:8.8.8.8' -C 'www.google.com:1.1.1.1' -C 'www.google.com:9.9.9.9' 

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initLog()
		ip = args[0]
		log.Info("Try to ping " + ip)
		ipList := strings.Split(ip, ",")
		pingRecords := []ping.PingRecord{}
		for _, ip := range ipList {
			if len(strings.TrimSpace(ip)) > 0 {
				pingRecords = append(pingRecords, ping.PingRecord{
					IPAddress: strings.TrimSpace(ip),
					PingRTT:   5 * time.Second,
				})
			}
		}
		option := ping.Options{
			PingRecords: pingRecords,
			Host:        host,
			OutPutType:  outputType,
		}
		ping.SetOptions(option)
		ping.Ping()
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	pingCmd.Flags().StringVarP(&outputType, "outputType", "o", "t", `output type, can be 't' or 'p', if the type is 'p', it will print "-C 'www.example.com:1.1.1.1' -C 'www.example.com:9.9.9.9'". If the type is 't', it will print 8.8.8.8 time: 62.17ms	`)
	pingCmd.Flags().StringVarP(&host, "domain", "d", "www.google.com", "Only avaliable when output type is 'p', domain to use output as a parameter, such as: -C 'www.google.com:1.1.1.1' -C 'www.google.com:8.8.8.8' ")
	pingCmd.Flags().StringVarP(&logLevel, "logLevel", "l", "Info", "logging level on a Logger,logging levels: Trace, Debug, Info, Warning, Error, Fatal and Panic.")
}
