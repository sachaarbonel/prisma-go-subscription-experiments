package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type PayloadStruct struct {
	Query string `json:"query"`
}

type ResponseStruct struct {
	Data struct {
		Users []struct {
			ID string `json:"id"`
		} `json:"users"`
	} `json:"data"`
}

// Call represents an active request
type Call struct {
	Req   PayloadStruct
	Res   ResponseStruct
	Done  chan bool
	Error error
}

func NewCall(req PayloadStruct) *Call {
	done := make(chan bool)
	return &Call{
		Req:  req,
		Done: done,
	}
}

type WSClient struct {
	conn *websocket.Conn
}

func New() *WSClient {
	return &WSClient{}
}

func (c *WSClient) read() {
	for {
		var res ResponseStruct
		err := c.conn.ReadJSON(&res)
		log.Printf("Receive : %v", res)
		if err != nil {
			log.Printf("error : %v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error : %v", err)
			}
			break
		}
	}
}

func (c *WSClient) Connect(url string) error {
	conn, _, err := websocket.DefaultDialer.Dial(url, http.Header{
		"User-Agent": []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36"},
	})
	if err != nil {
		return err
	}
	c.conn = conn
	go c.read()
	return nil
}

func (c *WSClient) Request(payload interface{}) (interface{}, error) {

	req := payload.(PayloadStruct)
	call := NewCall(req)

	err := c.conn.WriteJSON(&req)
	if err != nil {
		return nil, err
	}

	select {
	case <-call.Done:
	case <-time.After(2 * time.Second):
		call.Error = errors.New("request timeout")
	}

	if call.Error != nil {
		return nil, call.Error
	}
	return call.Res.Data, nil
}

func (c *WSClient) Close() error {
	return c.conn.Close()
}

func main() {
	client := New()
	err := client.Connect("ws://localhost:4466/")
	if err != nil {
		panic(err)
	}

	go func() {
		want := PayloadStruct{
			Query: "query {users {id }}",
		}
		_, err := client.Request(want)
		if err != nil {
			fmt.Println(err)
		}
	}()

	defer func() {
		err = client.Close()
		if err != nil {
			panic(err)
		}
	}()
}
