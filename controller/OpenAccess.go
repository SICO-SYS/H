/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"golang.org/x/net/context"
	"net/http"
	// "time"
	// "math"

	"github.com/SiCo-DevOps/Pb"
	"github.com/SiCo-DevOps/dao"
	. "github.com/SiCo-DevOps/log"
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
	err = dao.RedisSetShort(key, config.OpenAccess.TokenValid, config.OpenAccess.TokenExpired)
	rspdata := &ResponseData{}
	if err != nil {
		rspdata = ResponseErrmsg(126)
		LogProduce("error", "Cannot Set key")
	} else {
		rspdata = &ResponseData{0, &OpenToken{key}}
		LogProduce("info", "Sucess")
	}
	rsp, _ := json.Marshal(rspdata)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(rsp)
}

func GetAPIToken(rw http.ResponseWriter, req *http.Request) {
	if AuthOpenToken(req) {
		// key := GenerateRand()
		// secret := GenerateRand()
		defer func() {
			if rcv := recover(); rcv != nil {
				LogProduce("error", "gRPC connect error")
			}
		}()
		cc := dao.RpcConn("He")
		defer cc.Close()
		c := pb.NewOpenClient(cc)
		r, err := c.RegUser(context.Background(), &pb.OpenRequest{"reg"})
		if err != nil {
			LogErrMsg(50, "controller.GetAPIToken")
		}
		rsp, _ := json.Marshal(&SecretToken{Key: r.Key, Secret: r.Secret})
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(rsp)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Permission Denied"))
	}

}

func AuthOpenToken(req *http.Request) bool {
	k := req.URL.Query().Get("token")
	data, err1, err2 := dao.RedisGetValue(k)
	if err1 != nil {
		LogErrMsg(1, "controller.AuthOpenToken")
		return false
	}
	if err2 != nil {
		LogErrMsg(19, "controller.AuthOpenToken")
		return false
	}
	ok, err := dao.RedisBool(data)
	if err != nil {
		LogErrMsg(11, "controller.AuthOpenToken")
		return false
	}
	return ok
}

// func printTS(rw http.ResponseWriter, req *http.Request) {
// 	 := int64(math.Floor(float64(time.Now().Unix() / 30)))

// }
