package tentacle

import (
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
)

func (self *LocalConnT) Init(conn *websocket.Conn) *LocalConnT {
	self.Conn = conn
	self.ChMsg = make(MsgChannel, 100)
	//self.ConnId = CID(uuid.Must(uuid.NewV4()).String())
	self.ConnId = GetCID()
	return self
}

func (self LocalConnT) RecvChannel() MsgChannel {
	return self.ChMsg
}

func (self LocalConnT) CanBroadcast() bool {
	return true
}

func (self *LocalConnT) Close() {
	Tentacle().Router.Leave(self.ConnId)
}

func (self *LocalConnT) Start() {
	// register connection
	Tentacle().Router.JoinConn(self.ConnId, self)

	for {
		select {
		case msg, more := <- self.ChMsg:
			if more {
				if writeErr := self.writeJSON(msg.Raw); writeErr != nil {
					return
				}
			} else {
				// log.
				return
			}
		}
	}
}

func (self *LocalConnT) writeJSON(data *simplejson.Json) error {
	// send to self
	bytes, err := data.MarshalJSON()
	if err != nil {
		return err
	}
	self.Conn.WriteMessage(websocket.TextMessage, bytes)
	return nil
}
