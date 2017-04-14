/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	// "math"
	"net/http"
	// "time"

	. "github.com/SiCo-DevOps/H/log"
)

type OpenToken struct {
	Key string `json:"key"`
}

type SecretToken struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type TransMessage struct {
	Key    string `json:"key"`
	Token  string `json:"token"`
	Action string `json:"action"`
}

func GetOpenToken(rw http.ResponseWriter, req *http.Request) {
	key := GenerateRand()
	err = RedisSetShort(key, config.OpenAccess.TokenValid, config.OpenAccess.TokenExpired)
	rspdata := &ResponseData{}
	if err != nil {
		rspdata = ResponseErrmsg(126)
		WriteLog("error", "Cannot Set key")
	} else {
		rspdata = &ResponseData{0, &OpenToken{key}}
		WriteLog("info", "Sucess")
	}
	rsp, _ := json.Marshal(rspdata)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(rsp)
}

func GetAPIToken(rw http.ResponseWriter, req *http.Request) {
	if AuthOpenToken(req) {
		key := GenerateRand()
		secret := GenerateRand()
		// rsconn := RedisPool.Get()
		// defer rsconn.Close()
		// rsconn.Do("SET", key, secret)
		rsp, _ := json.Marshal(&SecretToken{Key: key, Secret: secret})
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(rsp)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Permission Denied"))
	}

}

func AuthOpenToken(req *http.Request) bool {
	k := req.URL.Query().Get("token")
	data, err1, err2 := RedisGetValue(k)
	if err1 != nil {
		WriteLog("error", "AuthOpenToken: connection error")
		// WriteLog("error", err1.Error())
		return false
	}
	if err2 != nil {
		WriteLog("error", "AuthOpenToken: Cannot Exec GET，I cannot procedure this error, maybe a large value")
		// WriteLog("error", err2.Error())
		return false
	}
	ok, err := RedisBool(data)
	if err != nil {
		WriteLog("error", "AuthOpenToken: Key parse error")
		// WriteLog("error", err.Error())
		return false
	}
	return ok
}

// func printTS(rw http.ResponseWriter, req *http.Request) {
// 	 := int64(math.Floor(float64(time.Now().Unix() / 30)))

// }
