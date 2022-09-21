package utils

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func MkAllDir(dirs string) {
	if !IsExist(dirs) {
		err := os.MkdirAll(dirs, os.ModePerm)
		check(err)
	}
}

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func RemoveAllDir(dirs string) {
	os.RemoveAll(strings.Split(dirs, "/")[0])
}

func IsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		log.Trace("Warn: " + filePath + " is not exist, error: " + err.Error())
		return false
	}
	return true
}

func WriteFile(fileName string, data []byte) {
	err := os.WriteFile(fileName, data, os.ModePerm)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Error("error :" + err.Error())
		panic(err)
	}
}

func FileAbs(path string) string {
	pwd, _ := os.Getwd()
	if absPath, err := filepath.Abs(filepath.Join(pwd, path)); err == nil {
		return absPath
	}
	return ""
}
