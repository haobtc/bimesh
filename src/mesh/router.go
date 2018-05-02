package mesh

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

func NewRouter() *Router {
	return new(Router).Init()
}

func GetConnId(c *websocket.Conn) string {
	return c.UnderlyingConn().RemoteAddr().String()
}

func RemoveElement(slice []CID, elems CID) []CID {
	for i := range slice {
		if slice[i] == elems {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (self *Router) Init() *Router {
	self.ChJoin = make(chan JoinCommand, 1000)
	self.ChLeave = make(chan LeaveCommand, 1000)
	self.ChMsg = make(MsgChannel, 10000)
	self.ChBroadcast = make(MsgChannel, 10000)

	self.serviceLock = new(sync.RWMutex)
	self.ServiceConnMap = make(map[string]([]CID))
	self.ConnServiceMap = make(map[CID]([]string))
	self.ConnMap = make(map[CID]ConnT)
	return self
}

func (self *Router) registerConn(connId CID, ch MsgChannel, intent string) {
	self.ConnMap[connId] = ConnT{RecvChannel: ch, Intent: intent}
	// register connId as a service name
	//self.RegisterService(connId, connId)
}

func (self *Router) RegisterService(connId CID, serviceName string) error {
	self.serviceLock.Lock()
	defer self.serviceLock.Unlock()

	// bi direction map
	cidArr, ok := self.ServiceConnMap[serviceName]
	if ok {
		cidArr = append(cidArr, connId)
	} else {
		var a []CID
		cidArr = append(a, connId)
	}
	self.ServiceConnMap[serviceName] = cidArr

	snArr, ok := self.ConnServiceMap[connId]
	if ok {
		snArr = append(snArr, serviceName)
	} else {
		var a []string
		snArr = append(a, serviceName)
	}
	self.ConnServiceMap[connId] = snArr

	return nil
}

func (self *Router) unregisterConn(connId CID) {
	self.ClearPending(connId)

	self.serviceLock.Lock()
	defer self.serviceLock.Unlock()

	serviceNames, ok := self.ConnServiceMap[connId]
	if ok {
		for _, serviceName := range serviceNames {
			connIds, ok := self.ServiceConnMap[serviceName]
			if !ok {
				continue
			}
			connIds = RemoveElement(connIds, connId)
			if len(connIds) > 0 {
				self.ServiceConnMap[serviceName] = connIds
			} else {
				delete(self.ServiceConnMap, serviceName)
			}
		}
		delete(self.ConnServiceMap, connId)
	}

	ct, ok := self.ConnMap[connId]
	if ok {
		delete(self.ConnMap, connId)
		close(ct.RecvChannel)
	}
}

func (self *Router) SelectConn(serviceName string) (CID, bool) {
	self.serviceLock.RLock()
	defer self.serviceLock.RUnlock()

	connIds, ok := self.ServiceConnMap[serviceName]
	if ok && len(connIds) > 0 {
		// or random or round-robin
		return connIds[0], true
	}
	return 0, false
}

func (self *Router) GetServices(connId CID) []string {
	self.serviceLock.RLock()
	defer self.serviceLock.RUnlock()
	return self.ConnServiceMap[connId]
}

func (self *Router) ClearTimeoutRequests() {
	now := time.Now()
	tmpMap := make(map[PendingKey]PendingValue)

	for pKey, pValue := range self.PendingMap {
		if now.After(pValue.Expire) {
			errMsg := NewErrorMessage(pKey.MsgId, 408, "request timeout")
			_ = self.deliverMessage(pKey.ConnId, errMsg)
		} else {
			tmpMap[pKey] = pValue
		}
	}
	self.PendingMap = tmpMap
}

func (self *Router) ClearPending(connId CID) {
	for pKey, pValue := range self.PendingMap {
		if pKey.ConnId == connId || pValue.ConnId == connId {
			delete(self.PendingMap, pKey)
		}
	}
}

func (self *Router) deletePending(pKey PendingKey) {
	delete(self.PendingMap, pKey)
}

func (self *Router) setPending(pKey PendingKey, pValue PendingValue) {
	self.PendingMap[pKey] = pValue
}

func (self *Router) routeMessage(msg RPCMessage) error {
	fromConnId := msg.FromConnId
	if msg.IsRequest() {
		toConnId, found := self.SelectConn(msg.ServiceName)
		if found {
			pKey := PendingKey{ConnId: fromConnId, MsgId: msg.Id}
			expireTime := time.Now().Add(DefaultRequestTimeout)
			pValue := PendingValue{ConnId: toConnId, Expire: expireTime}

			self.setPending(pKey, pValue)
			return self.deliverMessage(toConnId, msg)
		} else {
			errMsg := NewErrorMessage(msg.Id, 404, "service not found")
			return self.deliverMessage(fromConnId, errMsg)
		}
	} else if msg.IsNotify() {
		toConnId, found := self.SelectConn(msg.ServiceName)
		if found {
			return self.deliverMessage(toConnId, msg)
		} else {
			errMsg := NewErrorMessage(msg.Id, 404, "service not found")
			return self.deliverMessage(fromConnId, errMsg)
		}
	} else if msg.IsResultOrError() {
		for pKey, pValue := range self.PendingMap {
			if pKey.MsgId == msg.Id && pValue.ConnId == fromConnId {
				// delete key within a range loop is safe
				// refer to https://stackoverflow.com/questions/23229975/is-it-safe-to-remove-selected-keys-from-golang-map-within-a-range-loop
				self.deletePending(pKey)
				return self.deliverMessage(pKey.ConnId, msg)
			}
		} // end of for
	}
	return nil
}

func (self *Router) broadcastNotify(notify RPCMessage) error {
	if !notify.IsNotify() {
		errMsg := NewErrorMessage(notify.Id, 400, "only notify can be broadcasted")
		self.deliverMessage(notify.FromConnId, errMsg)
		return nil
	}
	for connId, connT := range self.ConnMap {
		if connT.Intent == "actor" {
			err := self.deliverMessage(connId, notify)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (self *Router) deliverMessage(connId CID, msg RPCMessage) error {
	ct, ok := self.ConnMap[connId]
	if ok {
		ct.RecvChannel <- msg
	}
	return nil
}

func (self *Router) Start() {
	for {
		select {
		case openCmd := <-self.ChJoin:
			self.registerConn(openCmd.ConnId, openCmd.Channel, openCmd.Intent)
		case msg := <-self.ChMsg:
			self.routeMessage(msg)
		case notify := <-self.ChBroadcast:
			self.broadcastNotify(notify)
		case closeCmd := <-self.ChLeave:
			self.unregisterConn(CID(closeCmd))
		}
	}
}

// commands
func (self *Router) RouteMessage(msg RPCMessage, fromConnId CID) {
	msg.FromConnId = fromConnId
	self.ChMsg <- msg
}

func (self *Router) BroadcastNotify(notify RPCMessage, fromConnId CID) {
	notify.FromConnId = fromConnId
	self.ChBroadcast <- notify
}

func (self *Router) Join(connId CID, ch MsgChannel, intent string) {
	self.ChJoin <- JoinCommand{ConnId: connId, Channel: ch, Intent: intent}
}

func (self *Router) Leave(connId CID) {
	self.ChLeave <- LeaveCommand(connId)
}
