/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/SiCo-DevOps/H/controller"
)

func Cloud() {
	v1 := Handler.PathPrefix("/v1/cloud").Subrouter()
	v1.HandleFunc("/{cloud}/{bsns}/op/{action}", controller.Cloud_call).Methods("GET")
	v1.HandleFunc("/{cloud}/{bsns}/template", controller.Cloud_template).Methods("GET")
	v1.HandleFunc("/{cloud}/{bsns}", controller.Cloud_rawCall).Methods("POST")
}
