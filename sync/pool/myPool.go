package pool

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type MyPool struct {
	p      sync.Pool
	macCnt int32
	cnt    int32
}

func (p *MyPool) Get() any {
	return p.p.Get()
}

func (p *MyPool) Put(val any) {
	// 大对象不放回去
	if unsafe.Sizeof(val) > 1024 {
		return
	}

	// 超过数量限制
	cnt := atomic.AddInt32(&p.cnt, 1)
	if cnt >= p.macCnt {
		atomic.AddInt32(&p.cnt, -1)
		return
	}

	p.p.Put(val)
}
