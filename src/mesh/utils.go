package mesh

import (
	"sync/atomic"
)

var counter uint64 = 10000

func GetUID() uint64 {
	return atomic.AddUint64(&counter, 1)
}

func GetCID() CID {
	return CID(GetUID())
}
