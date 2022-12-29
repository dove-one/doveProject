package main

import "fmt"

type Message string

type Greeter struct {
	Message Message
}

type Event struct {
	Greeter Greeter
}

func NewMessage(msg string) Message {
	return Message(msg)
}

func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}

func (g Greeter) Greet() Message {
	return g.Message
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

func main() {
	event := InitializeEvent("hello_world")
	event.Start()
}
