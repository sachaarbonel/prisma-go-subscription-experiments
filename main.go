package main

import (
	"log"
	"net/http"

	websocket "github.com/gorilla/websocket"
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

type WSClient struct {
	conn *websocket.Conn
}

func New() *WSClient {
	return &WSClient{}
}

func (c *WSClient) Connect(url string) error {

	//do better than strings
	headers := http.Header{
		"Host":                     []string{"localhost:4466"},
		"User-Agent":               []string{"Mozilla/5.0 (X11; Linux x86_64; rv:62.0) Gecko/20100101 Firefox/62.0"},
		"Accept":                   []string{"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"Accept-Language":          []string{"en-US,en;q=0.5"},
		"Accept-Encoding":          []string{"gzip, deflate"},
		"Sec-WebSocket-Version":    []string{"13"},
		"Origin":                   []string{"http://localhost:4466"},
		"Sec-WebSocket-Protocol":   []string{"graphql-ws"},
		"Sec-WebSocket-Extensions": []string{"permessage-deflate"},
		"Sec-WebSocket-Key":        []string{"/U1W39A3aqbOjTV7yFUyhg=="},
		"Pragma":                   []string{"no-cache"},
		"Cache-Control":            []string{"no-cache"},
	}
	conn, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		return err
	}
	c.conn = conn
	//go c.read()
	return nil
}

func main() {
	client := New()
	err := client.Connect("ws://localhost:4466")
	if err != nil {
		log.Println(err)
	}

}
