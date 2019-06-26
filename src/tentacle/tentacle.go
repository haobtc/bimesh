package tentacle

var (
	tentacle *TentacleT
)

func Tentacle() *TentacleT {
	if tentacle == nil {
		tentacle = new(TentacleT).Init()
	}
	return tentacle
}

func (self *TentacleT) Init() *TentacleT {
	self.Router = NewRouter()
	self.ServiceManager = new(ServiceManager).Init()
	return self
}

func (self *TentacleT) Start() {
	//go self.Router.Start()
	self.ServiceManager.Start()
}

