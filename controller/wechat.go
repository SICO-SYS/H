/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/cloud-go-sdk/wechat"
	"github.com/SiCo-Ops/dao/grpc"
	"github.com/SiCo-Ops/public"
)

func WechatValidateServer(rw http.ResponseWriter, req *http.Request) {
	echostr := req.URL.Query().Get("echostr")
	token := config.WechatToken
	nonce, timestamp, signature := wechat.GetValidation(req)
	isValid := wechat.ValidateServer(token, nonce, timestamp, signature)
	if !isValid {
		httpResponse("json", rw, responseErrMsg(10))
		return
	}
	rw.Write([]byte(echostr))
	return
}

func WechatReceiveMessage(rw http.ResponseWriter, req *http.Request) {
	token := config.WechatToken
	nonce, timestamp, signature := wechat.GetValidation(req)
	isValid := wechat.ValidateServer(token, nonce, timestamp, signature)
	if !isValid {
		httpResponse("json", rw, responseErrMsg(10))
		return
	}
	data, _ := ioutil.ReadAll(req.Body)
	v := wechat.Parse(data)
	var (
		msgtype string
		content string
	)
	if v.MsgType == "event" && v.Event == "subscribe" {
		msgtype = "text"
		content = "Welcome to use SiCo \nType ? for help"
	}
	if v.MsgType == "voice" {
		msgtype = "text"
		command := v.Recognition
		if strings.Contains(command, "注册") {
			in := &pb.AAAGenerateTokenCall{Email: v.FromUserName + "@wechat"}
			cc, err := rpc.Conn(config.RpcHeHost, config.RpcHePort)
			if err != nil {
				content = errorMessage(301)
			} else {
				r := rpc.AAATokenGenerateRPC(cc, in)
				if r.Code != 0 {
					content = errorMessage(r.Code)
				} else {
					content = "openid:\n" + v.FromUserName + "\n\n" + "SecretID:\n" + r.Id + "\n\n" + "SecretKey:\n" + r.Key + "\n\n" + "Save this Info and delete this message for safe"
				}
			}
		}
	}
	if v.MsgType == "text" {
		msgtype = "text"
		command := strings.Split(v.Content, " ")
		switch command[0] {
		case "?":
			content = "#signup" + "\n\n" + "#Signin TOKENIN SIGNATURE" + "\n\n" + "Documentation at https://docs.sico.io"
		case "？":
			content = "#signup" + "\n\n" + "#Signin TOKENIN SIGNATURE" + "\n\n" + "Documentation at https://docs.sico.io"
		case "#signup":
			in := &pb.AAAGenerateTokenCall{Email: v.FromUserName + "@wechat"}
			cc, err := rpc.Conn(config.RpcHeHost, config.RpcHePort)
			if err != nil {
				content = errorMessage(301)
				break
			}
			r := rpc.AAATokenGenerateRPC(cc, in)
			if r.Code != 0 {
				content = errorMessage(r.Code)
			} else {
				content = "openid:\n" + v.FromUserName + "\n\n" + "SecretID:\n" + r.Id + "\n\n" + "SecretKey:\n" + r.Key + "\n\n" + "Save this Info and delete this message for safe"
			}
		default:
			content = "Command error. Type ? for help"
		}
	}
	rsp := wechat.Marshal(v.FromUserName, v.ToUserName, public.CurrentTimeStamp(), msgtype, content)
	httpResponse("xml", rw, rsp)
	return
}
