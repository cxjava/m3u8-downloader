package utils

import (
	"encoding/hex"

	log "github.com/sirupsen/logrus"
)

func HexDecode(s string) []byte {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Error("HexDecode failed:" + err.Error())
		return nil
	}
	log.Tracef("HexDecode %s", decoded)
	return decoded
}

func HexEncode(s []byte) string {
	return hex.EncodeToString(s)
}
