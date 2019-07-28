package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type Msg struct {
	Topic   string
	Content string
	After   func(error)
}

type ProducerPool struct {
	C chan Msg // never close

	Addr   string
	Config *nsq.Config
	Size   int
}

func New(addr string, tls bool, size int) *ProducerPool {
	cfg := nsq.NewConfig()
	cfg.TlsV1 = tls

	p := &ProducerPool{
		C:      make(chan Msg, 100),
		Addr:   addr,
		Config: cfg,
		Size:   size,
	}
	go p.loop()
	return p
}

func (ppool *ProducerPool) loop() {
	if ppool.Size < 1 {
		ppool.Size = 1
	}

	start_ch := make(chan bool, ppool.Size)
	for {
		start_ch <- true
		go func() {
			defer func() {
				<-start_ch
			}()
			p, err := nsq.NewProducer(ppool.Addr, ppool.Config)
			if err != nil {
				return
			}
			for m := range ppool.C {
				err = p.Publish(m.Topic, []byte(m.Content))
				if m.After != nil {
					m.After(err)
				}
				if err != nil {
					return
				}
			}
		}()
	}
}

func client() {
	start := time.Now()

	p := New("ws.mofon.top:8911", true, 10)
	ch := p.C

	data, err := json.Marshal(HTTPReq{
		Topic:  "http_resp",
		Method: "GET",
		URL:    "https://httpbin.org/get",
	})
	if err != nil {
		panic(err)
	}

	wg := new(sync.WaitGroup)
	msg := Msg{
		Topic:   "http",
		Content: string(data),
		After: func(err error) {
			wg.Done()
		},
	}

	for i := 0; i < 1; i++ {
		wg.Add(1)
		ch <- msg
	}
	wg.Wait()
	end := time.Now()
	fmt.Println("used time:", end.Sub(start))
}
