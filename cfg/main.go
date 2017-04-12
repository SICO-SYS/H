/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	Config *Main
	err    error
)

type Main struct {
	Version    string     `json:"version"`
	Redis      Redis      `json:"redis"`
	Mysql      Mysql      `json:"mysql"`
	Mongo      Mongo      `json:"mongo"`
	Rabbit     Rabbit     `json:"rabbit"`
	OpenAccess OpenAccess `json:OpenAccess`
	Log        Log        `json:"log"`
}

func (c *Main) ReadConfig(path string) (*Main, error) {

	// We'll change to Call configAPI to parse the config
	// ioutil.ReadFile("https://cfg.siner.io/production")
	d, err := ioutil.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(d, &c)
	}

	return c, err
}

func init() {
	defer func() {
		if emsg := recover(); emsg != nil {
			log.Println(emsg)
		}
	}()
	Config, err = (&Main{}).ReadConfig("config.json")
	if err != nil {
		log.Println(err)
		log.Fatalln("Parse config.json problem, make sure file with correct syntax")
	}
	_, err = os.Stat(Config.Log.Logpath)
	if err != nil {
		err = os.MkdirAll(Config.Log.Logpath, 0666)
		if err != nil {
			log.Println(err)
			log.Panicln("Failed to create folder")
		}
	}
}
