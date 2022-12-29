package main

import (
	"fmt"
	"sync"
	"testing"
)

var memoryMu sync.Mutex
var value int

func TestWire(t *testing.T) {
	go func() {
		memoryMu.Lock()
		value++
		memoryMu.Unlock()
	}()
	memoryMu.Lock()
	if value == 0 {
		fmt.Println(0, value)
	} else {
		fmt.Println(1, value)
	}
	memoryMu.Unlock()

}
