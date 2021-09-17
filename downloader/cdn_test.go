package downloader

import (
	"fmt"
	"testing"
)

func Test_parseCDN(t *testing.T) {
	type args struct {
		cdns []string
	}
	tests := []struct {
		name string
		args args
		want map[string]CDNS
	}{
		{
			name: "get the original URI",
			args: args{
				cdns: []string{"www.google1.com:8.8.8.8", "www.google2.com:1.1.1.1", "www.google1.com:9.9.9.9", "www.google2.com:7.9.9.9", "www.google4.com:7.7.9.9"},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCDN(tt.args.cdns)
			if len(got) != 3 {
				t.Errorf("parseCDN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addCDN(t *testing.T) {
	type args struct {
		cdnmap map[string]CDNS
	}
	tests := []struct {
		name string
		args args
		want map[string]chan string
	}{
		{
			name: "get the cdn",
			args: args{
				cdnmap: map[string]CDNS{
					"1": {
						Domain: "www.google.com",
						IPs:    []string{"1.1.1.1", "2.2.2.2", "5.5.5.5", "6.6.6.6"},
					},
					"3": {
						Domain: "www.google2.com",
						IPs:    []string{"3.3.3.3", "4.4.4.4"},
					},
					"5": {
						Domain: "www.baidu.com",
						IPs:    []string{"bbb", "ddd", "eeee", "ffff"},
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addCDN(tt.args.cdnmap)
			for i := 0; i < 10; i++ {
				if v, ok := got["www.baidu.com"]; ok {
					fmt.Println(<-v)
				}
			}
			t.Logf("Done!")
		})
	}
}
