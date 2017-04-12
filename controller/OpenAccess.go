/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"H/connection"
	"encoding/json"
	// "math"
	"net/http"
	// "time"
)

var ()

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
		rspdata = &ResponseData{126, "sys err"}
	} else {
		rspdata = &ResponseData{0, &OpenToken{key}}
	}
	rsp, _ := json.Marshal(rspdata)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(rsp)

}

func GetSecretToken(rw http.ResponseWriter, req *http.Request) {
	if AuthOpenToken(req) {
		key := GenerateRand()
		secret := GenerateRand()
		rsconn := RedisPool.Get()
		defer rsconn.Close()
		rsconn.Do("SET", key, secret)
		rsp, _ := json.Marshal(&SecretToken{Key: key, Secret: secret})
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(rsp)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Permission Denied"))
	}

}

func AuthOpenToken(req *http.Request) bool {
	key := req.URL.Query().Get("token")
	rsconn := RedisPool.Get()
	defer rsconn.Close()
	data, _ := rsconn.Do("GET", key)
	ok, _ := connection.GetRedisValue(data)
	return ok
}

// func printTS(rw http.ResponseWriter, req *http.Request) {
// 	 := int64(math.Floor(float64(time.Now().Unix() / 30)))

// }
