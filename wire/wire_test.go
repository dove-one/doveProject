package wire

import (
	"fmt"
	"testing"
)

type Message string

type Channel struct {
	Message Message
}

type BroadCast struct {
	Channel Channel
}

func NewMessage() Message {
	return Message("Hello wire!")
}

func NewChannel(m Message) Channel {
	return Channel{
		Message: m,
	}
}

func NewBroadCast(c Channel) *BroadCast {
	return &BroadCast{
		Channel: c,
	}
}

func (c *Channel) GetMsg() Message {
	return c.Message
}

func (b *BroadCast) Start() {
	msg := b.Channel.GetMsg()
	fmt.Println(msg)
}

// 不使用 wire
func TestNoWire(t *testing.T) {
	message := NewMessage()
	channel := NewChannel(message)
	broadCast := NewBroadCast(channel)
	broadCast.Start()
}

// 使用 wire
func TestWire(t *testing.T) {
	b := InitializeBroadCast()
	b.Start()
}
