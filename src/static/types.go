package static

type StaticEndpoint struct {
	Url string
	Cert string
	serviceNames []string

	// for future usages
	ServiceType string

	// infix between service name and service method, default is '::'
	ServiceInfix string
}
