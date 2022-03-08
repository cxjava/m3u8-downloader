package utils

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func ExecUnixShell(s string) error {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Error("ExecUnixShell Error: " + err.Error())
		return err
	}
	outStr := out.String()
	if len(outStr) > 0 {
		log.Info(outStr)
	}
	return nil
}

func ExecWinShell(s string) error {
	cmd := exec.Command("cmd", "/C", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Error("ExecWinShell Error: " + err.Error())
		return err
	}
	outStr := out.String()
	if len(outStr) > 0 {
		log.Info(outStr)
	}
	return nil
}

//windows合并文件
func WinMergeFile(path string, del bool) {
	err := os.Chdir(path)
	check(err)
	log.Info("Copy all ts files to merged.tmp")
	err = ExecWinShell("copy /b *.ts merged.tmp")
	check(err)
	if del {
		log.Warn("Delete all ts files")
		err = ExecWinShell("del /Q *.ts")
		check(err)
	}
}

//unix合并文件
func UnixMergeFile(path string, del bool) {
	err := os.Chdir(path)
	check(err)
	log.Info("Copy all ts files to merged.tmp")
	cmd := `cat *.ts >> merged.tmp`
	err = ExecUnixShell(cmd)
	check(err)
	if del {
		log.Warn("Delete all ts files")
		err = ExecUnixShell("rm -rf *.ts")
		check(err)
	}
}

func FFmpegMergeFile(path string, del bool) {
	err := os.Chdir(path)
	check(err)

	// generate list file then invoke ffmpeg
	// https://trac.ffmpeg.org/wiki/Concatenate
	switch runtime.GOOS {
	case "windows":
		err = ExecWinShell("(for %i in (*.ts) do @echo file '%i') > templist.txt")
		check(err)
		err = ExecWinShell("ffmpeg -f concat -i templist.txt -c copy merged.mp4")
		check(err)
		if del {
			log.Warn("Delete all ts files")
			err = ExecWinShell("del /Q *.ts")
			check(err)
		}
	default:
		err = ExecUnixShell("for f in *.ts; do echo \"file '$f'\" >> templist.txt; done")
		check(err)
		err = ExecUnixShell("ffmpeg -f concat -i templist.txt -c copy merged.mp4")
		check(err)
		if del {
			log.Warn("Delete all ts files")
			err = ExecUnixShell("rm -rf *.ts")
			check(err)
		}
	}
}
