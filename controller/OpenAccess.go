/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"net/http"

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
	defer func() {
		recover()
		LogProduce("error", "gRPC connect error")
	}()
	if AuthOpenToken(req) {
		fmt.Println("1")
		cc := dao.RpcConn("He")
		fmt.Println("1")
		defer cc.Close()
		fmt.Println("1")
		c := pb.NewAAA_OpenClient(cc)
		fmt.Println("1")
		r, err := c.AAA_RegUser(context.Background(), &pb.AAA_OpenRequest{"reg"})
		fmt.Println("1")
		if err != nil {
			LogErrMsg(50, "controller.GetAPIToken")
		}
		if r != nil {
			rsp, _ := json.Marshal(&SecretToken{Key: r.Key, Secret: r.Secret})
			rw.Header().Add("Content-Type", "application/json")
			rw.Write(rsp)
			return
		}
		rsp, _ := json.Marshal(&SecretToken{Key: "", Secret: ""})
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
