package v2

import (
	"net/rpc"
)

var (
	InitialCap = 5
	MaxIdleCap = 10
	MaximumCap = 30
	network    = "tcp"
	address    = "127.0.0.1:7777"
	factory    = func() (interface{}, error) {
		return rpc.DialHTTP("tcp", address)
	}
	closeFac = func(v interface{}) error {
		nc := v.(*rpc.Client)
		return nc.Close()
	}
)
