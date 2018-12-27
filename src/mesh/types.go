package mesh

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"jsonrpc"
)


// 5 seconds
const DefaultRequestTimeout time.Duration = 1000000 * 5

// Commands
type MsgChannel chan jsonrpc.RPCMessage

type JoinCommand struct {
	ConnId  jsonrpc.CID
	Channel MsgChannel
	Intent  string
}

type LeaveCommand jsonrpc.CID

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
type ConnT struct {
	RecvChannel MsgChannel
	Intent      string
}

type Router struct {
	// channels
	ChJoin      chan JoinCommand
	ChLeave     chan LeaveCommand
	ChMsg       MsgChannel
	ChBroadcast MsgChannel

	serviceLock    *sync.RWMutex
	ServiceConnMap map[string]([]jsonrpc.CID)
	ConnServiceMap map[jsonrpc.CID]([]string)

	ConnMap    map[jsonrpc.CID](ConnT)
	PendingMap map[PendingKey]PendingValue
}

// An ConnActor manage a websocket connection and handles incoming messages
type Actor struct {
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

type ContextT struct {
	Router         *Router
	ServiceManager *ServiceManager
}
