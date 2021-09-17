package downloader

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

var (
	client         = resty.New()
	defaultHeaders = map[string]string{
		"dnt":             "1",
		"Accept":          "*/*",
		"Accept-Charset":  "UTF-8,*;q=0.5",
		"Accept-Encoding": "gzip,deflate,sdch",
		"Accept-Language": "en-US,en;q=0.8",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 Edg/91.0.864.54",
	}
)

func download(url string) ([]byte, error) {
	log.Trace("Start download " + url)
	resp, err := client.R().Get(url)
	log.Trace("Finish download " + url)
	log.Tracef("Error      :%+v", err)
	log.Trace("Status     :" + resp.Status())
	return resp.Body(), err
}

func initHttpClient(proxy string, headers []string) {
	log.Trace("Init http client")
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// addr is "google.com:443"
			domainAndPort := strings.Split(addr, ":")
			domain := domainAndPort[0]
			if ipch, ok := cdnChan[domain]; ok {
				ip := <-ipch
				port := domainAndPort[1]
				addr = fmt.Sprintf("%s:%s", ip, port)
				log.Debug("Use cdn address :", addr)
			}
			return dialer.DialContext(ctx, network, addr)
		},
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		DisableKeepAlives:     true, // for CDN
		MaxIdleConns:          100,
		IdleConnTimeout:       5 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	client.SetTransport(transport)

	client.SetTimeout(5 * time.Second)

	if len(proxy) > 0 {
		log.Info("Use proxy : ", proxy)
		client.SetProxy(proxy)
	}

	if len(headers) > 0 {
		log.Info("Set customization header")
		for _, header := range headers {
			h := strings.SplitN(header, ":", 2)
			log.Debug("customization header =>" + h[0] + ":" + h[1])
			client.SetHeader(h[0], strings.TrimSpace(h[1]))
		}
	} else {
		log.Debug("Use default header")
		log.Tracef("Use default header :%+v", defaultHeaders)
		client.SetHeaders(defaultHeaders)
	}
	client.
		SetRetryCount(5).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second)

	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(10))

}
