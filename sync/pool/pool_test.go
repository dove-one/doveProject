package pool

import (
	"fmt"
	"sync"
	"testing"
)

type MyCache struct {
	pool sync.Pool
}

func NewMyCache() *MyCache {
	return &MyCache{
		pool: sync.Pool{
			New: func() any {
				fmt.Println("hello")
				return []byte{}
			}}}
}

func TestPool(t *testing.T) {
	cache := NewMyCache()
	val := cache.pool.Get()
	//
	cache.pool.Put(val)
}

type User struct {
	ID   uint64
	Name string
}

func (u *User) Reset() {
	u.Name = ""
	u.ID = 0
}

func TestPool2(t *testing.T) {
	pool := sync.Pool{
		New: func() any {
			return &User{}
		},
	}
	u1 := pool.Get().(*User)
	u1.ID = 1
	u1.Name = "doveOne"
	// 重置掉后放回去 例-beego中context的重置
	u1.Reset()
	pool.Put(u1)

	u2 := pool.Get().(*User)
	fmt.Println(u2)
}
