package mesh

var (
	context *ContextT
)

func Context() *ContextT {
	if context == nil {
		context = new(ContextT).Init()
	}
	return context
}

func (self *ContextT) Init() *ContextT {
	self.Router = NewRouter()
	self.ServiceManager = new(ServiceManager).Init()
	return self
}

func (self *ContextT) Start() {
	go self.Router.Start()
	go self.ServiceManager.Start()
}
