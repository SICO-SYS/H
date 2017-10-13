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

	"github.com/SiCo-Ops/cloud-go-sdk/wechat"
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
	// token := config.WechatToken
	// nonce, timestamp, signature := wechat.GetValidation(req)
	// isValid := wechat.ValidateServer(token, nonce, timestamp, signature)
	// if !isValid {
	// 	httpResponse("json", rw, responseErrMsg(10))
	// 	return
	// }
	data, _ := ioutil.ReadAll(req.Body)
	v := wechat.Parse(data)
	var (
		msgtype string
		content string
	)
	if v.MsgType == "event" && v.Event == "subscribe" {
		msgtype = "text"
		content = "Welcome to use SiCo \n Type #signup to registry \n Type #Signin TOKENIN SIGNATURE to bind an exist token"
	}
	if v.MsgType == "text" {
		// msgtype = "text"
		command := strings.Split(v.Content, " ")
		switch command[0] {
		case "?":
			content = "Type #signup to registry \n Type #Signin TOKENIN SIGNATURE to bind an exist token"
		default:
			content = "Command error. Type ? for help"
		}
	}
	rsp := wechat.Marshal(v.FromUserName, v.ToUserName, public.CurrentTimeStamp(), msgtype, content)
	httpResponse("xml", rw, rsp)
	return
}
