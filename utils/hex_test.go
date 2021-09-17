package utils

import (
	"reflect"
	"testing"
)

func TestHexDecode(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "decode hello",
			args: args{
				s: "68656c6c6f",
			},
			want: []byte{104, 101, 108, 108, 111},
		},

		{
			name: "decode m3u8Golang",
			args: args{
				s: "6d337538476f6c616e67",
			},
			want: []byte("m3u8Golang"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HexDecode(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HexDecode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHexEncode(t *testing.T) {
	type args struct {
		s []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "encode hello",
			args: args{
				s: []byte("hello"),
			},
			want: "68656c6c6f",
		},
		{
			name: "encode m3u8Golang",
			args: args{
				s: []byte("m3u8Golang"),
			},
			want: "6d337538476f6c616e67",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HexEncode(tt.args.s); got != tt.want {
				t.Errorf("HexEncode() = %v, want %v", got, tt.want)
			}
		})
	}
}
