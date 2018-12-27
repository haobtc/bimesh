package mesh

import (
	"sync"
	"jsonrpc"
)

type EndpointSource int
type MsgChannel chan jsonrpc.RPCMessage

const (
	SourceTentacle EndpointSource = 1 + iota
	SourceStatic
	SourceBbox
)

type EndpointRef struct {
	Source EndpointSource
	Url string
}

type Endpoint struct {
	EndpointRef
	
	Cert string
	ServiceType string
	ServiceInfix string
	ServiceNames []string

	ChMsg MsgChannel
}

type Mesh struct {
	serviceLock *sync.RWMutex
	urlEndpointMap map[string](*Endpoint)
	serviceEndpointsMap map[string]([]*Endpoint)
}
