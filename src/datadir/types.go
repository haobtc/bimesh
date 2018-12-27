package datadir

type Bind struct {
	Host string `json:"host,omniempty"`
	Port int `json:"port,omniempty"`
}

type Config struct {
	Version string `json:"version,omniempty"`
	Bind Bind `json:"bind,omniempty"`
}

type BoxRef struct {
	Boxid string `json:"boxid"`
	Endpoint string `json:"endpoint"`
	Cert string `json:"cert,omniempty"`
	ServiceNames []string `json:"service_names"`
}

type StaticRouter struct {
	Version string `json:"version,omniempty"`
	Boxes []BoxRef `json:"boxes"`
}
