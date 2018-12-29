package bbox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"datadir"
	"go.etcd.io/etcd/clientv3"
	"mesh"
)
	
func GetTicket() Ticket {
	config := datadir.GetConfig()
	return Ticket{config.Bbox.Prefix, config.Bbox.Etcd}
}

func (self *BboxClient) Init() *BboxClient {
	self.ticket = GetTicket()
	self.etcdClient = clientv3.New(clientv3.Config{
		Endpoints: self.ticket.Etcd,
		DialTimeout: 5 * time.Second})
}

func (self *BboxClient) EtcdPath(key string) string {
	return fmt.Sprintf("/%s/%s", self.ticket.Prefix, key)
}

func (self *BboxClient) GetBoxes() ([]BoxInfo, error) {
	kv := clientv3.NewKv(self.etcdClient)
	resp, err := kv.Get(context.TODO(), self.EtcdPath("boxes"), clientv3.WithPrefix())
	var arr []BoxInfo
	for _, kv := range resp.Kvs() {
		var boxInfo BoxInfo
		err := json.Unmarshal(kv.Value, &boxInfo)
		if err != nil {
			return nil, err
		}
		arr = append(arr, boxInfo)
	}
	return arr, nil
}

func (self *BboxClient) WatchBoxes() error {
	boxInfos, err := self.GetBoxes()
	if err != nil {
		return err
	}
	self.JoinBoxes(boxInfos)	
	chBox := clientv3.Watch(self.EtcdPath("boxes"))
	for true {
		resp, done := <- chBox
		if done {
			break
		}
		boxInfos, err := self.GetBoxes()
		if err != nil {
			return err
		}
		self.JoinBoxes(boxInfos)
	}
	return nil
}

func (self BboxClient) JoinBoxes(boxInfos []BoxInfo) {
	m := mesh.GetMesh()
	for _, boxInfo := range boxInfos {
		var endpoint = BboxEndpoint{BoxInfo: boxInfo}
		m.Update(&endpoint)
	}
}
