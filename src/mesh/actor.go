package mesh

import (
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
	//"github.com/satori/go.uuid"
)

func (self *Actor) Init(conn *websocket.Conn) *Actor {
	self.Conn = conn
	self.ChMsg = make(MsgChannel, 100)
	//self.ConnId = CID(uuid.Must(uuid.NewV4()).String())
	self.ConnId = GetCID()
	return self
}

func (self *Actor) Close() {
	Context().Router.Leave(self.ConnId)
}

func (self *Actor) Start() {
	// register connection
	Context().Router.Join(self.ConnId, self.ChMsg, "actor")

	for {
		select {
		case msg, more := <-self.ChMsg:
			if more {
				self.writeJSON(msg.Raw)
			} else {
				// log.
				return
			}
		}
	}
}

func (self *Actor) writeJSON(data *simplejson.Json) error {
	// send to self
	bytes, err := data.MarshalJSON()
	if err != nil {
		return err
	}
	self.Conn.WriteMessage(websocket.TextMessage, bytes)
	return nil
}
