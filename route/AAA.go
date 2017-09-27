/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/SiCo-Ops/H/controller"
)

func AAA() {
	v1 := HTTPHandler.PathPrefix("/v1/AAA").Subrouter()
	v1.Path("/token").HandlerFunc(controller.AAAGenerateToken).Methods("POST")
	v1.Path("/authentication").HandlerFunc(controller.AAAAuthentication).Methods("POST")
}
