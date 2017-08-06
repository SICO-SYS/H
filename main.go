/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package main

import (
	"net/http"

	"github.com/SiCo-Ops/H/route"
)

func Run() {
	http.ListenAndServe("0.0.0.0:2048", route.Handler)
}

func main() {
	Run()
}
