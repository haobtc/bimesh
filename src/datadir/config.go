package datadir

import (
	"os"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var cfg Config
var cfgParsed bool = false

func GetConfig() Config {
	if !cfgParsed {
		err := cfg.ParseConfig()
		if err != nil {
			panic(err)
		}
		cfgParsed = true
	}
	return cfg
}

func (self *Config) ParseConfig() (err error) {
	cfgPath := DataPath("config.yml")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		self.FillDefaultValues()
		return nil
	}
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, self)
	if err != nil {
		return err
	}
	self.FillDefaultValues()
	return nil
}

func (self *Config) FillDefaultValues() {
	if self.Version == "" {
		self.Version = "1.0"
	}
	if self.Bind.Host == "" {
		self.Bind.Host = "127.0.0.1"
	}
	if self.Bind.Port <= 0 || self.Bind.Port > 65535 {
		self.Bind.Port = 18666
	}
}

