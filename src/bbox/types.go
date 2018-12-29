package bbox

import (
	"go.etcd.io/etcd/clientv3"
)


type Ticket struct {
	Prefix string
	Etcd []string
}

type BboxClient struct {
	ticket Ticket
	etcdClient *clientv3.Client
}

type BoxInfo struct {
	BoxId string `json:"boxid"`
	Bind string `json:"bind"`
	Ssl string `json:"ssl"`
	ServiceNames []string `json:"services"`
}

type BboxEndpoint struct {
	BoxInfo BoxInfo
}
