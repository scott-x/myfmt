package util

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
)

func MD5(file string) string {
	bs, _ := ioutil.ReadFile(file)
	return fmt.Sprintf("%x", md5.Sum(bs))
}
