package datadir

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

var router Router
var routerParsed bool = false

func GetRouter() Router {
	if !routerParsed {
		err := router.Parse()
		if err != nil {
			panic(err)
		}
		routerParsed = true
	}
	return router
}

func (self *Router) Parse() (err error) {
	routerPath := DataPath("router.json")
	if _, err := os.Stat(routerPath); os.IsNotExist(err) {
		self.FillDefaultValues()
		return nil
	}
	data, err := ioutil.ReadFile(routerPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, self)
	if err != nil {
		return err
	}
	self.FillDefaultValues()
	return nil
}

func (self *Router) FillDefaultValues() {
	if self.Version == "" {
		self.Version = "1.0"
	}
}

