/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/gorilla/mux"
)

var (
	HTTPHandler *mux.Router // Define HTTPHandler for http handler
)

func init() {
	HTTPHandler = mux.NewRouter()
	HTTPHandler.StrictSlash(true)
	Public()
	AAA()
	Cloud()
	Asset()
	Config()
}
