package mesh

import (
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type CID uint64

type RPCMessage struct {
	Initialized bool
	FromConnId  CID
	Id          interface{}
	ServiceName string
	Method      string
	Params      *simplejson.Json
	Result      *simplejson.Json
	Error       *simplejson.Json
	Raw         *simplejson.Json
}

// 5 seconds
const DefaultRequestTimeout time.Duration = 1000000 * 5

// Commands
type MsgChannel chan RPCMessage

type JoinCommand struct {
	ConnId  CID
	Channel MsgChannel
	Intent  string
}

type LeaveCommand CID

// Pending Struct
type PendingKey struct {
	ConnId CID
	MsgId  interface{}
}

type PendingValue struct {
	ConnId CID
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
	ServiceConnMap map[string]([]CID)
	ConnServiceMap map[CID]([]string)

	ConnMap    map[CID](ConnT)
	PendingMap map[PendingKey]PendingValue
}

// An ConnActor manage a websocket connection and handles incoming messages
type Actor struct {
	ChMsg  MsgChannel
	ConnId CID
	Conn   *websocket.Conn
}

type Requester struct {
	ChMsg  MsgChannel
	ConnId CID
}

// builtin services
type ServiceManager struct {
	ChMsg  MsgChannel
	ConnId CID
}

type ContextT struct {
	Router         *Router
	ServiceManager *ServiceManager
}
