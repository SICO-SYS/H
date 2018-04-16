/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/SiCo-Ops/H/controller"
)

func Wechat() {
	v1 := HTTPHandler.PathPrefix("/v1/wechat").Subrouter()
	v1.Path("/").HandlerFunc(controller.WechatValidateServer).Methods("GET")
	v1.Path("/").HandlerFunc(controller.WechatReceiveMessage).Methods("POST")
}
