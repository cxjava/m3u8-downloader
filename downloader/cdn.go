package downloader

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
)

type CDNS struct {
	Domain string
	IPs    []string
}

func (d *Downloader) initCDN() {
	cdnList := parseCDN(d.cdns)
	ctx, cancel := context.WithCancel(context.Background())
	d.cdnCancel = cancel
	d.cdnChan = addCDN(ctx, cdnList)
}

func (d *Downloader) StopCDN() {
	if d.cdnCancel != nil {
		d.cdnCancel()
	}
}

func parseCDN(cdns []string) map[string]CDNS {
	log.Trace("Parse CDN")
	fastestCDN := make(map[string]CDNS, 2)

	for _, v := range cdns {
		dp := strings.Split(v, ":")
		domain := dp[0]
		ip := dp[1]
		if cdns, ok := fastestCDN[domain]; ok {
			log.Trace("Exist CDNS for domain: " + domain + ", IP: " + ip)
			cdns.IPs = append(cdns.IPs, ip)
			fastestCDN[domain] = cdns
		} else {
			log.Trace("Create CDNS for domain: " + domain + ", IP: " + ip)
			fastestCDN[domain] = CDNS{
				Domain: domain,
				IPs:    []string{ip},
			}
		}
	}
	return fastestCDN
}

func addCDN(ctx context.Context, cmap map[string]CDNS) map[string]chan string {
	log.Trace("Add CDN")
	cdnMap := make(map[string]chan string, len(cmap))
	for _, cdn := range cmap {
		ipChan := make(chan string, 10)
		go func(ctx context.Context, ipChan chan string, ips []string) {
			for {
				for _, v := range ips {
					select {
					case ipChan <- v:
					case <-ctx.Done():
						log.Trace("CDN goroutine stopped for domain")
						return
					}
				}
			}
		}(ctx, ipChan, cdn.IPs)
		cdnMap[cdn.Domain] = ipChan
	}
	return cdnMap
}
