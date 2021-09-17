package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMkAllDir(t *testing.T) {
	type args struct {
		dirs string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "create test folder successfully",
			args: args{
				dirs: "test/test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MkAllDir(tt.args.dirs)
			RemoveAllDir(tt.args.dirs)
		})
	}
}

func TestRemoveAllDir(t *testing.T) {
	type args struct {
		dirs string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "remove folder successfully",
			args: args{
				dirs: "test/test/testaa",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MkAllDir(tt.args.dirs)
			RemoveAllDir(tt.args.dirs)
			if IsExist(tt.args.dirs) {
				log.Fatalln("remove failed")
			}
		})
	}
	pwd, _ := os.Getwd()
	fmt.Println(filepath.Join(pwd + "/./abc"))
	fmt.Println(filepath.Join(pwd + "/../../def"))
	fmt.Println(filepath.Join(pwd + "/abc"))
	fmt.Println(filepath.Join(pwd + "abcd"))
}

func TestIsExist(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Not exist such path",
			args: args{
				filePath: "test/test",
			},
			want: false,
		},
		{
			name: "exist such file",
			args: args{
				filePath: "file.go",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExist(tt.args.filePath); got != tt.want {
				t.Errorf("IsExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileAbs(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "exist such file.go",
			args: args{
				path: "../file.go",
			},
			want: "m3u8-downloader/file.go",
		},
		{
			name: "exist such file abc",
			args: args{
				path: "./../abc",
			},
			want: "m3u8-downloader/abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileAbs(tt.args.path); !strings.HasSuffix(got, tt.want) {
				t.Errorf("FileAbs() = %v, want end with %v", got, tt.want)
			}
		})
	}
}

func TestFilenameWithoutExtension(t *testing.T) {
	type args struct {
		fn string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get same file name",
			args: args{
				fn: "../.././../abc",
			},
			want: "../.././../abc",
		},
		{
			name: "get the same file name",
			args: args{
				fn: "afasdf-adf_adfas",
			},
			want: "afasdf-adf_adfas",
		},
		{
			name: "get def.exe",
			args: args{
				fn: "def.exe.def",
			},
			want: "def.exe",
		},
		{
			name: "get abc",
			args: args{
				fn: "abc.txt",
			},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilenameWithoutExtension(tt.args.fn); got != tt.want {
				t.Errorf("FilenameWithoutExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}
