package datadir

type Bind struct {
	Host string `yaml:"host,omitempty"`
	Port int `yaml:"port,omitempty"`
}

type StaticEndpointConfig struct {
	Url string `yaml:"url"`
	Cert string `yaml:"cert,omitempty"`

	// enum{jsonrpc}, currently only jsonrpc is allowd, reserved
	// for future usages
	ServiceType string `yaml:"service_type,omitempty"`

	// infix between service name and service method, default is '::'
	ServiceInfix string `yaml:"service_infix,omitempty"`

	// list serice names
	ServiceNames []string `yaml:"service_names"`
}

type Config struct {
	Version string `yaml:"version,omitempty"`
	Bind Bind `yaml:"bind,omitempty"`
	StaticEndpoints []StaticEndpointConfig `yaml:"static_endpoints,flow"`
}
