package mesh

import (
	"errors"
)

func (self *Requester) Init() *Requester {
	self.ChMsg = make(MsgChannel, 100)
	self.ConnId = GetCID()
	return self
}

func (self *Requester) Close() {
	Context().Router.Leave(self.ConnId)
}

func (self *Requester) RequestAndWait(msg RPCMessage) (RPCMessage, error) {
	// register connection
	Context().Router.Join(self.ConnId, self.ChMsg, "requester")
	Context().Router.RouteMessage(msg, self.ConnId)

	for {
		select {
		case res, more := <-self.ChMsg:
			if more {
				if res.IsResultOrError() {
					return res, nil
				}
			} else {
				// log.
				return RPCMessage{}, errors.New("connection closed")
			}
		}
	}
	return RPCMessage{}, nil
}
