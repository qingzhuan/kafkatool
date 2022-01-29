package common

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)


func MD5Value(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func DirIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err){
			return false
		}
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}