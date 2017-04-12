/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"fmt"
	"net/http"
	"os"
)

func GetCfgVersion(rw http.ResponseWriter, req *http.Request) {
	// rw.Header().Add("content-type", "application/json")
	rw.Write([]byte("[Success] config version  === " + config.Version))
}

func GenerateRand() string {
	data, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	defer data.Close()
	buf := make([]byte, 16)
	data.Read(buf)
	v := fmt.Sprintf("%X", buf)
	return v
}

func Sha256Encrypt(v interface{}) string {
	return "hello"
}
