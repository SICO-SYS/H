/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package log

import (
	"H/cfg"
	"log"
	"os"
)

var (
	filename    string
	warningFile = cfg.Config.Log.WARNING
	errorFile   = cfg.Config.Log.ERROR
	fatalFile   = cfg.Config.Log.FATAL
	logFileDir  = cfg.Config.Log.Logpath
)

func WriteLog(level string, msg string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("File cannot Write in " + filename)
		}
	}()
	switch level {
	case "warning":
		filename = warningFile
	case "error":
		filename = errorFile
	case "fatal":
		filename = fatalFile
	default:
		filename = "unknown.log"
	}
	fd, err := os.OpenFile(logFileDir+filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	defer fd.Close()
	if err != nil {
		log.Panicln(filename + " :Open log file Failed")
	}

	logger := log.New(fd, "[SiCo]", log.Lshortfile)
	logger.Println(msg)
}
