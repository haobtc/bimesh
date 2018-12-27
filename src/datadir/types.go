package datadir

type Bind struct {
	Host string `yaml:"host,omitempty"`
	Port int `yaml:"port,omitempty"`
}


type BoxRef struct {
	Boxid string `yaml:"boxid"`
	Endpoint string `yaml:"endpoint"`
	Cert string `yaml:"cert,omitempty"`
	ServiceNames []string `yaml:"service_names"`
}

type Config struct {
	Version string `yaml:"version,omitempty"`
	Bind Bind `yaml:"bind,omitempty"`
	StaticBoxes []BoxRef `yaml:"static_boxes,omitempty"`
}
