package downloader

import (
	"net/url"
	"testing"

	"github.com/grafov/m3u8"
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

func Test_resolveKeyURL(t *testing.T) {
	globalKey := &m3u8.Key{URI: "https://example.com/key", IV: "0x0102030405060708090a0b0c0d0e0f10"}
	segmentKey := &m3u8.Key{URI: "https://segment.key", IV: "0xdeadbeef"}

	tests := []struct {
		name       string
		segment    *m3u8.MediaSegment
		globalKey  *m3u8.Key
		wantKeyURL string
		wantIV     string
	}{
		{
			name:       "segment key takes priority",
			segment:    &m3u8.MediaSegment{Key: segmentKey},
			globalKey:  globalKey,
			wantKeyURL: "https://segment.key",
			wantIV:     "0xdeadbeef",
		},
		{
			name:       "fallback to global key",
			segment:    &m3u8.MediaSegment{},
			globalKey:  globalKey,
			wantKeyURL: "https://example.com/key",
			wantIV:     "0x0102030405060708090a0b0c0d0e0f10",
		},
		{
			name:       "no keys available",
			segment:    &m3u8.MediaSegment{},
			globalKey:  nil,
			wantKeyURL: "",
			wantIV:     "",
		},
		{
			name:       "segment key URI empty string",
			segment:    &m3u8.MediaSegment{Key: &m3u8.Key{URI: ""}},
			globalKey:  globalKey,
			wantKeyURL: "https://example.com/key",
			wantIV:     "0x0102030405060708090a0b0c0d0e0f10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeyURL, gotIV := resolveKeyURL(tt.segment, tt.globalKey)
			if gotKeyURL != tt.wantKeyURL {
				t.Errorf("resolveKeyURL() keyURL = %v, want %v", gotKeyURL, tt.wantKeyURL)
			}
			if gotIV != tt.wantIV {
				t.Errorf("resolveKeyURL() iv = %v, want %v", gotIV, tt.wantIV)
			}
		})
	}
}

func Test_parseKey(t *testing.T) {
	d := &Downloader{}
	tests := []struct {
		name     string
		format   string
		keyStr   string
		want     string
	}{
		{
			name:     "original format",
			format:   "original",
			keyStr:   "mysecretkey12345",
			want:     "mysecretkey12345",
		},
		{
			name:     "hex format",
			format:   "hex",
			keyStr:   "6865786b6579",
			want:     "hexkey",
		},
		{
			name:     "base64 format",
			format:   "base64",
			keyStr:   "bXkza2V5MTIz",
			want:     "my3key123",
		},
		{
			name:     "default to original for unknown format",
			format:   "unknown",
			keyStr:   "fallbackkey123",
			want:     "fallbackkey123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d.keyFormat = tt.format
			d.keyStr = tt.keyStr
			got := d.parseKey()
			if string(got) != tt.want {
				t.Errorf("parseKey() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func Test_parseIV(t *testing.T) {
	d := &Downloader{}
	tests := []struct {
		name  string
		ivStr string
		index int
		want  []byte
	}{
		{
			name:  "parse iv with 0x prefix",
			ivStr: "0x0102030405060708090a0b0c0d0e0f10",
			index: 0,
			want:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
		{
			name:  "parse iv without 0x prefix",
			ivStr: "ffeeddccbbaa99887766554433221100",
			index: 0,
			want:  []byte{255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0},
		},
		{
			name:  "default iv based on index",
			ivStr: "",
			index: 5,
			want:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5},
		},
		{
			name:  "default iv index 255",
			ivStr: "",
			index: 255,
			want:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.parseIV(tt.ivStr, tt.index)
			if len(got) != 16 {
				t.Errorf("parseIV() length = %d, want 16", len(got))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseIV() [%d] = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func Test_removeSyncBytes(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		enabled bool
		want    []byte
	}{
		{
			name:    "disabled returns original data",
			data:    []byte{0x47, 0x01, 0x02},
			enabled: false,
			want:    []byte{0x47, 0x01, 0x02},
		},
		{
			name:    "starts with sync byte",
			data:    []byte{0x47, 0x01, 0x02},
			enabled: true,
			want:    []byte{0x47, 0x01, 0x02},
		},
		{
			name:    "remove bytes before sync byte",
			data:    []byte{0x00, 0x00, 0x47, 0x01, 0x02},
			enabled: true,
			want:    []byte{0x47, 0x01, 0x02},
		},
		{
			name:    "no sync byte found returns all",
			data:    []byte{0x00, 0x01, 0x02},
			enabled: true,
			want:    []byte{0x00, 0x01, 0x02},
		},
		{
			name:    "empty data",
			data:    []byte{},
			enabled: true,
			want:    []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeSyncBytes(tt.data, tt.enabled)
			if len(got) != len(tt.want) {
				t.Errorf("removeSyncBytes() length = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("removeSyncBytes() [%d] = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func parseUrl(u string) *url.URL {
	ur, _ := url.Parse(u)
	return ur
}
