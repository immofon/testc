package main

import (
	"fmt"
	"time"

	"github.com/immofon/ebus"
)

func main() {
	fmt.Println("hello ci")

	defer func() {
		recover()
		fmt.Println("err")
	}()
	c := ebus.NewClient("ws://39.105.42.45:8100/", func(e ebus.Event) {
	})

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		c.Emit(ebus.Event{
			To:    "@record",
			Topic: "set",
			Data:  []string{"time", time.Now().String()},
		})
	}
}
