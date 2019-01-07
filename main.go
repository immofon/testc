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
	fmt.Println("1")
	c := ebus.NewClient("ws://39.105.42.45:8100/", func(e ebus.Event) {
	})
	fmt.Println("2")

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		fmt.Println("evemt")
		c.Emit(ebus.Event{
			To:    "@record",
			Topic: "set",
			Data:  []string{"time", time.Now().String()},
		})
	}
}
