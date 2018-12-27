package mesh

import (
	"sync"
	"errors"
)

var _mesh *Mesh = nil

func GetMesh() *Mesh {
	if _mesh == nil {
		_mesh = new(Mesh).Init()
	}
	return _mesh
}

func (self *Mesh) Init() *Mesh {
	self.serviceLock = new(sync.RWMutex)
	self.idEndpointMap = make(map[string](*Endpoint))
	self.serviceEndpointsMap = make(map[string]([]*Endpoint))
	return self
}


func (self * Mesh) Join(endpoint *Endpoint) error {
	self.serviceLock.Lock()
	defer self.serviceLock.Unlock()

	epId := (*endpoint).GetId()
	_, ok := self.idEndpointMap[epId]
	if ok {
		return errors.New("endpoint already exist")
	}

	self.idEndpointMap[epId] = endpoint

	for _, serviceName := range (*endpoint).GetServiceNames() {
		arr, ok := self.serviceEndpointsMap[serviceName]
		if ok {
			arr = append(arr, endpoint)
		} else {
			var emptyArr [](*Endpoint)
			arr = append(emptyArr, endpoint)
		}
		self.serviceEndpointsMap[serviceName] = arr
	}
	return nil
}

func (self *Mesh) Leave(endpointId string) (*Endpoint) {
	self.serviceLock.Lock()
	defer self.serviceLock.Unlock()

	endpoint, ok := self.idEndpointMap[endpointId]
	if !ok {
		return nil
	}

	for _, serviceName := range (*endpoint).GetServiceNames() {
		arr, ok := self.serviceEndpointsMap[serviceName]
		if !ok {
			// FIXME: it almost not happend
			continue
		}
		var newArr [](*Endpoint)
		for _, ep := range arr {
			epId := (*ep).GetId()
			if epId != endpointId {
				newArr = append(newArr, ep)
			}
		}
		if len(newArr) > 0 {
			self.serviceEndpointsMap[serviceName] = newArr
		} else {
			delete(self.serviceEndpointsMap, serviceName)
		}
	}
	delete(self.idEndpointMap, endpointId)
	return endpoint
}

func (self *Mesh) GetEndpoint(serviceName string) *Endpoint {
	self.serviceLock.RLock()
	defer self.serviceLock.RUnlock()

	arr, ok := self.serviceEndpointsMap[serviceName]
	if ok && len(arr) > 0 {
		return arr[0]
	}  else {
		return nil
	}
}

func (self *Mesh) GetEndpoints(serviceName string) ([](*Endpoint), bool) {
	self.serviceLock.RLock()
	defer self.serviceLock.RUnlock()
	
	arr, ok := self.serviceEndpointsMap[serviceName]
	return arr, ok
}
