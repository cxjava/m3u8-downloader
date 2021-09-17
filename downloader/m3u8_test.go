package downloader

import (
	"net/url"
	"testing"
)

func Test_formatURI(t *testing.T) {
	type args struct {
		base *url.URL
		u    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "get the format URI",
			args: args{
				base: parseUrl("https://res001.geekbang.org/media/audio/51/02/51"),
				u:    "/124.ts",
			},
			want:    "https://res001.geekbang.org/124.ts",
			wantErr: false,
		},
		{
			name: "get the original URI",
			args: args{
				base: parseUrl("https://res001.geekbang.org/media/audio/51/02/51"),
				u:    "https://www.baidu.com/124.ts",
			},
			want:    "https://www.baidu.com/124.ts",
			wantErr: false,
		},
		{
			name: "throw error when base is nil",
			args: args{
				base: nil,
				u:    "124.ts",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatURI(tt.args.base, tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("formatURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func parseUrl(u string) *url.URL {
	ur, _ := url.Parse(u)
	return ur
}
