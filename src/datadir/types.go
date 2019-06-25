package datadir

type Bind struct {
	Host string `yaml:"host,omitempty"`
	Port int `yaml:"port,omitempty"`
}

type Config struct {
	Version string `yaml:"version,omitempty"`
	Bind Bind `yaml:"bind,omitempty"`
}
