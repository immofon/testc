package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ccding/go-stun/stun"
	"github.com/immofon/ebus"
)

func main() {
	nat, host, err := stun.NewClient().Discover()
	fmt.Println(nat, host, err)

	fmt.Println("hello ci")

	defer func() {
		recover()
		fmt.Println("err")
	}()
	var c *ebus.Client
	c = ebus.NewClient("ws://39.105.42.45:8100/", func(e ebus.Event) {
		if strings.HasPrefix(e.From, "&") {
			url := e.Topic
			content := get(url)

			c.Emit(ebus.Event{
				To:    e.From,
				Topic: e.Topic,
				Data:  []string{content},
			})
		}
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

func get(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return err.Error()
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	return string(data)
}
