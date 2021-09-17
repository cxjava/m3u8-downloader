package ping

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/go-ping/ping"
	log "github.com/sirupsen/logrus"
)

var (
	pingRecords = []PingRecord{}
	outputType  = "p"
	host        = "www.google.com"
)

type PingRecord struct {
	IPAddress string
	PingRTT   time.Duration
}

type Options struct {
	PingRecords []PingRecord
	OutPutType  string
	Host        string
}

func SetOptions(opt Options) {
	pingRecords = opt.PingRecords
	if len(opt.OutPutType) > 0 {
		outputType = opt.OutPutType
	}
	if len(opt.Host) > 0 {
		host = opt.Host
	}
}
func Ping() {
	PingAndPrint(pingRecords, outputType)
}

func PingAndPrint(pingRecords []PingRecord, outputType string) {
	var wg sync.WaitGroup
	threadLimiter := make(chan struct{}, 10)

	var total = int(len(pingRecords))
	bar := pb.Full.Start(total)

	pingedRecords := []PingRecord{}
	for _, pingRecord := range pingRecords {
		wg.Add(1)
		threadLimiter <- struct{}{}
		go func(pingRecord PingRecord) {
			defer func() {
				bar.Increment()
				wg.Done()
				<-threadLimiter
				log.Trace("Finished ping " + pingRecord.IPAddress)
			}()
			pingedRecords = append(pingedRecords, pingIP(pingRecord))
		}(pingRecord)
	}
	wg.Wait()
	bar.Finish()

	sort.Slice(pingedRecords, func(i, j int) bool {
		return pingedRecords[i].PingRTT < pingedRecords[j].PingRTT
	})

	if outputType == "t" {
		outputTime(pingedRecords)
	} else if outputType == "p" {
		outputParameter(pingedRecords)
	}
}

func pingIP(pingRecord PingRecord) PingRecord {
	pinger, err := ping.NewPinger(pingRecord.IPAddress)
	if err != nil {
		log.Error("Ping IP Error: " + err.Error())
		panic(err)
	}
	pinger.Count = 2
	pinger.Timeout = 2 * time.Second
	err = pinger.Run() // blocks until finished
	if err != nil {
		log.Error("Ping Run Error: " + err.Error())
		panic(err)
	}
	stats := pinger.Statistics() // get send/receive/rtt stats
	pingRecord.PingRTT = stats.AvgRtt
	return pingRecord
}

func outputTime(pingRecords []PingRecord) {
	log.Trace("Print all ping time")
	fmt.Println("")
	for _, pingRecord := range pingRecords {
		if pingRecord.PingRTT > 0 {
			fmt.Printf("%s	time: %s\n", pingRecord.IPAddress, pingRecord.PingRTT)
		}
	}
	fmt.Println("")
}

func outputParameter(pingRecords []PingRecord) {
	log.Trace("Print as parameter for downloading")
	fmt.Println("")
	for _, pingRecord := range pingRecords {
		if pingRecord.PingRTT > 0 {
			fmt.Printf(`-C '%s:%s' `, host, pingRecord.IPAddress)
		}
	}
	fmt.Println("")
	fmt.Println("")
}
