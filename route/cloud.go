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
	v1.HandleFunc("/{cloud}/{bsns}/sync", controller.Cloud_SyncResourse).Methods("GET")
	v1.HandleFunc("/{cloud}/{bsns}", controller.Cloud_Call).Methods("POST")
}
