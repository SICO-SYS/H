/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package main

import (
	"net/http"

	"github.com/SiCo-DevOps/H/route"
	"github.com/SiCo-DevOps/dao"
)

func Run() {
	http.ListenAndServe("0.0.0.0:2048", route.Handler)
}

func main() {
	defer func() { recover(); dao.MgoUserConn.Close() }()
	Run()
}
