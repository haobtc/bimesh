package tentacle

import (
	"jsonrpc"
)

// core services manager
func (self *ServiceManager) Init() *ServiceManager {
	self.ChMsg = make(MsgChannel, 100)
	self.ConnId = GetCID()
	return self
}

func (self *ServiceManager) Close() {
	Tentacle().Router.Leave(self.ConnId)
}

func (self *ServiceManager) Start() {
	go self.Run()
}

func (self *ServiceManager) RecvChannel() MsgChannel {
	return self.ChMsg
}

func (self *ServiceManager) CanBroadcast() bool {
	return false
}

func (self *ServiceManager) Run() {
	//Tentacle().Router.Join(self.ConnId, self.ChMsg, "builtin")
	Tentacle().Router.JoinConn(self.ConnId, self)
	Tentacle().Router.RegisterService(self.ConnId, "core.services")

	for {
		select {
		case msg, more := <-self.ChMsg:
			if more {
				self.handleMessage(msg)
			} else {
				// log.
				return
			}
		}
	}
}

func (self *ServiceManager) registerServices(msg jsonrpc.RPCMessage) {
	tentacle = Tentacle()
	params, err := msg.Params.Array()
	if err != nil {
		errMsg := jsonrpc.NewErrorMessage(msg.Id, 400, "params must be array")
		tentacle.Router.RouteMessage(errMsg, self.ConnId)
		return
	}
	var serviceNames []string

	for _, v := range params {
		serviceName, ok := v.(string)
		if !ok {
			errMsg := jsonrpc.NewErrorMessage(msg.Id, 400, "service name must be string")
			tentacle.Router.RouteMessage(errMsg, self.ConnId)
			return
		}
		serviceNames = append(serviceNames, serviceName)
	}

	for _, serviceName := range serviceNames {
		tentacle.Router.RegisterService(msg.FromConnId, serviceName)
	}
	result := jsonrpc.NewResultMessage(msg.Id, "ok")
	tentacle.Router.RouteMessage(result, self.ConnId)
}

func (self *ServiceManager) unregisterServices(msg jsonrpc.RPCMessage) {
	tentacle = Tentacle()
	params, err := msg.Params.Array()
	if err != nil {
		errMsg := jsonrpc.NewErrorMessage(msg.Id, 400, "params must be array")
		tentacle.Router.RouteMessage(errMsg, self.ConnId)
		return
	}
	var serviceNames []string

	for _, v := range params {
		serviceName, ok := v.(string)
		if !ok {
			errMsg := jsonrpc.NewErrorMessage(msg.Id, 400, "service name must be string")
			tentacle.Router.RouteMessage(errMsg, self.ConnId)
			return
		}
		serviceNames = append(serviceNames, serviceName)
	}

	for _, serviceName := range serviceNames {
		tentacle.Router.UnRegisterService(msg.FromConnId, serviceName)
	}
	result := jsonrpc.NewResultMessage(msg.Id, "ok")
	tentacle.Router.RouteMessage(result, self.ConnId)
}

func (self *ServiceManager) handleMessage(msg jsonrpc.RPCMessage) {
	switch msg.Method {
	case "register":
		self.registerServices(msg)
	case "unregister":
		self.unregisterServices(msg)
	case "getServices":
		serviceNames := Tentacle().Router.GetServices(self.ConnId)
		result := jsonrpc.NewResultMessage(msg.Id, serviceNames)
		Tentacle().Router.RouteMessage(result, self.ConnId)
	case "getId":
		result := jsonrpc.NewResultMessage(msg.Id, msg.FromConnId)
		Tentacle().Router.RouteMessage(result, self.ConnId)
	case "ping":
		result := jsonrpc.NewResultMessage(msg.Id, "pong")
		Tentacle().Router.RouteMessage(result, self.ConnId)
	default:
		errorMsg := jsonrpc.NewErrorMessage(msg.Id, 404, "method not found")
		Tentacle().Router.RouteMessage(errorMsg, self.ConnId)
	}
}
