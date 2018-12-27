package tentacle

import (
	"sync/atomic"
	"jsonrpc"
)

var counter uint64 = 10000

func GetUID() uint64 {
	return atomic.AddUint64(&counter, 1)
}

func GetCID() jsonrpc.CID {
	return jsonrpc.CID(GetUID())
}
