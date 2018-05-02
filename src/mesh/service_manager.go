package mesh

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

func (self *ServiceManager) handleMessage(msg RPCMessage) {
	context = Context()
	if msg.Method == "register" {
		params, err := msg.Params.Array()
		if err != nil {
			errMsg := NewErrorMessage(msg.Id, 400, "params must be array")
			context.Router.RouteMessage(errMsg, self.ConnId)
			return
		}
		var serviceNames []string

		for _, v := range params {
			serviceName, ok := v.(string)
			if !ok {
				errMsg := NewErrorMessage(msg.Id, 400, "service name must be string")
				context.Router.RouteMessage(errMsg, self.ConnId)
				return
			}
			serviceNames = append(serviceNames, serviceName)
		}

		for _, serviceName := range serviceNames {
			context.Router.RegisterService(msg.FromConnId, serviceName)
		}
		result := NewResultMessage(msg.Id, "ok")
		context.Router.RouteMessage(result, self.ConnId)
	} else if msg.Method == "getServices" {
		serviceNames := context.Router.GetServices(self.ConnId)
		result := NewResultMessage(msg.Id, serviceNames)
		context.Router.RouteMessage(result, self.ConnId)
	} else if msg.Method == "getId" {
		result := NewResultMessage(msg.Id, msg.FromConnId)
		context.Router.RouteMessage(result, self.ConnId)
	} else if msg.Method == "ping" {
		result := NewResultMessage(msg.Id, "pong")
		context.Router.RouteMessage(result, self.ConnId)
	}
}
