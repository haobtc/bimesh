package bbox

import (
	"time"
	"net/http"
	"bytes"
	"jsonrpc"
)

func (self BboxEndpoint) GetId() string {
	return self.BoxInfo.Bind
}

func (self BboxEndpoint) GetServiceNames() []string {
	return self.BoxInfo.ServiceNames
}

func (self *BboxEndpoint) Request(msg jsonrpc.RPCMessage) (jsonrpc.RPCMessage, error) {
	data, err := msg.Raw.MarshalJSON()
	if err != nil {
		return jsonrpc.RPCMessage{}, err
	}
	req, err := http.NewRequest("POST", self.BoxInfo.Bind, bytes.NewBuffer(data))
	req.Header.Set("X-Framework", "bimesh")
	// TODO: other headers
	client := &http.Client{Timeout: 10 * time.Second}
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

