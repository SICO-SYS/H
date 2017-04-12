package route

import (
	"H/controller"
)

func OpenAPI() {
	Handler.HandleFunc("/open/config", controller.GetCfgVersion).Methods("GET")
	Handler.HandleFunc("/open/Token", controller.GetOpenToken).Methods("GET")
	Handler.HandleFunc("/open/APIToken", controller.GetSecretToken).Methods("GET")
}
