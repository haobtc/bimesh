package tentacle

import (
	"sync"
	"time"
	"github.com/gorilla/websocket"
	"jsonrpc"
)

func NewRouter() *Router {
	return new(Router).Init()
}

func GetConnId(c *websocket.Conn) string {
	return c.UnderlyingConn().RemoteAddr().String()
}

func RemoveElement(slice []jsonrpc.CID, elems jsonrpc.CID) []jsonrpc.CID {
	for i := range slice {
		if slice[i] == elems {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (self *Router) Init() *Router {
	self.serviceLock = new(sync.RWMutex)
	self.ServiceConnMap = make(map[string]([]jsonrpc.CID))
	self.ConnServiceMap = make(map[jsonrpc.CID]([]string))
	self.ConnMap = make(map[jsonrpc.CID](ConnT))
	self.PendingMap = make(map[PendingKey]PendingValue)
	return self
}

func (self *Router) registerConn(connId jsonrpc.CID, conn ConnT) {
	self.ConnMap[connId] = conn
	// register connId as a service name
}

func (self *Router) RegisterService(connId jsonrpc.CID, serviceName string) error {
	self.serviceLock.Lock()
	defer self.serviceLock.Unlock()

	// bi direction map
	cidArr, ok := self.ServiceConnMap[serviceName]
	if ok {
		cidArr = append(cidArr, connId)
	} else {
		var a []jsonrpc.CID
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

func (self *Router) UnRegisterService(connId jsonrpc.CID, serviceName string) error {
	self.serviceLock.Lock()
	defer self.serviceLock.Unlock()

	serviceNames, ok := self.ConnServiceMap[connId]
	if ok {
		var tmpServiceNames []string

		for _, sname := range serviceNames {
			if sname != serviceName {
				tmpServiceNames = append(tmpServiceNames, sname)
			}
		}
		if len(tmpServiceNames) > 0 {
			self.ConnServiceMap[connId] = tmpServiceNames
		} else {
			delete(self.ConnServiceMap, connId)
		}
	}

	connIds, ok := self.ServiceConnMap[serviceName]
	if ok {
		var tmpConnIds []jsonrpc.CID
		for _, cid := range connIds {
			if cid != connId {
				tmpConnIds = append(tmpConnIds, cid)
			}

			if len(tmpConnIds) > 0 {
				self.ServiceConnMap[serviceName] = tmpConnIds
			} else {
				delete(self.ServiceConnMap, serviceName)
			}
		}
	}


	ct, ok := self.ConnMap[connId]
	if ok {
		delete(self.ConnMap, connId)
		close(ct.RecvChannel())
	}
	return nil
}

func (self *Router) unregisterConn(connId jsonrpc.CID) {
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
		close(ct.RecvChannel())
	}
}

func (self *Router) SelectConn(serviceName string) (jsonrpc.CID, bool) {
	self.serviceLock.RLock()
	defer self.serviceLock.RUnlock()

	connIds, ok := self.ServiceConnMap[serviceName]
	if ok && len(connIds) > 0 {
		// or random or round-robin
		return connIds[0], true
	}
	return 0, false
}

func (self *Router) GetServices(connId jsonrpc.CID) []string {
	self.serviceLock.RLock()
	defer self.serviceLock.RUnlock()
	return self.ConnServiceMap[connId]
}

func (self *Router) ClearTimeoutRequests() {
	now := time.Now()
	tmpMap := make(map[PendingKey]PendingValue)

	for pKey, pValue := range self.PendingMap {
		if now.After(pValue.Expire) {
			errMsg := jsonrpc.NewErrorMessage(pKey.MsgId, 408, "request timeout")
			_ = self.deliverMessage(pKey.ConnId, errMsg)
		} else {
			tmpMap[pKey] = pValue
		}
	}
	self.PendingMap = tmpMap
}

func (self *Router) ClearPending(connId jsonrpc.CID) {
	for pKey, pValue := range self.PendingMap {
		if pKey.ConnId == connId || pValue.ConnId == connId {
			self.deletePending(pKey)
		}
	}
}

func (self *Router) deletePending(pKey PendingKey) {
	delete(self.PendingMap, pKey)
}

func (self *Router) setPending(pKey PendingKey, pValue PendingValue) {
	self.PendingMap[pKey] = pValue
}

func (self *Router) routeMessage(msg jsonrpc.RPCMessage) *ConnT {
	fromConnId := msg.FromConnId
	if msg.IsRequest() {
		toConnId, found := self.SelectConn(msg.ServiceName)
		if found {
			pKey := PendingKey{ConnId: fromConnId, MsgId: msg.Id}
			expireTime := time.Now().Add(DefaultRequestTimeout)
			pValue := PendingValue{ConnId: toConnId, Expire: expireTime}

			self.setPending(pKey, pValue)
			return self.deliverMessage(toConnId, msg)
		} /*else {
			errMsg := jsonrpc.NewErrorMessage(msg.Id, 404, "service not found")
			return self.deliverMessage(fromConnId, errMsg)
		}*/
	} else if msg.IsNotify() {
		toConnId, found := self.SelectConn(msg.ServiceName)
		if found {
			return self.deliverMessage(toConnId, msg)
		}
		/* else {
			errMsg := jsonrpc.NewErrorMessage(msg.Id, 404, "service not found")
			return self.deliverMessage(fromConnId, errMsg)
		} */
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

func (self *Router) broadcastNotify(notify jsonrpc.RPCMessage) (int, error) {
	if !notify.IsNotify() {
		/*errMsg := jsonrpc.NewErrorMessage(notify.Id, 400, "only notify can be broadcasted")
		self.deliverMessage(notify.FromConnId, errMsg)
		return nil */
		return 0, ErrNotNotify
	}
	cntDeliver := 0
	for connId, conn := range self.ConnMap {
		if conn.CanBroadcast() { // == IntentLocal {
			self.deliverMessage(connId, notify)
			cntDeliver += 1
		}
	}
	return cntDeliver, nil
}

func (self *Router) deliverMessage(connId jsonrpc.CID, msg jsonrpc.RPCMessage) *ConnT {
	ct, ok := self.ConnMap[connId]
	if ok {
		ct.RecvChannel() <- msg
		return &ct
	}
	return nil
}

//func (self *Router) Start() {
/*	for {
		select {
		case cmdOpen := <-self.ChJoin:
			//self.registerConn(cmdOpen.ConnId, cmdOpen.Channel, cmdOpen.Intent)
		case msg := <-self.ChMsg:
			self.routeMessage(msg)
		case notify := <-self.ChBroadcast:
			self.broadcastNotify(notify)
		case cmdClose := <-self.ChLeave:
			self.unregisterConn(jsonrpc.CID(cmdClose))
		}
	} */
//}

// commands
func (self *Router) RouteMessage(msg jsonrpc.RPCMessage, fromConnId jsonrpc.CID) *ConnT {
	self.serviceLock.RLock()
	defer self.serviceLock.RUnlock()
	
	msg.FromConnId = fromConnId
	//self.ChMsg <- msg
	return self.routeMessage(msg)
}

func (self *Router) BroadcastNotify(notify jsonrpc.RPCMessage, fromConnId jsonrpc.CID) (int, error) {
	self.serviceLock.RLock()
	defer self.serviceLock.RUnlock()
	
	notify.FromConnId = fromConnId
	//self.ChBroadcast <- notify
	return self.broadcastNotify(notify)
}

/*func (self *Router) Join(connId jsonrpc.CID, ch MsgChannel, intent string) {
	conn := &ConnT{RecvChannel: ch, Intent: intent}
	//self.registerConn(cmdOpen.ConnId, conn)
	self.JoinConn(connId, conn)
}*/

func (self *Router) JoinConn(connId jsonrpc.CID, conn ConnT) {
	self.serviceLock.Lock()
	defer self.serviceLock.Unlock()
	self.registerConn(connId, conn)
}

func (self *Router) Leave(connId jsonrpc.CID) {
	//self.ChLeave <- LeaveCommand(connId)
	self.serviceLock.Lock()
	defer self.serviceLock.Unlock()
	self.unregisterConn(connId)
}
