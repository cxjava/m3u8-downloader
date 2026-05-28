package utils

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func concatFiles(dir string, output string, del bool) string {
	files, err := filepath.Glob(filepath.Join(dir, "*.ts"))
	check(err)

	outPath := filepath.Join(dir, output)
	outFile, err := os.Create(outPath)
	check(err)
	defer func() {
		check(outFile.Close())
	}()

	for _, file := range files {
		inFile, err := os.Open(file)
		check(err)
		_, err = io.Copy(outFile, inFile)
		closeErr := inFile.Close()
		check(err)
		check(closeErr)
		if del {
			check(os.Remove(file))
		}
	}

	return output
}

//windows合并文件
func WinMergeFile(path string, del bool) string {
	log.Info("Copy all ts files to merged.tmp")
	return concatFiles(path, "merged.tmp", del)
}

//unix合并文件
func UnixMergeFile(path string, del bool) string {
	log.Info("Copy all ts files to merged.tmp")
	return concatFiles(path, "merged.tmp", del)
}

func FFmpegMergeFile(path string, del bool) string {
	if err := exec.Command("ffmpeg", "-L").Run(); err != nil {
		log.Warn("Check ffmpeg failed, fallback to merge by copy")
		return concatFiles(path, "merged.tmp", del)
	}

	files, err := filepath.Glob(filepath.Join(path, "*.ts"))
	check(err)

	var list strings.Builder
	for _, file := range files {
		list.WriteString("file '")
		list.WriteString(filepath.Base(file))
		list.WriteString("'\n")
	}

	listPath := filepath.Join(path, "templist.txt")
	check(os.WriteFile(listPath, []byte(list.String()), 0644))

	cmd := exec.Command("ffmpeg", "-f", "concat", "-i", "templist.txt", "-c", "copy", "merged.mp4")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		log.Warn("Check ffmpeg failed, fallback to merge by copy")
		check(os.Remove(listPath))
		return concatFiles(path, "merged.tmp", del)
	}

	if del {
		check(os.Remove(listPath))
		for _, file := range files {
			check(os.Remove(file))
		}
	}
	return "merged.mp4"
}
