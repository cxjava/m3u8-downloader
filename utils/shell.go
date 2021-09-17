package utils

import (
	"bytes"
	"os"
	"os/exec"

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
