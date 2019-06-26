package tentacle

import (
	"github.com/gorilla/websocket"
	"errors"
	"sync"
	"time"
	"jsonrpc"
)


// 5 seconds
const (
	DefaultRequestTimeout time.Duration = 1000000 * 5
	
	IntentLocal string = "local"
)

var (
	ErrNotNotify = errors.New("json message is not notify")
)

// Commands
type MsgChannel chan jsonrpc.RPCMessage

// Pending Struct
type PendingKey struct {
	ConnId jsonrpc.CID
	MsgId  interface{}
}

type PendingValue struct {
	ConnId jsonrpc.CID
	Expire time.Time
}

// ConnT
type ConnT interface {
	RecvChannel() MsgChannel
	CanBroadcast() bool
}


type Router struct {
	// channels
	serviceLock    *sync.RWMutex
	ServiceConnMap map[string]([]jsonrpc.CID)
	ConnServiceMap map[jsonrpc.CID]([]string)

	ConnMap    map[jsonrpc.CID](ConnT)
	PendingMap map[PendingKey]PendingValue
}

// An ConnActor manage a websocket connection and handles incoming messages
type LocalConnT struct {
	ChMsg  MsgChannel
	ConnId jsonrpc.CID
	Conn   *websocket.Conn
}

type Requester struct {
	ChMsg  MsgChannel
	ConnId jsonrpc.CID
}

// builtin services
type ServiceManager struct {
	ChMsg  MsgChannel
	ConnId jsonrpc.CID
}

type TentacleT struct {
	Router         *Router
	ServiceManager *ServiceManager
}
