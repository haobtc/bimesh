package static

import (
	"bytes"
	"net/http"
	"mesh"
	"jsonrpc"
	"datadir"
)

func FromConfig(config datadir.StaticEndpointConfig) StaticEndpoint {
	endpoint := StaticEndpoint{}
	endpoint.Url = config.Url
	endpoint.Cert = config.Cert
	endpoint.ServiceType = config.ServiceType
	endpoint.ServiceInfix = config.ServiceInfix
	endpoint.serviceNames = config.ServiceNames
	return endpoint
}

func (self StaticEndpoint) GetId() string {
	return self.Url
}

func (self StaticEndpoint) GetServiceNames() []string {
	return self.serviceNames
}

func (self StaticEndpoint) Request(msg jsonrpc.RPCMessage) (jsonrpc.RPCMessage, error) {
	data, err := msg.Raw.MarshalJSON()
	if err != nil {
		return jsonrpc.RPCMessage{}, err
	}
	req, err := http.NewRequest("POST", self.Url, bytes.NewBuffer(data))
	req.Header.Set("X-Framework", "bimesh")
	// TODO: other headers
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return jsonrpc.RPCMessage{}, err
	}
	// TODO: check response headers
	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return jsonrpc.RPCMessage{}, err
	}
	respMsg, err := jsonrpc.ParseMessage(buffer.Bytes())
	if err != nil {
		return jsonrpc.RPCMessage{}, err
	}
	return respMsg, err
}

func JoinMesh() {
	mesh := mesh.GetMesh()
	config := datadir.GetConfig()
	for _, epConfig := range config.StaticEndpoints {
		endpoint := FromConfig(epConfig)
		mesh.Join(endpoint)
	}
}
