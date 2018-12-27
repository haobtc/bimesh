package datadir

import (
	"os"
	"io/ioutil"
	"errors"

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

func (self *Config) ParseConfig() error {
	cfgPath := DataPath("config.yml")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		err = self.ValidateValues()
		if err != nil {
			return err
		}
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
	return self.ValidateValues()
}

func (self *Config) ValidateValues() error {
	if self.Version == "" {
		self.Version = "1.0"
	}
	if self.Bind.Host == "" {
		self.Bind.Host = "127.0.0.1"
	}

	if self.Bind.Port <= 0 || self.Bind.Port > 65535 {
		// should be a legal port
		self.Bind.Port = 18666
	}

	for _, endp := range self.StaticEndpoints {
		if endp.ServiceType == "" {
			endp.ServiceType = "jsonrpc"
		} else if endp.ServiceType != "jsonrpc" {
			// currently we only support jsonrpc
			return errors.New("endpoint service type is not jsonrpc")
		}

		if endp.ServiceInfix == "" {
			endp.ServiceInfix = "::"
		}
	}
	return nil
}

