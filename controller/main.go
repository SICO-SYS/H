/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"github.com/getsentry/raven-go"

	"github.com/SiCo-Ops/cfg"
)

var (
	config  = cfg.Config
	RPCAddr = map[string]string{
		"He": "He.SiCo" + config.RPCPort.He,
		"Li": "Li.SiCo" + config.RPCPort.Li,
		"Be": "Be.SiCo" + config.RPCPort.Be,
		"B":  "B.SiCo" + config.RPCPort.B,
		"C":  "C.SiCo" + config.RPCPort.C,
		"N":  "N.SiCo" + config.RPCPort.N,
	}
)

func init() {
	if config.Sentry.Enable {
		raven.SetDSN(config.Sentry.DSN)
	}
}
