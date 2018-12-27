package mesh

import (
	"sync"
	"jsonrpc"
)

type EndpointSource int

/*const (
	SourceTentacle EndpointSource = 1 + iota
	SourceStatic
	SourceBbox
)
*/
type Endpoint interface {
	GetId() string
	GetServiceNames() []string
	Request(msg jsonrpc.RPCMessage) (jsonrpc.RPCMessage, error)
}

type Mesh struct {
	serviceLock *sync.RWMutex
	idEndpointMap map[string](Endpoint)
	serviceEndpointsMap map[string]([]Endpoint)
}
