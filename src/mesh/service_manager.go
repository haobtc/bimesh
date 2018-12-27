package mesh

import "jsonrpc"

// builtin services manager

func (self *ServiceManager) Init() *ServiceManager {
	self.ChMsg = make(MsgChannel, 100)
	//self.ConnId = CID(uuid.Must(uuid.NewV4()).String())
	self.ConnId = GetCID()
	return self
}

func (self *ServiceManager) Close() {
	Context().Router.Leave(self.ConnId)
}

func (self *ServiceManager) Start() {
	Context().Router.Join(self.ConnId, self.ChMsg, "builtin")
	Context().Router.RegisterService(self.ConnId, "builtin-services")

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
	context = Context()
	params, err := msg.Params.Array()
	if err != nil {
		errMsg := jsonrpc.NewErrorMessage(msg.Id, 400, "params must be array")
		context.Router.RouteMessage(errMsg, self.ConnId)
		return
	}
	var serviceNames []string

	for _, v := range params {
		serviceName, ok := v.(string)
		if !ok {
			errMsg := jsonrpc.NewErrorMessage(msg.Id, 400, "service name must be string")
			context.Router.RouteMessage(errMsg, self.ConnId)
			return
		}
		serviceNames = append(serviceNames, serviceName)
	}

	for _, serviceName := range serviceNames {
		context.Router.RegisterService(msg.FromConnId, serviceName)
	}
	result := jsonrpc.NewResultMessage(msg.Id, "ok")
	context.Router.RouteMessage(result, self.ConnId)
}

func (self *ServiceManager) unregisterServices(msg jsonrpc.RPCMessage) {
	context = Context()
	params, err := msg.Params.Array()
	if err != nil {
		errMsg := jsonrpc.NewErrorMessage(msg.Id, 400, "params must be array")
		context.Router.RouteMessage(errMsg, self.ConnId)
		return
	}
	var serviceNames []string

	for _, v := range params {
		serviceName, ok := v.(string)
		if !ok {
			errMsg := jsonrpc.NewErrorMessage(msg.Id, 400, "service name must be string")
			context.Router.RouteMessage(errMsg, self.ConnId)
			return
		}
		serviceNames = append(serviceNames, serviceName)
	}

	for _, serviceName := range serviceNames {
		context.Router.UnRegisterService(msg.FromConnId, serviceName)
	}
	result := jsonrpc.NewResultMessage(msg.Id, "ok")
	context.Router.RouteMessage(result, self.ConnId)
}

func (self *ServiceManager) handleMessage(msg jsonrpc.RPCMessage) {
	switch msg.Method {
	case "register":
		self.registerServices(msg)
	case "unregister":
		self.unregisterServices(msg)
	case "getServices":
		serviceNames := Context().Router.GetServices(self.ConnId)
		result := jsonrpc.NewResultMessage(msg.Id, serviceNames)
		Context().Router.RouteMessage(result, self.ConnId)
	case "getId":
		result := jsonrpc.NewResultMessage(msg.Id, msg.FromConnId)
		Context().Router.RouteMessage(result, self.ConnId)
	case "ping":
		result := jsonrpc.NewResultMessage(msg.Id, "pong")
		Context().Router.RouteMessage(result, self.ConnId)
	default:
		errorMsg := jsonrpc.NewErrorMessage(msg.Id, 404, "method not found")
		Context().Router.RouteMessage(errorMsg, self.ConnId)
	}
}
