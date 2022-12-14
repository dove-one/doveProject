package v2

import (
	"errors"
	"sync"
	"time"
)

var (
	//ErrMaxActiveConnReached 连接池超限
	ErrMaxActiveConnReached = errors.New("MaxActiveConnReached")
)

// Config 连接池相关配置
type Config struct {
	initialCap  int
	MaxCap      int
	MaxIdle     int
	Factory     func() (any, error)
	Close       func(any) error
	Ping        func(any) error
	IdleTimeout time.Duration
}

type connReq struct {
	idleConn *idleConn
}

// channelPool 实现了基于缓冲通道的 Pool 接口
type channelPool struct {
	mu           sync.RWMutex
	conns        chan *idleConn
	factory      func() (any, error)
	close        func(any) error
	ping         func(any) error
	idleTimeout  time.Duration
	waitTimeout  time.Duration
	maxActive    int
	openingConns int
	connReqs     []chan connReq
}

type idleConn struct {
	conn any
	t    time.Time
}

func NewChannelPool(poolConfig *Config) (Pool, error) {
	//TODO
	panic("implement me")
}

func (c *channelPool) getConns() chan *idleConn {
	c.mu.Lock()
	defer c.mu.Unlock()
	conns := c.conns
	return conns
}
func (c *channelPool) Get() (any, error) {
	//TODO implement me
	panic("implement me")
}

func (c *channelPool) Put(a any) error {
	//TODO implement me
	panic("implement me")
}

func (c *channelPool) Close(conn any) error {
	if conn == nil {
		return errors.New("connection is nil.rejecting")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.close == nil {
		return nil
	}
	c.openingConns--
	return c.close(conn)
}

func (c *channelPool) Ping(conn any) error {
	if conn == nil {
		return errors.New("connection is nil.rejecting")
	}
	return c.ping(conn)
}

func (c *channelPool) Len() int {
	return len(c.getConns())
}

func (c *channelPool) Release() {
	//TODO implement me
	panic("implement me")
}
