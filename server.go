package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type HTTPReq struct {
	Topic string `json:"topic"` // MUST
	ID    string `json:"id"`

	Method  string            `json:"method"` // MUST
	URL     string            `json:"url"`    // MUST
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
}

type HTTPResp struct {
	ID string

	Method  string
	URL     string
	Headers map[string]string
	Status  int
	Body    string
}

func server() {
	p := New("ws.mofon.top:8911", true, 10)

	cfg := nsq.NewConfig()
	cfg.TlsV1 = true
	consumer, err := nsq.NewConsumer("http", "process", cfg)
	if err != nil {
		panic(err)
	}

	consumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) (err error) {
		defer func() {
			if err != nil {
				p.C <- Msg{Topic: "log/error/http", Content: err.Error()}
			}
			err = nil
		}()

		var r HTTPReq
		err = json.Unmarshal(msg.Body, &r)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(r.Method, r.URL, strings.NewReader(r.Body))
		if err != nil {
			return err
		}
		for k, v := range r.Headers {
			req.Header.Add(k, v)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var ret HTTPResp
		ret.Method = r.Method
		ret.ID = r.ID
		ret.URL = r.URL
		ret.Headers = make(map[string]string)
		ret.Status = resp.StatusCode

		for k, _ := range resp.Header {
			ret.Headers[k] = resp.Header.Get(k)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		ret.Body = string(data)

		raw, err := json.Marshal(ret)
		if err != nil {
			return err
		}

		p.C <- Msg{
			Topic:   r.Topic,
			Content: string(raw),
		}

		return nil
	}))

	//err = consumer.ConnectToNSQD("localhost:4150")
	err = consumer.ConnectToNSQD("ws.mofon.top:8911")
	if err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Hour)
	<-consumer.StopChan
}
