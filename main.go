package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Session struct {
	ws      *websocket.Conn
	errChan chan error
}

const (
	connectionInitMsg      = "connection_init"      // Client -> Server
	connectionTerminateMsg = "connection_terminate" // Client -> Server
	startMsg               = "start"                // Client -> Server
	stopMsg                = "stop"                 // Client -> Server
	connectionAckMsg       = "connection_ack"       // Server -> Client
	connectionErrorMsg     = "connection_error"     // Server -> Client
	dataMsg                = "data"                 // Server -> Client
	errorMsg               = "error"                // Server -> Client
	completeMsg            = "complete"             // Server -> Client
	//connectionKeepAliveMsg = "ka"                 // Server -> Client  TODO: keepalives
)

type operationMessage struct {
	Payload json.RawMessage `json:"payload,omitempty"`
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
}

type PostSubscriptionResponse struct {
	Data struct {
		Post struct {
			Node struct {
				ID    string `json:"id"`
				Title string `json:"title"`
			} `json:"node"`
		} `json:"post"`
	} `json:"data"`
}

func wsConnect(url string) *websocket.Conn {
	headers := make(http.Header)
	headers.Add("Sec-Websocket-Protocol", "graphql-ws")
	c, _, err := websocket.DefaultDialer.Dial(url, headers)

	if err != nil {
		panic(err)
	}
	return c
}

func (s *Session) ReadOp() (operationMessage, error) {
	var msg operationMessage
	err := s.ws.ReadJSON(&msg)
	if err != nil {
		panic(err)
	}
	return msg, err
}

func (s *Session) Subscribe(query string) (<-chan string, <-chan error) {

	channel := make(chan string)

	s.ws.WriteJSON(&operationMessage{
		Type:    startMsg,
		ID:      "test_1", // Do I need to generate a random ID here
		Payload: json.RawMessage(query),
	})

	go func() {
		for {

			msg, err := s.ReadOp()
			if err != nil {
				s.errChan <- err
			}
			rawPayload := json.RawMessage(msg.Payload)
			strPayload := string(rawPayload[:])
			channel <- strPayload

		}
		close(channel)
		close(s.errChan)
	}()

	return channel, s.errChan
}

func main() {

	c := wsConnect("ws://localhost:4466")
	defer c.Close()

	session := &Session{
		ws: c,
	}
	session.ws.WriteJSON(&operationMessage{Type: connectionInitMsg})
	msg, _ := session.ReadOp()
	log.Println(msg.Type)

	query := string(`{"query": "subscription { post { node { id title } } }"}`)
	subscription_future, _ := session.Subscribe(query)
	for subscription := range subscription_future {
		fmt.Println(subscription)
	}

}
