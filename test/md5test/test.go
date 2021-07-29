package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func main() {
	fmt.Println(GetMd5String("1"))
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
